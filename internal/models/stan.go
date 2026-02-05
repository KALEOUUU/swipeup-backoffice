package models

import (
	"time"

	"gorm.io/gorm"
)

type Stan struct {
	ID              uint           `json:"id" gorm:"column:id;primaryKey"`
	NamaStan        string         `json:"nama_stan" gorm:"column:nama_stan;type:varchar(100);not null"`
	NamaPemilik     string         `json:"nama_pemilik" gorm:"column:nama_pemilik;type:varchar(100);not null"`
	Telp            string         `json:"telp" gorm:"column:telp;type:varchar(20)"`
	Foto            string         `json:"foto" gorm:"column:foto;type:varchar(255)"`
	QrisImage       string         `json:"qris_image" gorm:"column:qris_image;type:varchar(255)"`
	AcceptCash      bool           `json:"accept_cash" gorm:"column:accept_cash;default:true"`
	AcceptQris      bool           `json:"accept_qris" gorm:"column:accept_qris;default:false"`
	IDUser          uint           `json:"id_user" gorm:"column:id_user;not null"`
	CreatedBy       string         `json:"created_by" gorm:"column:created_by"`
	UpdatedBy       string         `json:"updated_by" gorm:"column:updated_by"`
	CreatedAt       time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt       time.Time      `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;index"`
	
	// Relations
	User        User         `json:"user" gorm:"foreignKey:IDUser;constraint:OnDelete:CASCADE"`
	Menu        []Menu       `json:"menu,omitempty" gorm:"foreignKey:IDStan"`
	Transaksi   []Transaksi  `json:"transaksi,omitempty" gorm:"foreignKey:IDStan"`
}
