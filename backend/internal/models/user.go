package models

import "time"

type User struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	Username     string    `gorm:"size:50;not null;uniqueIndex" json:"username"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"`
	DisplayName  string    `gorm:"size:100;not null" json:"display_name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
