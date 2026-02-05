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
	ID               uint           `json:"id" gorm:"column:id;primaryKey"`
	NamaDiskon       string         `json:"nama_diskon" gorm:"column:nama_diskon;type:varchar(100);not null"`
	PersentaseDiskon float64        `json:"persentase_diskon" gorm:"column:persentase_diskon;not null"`
	TanggalAwal      time.Time      `json:"tanggal_awal" gorm:"column:tanggal_awal;not null"`
	TanggalAkhir     time.Time      `json:"tanggal_akhir" gorm:"column:tanggal_akhir;not null"`
	TipeDiskon       TipeDiskon     `json:"tipe_diskon" gorm:"column:tipe_diskon;type:varchar(20);not null;default:'global'"`
	IDStan           *uint          `json:"id_stan" gorm:"column:id_stan;index"` // NULL untuk global, berisi ID untuk diskon stan
	CreatedBy        string         `json:"created_by" gorm:"column:created_by"`
	UpdatedBy        string         `json:"updated_by" gorm:"column:updated_by"`
	CreatedAt        time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt        time.Time      `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt        gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;index"`
	
	// Relations
	Stan       *Stan        `json:"stan,omitempty" gorm:"foreignKey:IDStan"`
	MenuDiskon []MenuDiskon `json:"menu_diskon,omitempty" gorm:"foreignKey:IDDiskon"`
}

