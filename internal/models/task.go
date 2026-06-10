package models

import "time"

type TaskStatus string

const (
	StatusTodo       TaskStatus = "todo"
	StatusInProgress TaskStatus = "in_progress"
	StatusDone       TaskStatus = "done"
)

type TaskPriority string

const (
	PriorityLow    TaskPriority = "low"
	PriorityMedium TaskPriority = "medium"
	PriorityHigh   TaskPriority = "high"
)

type Task struct {
	ID          uint         `gorm:"primaryKey" json:"id"`
	UserID      uint         `gorm:"index;not null" json:"user_id"`
	Title       string       `gorm:"size:200;not null" json:"title"`
	Description string       `gorm:"size:2000" json:"description"`
	Status      TaskStatus   `gorm:"size:20;not null;default:todo;index" json:"status"`
	Priority    TaskPriority `gorm:"size:20;not null;default:medium" json:"priority"`
	DueDate     *time.Time   `json:"due_date"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}
