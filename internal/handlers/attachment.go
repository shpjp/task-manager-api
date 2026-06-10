package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"task-manager-api/internal/models"
	"task-manager-api/internal/repository"
	"task-manager-api/internal/services"

	"github.com/gin-gonic/gin"
)

type AttachmentHandler struct {
	service *services.AttachmentService
}

func NewAttachmentHandler(service *services.AttachmentService) *AttachmentHandler {
	return &AttachmentHandler{service: service}
}

func (h *AttachmentHandler) Upload(c *gin.Context) {
	taskID, ok := parseID(c)
	if !ok {
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_BODY", `Send the file as multipart form field "file"`)
		return
	}

	attachment, err := h.service.Upload(actor(c), taskID, file)
	switch {
	case errors.Is(err, repository.ErrNotFound):
		respondError(c, http.StatusNotFound, "NOT_FOUND", "Task not found")
	case errors.Is(err, services.ErrFileTooLarge):
		respondError(c, http.StatusRequestEntityTooLarge, "FILE_TOO_LARGE",
			fmt.Sprintf("File exceeds the %d MB limit", h.service.MaxBytes()/(1<<20)))
	case errors.Is(err, services.ErrUnsupportedType):
		respondError(c, http.StatusUnsupportedMediaType, "UNSUPPORTED_TYPE",
			"Only images and common document formats are allowed")
	case err != nil:
		respondError(c, http.StatusInternalServerError, "INTERNAL", "Failed to upload file")
	default:
		respondData(c, http.StatusCreated, attachment)
	}
}

func (h *AttachmentHandler) List(c *gin.Context) {
	taskID, ok := parseID(c)
	if !ok {
		return
	}

	attachments, err := h.service.List(actor(c), taskID)
	if errors.Is(err, repository.ErrNotFound) {
		respondError(c, http.StatusNotFound, "NOT_FOUND", "Task not found")
		return
	}
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL", "Failed to list attachments")
		return
	}
	if attachments == nil {
		attachments = []models.Attachment{}
	}
	respondData(c, http.StatusOK, attachments)
}

func (h *AttachmentHandler) Download(c *gin.Context) {
	taskID, ok := parseID(c)
	if !ok {
		return
	}
	attachmentID, ok := parseAttachmentID(c)
	if !ok {
		return
	}

	attachment, path, err := h.service.Resolve(actor(c), taskID, attachmentID)
	if errors.Is(err, repository.ErrNotFound) {
		respondError(c, http.StatusNotFound, "NOT_FOUND", "Attachment not found")
		return
	}
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL", "Failed to fetch attachment")
		return
	}
	c.FileAttachment(path, attachment.FileName)
}

func (h *AttachmentHandler) Delete(c *gin.Context) {
	taskID, ok := parseID(c)
	if !ok {
		return
	}
	attachmentID, ok := parseAttachmentID(c)
	if !ok {
		return
	}

	err := h.service.Delete(actor(c), taskID, attachmentID)
	if errors.Is(err, repository.ErrNotFound) {
		respondError(c, http.StatusNotFound, "NOT_FOUND", "Attachment not found")
		return
	}
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL", "Failed to delete attachment")
		return
	}
	c.Status(http.StatusNoContent)
}

func parseAttachmentID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("attachmentID"), 10, 64)
	if err != nil || id == 0 {
		respondError(c, http.StatusBadRequest, "INVALID_ID", "Attachment id must be a positive integer")
		return 0, false
	}
	return uint(id), true
}
