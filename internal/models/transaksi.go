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
	ID              uint            `json:"id" gorm:"primaryKey"`
	Tanggal         time.Time       `json:"tanggal" gorm:"not null"`
	IDStan          uint            `json:"id_stan" gorm:"not null"`
	IDSiswa         uint            `json:"id_siswa" gorm:"not null"`
	Status          StatusTransaksi `json:"status" gorm:"type:varchar(20);not null;default:'belum dikonfirm'"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt  `gorm:"index"`
	
	// Relations
	Stan            Stan              `json:"stan" gorm:"foreignKey:IDStan;constraint:OnDelete:CASCADE"`
	Siswa           Siswa             `json:"siswa" gorm:"foreignKey:IDSiswa;constraint:OnDelete:CASCADE"`
	DetailTransaksi []DetailTransaksi `json:"detail_transaksi,omitempty" gorm:"foreignKey:IDTransaksi"`
}
