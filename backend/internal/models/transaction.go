package models

import "time"

type Transaction struct {
	ID               uint      `gorm:"primarykey" json:"id"`
	CustomerID       uint      `gorm:"not null;index:idx_customer_time,priority:1" json:"customer_id"`
	Amount           float64   `gorm:"type:decimal(10,2);not null" json:"amount"`
	Channel          string    `gorm:"size:20;not null" json:"channel"`
	CommissionRate   float64   `gorm:"type:decimal(5,4);not null" json:"commission_rate"`
	CommissionAmount float64   `gorm:"type:decimal(10,2);not null" json:"commission_amount"`
	IsNewCustomer    bool      `gorm:"not null" json:"is_new_customer"`
	RecordedBy       uint      `gorm:"not null;index:idx_recorder_time,priority:1" json:"recorded_by"`
	CreatedAt        time.Time `gorm:"index:idx_created_at,sort:desc;index:idx_customer_time,priority:2,sort:desc;index:idx_recorder_time,priority:2,sort:desc" json:"created_at"`

	Customer Customer `gorm:"foreignKey:CustomerID;constraint:OnDelete:RESTRICT" json:"customer,omitempty"`
	Recorder User     `gorm:"foreignKey:RecordedBy;constraint:OnDelete:RESTRICT" json:"recorder,omitempty"`
}
