package repository

import (
	"task-manager-api/internal/models"

	"gorm.io/gorm"
)

type ActivityRepository struct {
	db *gorm.DB
}

func NewActivityRepository(db *gorm.DB) *ActivityRepository {
	return &ActivityRepository{db: db}
}

func (r *ActivityRepository) Create(activity *models.TaskActivity) error {
	return r.db.Create(activity).Error
}

func (r *ActivityRepository) ListByTask(taskID uint) ([]models.TaskActivity, error) {
	var activities []models.TaskActivity
	err := r.db.
		Preload("User").
		Where("task_id = ?", taskID).
		Order("created_at DESC, id DESC").
		Limit(100).
		Find(&activities).Error
	return activities, err
}
