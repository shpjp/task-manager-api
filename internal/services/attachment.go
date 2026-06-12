package services

import (
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"task-manager-api/internal/models"
	"task-manager-api/internal/repository"
	"task-manager-api/internal/storage"
)

var (
	ErrFileTooLarge    = errors.New("file is too large")
	ErrUnsupportedType = errors.New("unsupported file type")
)

// Allowed upload extensions: images and common documents.
var allowedExtensions = map[string]bool{
	".png": true, ".jpg": true, ".jpeg": true, ".gif": true, ".webp": true,
	".pdf": true, ".txt": true, ".md": true, ".csv": true,
	".doc": true, ".docx": true, ".xls": true, ".xlsx": true,
}

type AttachmentService struct {
	attachments *repository.AttachmentRepository
	tasks       *repository.TaskRepository
	taskService *TaskService
	storage     storage.Provider
	maxBytes    int64
}

func NewAttachmentService(
	attachments *repository.AttachmentRepository,
	tasks *repository.TaskRepository,
	taskService *TaskService,
	store storage.Provider,
	maxBytes int64,
) *AttachmentService {
	return &AttachmentService{
		attachments: attachments,
		tasks:       tasks,
		taskService: taskService,
		storage:     store,
		maxBytes:    maxBytes,
	}
}

func (s *AttachmentService) MaxBytes() int64 {
	return s.maxBytes
}

func (s *AttachmentService) taskForRead(actor Actor, taskID uint) (*models.Task, error) {
	if actor.Admin {
		return s.tasks.FindByIDAny(taskID)
	}
	return s.tasks.FindByID(actor.UserID, taskID)
}

func (s *AttachmentService) Upload(actor Actor, taskID uint, file *multipart.FileHeader) (*models.Attachment, error) {
	task, err := s.tasks.FindByID(actor.UserID, taskID)
	if err != nil {
		return nil, err
	}

	if file.Size > s.maxBytes {
		return nil, ErrFileTooLarge
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExtensions[ext] {
		return nil, ErrUnsupportedType
	}

	uploaded, err := s.storage.Upload(file)
	if err != nil {
		return nil, err
	}

	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	size := file.Size
	if uploaded.Bytes > 0 {
		size = uploaded.Bytes
	}

	attachment := &models.Attachment{
		TaskID:             task.ID,
		UserID:             actor.UserID,
		FileName:           filepath.Base(file.Filename),
		ContentType:        contentType,
		Size:               size,
		URL:                uploaded.URL,
		CloudinaryPublicID: uploaded.PublicID,
		ResourceType:       uploaded.ResourceType,
	}
	if err := s.attachments.Create(attachment); err != nil {
		_ = s.storage.Delete(uploaded.PublicID, uploaded.ResourceType)
		return nil, err
	}

	s.taskService.RecordActivity(task.ID, actor.UserID, "attachment_added",
		fmt.Sprintf("attached %q", attachment.FileName))
	s.taskService.PublishTaskUpdated(task)
	return attachment, nil
}

func (s *AttachmentService) List(actor Actor, taskID uint) ([]models.Attachment, error) {
	if _, err := s.taskForRead(actor, taskID); err != nil {
		return nil, err
	}
	return s.attachments.ListByTask(taskID)
}

func (s *AttachmentService) Delete(actor Actor, taskID, attachmentID uint) error {
	task, err := s.tasks.FindByID(actor.UserID, taskID)
	if err != nil {
		return err
	}
	attachment, err := s.attachments.FindByID(taskID, attachmentID)
	if err != nil {
		return err
	}
	if err := s.attachments.Delete(attachment.ID); err != nil {
		return err
	}
	_ = s.storage.Delete(attachment.CloudinaryPublicID, attachment.ResourceType)

	s.taskService.RecordActivity(task.ID, actor.UserID, "attachment_removed",
		fmt.Sprintf("removed attachment %q", attachment.FileName))
	s.taskService.PublishTaskUpdated(task)
	return nil
}
