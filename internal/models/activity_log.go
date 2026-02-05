package models

import (
	"time"
)

type ActivityLog struct {
	ID          uint      `json:"id" gorm:"column:id;primaryKey"`
	IDUser      uint      `json:"id_user" gorm:"column:id_user;not null"`
	Action      string    `json:"action" gorm:"column:action;not null;size:100"`
	Description string    `json:"description" gorm:"column:description;type:text"`
	IPAddress   string    `json:"ip_address" gorm:"column:ip_address;type:inet"`
	UserAgent   string    `json:"user_agent" gorm:"column:user_agent;type:text"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at"`

	// Relations
	User User `json:"user" gorm:"foreignKey:IDUser;constraint:OnDelete:CASCADE"`
}