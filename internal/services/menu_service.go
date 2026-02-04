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

// UpdateStock updates the stock of a menu item
func (s *MenuService) UpdateStock(id uint, stock int) error {
	isAvailable := stock > 0
	return s.GetDB().Model(&models.Menu{}).Where("id = ?", id).Updates(map[string]interface{}{
		"stock":        stock,
		"is_available": isAvailable,
	}).Error
}

// AdjustStock adjusts stock by delta (positive or negative)
func (s *MenuService) AdjustStock(id uint, delta int) error {
	return s.GetDB().Transaction(func(tx *gorm.DB) error {
		var menu models.Menu
		if err := tx.First(&menu, id).Error; err != nil {
			return err
		}
		
		newStock := menu.Stock + delta
		if newStock < 0 {
			newStock = 0
		}
		
		isAvailable := newStock > 0
		return tx.Model(&menu).Updates(map[string]interface{}{
			"stock":        newStock,
			"is_available": isAvailable,
		}).Error
	})
}

// GetAvailableMenuByStanID gets only available (in-stock) menu items
func (s *MenuService) GetAvailableMenuByStanID(stanID uint) ([]models.Menu, error) {
	var menus []models.Menu
	err := s.GetDB().Preload("Stan").Where("id_stan = ? AND is_available = ? AND stock > 0", stanID, true).Find(&menus).Error
	return menus, err
}
