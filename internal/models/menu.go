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
	ID           uint           `json:"id" gorm:"column:id;primaryKey"`
	NamaMakanan  string         `json:"nama_makanan" gorm:"column:nama_makanan;type:varchar(100);not null"`
	Harga        float64        `json:"harga" gorm:"column:harga;not null"`
	Jenis        JenisMenu      `json:"jenis" gorm:"column:jenis;type:varchar(20);not null"`
	Foto         string         `json:"foto" gorm:"column:foto;type:varchar(255)"`
	Deskripsi    string         `json:"deskripsi" gorm:"column:deskripsi;type:text"`
	Stock        int            `json:"stock" gorm:"column:stock;default:0"`
	IsAvailable  bool           `json:"is_available" gorm:"column:is_available;default:true"`
	IDStan       uint           `json:"id_stan" gorm:"column:id_stan;not null"`
	CreatedBy    string         `json:"created_by" gorm:"column:created_by"`
	UpdatedBy    string         `json:"updated_by" gorm:"column:updated_by"`
	CreatedAt    time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;index"`
	
	// Relations
	Stan            Stan              `json:"stan" gorm:"foreignKey:IDStan;constraint:OnDelete:CASCADE"`
	DetailTransaksi []DetailTransaksi `json:"detail_transaksi,omitempty" gorm:"foreignKey:IDMenu"`
	MenuDiskon      []MenuDiskon      `json:"menu_diskon,omitempty" gorm:"foreignKey:IDMenu"`
}
