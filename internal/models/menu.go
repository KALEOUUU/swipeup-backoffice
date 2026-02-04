package models

import (
	"time"

	"gorm.io/gorm"
)

type JenisMenu string

const (
	JenisMakanan  JenisMenu = "makanan"
	JenisMinuman  JenisMenu = "minuman"
)

type Menu struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	NamaMakanan  string    `json:"nama_makanan" gorm:"type:varchar(100);not null"`
	Harga        float64   `json:"harga" gorm:"not null"`
	Jenis        JenisMenu `json:"jenis" gorm:"type:varchar(20);not null"`
	Foto         string    `json:"foto" gorm:"type:varchar(255)"`
	Deskripsi    string    `json:"deskripsi" gorm:"type:text"`
	Stock        int       `json:"stock" gorm:"default:0"`
	IsAvailable  bool      `json:"is_available" gorm:"default:true"`
	IDStan       uint      `json:"id_stan" gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	
	// Relations
	Stan            Stan              `json:"stan" gorm:"foreignKey:IDStan;constraint:OnDelete:CASCADE"`
	DetailTransaksi []DetailTransaksi `json:"detail_transaksi,omitempty" gorm:"foreignKey:IDMenu"`
	MenuDiskon      []MenuDiskon      `json:"menu_diskon,omitempty" gorm:"foreignKey:IDMenu"`
}
