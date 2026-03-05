package models

import (
	"gorm.io/gorm"
	"time"
)

type KeySpec struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	Name         string         `gorm:"size:100;not null;uniqueIndex" json:"name"`
	Description  string         `gorm:"type:text" json:"description"`
	DisplayOrder int            `gorm:"not null;default:0;index" json:"display_order"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
