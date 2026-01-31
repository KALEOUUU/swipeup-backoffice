package models

import (
	"time"
)

type ActivityLog struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	IDUser      uint      `json:"id_user" gorm:"not null"`
	Action      string    `json:"action" gorm:"not null;size:100"`
	Description string    `json:"description" gorm:"type:text"`
	IPAddress   string    `json:"ip_address" gorm:"type:inet"`
	UserAgent   string    `json:"user_agent" gorm:"type:text"`
	CreatedAt   time.Time `json:"created_at"`

	// Relations
	User User `json:"user" gorm:"foreignKey:IDUser;constraint:OnDelete:CASCADE"`
}