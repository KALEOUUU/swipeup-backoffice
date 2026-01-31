package services

import (
	"swipeup-be/internal/models"
	"time"

	"gorm.io/gorm"
)

type MenuService struct {
	*BaseService[models.Menu]
}

func NewMenuService(db *gorm.DB) *MenuService {
	return &MenuService{
		BaseService: NewBaseService[models.Menu](db),
	}
}

func (s *MenuService) GetByStanID(stanID uint) ([]models.Menu, error) {
	return s.FindWithCondition(map[string]interface{}{"id_stan": stanID}, "Stan")
}

func (s *MenuService) UpdateFields(id uint, updates map[string]interface{}) error {
	return s.GetDB().Model(&models.Menu{}).Where("id = ?", id).Updates(updates).Error
}

func (s *MenuService) GetByJenis(jenis models.JenisMenu) ([]models.Menu, error) {
	return s.FindWithCondition(map[string]interface{}{"jenis": jenis}, "Stan")
}

func (s *MenuService) GetWithDiskon(id uint) (*models.Menu, error) {
	return s.FindByID(id, "Stan", "MenuDiskon", "MenuDiskon.Diskon")
}

func (s *MenuService) GetMenuWithActiveDiskon(id uint) (*models.Menu, error) {
	var menu models.Menu
	now := time.Now()
	err := s.GetDB().Preload("Stan").
		Preload("MenuDiskon", "deleted_at IS NULL").
		Preload("MenuDiskon.Diskon", "tanggal_awal <= ? AND tanggal_akhir >= ?", now, now).
		First(&menu, id).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

func (s *MenuService) SearchByName(name string) ([]models.Menu, error) {
	var menus []models.Menu
	err := s.GetDB().Preload("Stan").Where("nama_makanan ILIKE ?", "%"+name+"%").Find(&menus).Error
	return menus, err
}
