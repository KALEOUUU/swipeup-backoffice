package services

import (
	"swipeup-be/internal/models"
	"time"

	"gorm.io/gorm"
)

type DiskonService struct {
	*BaseService[models.Diskon]
}

func NewDiskonService(db *gorm.DB) *DiskonService {
	return &DiskonService{
		BaseService: NewBaseService[models.Diskon](db),
	}
}

func (s *DiskonService) GetActiveDiskon() ([]models.Diskon, error) {
	var diskon []models.Diskon
	now := time.Now()
	err := s.GetDB().Where("tanggal_awal <= ? AND tanggal_akhir >= ?", now, now).Preload("Stan").Find(&diskon).Error
	return diskon, err
}

func (s *DiskonService) GetActiveDiskonByStan(stanID uint) ([]models.Diskon, error) {
	var diskon []models.Diskon
	now := time.Now()
	// Get global discounts (id_stan is NULL) OR discounts for specific stan
	err := s.GetDB().Where("tanggal_awal <= ? AND tanggal_akhir >= ? AND (id_stan IS NULL OR id_stan = ?)", 
		now, now, stanID).Preload("Stan").Find(&diskon).Error
	return diskon, err
}

func (s *DiskonService) GetByStanID(stanID uint) ([]models.Diskon, error) {
	var diskon []models.Diskon
	err := s.GetDB().Where("id_stan = ?", stanID).Preload("Stan").Find(&diskon).Error
	return diskon, err
}

func (s *DiskonService) GetGlobalDiskon() ([]models.Diskon, error) {
	var diskon []models.Diskon
	err := s.GetDB().Where("tipe_diskon = ? OR id_stan IS NULL", models.DiskonGlobal).Find(&diskon).Error
	return diskon, err
}

func (s *DiskonService) UpdateFields(id uint, updates map[string]interface{}) error {
	return s.GetDB().Model(&models.Diskon{}).Where("id = ?", id).Updates(updates).Error
}

func (s *DiskonService) GetByDateRange(startDate, endDate time.Time) ([]models.Diskon, error) {
	var diskon []models.Diskon
	err := s.GetDB().Where("(tanggal_awal BETWEEN ? AND ?) OR (tanggal_akhir BETWEEN ? AND ?)", 
		startDate, endDate, startDate, endDate).Find(&diskon).Error
	return diskon, err
}

func (s *DiskonService) AssignToMenu(diskonID, menuID uint) error {
	menuDiskon := models.MenuDiskon{
		IDMenu:   menuID,
		IDDiskon: diskonID,
	}
	return s.GetDB().Create(&menuDiskon).Error
}

func (s *DiskonService) RemoveFromMenu(diskonID, menuID uint) error {
	return s.GetDB().Where("id_menu = ? AND id_diskon = ?", menuID, diskonID).Delete(&models.MenuDiskon{}).Error
}
