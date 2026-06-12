package storage

import (
	"fmt"
	"mime/multipart"
)

// NoopProvider is used in tests. Returns fake URLs without calling Cloudinary.
type NoopProvider struct{}

func (NoopProvider) Upload(file *multipart.FileHeader) (*UploadResult, error) {
	return &UploadResult{
		URL:          "https://res.cloudinary.com/demo/image/upload/sample.jpg",
		PublicID:     "test/" + file.Filename,
		ResourceType: "image",
		Bytes:        file.Size,
	}, nil
}

func (NoopProvider) Delete(_ string, _ string) error { return nil }

// ErrNotConfigured is returned when Cloudinary env vars are missing.
var ErrNotConfigured = fmt.Errorf("cloudinary is not configured")
