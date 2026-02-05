package models

import (
	"time"
)

type Cart struct {
	ID        uint      `json:"id" gorm:"column:id;primaryKey"`
	IDSiswa   uint      `json:"id_siswa" gorm:"column:id_siswa;not null"`
	IDMenu    uint      `json:"id_menu" gorm:"column:id_menu;not null"`
	Qty       int       `json:"qty" gorm:"column:qty;not null;check:qty > 0"`
	CreatedBy string    `json:"created_by" gorm:"column:created_by"`
	UpdatedBy string    `json:"updated_by" gorm:"column:updated_by"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`

	// Relations
	Siswa Siswa `json:"siswa" gorm:"foreignKey:IDSiswa;constraint:OnDelete:CASCADE"`
	Menu  Menu  `json:"menu" gorm:"foreignKey:IDMenu;constraint:OnDelete:CASCADE"`
}