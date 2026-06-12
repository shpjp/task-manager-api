package storage

import "mime/multipart"

// UploadResult holds metadata returned after a successful cloud upload.
type UploadResult struct {
	URL          string
	PublicID     string
	ResourceType string
	Bytes        int64
}

// Provider uploads and deletes files in external object storage (Cloudinary).
type Provider interface {
	Upload(file *multipart.FileHeader) (*UploadResult, error)
	Delete(publicID, resourceType string) error
}
