package models

import (
	"time"

	"gorm.io/gorm"
)

type MenuDiskon struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	IDMenu    uint      `json:"id_menu" gorm:"not null"`
	IDDiskon  uint      `json:"id_diskon" gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	
	// Relations
	Menu   Menu   `json:"menu" gorm:"foreignKey:IDMenu;constraint:OnDelete:CASCADE"`
	Diskon Diskon `json:"diskon" gorm:"foreignKey:IDDiskon;constraint:OnDelete:CASCADE"`
}
