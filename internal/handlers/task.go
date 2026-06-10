package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"task-manager-api/internal/middleware"
	"task-manager-api/internal/models"
	"task-manager-api/internal/repository"
	"task-manager-api/internal/services"

	"github.com/gin-gonic/gin"
)

const (
	defaultPageSize = 10
	maxPageSize     = 100
)

type TaskHandler struct {
	service *services.TaskService
}

func NewTaskHandler(service *services.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

// nullableTime distinguishes `"due_date": null` (clear the date) from the
// field being absent (leave unchanged) in PATCH bodies.
type nullableTime struct {
	Set   bool
	Value *time.Time
}

func (n *nullableTime) UnmarshalJSON(b []byte) error {
	n.Set = true
	if string(b) == "null" {
		return nil
	}
	return json.Unmarshal(b, &n.Value)
}

type createTaskRequest struct {
	Title       string     `json:"title" binding:"required,min=1,max=200"`
	Description string     `json:"description" binding:"omitempty,max=2000"`
	Status      string     `json:"status" binding:"omitempty,oneof=todo in_progress done"`
	Priority    string     `json:"priority" binding:"omitempty,oneof=low medium high"`
	DueDate     *time.Time `json:"due_date"`
}

type updateTaskRequest struct {
	Title       *string      `json:"title" binding:"omitempty,min=1,max=200"`
	Description *string      `json:"description" binding:"omitempty,max=2000"`
	Status      *string      `json:"status" binding:"omitempty,oneof=todo in_progress done"`
	Priority    *string      `json:"priority" binding:"omitempty,oneof=low medium high"`
	DueDate     nullableTime `json:"due_date"`
}

type listMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

func (h *TaskHandler) Create(c *gin.Context) {
	var req createTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBindingError(c, err)
		return
	}

	task, err := h.service.Create(middleware.CurrentUserID(c), services.CreateTaskInput{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Priority:    req.Priority,
		DueDate:     req.DueDate,
	})
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL", "Failed to create task")
		return
	}
	respondData(c, http.StatusCreated, task)
}

func (h *TaskHandler) List(c *gin.Context) {
	filter, ok := parseTaskFilter(c)
	if !ok {
		return
	}

	tasks, total, err := h.service.List(middleware.CurrentUserID(c), filter)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL", "Failed to list tasks")
		return
	}

	totalPages := int((total + int64(filter.Limit) - 1) / int64(filter.Limit))
	if tasks == nil {
		tasks = []models.Task{}
	}
	c.JSON(http.StatusOK, gin.H{
		"data": tasks,
		"meta": listMeta{Page: filter.Page, Limit: filter.Limit, Total: total, TotalPages: totalPages},
	})
}

func (h *TaskHandler) Get(c *gin.Context) {
	taskID, ok := parseID(c)
	if !ok {
		return
	}

	task, err := h.service.Get(middleware.CurrentUserID(c), taskID)
	if errors.Is(err, repository.ErrNotFound) {
		respondError(c, http.StatusNotFound, "NOT_FOUND", "Task not found")
		return
	}
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL", "Failed to fetch task")
		return
	}
	respondData(c, http.StatusOK, task)
}

func (h *TaskHandler) Update(c *gin.Context) {
	taskID, ok := parseID(c)
	if !ok {
		return
	}

	var req updateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBindingError(c, err)
		return
	}

	input := services.UpdateTaskInput{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Priority:    req.Priority,
	}
	if req.DueDate.Set {
		if req.DueDate.Value == nil {
			input.ClearDue = true
		} else {
			input.DueDate = req.DueDate.Value
		}
	}

	task, err := h.service.Update(middleware.CurrentUserID(c), taskID, input)
	if errors.Is(err, repository.ErrNotFound) {
		respondError(c, http.StatusNotFound, "NOT_FOUND", "Task not found")
		return
	}
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL", "Failed to update task")
		return
	}
	respondData(c, http.StatusOK, task)
}

func (h *TaskHandler) Delete(c *gin.Context) {
	taskID, ok := parseID(c)
	if !ok {
		return
	}

	err := h.service.Delete(middleware.CurrentUserID(c), taskID)
	if errors.Is(err, repository.ErrNotFound) {
		respondError(c, http.StatusNotFound, "NOT_FOUND", "Task not found")
		return
	}
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL", "Failed to delete task")
		return
	}
	c.Status(http.StatusNoContent)
}

func parseID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		respondError(c, http.StatusBadRequest, "INVALID_ID", "Task id must be a positive integer")
		return 0, false
	}
	return uint(id), true
}

func parseTaskFilter(c *gin.Context) (repository.TaskFilter, bool) {
	filter := repository.TaskFilter{
		Status: c.Query("status"),
		Search: c.Query("search"),
		SortBy: c.DefaultQuery("sort_by", "created_at"),
		Order:  c.DefaultQuery("order", "desc"),
		Page:   1,
		Limit:  defaultPageSize,
	}

	if filter.Status != "" && filter.Status != "todo" && filter.Status != "in_progress" && filter.Status != "done" {
		respondError(c, http.StatusBadRequest, "INVALID_QUERY", "status must be one of: todo, in_progress, done")
		return filter, false
	}
	if filter.SortBy != "due_date" && filter.SortBy != "priority" && filter.SortBy != "created_at" {
		respondError(c, http.StatusBadRequest, "INVALID_QUERY", "sort_by must be one of: due_date, priority, created_at")
		return filter, false
	}
	if filter.Order != "asc" && filter.Order != "desc" {
		respondError(c, http.StatusBadRequest, "INVALID_QUERY", "order must be one of: asc, desc")
		return filter, false
	}

	if raw := c.Query("page"); raw != "" {
		page, err := strconv.Atoi(raw)
		if err != nil || page < 1 {
			respondError(c, http.StatusBadRequest, "INVALID_QUERY", "page must be a positive integer")
			return filter, false
		}
		filter.Page = page
	}
	if raw := c.Query("limit"); raw != "" {
		limit, err := strconv.Atoi(raw)
		if err != nil || limit < 1 || limit > maxPageSize {
			respondError(c, http.StatusBadRequest, "INVALID_QUERY", "limit must be between 1 and 100")
			return filter, false
		}
		filter.Limit = limit
	}
	return filter, true
}
