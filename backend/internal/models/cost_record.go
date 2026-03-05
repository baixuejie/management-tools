package models

import "time"

type CostRecord struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	Amount     float64   `gorm:"type:decimal(10,2);not null" json:"amount"`
	Note       string    `gorm:"type:text" json:"note"`
	RecordedBy uint      `gorm:"not null;index" json:"recorded_by"`
	CreatedAt  time.Time `gorm:"index:idx_cost_created_at,sort:desc" json:"created_at"`

	Recorder User `gorm:"foreignKey:RecordedBy;constraint:OnDelete:RESTRICT" json:"recorder,omitempty"`
}
