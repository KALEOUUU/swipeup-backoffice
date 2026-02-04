package models

import (
	"time"

	"gorm.io/gorm"
)

type Stan struct {
	ID              uint   `json:"id" gorm:"primaryKey"`
	NamaStan        string `json:"nama_stan" gorm:"type:varchar(100);not null"`
	NamaPemilik     string `json:"nama_pemilik" gorm:"type:varchar(100);not null"`
	Telp            string `json:"telp" gorm:"type:varchar(20)"`
	Foto            string `json:"foto" gorm:"type:varchar(255)"`
	QrisImage       string `json:"qris_image" gorm:"type:varchar(255)"`
	AcceptCash      bool   `json:"accept_cash" gorm:"default:true"`
	AcceptQris      bool   `json:"accept_qris" gorm:"default:false"`
	IDUser          uint   `json:"id_user" gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	
	// Relations
	User        User         `json:"user" gorm:"foreignKey:IDUser;constraint:OnDelete:CASCADE"`
	Menu        []Menu       `json:"menu,omitempty" gorm:"foreignKey:IDStan"`
	Transaksi   []Transaksi  `json:"transaksi,omitempty" gorm:"foreignKey:IDStan"`
}
