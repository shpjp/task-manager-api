package repository

import (
	"errors"

	"task-manager-api/internal/models"

	"gorm.io/gorm"
)

type AttachmentRepository struct {
	db *gorm.DB
}

func NewAttachmentRepository(db *gorm.DB) *AttachmentRepository {
	return &AttachmentRepository{db: db}
}

func (r *AttachmentRepository) Create(attachment *models.Attachment) error {
	return r.db.Create(attachment).Error
}

func (r *AttachmentRepository) ListByTask(taskID uint) ([]models.Attachment, error) {
	var attachments []models.Attachment
	err := r.db.Where("task_id = ?", taskID).Order("created_at DESC").Find(&attachments).Error
	return attachments, err
}

func (r *AttachmentRepository) FindByID(taskID, attachmentID uint) (*models.Attachment, error) {
	var attachment models.Attachment
	err := r.db.Where("id = ? AND task_id = ?", attachmentID, taskID).First(&attachment).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &attachment, nil
}

func (r *AttachmentRepository) Delete(id uint) error {
	return r.db.Delete(&models.Attachment{}, id).Error
}
