package models

import (
	"time"

	"gorm.io/gorm"
)

type Siswa struct {
	ID         uint           `json:"id" gorm:"column:id;primaryKey"`
	NamaSiswa  string         `json:"nama_siswa" gorm:"column:nama_siswa;type:varchar(100);not null"`
	Alamat     string         `json:"alamat" gorm:"column:alamat;type:text"`
	Telp       string         `json:"telp" gorm:"column:telp;type:varchar(20)"`
	IDUser     uint           `json:"id_user" gorm:"column:id_user;not null"`
	Foto       string         `json:"foto" gorm:"column:foto;type:varchar(255)"`
	CreatedBy  string         `json:"created_by" gorm:"column:created_by"`
	UpdatedBy  string         `json:"updated_by" gorm:"column:updated_by"`
	CreatedAt  time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt  time.Time      `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;index"`
	
	// Relations
	User        User         `json:"user" gorm:"foreignKey:IDUser;constraint:OnDelete:CASCADE"`
	Transaksi   []Transaksi  `json:"transaksi,omitempty" gorm:"foreignKey:IDSiswa"`
}
