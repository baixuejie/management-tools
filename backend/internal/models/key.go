package models

import (
	"time"
	"gorm.io/gorm"
)

type Key struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	SpecID    uint           `gorm:"not null;index:idx_spec_used_time;index:idx_spec_created,priority:1" json:"spec_id"`
	KeyValue  string         `gorm:"type:text;not null" json:"key_value"`
	IsUsed    bool           `gorm:"default:false;index:idx_spec_used_time" json:"is_used"`
	UsedAt    *time.Time     `gorm:"index:idx_spec_used_time" json:"used_at"`
	CreatedAt time.Time      `gorm:"not null;index:idx_spec_created,priority:2" json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Spec      KeySpec        `gorm:"foreignKey:SpecID;constraint:OnDelete:RESTRICT" json:"spec,omitempty"`
}
