package models

import "time"

type Attachment struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	TaskID             uint      `gorm:"index;not null" json:"task_id"`
	UserID             uint      `gorm:"not null" json:"user_id"`
	FileName           string    `gorm:"size:255;not null" json:"file_name"`
	ContentType        string    `gorm:"size:100;not null" json:"content_type"`
	Size               int64     `gorm:"not null" json:"size"`
	URL                string    `gorm:"size:512;not null" json:"url"`
	CloudinaryPublicID string    `gorm:"size:255;not null" json:"-"`
	ResourceType       string    `gorm:"size:32;not null;default:image" json:"-"`
	CreatedAt          time.Time `json:"created_at"`

	Task Task `gorm:"constraint:OnDelete:CASCADE" json:"-"`
}
