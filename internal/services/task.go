package services

import (
	"time"

	"task-manager-api/internal/models"
	"task-manager-api/internal/repository"
)

type TaskService struct {
	tasks *repository.TaskRepository
}

func NewTaskService(tasks *repository.TaskRepository) *TaskService {
	return &TaskService{tasks: tasks}
}

type CreateTaskInput struct {
	Title       string
	Description string
	Status      string
	Priority    string
	DueDate     *time.Time
}

type UpdateTaskInput struct {
	Title       *string
	Description *string
	Status      *string
	Priority    *string
	DueDate     *time.Time
	ClearDue    bool
}

func (s *TaskService) Create(userID uint, in CreateTaskInput) (*models.Task, error) {
	task := &models.Task{
		UserID:      userID,
		Title:       in.Title,
		Description: in.Description,
		Status:      models.StatusTodo,
		Priority:    models.PriorityMedium,
		DueDate:     in.DueDate,
	}
	if in.Status != "" {
		task.Status = models.TaskStatus(in.Status)
	}
	if in.Priority != "" {
		task.Priority = models.TaskPriority(in.Priority)
	}
	if err := s.tasks.Create(task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskService) Get(userID, taskID uint) (*models.Task, error) {
	return s.tasks.FindByID(userID, taskID)
}

func (s *TaskService) Update(userID, taskID uint, in UpdateTaskInput) (*models.Task, error) {
	task, err := s.tasks.FindByID(userID, taskID)
	if err != nil {
		return nil, err
	}

	if in.Title != nil {
		task.Title = *in.Title
	}
	if in.Description != nil {
		task.Description = *in.Description
	}
	if in.Status != nil {
		task.Status = models.TaskStatus(*in.Status)
	}
	if in.Priority != nil {
		task.Priority = models.TaskPriority(*in.Priority)
	}
	if in.ClearDue {
		task.DueDate = nil
	} else if in.DueDate != nil {
		task.DueDate = in.DueDate
	}

	if err := s.tasks.Update(task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskService) Delete(userID, taskID uint) error {
	return s.tasks.Delete(userID, taskID)
}

func (s *TaskService) List(userID uint, filter repository.TaskFilter) ([]models.Task, int64, error) {
	return s.tasks.List(userID, filter)
}
