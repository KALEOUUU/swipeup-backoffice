package services

import (
	"swipeup-be/internal/models"

	"gorm.io/gorm"
)

type StanService struct {
	*BaseService[models.Stan]
}

func NewStanService(db *gorm.DB) *StanService {
	return &StanService{
		BaseService: NewBaseService[models.Stan](db),
	}
}

func (s *StanService) GetByUserID(userID uint) (*models.Stan, error) {
	var stan models.Stan
	err := s.GetDB().Where("id_user = ?", userID).Preload("User").First(&stan).Error
	if err != nil {
		return nil, err
	}
	return &stan, nil
}

func (s *StanService) UpdateFields(id uint, updates map[string]interface{}) error {
	return s.GetDB().Model(&models.Stan{}).Where("id = ?", id).Updates(updates).Error
}

func (s *StanService) GetWithMenu(id uint) (*models.Stan, error) {
	return s.FindByID(id, "User", "Menu")
}

func (s *StanService) GetWithTransaksi(id uint) (*models.Stan, error) {
	return s.FindByID(id, "User", "Transaksi", "Transaksi.Siswa", "Transaksi.DetailTransaksi")
}
