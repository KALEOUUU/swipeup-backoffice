package models

import (
	"time"
)

type Cart struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	IDSiswa   uint      `json:"id_siswa" gorm:"not null"`
	IDMenu    uint      `json:"id_menu" gorm:"not null"`
	Qty       int       `json:"qty" gorm:"not null;check:qty > 0"`
	CreatedAt time.Time
	UpdatedAt time.Time

	// Relations
	Siswa Siswa `json:"siswa" gorm:"foreignKey:IDSiswa;constraint:OnDelete:CASCADE"`
	Menu  Menu  `json:"menu" gorm:"foreignKey:IDMenu;constraint:OnDelete:CASCADE"`
}