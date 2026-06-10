package repository

import (
	"errors"
	"fmt"

	"task-manager-api/internal/models"

	"gorm.io/gorm"
)

type TaskFilter struct {
	Status string
	Search string
	SortBy string // due_date | priority | created_at
	Order  string // asc | desc
	Page   int
	Limit  int
}

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(task *models.Task) error {
	return r.db.Create(task).Error
}

func (r *TaskRepository) FindByID(userID, taskID uint) (*models.Task, error) {
	var task models.Task
	err := r.db.Where("id = ? AND user_id = ?", taskID, userID).First(&task).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *TaskRepository) Update(task *models.Task) error {
	return r.db.Save(task).Error
}

func (r *TaskRepository) Delete(userID, taskID uint) error {
	res := r.db.Where("id = ? AND user_id = ?", taskID, userID).Delete(&models.Task{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *TaskRepository) List(userID uint, filter TaskFilter) ([]models.Task, int64, error) {
	query := r.db.Model(&models.Task{}).Where("user_id = ?", userID)

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Search != "" {
		query = query.Where(`LOWER(title) LIKE LOWER(?) ESCAPE '\'`, "%"+escapeLike(filter.Search)+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query = query.Order(orderClause(filter.SortBy, filter.Order))

	offset := (filter.Page - 1) * filter.Limit
	var tasks []models.Task
	if err := query.Offset(offset).Limit(filter.Limit).Find(&tasks).Error; err != nil {
		return nil, 0, err
	}
	return tasks, total, nil
}

// escapeLike escapes LIKE wildcards so user input is matched literally.
func escapeLike(s string) string {
	out := make([]rune, 0, len(s))
	for _, r := range s {
		if r == '%' || r == '_' || r == '\\' {
			out = append(out, '\\')
		}
		out = append(out, r)
	}
	return string(out)
}

func orderClause(sortBy, order string) string {
	dir := "ASC"
	if order == "desc" {
		dir = "DESC"
	}

	switch sortBy {
	case "due_date":
		// Keep tasks without a due date at the end regardless of direction.
		return fmt.Sprintf("due_date IS NULL, due_date %s, id %s", dir, dir)
	case "priority":
		return fmt.Sprintf("CASE priority WHEN 'high' THEN 3 WHEN 'medium' THEN 2 ELSE 1 END %s, id %s", dir, dir)
	case "created_at":
		return fmt.Sprintf("created_at %s, id %s", dir, dir)
	default:
		return "created_at DESC, id DESC"
	}
}
