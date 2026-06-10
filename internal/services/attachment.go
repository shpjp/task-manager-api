package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"task-manager-api/internal/models"
	"task-manager-api/internal/repository"
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
	uploadDir   string
	maxBytes    int64
}

func NewAttachmentService(
	attachments *repository.AttachmentRepository,
	tasks *repository.TaskRepository,
	taskService *TaskService,
	uploadDir string,
	maxBytes int64,
) *AttachmentService {
	return &AttachmentService{
		attachments: attachments,
		tasks:       tasks,
		taskService: taskService,
		uploadDir:   uploadDir,
		maxBytes:    maxBytes,
	}
}

func (s *AttachmentService) MaxBytes() int64 {
	return s.maxBytes
}

// taskForRead resolves a task the actor may view (owner, or admin read-only).
func (s *AttachmentService) taskForRead(actor Actor, taskID uint) (*models.Task, error) {
	if actor.Admin {
		return s.tasks.FindByIDAny(taskID)
	}
	return s.tasks.FindByID(actor.UserID, taskID)
}

func (s *AttachmentService) Upload(actor Actor, taskID uint, file *multipart.FileHeader) (*models.Attachment, error) {
	// Only the owner may modify a task's attachments.
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

	if err := os.MkdirAll(s.uploadDir, 0o755); err != nil {
		return nil, err
	}

	storedName, err := randomName(ext)
	if err != nil {
		return nil, err
	}

	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	dstPath := filepath.Join(s.uploadDir, storedName)
	dst, err := os.Create(dstPath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, io.LimitReader(src, s.maxBytes+1)); err != nil {
		os.Remove(dstPath)
		return nil, err
	}

	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	attachment := &models.Attachment{
		TaskID:      task.ID,
		UserID:      actor.UserID,
		FileName:    filepath.Base(file.Filename),
		ContentType: contentType,
		Size:        file.Size,
		StoredName:  storedName,
	}
	if err := s.attachments.Create(attachment); err != nil {
		os.Remove(dstPath)
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

// Resolve returns the attachment metadata and its absolute file path.
func (s *AttachmentService) Resolve(actor Actor, taskID, attachmentID uint) (*models.Attachment, string, error) {
	if _, err := s.taskForRead(actor, taskID); err != nil {
		return nil, "", err
	}
	attachment, err := s.attachments.FindByID(taskID, attachmentID)
	if err != nil {
		return nil, "", err
	}
	return attachment, filepath.Join(s.uploadDir, attachment.StoredName), nil
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
	os.Remove(filepath.Join(s.uploadDir, attachment.StoredName))

	s.taskService.RecordActivity(task.ID, actor.UserID, "attachment_removed",
		fmt.Sprintf("removed attachment %q", attachment.FileName))
	s.taskService.PublishTaskUpdated(task)
	return nil
}

func randomName(ext string) (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf) + ext, nil
}
