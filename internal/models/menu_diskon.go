package models

import (
	"time"

	"gorm.io/gorm"
)

type MenuDiskon struct {
	ID        uint           `json:"id" gorm:"column:id;primaryKey"`
	IDMenu    uint           `json:"id_menu" gorm:"column:id_menu;not null"`
	IDDiskon  uint           `json:"id_diskon" gorm:"column:id_diskon;not null"`
	CreatedBy string         `json:"created_by" gorm:"column:created_by"`
	UpdatedBy string         `json:"updated_by" gorm:"column:updated_by"`
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;index"`
	
	// Relations
	Menu   Menu   `json:"menu" gorm:"foreignKey:IDMenu;constraint:OnDelete:CASCADE"`
	Diskon Diskon `json:"diskon" gorm:"foreignKey:IDDiskon;constraint:OnDelete:CASCADE"`
}
