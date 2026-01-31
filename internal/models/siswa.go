package models

import (
	"time"

	"gorm.io/gorm"
)

type Siswa struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	NamaSiswa  string `json:"nama_siswa" gorm:"type:varchar(100);not null"`
	Alamat     string `json:"alamat" gorm:"type:text"`
	Telp       string `json:"telp" gorm:"type:varchar(20)"`
	IDUser     uint   `json:"id_user" gorm:"not null"`
	Foto       string `json:"foto" gorm:"type:varchar(255)"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	
	// Relations
	User        User         `json:"user" gorm:"foreignKey:IDUser;constraint:OnDelete:CASCADE"`
	Transaksi   []Transaksi  `json:"transaksi,omitempty" gorm:"foreignKey:IDSiswa"`
}
