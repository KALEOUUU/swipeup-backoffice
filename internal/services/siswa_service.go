package services

import (
	"swipeup-be/internal/models"

	"gorm.io/gorm"
)

type SiswaService struct {
	*BaseService[models.Siswa]
}

func NewSiswaService(db *gorm.DB) *SiswaService {
	return &SiswaService{
		BaseService: NewBaseService[models.Siswa](db),
	}
}

func (s *SiswaService) GetByUserID(userID uint) (*models.Siswa, error) {
	var siswa models.Siswa
	err := s.GetDB().Where("id_user = ?", userID).Preload("User").First(&siswa).Error
	if err != nil {
		return nil, err
	}
	return &siswa, nil
}

func (s *SiswaService) UpdateFields(id uint, updates map[string]interface{}) error {
	return s.GetDB().Model(&models.Siswa{}).Where("id = ?", id).Updates(updates).Error
}

func (s *SiswaService) GetWithTransaksi(id uint) (*models.Siswa, error) {
	return s.FindByID(id, "User", "Transaksi", "Transaksi.DetailTransaksi")
}
