package services

import (
	"swipeup-be/internal/models"
	"time"

	"gorm.io/gorm"
)

// StanAdminService provides stan_admin-specific operations
type StanAdminService struct {
	db            *gorm.DB
	stanService   *StanService
	menuService   *MenuService
	diskonService *DiskonService
}

func NewStanAdminService(
	db *gorm.DB,
	stanService *StanService,
	menuService *MenuService,
	diskonService *DiskonService,
) *StanAdminService {
	return &StanAdminService{
		db:            db,
		stanService:   stanService,
		menuService:   menuService,
		diskonService: diskonService,
	}
}

// GetStanByUserID retrieves the stan owned by the user
func (s *StanAdminService) GetStanByUserID(userID uint) (*models.Stan, error) {
	return s.stanService.GetByUserID(userID)
}

// UpdateStanProfile updates the stan profile for the authenticated user
func (s *StanAdminService) UpdateStanProfile(userID uint, updates map[string]interface{}) error {
	stan, err := s.stanService.GetByUserID(userID)
	if err != nil {
		return err
	}
	return s.stanService.UpdateFields(stan.ID, updates)
}

// CreateMenu creates a new menu item for the stan
func (s *StanAdminService) CreateMenu(userID uint, menu *models.Menu) error {
	stan, err := s.stanService.GetByUserID(userID)
	if err != nil {
		return err
	}
	menu.IDStan = stan.ID
	return s.menuService.Create(menu)
}

// UpdateMenu updates a menu item owned by the stan
func (s *StanAdminService) UpdateMenu(userID uint, menuID uint, updates map[string]interface{}) error {
	// Verify menu belongs to stan
	var menu models.Menu
	if err := s.db.First(&menu, menuID).Error; err != nil {
		return err
	}

	stan, err := s.stanService.GetByUserID(userID)
	if err != nil {
		return err
	}

	if menu.IDStan != stan.ID {
		return gorm.ErrRecordNotFound
	}

	return s.menuService.UpdateFields(menuID, updates)
}

// DeleteMenu deletes a menu item owned by the stan
func (s *StanAdminService) DeleteMenu(userID uint, menuID uint) error {
	// Verify menu belongs to stan
	var menu models.Menu
	if err := s.db.First(&menu, menuID).Error; err != nil {
		return err
	}

	stan, err := s.stanService.GetByUserID(userID)
	if err != nil {
		return err
	}

	if menu.IDStan != stan.ID {
		return gorm.ErrRecordNotFound
	}

	return s.menuService.Delete(menuID)
}

// GetMenusByStan retrieves all menu items for the stan
func (s *StanAdminService) GetMenusByStan(userID uint) ([]models.Menu, error) {
	stan, err := s.stanService.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	return s.menuService.GetByStanID(stan.ID)
}

// CreateStanDiscount creates a new stan-level discount
func (s *StanAdminService) CreateStanDiscount(userID uint, diskon *models.Diskon) error {
	stan, err := s.stanService.GetByUserID(userID)
	if err != nil {
		return err
	}
	diskon.TipeDiskon = models.DiskonStan
	diskon.IDStan = &stan.ID
	return s.diskonService.Create(diskon)
}

// CreateMenuDiscount creates a new menu-level discount
func (s *StanAdminService) CreateMenuDiscount(userID uint, diskon *models.Diskon, menuIDs []uint) error {
	stan, err := s.stanService.GetByUserID(userID)
	if err != nil {
		return err
	}

	// Verify all menus belong to stan
	for _, menuID := range menuIDs {
		var menu models.Menu
		if err := s.db.First(&menu, menuID).Error; err != nil {
			return err
		}
		if menu.IDStan != stan.ID {
			return gorm.ErrRecordNotFound
		}
	}

	diskon.TipeDiskon = models.DiskonMenu
	diskon.IDStan = &stan.ID
	if err := s.diskonService.Create(diskon); err != nil {
		return err
	}

	// Assign to menus
	for _, menuID := range menuIDs {
		if err := s.diskonService.AssignToMenu(diskon.ID, menuID); err != nil {
			return err
		}
	}

	return nil
}

// UpdateDiscount updates a discount owned by the stan
func (s *StanAdminService) UpdateDiscount(userID uint, diskonID uint, updates map[string]interface{}) error {
	// Verify discount belongs to stan
	var diskon models.Diskon
	if err := s.db.First(&diskon, diskonID).Error; err != nil {
		return err
	}

	stan, err := s.stanService.GetByUserID(userID)
	if err != nil {
		return err
	}

	if diskon.IDStan == nil || *diskon.IDStan != stan.ID {
		return gorm.ErrRecordNotFound
	}

	return s.diskonService.UpdateFields(diskonID, updates)
}

