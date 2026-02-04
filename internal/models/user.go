package models

import (
	"time"

	"gorm.io/gorm"
)

type UserRole string

const (
	RoleSuperAdmin UserRole = "superadmin"
	RoleAdminStan  UserRole = "admin_stan"
	RoleSiswa      UserRole = "siswa"
)

type User struct {
	ID           uint     `json:"id" gorm:"primaryKey"`
	Username     string   `json:"username" gorm:"type:varchar(100);unique;not null"`
	Password     string   `json:"-" gorm:"type:varchar(100);not null"`
	Role         UserRole `json:"role" gorm:"type:varchar(20);not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	
	// Relations
	Stan  *Stan  `json:"stan,omitempty" gorm:"foreignKey:IDUser"`
	Siswa *Siswa `json:"siswa,omitempty" gorm:"foreignKey:IDUser"`
}