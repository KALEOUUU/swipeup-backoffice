package models

import (
	"time"

	"gorm.io/gorm"
)

type DetailTransaksi struct {
	ID           uint           `json:"id" gorm:"column:id;primaryKey"`
	IDTransaksi  uint           `json:"id_transaksi" gorm:"column:id_transaksi;not null"`
	IDMenu       uint           `json:"id_menu" gorm:"column:id_menu;not null"`
	Qty          int            `json:"qty" gorm:"column:qty;not null"`
	HargaBeli    float64        `json:"harga_beli" gorm:"column:harga_beli;not null"`
	NamaDiskon   string         `json:"nama_diskon" gorm:"column:nama_diskon;type:varchar(100);default:''"`
	CreatedBy    string         `json:"created_by" gorm:"column:created_by"`
	UpdatedBy    string         `json:"updated_by" gorm:"column:updated_by"`
	CreatedAt    time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;index"`
	
	// Relations
	Transaksi Transaksi `json:"transaksi" gorm:"foreignKey:IDTransaksi;constraint:OnDelete:CASCADE"`
	Menu      Menu      `json:"menu" gorm:"foreignKey:IDMenu;constraint:OnDelete:CASCADE"`
}