// DeleteDiscount deletes a discount owned by the stan
func (s *StanAdminService) DeleteDiscount(userID uint, diskonID uint) error {
	// Verify discount belongs to stan
	var diskon models.Diskon
	if err := s.db.First(&diskon, diskonID).Error; err != nil {
		return err
	}

	stan, err := s.stanService.GetByUserID(userID)
	if err != nil {
		return err
	}

	if diskon.IDStan == nil || *diskon.IDStan != stan.ID {
		return gorm.ErrRecordNotFound
	}

	return s.diskonService.Delete(diskonID)
}

// GetDiscountsByStan retrieves all discounts for the stan
func (s *StanAdminService) GetDiscountsByStan(userID uint) ([]models.Diskon, error) {
	stan, err := s.stanService.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	return s.diskonService.GetByStanID(stan.ID)
}

// GetActiveDiscountsByStan retrieves active discounts for the stan
func (s *StanAdminService) GetActiveDiscountsByStan(userID uint) ([]models.Diskon, error) {
	stan, err := s.stanService.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	return s.diskonService.GetActiveDiskonByStan(stan.ID)
}

// GetTransactionsByStan retrieves all transactions for the stan
func (s *StanAdminService) GetTransactionsByStan(userID uint) ([]models.Transaksi, error) {
	stan, err := s.stanService.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	var transaksiService TransaksiService
	return transaksiService.GetByStanID(stan.ID)
}

// GetTransactionsByStanAndDateRange retrieves transactions for the stan within a date range
func (s *StanAdminService) GetTransactionsByStanAndDateRange(userID uint, startDate, endDate time.Time) ([]models.Transaksi, error) {
	stan, err := s.stanService.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	var transaksi []models.Transaksi
	query := s.db.Preload("Siswa").Preload("DetailTransaksi").Preload("DetailTransaksi.Menu").
		Where("id_stan = ?", stan.ID)

	if !startDate.IsZero() && !endDate.IsZero() {
		query = query.Where("tanggal BETWEEN ? AND ?", startDate, endDate)
	}

	err = query.Find(&transaksi).Error
	return transaksi, err
}

// UpdateTransactionStatus updates the status of a transaction for the stan
func (s *StanAdminService) UpdateTransactionStatus(userID uint, transaksiID uint, status models.StatusTransaksi) error {
	// Verify transaction belongs to stan
	var transaksi models.Transaksi
	if err := s.db.First(&transaksi, transaksiID).Error; err != nil {
		return err
	}

	stan, err := s.stanService.GetByUserID(userID)
	if err != nil {
		return err
	}

	if transaksi.IDStan != stan.ID {
		return gorm.ErrRecordNotFound
	}

	var transaksiService TransaksiService
	return transaksiService.UpdateStatus(transaksiID, status)
}

// GetStanRevenue retrieves revenue for the stan
func (s *StanAdminService) GetStanRevenue(userID uint, startDate, endDate time.Time) (float64, int, error) {
	stan, err := s.stanService.GetByUserID(userID)
	if err != nil {
		return 0, 0, err
	}

	var result struct {
		TotalRevenue float64
		TotalOrders  int
	}

	query := s.db.Model(&models.DetailTransaksi{}).
		Joins("JOIN transaksis ON detail_transaksis.id_transaksi = transaksis.id").
		Where("transaksis.id_stan = ?", stan.ID)

	if !startDate.IsZero() && !endDate.IsZero() {
		query = query.Where("transaksis.tanggal BETWEEN ? AND ?", startDate, endDate)
	}

	query.Select("COALESCE(SUM(detail_transaksis.qty * detail_transaksis.harga_beli), 0) as total_revenue, COUNT(DISTINCT transaksis.id) as total_orders").
		Scan(&result)

	return result.TotalRevenue, result.TotalOrders, nil
}

// UpdatePaymentSettings updates payment settings for the stan
func (s *StanAdminService) UpdatePaymentSettings(userID uint, acceptCash, acceptQris bool, qrisImage string) error {
	stan, err := s.stanService.GetByUserID(userID)
	if err != nil {
		return err
	}

	updates := map[string]interface{}{
		"accept_cash": acceptCash,
		"accept_qris": acceptQris,
	}

	if qrisImage != "" {
		updates["qris_image"] = qrisImage
	}

	return s.stanService.UpdateFields(stan.ID, updates)
}
