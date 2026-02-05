package models

import (
	"time"

	"gorm.io/gorm"
)

type StatusTransaksi string

const (
	StatusBelumDikonfirm StatusTransaksi = "belum dikonfirm"
	StatusDimasak        StatusTransaksi = "dimasak"
	StatusDiantar        StatusTransaksi = "diantar"
	StatusSampai         StatusTransaksi = "sampai"
)

type Transaksi struct {
	ID              uint            `json:"id" gorm:"column:id;primaryKey"`
	Tanggal         time.Time       `json:"tanggal" gorm:"column:tanggal;not null"`
	IDStan          uint            `json:"id_stan" gorm:"column:id_stan;not null"`
	IDSiswa         uint            `json:"id_siswa" gorm:"column:id_siswa;not null"`
	Status          StatusTransaksi `json:"status" gorm:"column:status;type:varchar(20);not null;default:'belum dikonfirm'"`
	CreatedBy       string          `json:"created_by" gorm:"column:created_by"`
	UpdatedBy       string          `json:"updated_by" gorm:"column:updated_by"`
	CreatedAt       time.Time       `json:"created_at" gorm:"column:created_at"`
	UpdatedAt       time.Time       `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt       gorm.DeletedAt  `json:"deleted_at" gorm:"column:deleted_at;index"`
	
	// Relations
	Stan            Stan              `json:"stan" gorm:"foreignKey:IDStan;constraint:OnDelete:CASCADE"`
	Siswa           Siswa             `json:"siswa" gorm:"foreignKey:IDSiswa;constraint:OnDelete:CASCADE"`
	DetailTransaksi []DetailTransaksi `json:"detail_transaksi,omitempty" gorm:"foreignKey:IDTransaksi"`
}
