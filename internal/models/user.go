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
	ID           uint           `json:"id" gorm:"column:id;primaryKey"`
	Username     string         `json:"username" gorm:"column:username;type:varchar(100);unique;not null"`
	Password     string         `json:"-" gorm:"column:password;type:varchar(100);not null"`
	Role         UserRole       `json:"role" gorm:"column:role;type:varchar(20);not null"`
	CreatedBy    string         `json:"created_by" gorm:"column:created_by"`
	UpdatedBy    string         `json:"updated_by" gorm:"column:updated_by"`
	CreatedAt    time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;index"`
	
	// Relations
	Stan  *Stan  `json:"stan,omitempty" gorm:"foreignKey:IDUser"`
	Siswa *Siswa `json:"siswa,omitempty" gorm:"foreignKey:IDUser"`
}