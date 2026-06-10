package handlers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"task-manager-api/internal/auth"
	"task-manager-api/internal/handlers"
	"task-manager-api/internal/models"
	"task-manager-api/internal/repository"
	"task-manager-api/internal/routes"
	"task-manager-api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func newTestServer(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory database: %v", err)
	}
	if err := db.AutoMigrate(&models.User{}, &models.Task{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	jwtManager := auth.NewJWTManager("test-secret", time.Hour)
	authService := services.NewAuthService(repository.NewUserRepository(db), jwtManager)
	taskService := services.NewTaskService(repository.NewTaskRepository(db))

	r := gin.New()
	routes.RegisterRoutes(
		r,
		handlers.NewAuthHandler(authService, jwtManager, false),
		handlers.NewTaskHandler(taskService),
		jwtManager,
	)
	return r
}

func doJSON(t *testing.T, r *gin.Engine, method, path, token, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func signup(t *testing.T, r *gin.Engine, email string) string {
	t.Helper()
	body := fmt.Sprintf(`{"name":"Test User","email":%q,"password":"password123"}`, email)
	w := doJSON(t, r, http.MethodPost, "/auth/signup", "", body)
	if w.Code != http.StatusCreated {
		t.Fatalf("signup failed with status %d: %s", w.Code, w.Body.String())
	}
	var resp struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode signup response: %v", err)
	}
	return resp.Data.Token
}

func TestSignupLoginAndDuplicateEmail(t *testing.T) {
	r := newTestServer(t)
	signup(t, r, "alice@example.com")

	// Duplicate signup must be rejected.
	w := doJSON(t, r, http.MethodPost, "/auth/signup", "",
		`{"name":"Alice Again","email":"alice@example.com","password":"password123"}`)
	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409 for duplicate email, got %d: %s", w.Code, w.Body.String())
	}

	// Login with correct credentials succeeds.
	w = doJSON(t, r, http.MethodPost, "/auth/login", "",
		`{"email":"alice@example.com","password":"password123"}`)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for valid login, got %d: %s", w.Code, w.Body.String())
	}

	// Login with a wrong password fails.
	w = doJSON(t, r, http.MethodPost, "/auth/login", "",
		`{"email":"alice@example.com","password":"wrong-password"}`)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for wrong password, got %d: %s", w.Code, w.Body.String())
	}
}

func TestTaskRoutesRequireAuth(t *testing.T) {
	r := newTestServer(t)

	for _, tc := range []struct{ method, path string }{
		{http.MethodPost, "/tasks"},
		{http.MethodGet, "/tasks"},
		{http.MethodGet, "/tasks/1"},
		{http.MethodPatch, "/tasks/1"},
		{http.MethodDelete, "/tasks/1"},
	} {
		w := doJSON(t, r, tc.method, tc.path, "", `{}`)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("%s %s: expected 401 without token, got %d", tc.method, tc.path, w.Code)
		}
	}
}

