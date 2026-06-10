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
	"task-manager-api/internal/realtime"
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
	if err := db.AutoMigrate(
		&models.User{},
		&models.Task{},
		&models.TaskActivity{},
		&models.Attachment{},
	); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	jwtManager := auth.NewJWTManager("test-secret", time.Hour)
	hub := realtime.NewHub()
	taskRepo := repository.NewTaskRepository(db)
	taskService := services.NewTaskService(taskRepo, repository.NewActivityRepository(db), hub)
	authService := services.NewAuthService(
		repository.NewUserRepository(db), jwtManager, []string{"admin@example.com"},
	)
	attachmentService := services.NewAttachmentService(
		repository.NewAttachmentRepository(db), taskRepo, taskService, t.TempDir(), 5<<20,
	)

	r := gin.New()
	routes.RegisterRoutes(r, routes.Handlers{
		Auth:        handlers.NewAuthHandler(authService, jwtManager, false),
		Tasks:       handlers.NewTaskHandler(taskService),
		Attachments: handlers.NewAttachmentHandler(attachmentService),
		Events:      handlers.NewEventsHandler(hub),
	}, jwtManager)
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

func TestAdminCanViewAllTasksButNotModifyThem(t *testing.T) {
	r := newTestServer(t)
	userToken := signup(t, r, "regular@example.com")
	adminToken := signup(t, r, "admin@example.com") // promoted via ADMIN_EMAILS

	w := doJSON(t, r, http.MethodPost, "/tasks", userToken, `{"title":"User's private task"}`)
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

	// A regular user may not list everyone's tasks.
	if w := doJSON(t, r, http.MethodGet, "/tasks?scope=all", userToken, ""); w.Code != http.StatusForbidden {
		t.Errorf("expected 403 for non-admin scope=all, got %d", w.Code)
	}

	// The admin sees the other user's task in the global listing with owner info.
	w = doJSON(t, r, http.MethodGet, "/tasks?scope=all", adminToken, "")
	if w.Code != http.StatusOK {
		t.Fatalf("admin scope=all failed: %d %s", w.Code, w.Body.String())
	}
	var list struct {
		Data []struct {
			Title string `json:"title"`
			Owner *struct {
				Email string `json:"email"`
			} `json:"owner"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &list); err != nil {
		t.Fatalf("failed to decode list: %v", err)
	}
	if len(list.Data) != 1 || list.Data[0].Owner == nil || list.Data[0].Owner.Email != "regular@example.com" {
		t.Errorf("expected the user's task with owner email in admin listing, got %+v", list.Data)
	}

	// Admin can read a single task but not modify or delete it.
	if w := doJSON(t, r, http.MethodGet, path, adminToken, ""); w.Code != http.StatusOK {
		t.Errorf("admin GET: expected 200, got %d", w.Code)
	}
	if w := doJSON(t, r, http.MethodPatch, path, adminToken, `{"title":"nope"}`); w.Code != http.StatusNotFound {
		t.Errorf("admin PATCH: expected 404 (read-only), got %d", w.Code)
	}
	if w := doJSON(t, r, http.MethodDelete, path, adminToken, ""); w.Code != http.StatusNotFound {
		t.Errorf("admin DELETE: expected 404 (read-only), got %d", w.Code)
	}
}

func TestActivityLogRecordsChanges(t *testing.T) {
	r := newTestServer(t)
	token := signup(t, r, "dana@example.com")

	w := doJSON(t, r, http.MethodPost, "/tasks", token, `{"title":"Track me"}`)
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
	if w := doJSON(t, r, http.MethodPatch, path, token, `{"status":"done"}`); w.Code != http.StatusOK {
		t.Fatalf("failed to update task: %d %s", w.Code, w.Body.String())
	}

	w = doJSON(t, r, http.MethodGet, path+"/activity", token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("failed to fetch activity: %d %s", w.Code, w.Body.String())
	}
	var activity struct {
		Data []struct {
			Action string `json:"action"`
			Detail string `json:"detail"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &activity); err != nil {
		t.Fatalf("failed to decode activity: %v", err)
	}
	if len(activity.Data) != 2 {
		t.Fatalf("expected 2 activity entries (created + updated), got %d: %+v", len(activity.Data), activity.Data)
	}
	// Newest first.
	if activity.Data[0].Action != "updated" || !strings.Contains(activity.Data[0].Detail, "todo → done") {
		t.Errorf("expected status-change entry first, got %+v", activity.Data[0])
	}
	if activity.Data[1].Action != "created" {
		t.Errorf("expected created entry, got %+v", activity.Data[1])
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
