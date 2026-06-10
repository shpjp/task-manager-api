package models

import "time"

type TaskActivity struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TaskID    uint      `gorm:"index;not null" json:"task_id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	Action    string    `gorm:"size:50;not null" json:"action"`
	Detail    string    `gorm:"size:1000" json:"detail"`
	CreatedAt time.Time `json:"created_at"`

	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Task Task  `gorm:"constraint:OnDelete:CASCADE" json:"-"`
}
