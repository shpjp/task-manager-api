package services

import (
	"fmt"
	"strings"
	"time"

	"task-manager-api/internal/models"
	"task-manager-api/internal/realtime"
	"task-manager-api/internal/repository"
)

// Actor identifies who is performing an operation.
type Actor struct {
	UserID uint
	Admin  bool
}

type TaskService struct {
	tasks      *repository.TaskRepository
	activities *repository.ActivityRepository
	hub        *realtime.Hub
}

func NewTaskService(
	tasks *repository.TaskRepository,
	activities *repository.ActivityRepository,
	hub *realtime.Hub,
) *TaskService {
	return &TaskService{tasks: tasks, activities: activities, hub: hub}
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

func (s *TaskService) Create(actor Actor, in CreateTaskInput) (*models.Task, error) {
	task := &models.Task{
		UserID:      actor.UserID,
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

	s.recordActivity(task.ID, actor.UserID, "created", fmt.Sprintf("created task %q", task.Title))
	s.publish(task.UserID, realtime.Event{Type: "task.created", TaskID: task.ID, Task: task})
	return task, nil
}

func (s *TaskService) Get(actor Actor, taskID uint) (*models.Task, error) {
	if actor.Admin {
		return s.tasks.FindByIDAny(taskID)
	}
	return s.tasks.FindByID(actor.UserID, taskID)
}

func (s *TaskService) Update(actor Actor, taskID uint, in UpdateTaskInput) (*models.Task, error) {
	// Updates are owner-only; admins have read-only access to others' tasks.
	task, err := s.tasks.FindByID(actor.UserID, taskID)
	if err != nil {
		return nil, err
	}

	var changes []string

	if in.Title != nil && *in.Title != task.Title {
		changes = append(changes, fmt.Sprintf("title %q → %q", task.Title, *in.Title))
		task.Title = *in.Title
	}
	if in.Description != nil && *in.Description != task.Description {
		changes = append(changes, "description updated")
		task.Description = *in.Description
	}
	if in.Status != nil && models.TaskStatus(*in.Status) != task.Status {
		changes = append(changes, fmt.Sprintf("status %s → %s", task.Status, *in.Status))
		task.Status = models.TaskStatus(*in.Status)
	}
	if in.Priority != nil && models.TaskPriority(*in.Priority) != task.Priority {
		changes = append(changes, fmt.Sprintf("priority %s → %s", task.Priority, *in.Priority))
		task.Priority = models.TaskPriority(*in.Priority)
	}
	if in.ClearDue {
		if task.DueDate != nil {
			changes = append(changes, "due date removed")
		}
		task.DueDate = nil
	} else if in.DueDate != nil && (task.DueDate == nil || !in.DueDate.Equal(*task.DueDate)) {
		changes = append(changes, fmt.Sprintf("due date set to %s", in.DueDate.Format("2006-01-02")))
		task.DueDate = in.DueDate
	}

	if err := s.tasks.Update(task); err != nil {
		return nil, err
	}

	if len(changes) > 0 {
		s.recordActivity(task.ID, actor.UserID, "updated", strings.Join(changes, "; "))
	}
	s.publish(task.UserID, realtime.Event{Type: "task.updated", TaskID: task.ID, Task: task})
	return task, nil
}

func (s *TaskService) Delete(actor Actor, taskID uint) error {
	task, err := s.tasks.FindByID(actor.UserID, taskID)
	if err != nil {
		return err
	}
	if err := s.tasks.Delete(actor.UserID, taskID); err != nil {
		return err
	}
	s.publish(task.UserID, realtime.Event{Type: "task.deleted", TaskID: taskID})
	return nil
}

func (s *TaskService) List(actor Actor, filter repository.TaskFilter) ([]models.Task, int64, error) {
	if !actor.Admin {
		filter.AllUsers = false
	}
	return s.tasks.List(actor.UserID, filter)
}

// ListActivity returns the change history of a task the actor may view.
func (s *TaskService) ListActivity(actor Actor, taskID uint) ([]models.TaskActivity, error) {
	if _, err := s.Get(actor, taskID); err != nil {
		return nil, err
	}
	return s.activities.ListByTask(taskID)
}

// RecordActivity lets collaborating services (e.g. attachments) add entries.
func (s *TaskService) RecordActivity(taskID, userID uint, action, detail string) {
	s.recordActivity(taskID, userID, action, detail)
}

// PublishTaskUpdated notifies subscribers that a task changed.
func (s *TaskService) PublishTaskUpdated(task *models.Task) {
	s.publish(task.UserID, realtime.Event{Type: "task.updated", TaskID: task.ID, Task: task})
}

func (s *TaskService) recordActivity(taskID, userID uint, action, detail string) {
	if s.activities == nil {
		return
	}
	// Activity logging is best-effort; it must not fail the main operation.
	_ = s.activities.Create(&models.TaskActivity{
		TaskID: taskID,
		UserID: userID,
		Action: action,
		Detail: detail,
	})
}

func (s *TaskService) publish(ownerID uint, event realtime.Event) {
	if s.hub == nil {
		return
	}
	s.hub.Publish(ownerID, event)
}