func TestCreateTaskValidation(t *testing.T) {
	r := newTestServer(t)
	token := signup(t, r, "bob@example.com")

	w := doJSON(t, r, http.MethodPost, "/tasks", token, `{"title":"","status":"bogus","priority":"urgent"}`)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid payload, got %d: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Error struct {
			Code   string            `json:"code"`
			Fields map[string]string `json:"fields"`
		} `json:"error"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}
	if resp.Error.Code != "VALIDATION_ERROR" {
		t.Fatalf("expected VALIDATION_ERROR code, got %q", resp.Error.Code)
	}
	for _, field := range []string{"title", "status", "priority"} {
		if _, ok := resp.Error.Fields[field]; !ok {
			t.Errorf("expected validation message for field %q, got fields: %v", field, resp.Error.Fields)
		}
	}
}

func TestUsersCannotAccessOthersTasks(t *testing.T) {
	r := newTestServer(t)
	aliceToken := signup(t, r, "alice@example.com")
	bobToken := signup(t, r, "bob@example.com")

	w := doJSON(t, r, http.MethodPost, "/tasks", aliceToken, `{"title":"Alice's secret task"}`)
	if w.Code != http.StatusCreated {
		t.Fatalf("failed to create task: %d %s", w.Code, w.Body.String())
	}
	var created struct {
		Data struct {
			ID uint `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
		t.Fatalf("failed to decode create response: %v", err)
	}

	path := fmt.Sprintf("/tasks/%d", created.Data.ID)

	if w := doJSON(t, r, http.MethodGet, path, bobToken, ""); w.Code != http.StatusNotFound {
		t.Errorf("GET: expected 404 for another user's task, got %d", w.Code)
	}
	if w := doJSON(t, r, http.MethodPatch, path, bobToken, `{"title":"hijacked"}`); w.Code != http.StatusNotFound {
		t.Errorf("PATCH: expected 404 for another user's task, got %d", w.Code)
	}
	if w := doJSON(t, r, http.MethodDelete, path, bobToken, ""); w.Code != http.StatusNotFound {
		t.Errorf("DELETE: expected 404 for another user's task, got %d", w.Code)
	}

	// The owner can still access it.
	if w := doJSON(t, r, http.MethodGet, path, aliceToken, ""); w.Code != http.StatusOK {
		t.Errorf("owner GET: expected 200, got %d", w.Code)
	}
}

func TestListFilterSearchSortAndPagination(t *testing.T) {
	r := newTestServer(t)
	token := signup(t, r, "carol@example.com")

	seed := []string{
		`{"title":"Write report","status":"todo","priority":"low","due_date":"2026-06-20T00:00:00Z"}`,
		`{"title":"Review report","status":"done","priority":"high","due_date":"2026-06-15T00:00:00Z"}`,
		`{"title":"Plan meeting","status":"todo","priority":"medium","due_date":"2026-06-18T00:00:00Z"}`,
		`{"title":"Send report email","status":"todo","priority":"high"}`,
	}
	for _, body := range seed {
		if w := doJSON(t, r, http.MethodPost, "/tasks", token, body); w.Code != http.StatusCreated {
			t.Fatalf("seed task failed: %d %s", w.Code, w.Body.String())
		}
	}

	type listResp struct {
		Data []struct {
			Title    string `json:"title"`
			Priority string `json:"priority"`
		} `json:"data"`
		Meta struct {
			Total      int64 `json:"total"`
			TotalPages int   `json:"total_pages"`
			Page       int   `json:"page"`
		} `json:"meta"`
	}
	list := func(query string) listResp {
		w := doJSON(t, r, http.MethodGet, "/tasks?"+query, token, "")
		if w.Code != http.StatusOK {
			t.Fatalf("list %q failed: %d %s", query, w.Code, w.Body.String())
		}
		var resp listResp
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to decode list response: %v", err)
		}
		return resp
	}

	// Status filter + title search combined.
	resp := list("status=todo&search=report")
	if resp.Meta.Total != 2 {
		t.Errorf("expected 2 todo tasks matching 'report', got %d", resp.Meta.Total)
	}

	// Sort by priority descending: high tasks first.
	resp = list("sort_by=priority&order=desc")
	if len(resp.Data) == 0 || resp.Data[0].Priority != "high" {
		t.Errorf("expected highest priority first, got %+v", resp.Data)
	}

	// Sort by due date ascending: earliest due date first, no-due-date last.
	resp = list("sort_by=due_date&order=asc")
	if len(resp.Data) != 4 || resp.Data[0].Title != "Review report" {
		t.Errorf("expected 'Review report' (earliest due) first, got %+v", resp.Data)
	}
	if resp.Data[3].Title != "Send report email" {
		t.Errorf("expected task without due date last, got %+v", resp.Data)
	}

	// Pagination.
	resp = list("limit=3&page=2")
	if resp.Meta.TotalPages != 2 || resp.Meta.Page != 2 || len(resp.Data) != 1 {
		t.Errorf("expected page 2 of 2 with 1 item, got meta=%+v len=%d", resp.Meta, len(resp.Data))
	}

	// Invalid query values are rejected.
	if w := doJSON(t, r, http.MethodGet, "/tasks?status=bogus", token, ""); w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid status filter, got %d", w.Code)
	}
}
