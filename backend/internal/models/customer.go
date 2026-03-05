package models

import "time"

type Customer struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Name      string    `gorm:"size:100;not null;index:idx_customer_name;uniqueIndex:uk_customer_name_owner" json:"name"`
	CreatedBy uint      `gorm:"not null;uniqueIndex:uk_customer_name_owner" json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Creator User `gorm:"foreignKey:CreatedBy;constraint:OnDelete:RESTRICT" json:"creator,omitempty"`
}
