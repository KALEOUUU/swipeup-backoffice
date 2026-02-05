package models

import (
	"time"

	"gorm.io/gorm"
)

type TipeDiskon string

const (
	DiskonGlobal TipeDiskon = "global" // Diatur oleh superadmin, berlaku semua stan
	DiskonStan   TipeDiskon = "stan"   // Diatur oleh admin stan, berlaku untuk stannya saja
	DiskonMenu   TipeDiskon = "menu"   // Diatur oleh admin stan, berlaku untuk menu tertentu
)

type Diskon struct {
	ID               uint       `json:"id" gorm:"primaryKey"`
	NamaDiskon       string     `json:"nama_diskon" gorm:"type:varchar(100);not null"`
	PersentaseDiskon float64    `json:"persentase_diskon" gorm:"not null"`
	TanggalAwal      time.Time  `json:"tanggal_awal" gorm:"not null"`
	TanggalAkhir     time.Time  `json:"tanggal_akhir" gorm:"not null"`
	TipeDiskon       TipeDiskon `json:"tipe_diskon" gorm:"type:varchar(20);not null;default:'global'"`
	IDStan           *uint      `json:"id_stan" gorm:"index"` // NULL untuk global, berisi ID untuk diskon stan
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
	
	// Relations
	Stan       *Stan        `json:"stan,omitempty" gorm:"foreignKey:IDStan"`
	MenuDiskon []MenuDiskon `json:"menu_diskon,omitempty" gorm:"foreignKey:IDDiskon"`
}

