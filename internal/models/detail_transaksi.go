package models

import (
	"time"

	"gorm.io/gorm"
)

type DetailTransaksi struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	IDTransaksi  uint      `json:"id_transaksi" gorm:"not null"`
	IDMenu       uint      `json:"id_menu" gorm:"not null"`
	Qty          int       `json:"qty" gorm:"not null"`
	HargaBeli    float64   `json:"harga_beli" gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	
	// Relations
	Transaksi Transaksi `json:"transaksi" gorm:"foreignKey:IDTransaksi;constraint:OnDelete:CASCADE"`
	Menu      Menu      `json:"menu" gorm:"foreignKey:IDMenu;constraint:OnDelete:CASCADE"`
}
