package services

import (
	"swipeup-be/internal/models"
	"time"

	"gorm.io/gorm"
)

// SuperadminService provides superadmin-specific operations
type SuperadminService struct {
	db *gorm.DB
}

func NewSuperadminService(db *gorm.DB) *SuperadminService {
	return &SuperadminService{db: db}
}

// StanRevenue represents revenue data for a stan
type StanRevenue struct {
	StanID       uint    `json:"stan_id"`
	NamaStan     string  `json:"nama_stan"`
	TotalRevenue float64 `json:"total_revenue"`
	TotalOrders  int     `json:"total_orders"`
}

// RevenueReport represents a comprehensive revenue report
type RevenueReport struct {
	TotalRevenue       float64       `json:"total_revenue"`
	TotalOrders        int           `json:"total_orders"`
	StanRevenues       []StanRevenue `json:"stan_revenues"`
	StartDate          time.Time     `json:"start_date"`
	EndDate            time.Time     `json:"end_date"`
}

// GetRevenueByStanID calculates total revenue for a specific stan
func (s *SuperadminService) GetRevenueByStanID(stanID uint, startDate, endDate time.Time) (*StanRevenue, error) {
	var stan models.Stan
	if err := s.db.First(&stan, stanID).Error; err != nil {
		return nil, err
	}

	var result struct {
		TotalRevenue float64
		TotalOrders  int
	}

	query := s.db.Model(&models.DetailTransaksi{}).
		Joins("JOIN transaksis ON detail_transaksis.id_transaksi = transaksis.id").
		Where("transaksis.id_stan = ?", stanID)

	if !startDate.IsZero() && !endDate.IsZero() {
		query = query.Where("transaksis.tanggal BETWEEN ? AND ?", startDate, endDate)
	}

	query.Select("COALESCE(SUM(detail_transaksis.qty * detail_transaksis.harga_beli), 0) as total_revenue, COUNT(DISTINCT transaksis.id) as total_orders").
		Scan(&result)

	return &StanRevenue{
		StanID:       stanID,
		NamaStan:     stan.NamaStan,
		TotalRevenue: result.TotalRevenue,
		TotalOrders:  result.TotalOrders,
	}, nil
}

// GetAllStanRevenue calculates revenue for all stans
func (s *SuperadminService) GetAllStanRevenue(startDate, endDate time.Time) ([]StanRevenue, error) {
	var stans []models.Stan
	if err := s.db.Find(&stans).Error; err != nil {
		return nil, err
	}

	var revenues []StanRevenue
	for _, stan := range stans {
		revenue, err := s.GetRevenueByStanID(stan.ID, startDate, endDate)
		if err != nil {
			continue
		}
		revenues = append(revenues, *revenue)
	}

	return revenues, nil
}

// GetRevenueReport generates a comprehensive revenue report
func (s *SuperadminService) GetRevenueReport(startDate, endDate time.Time) (*RevenueReport, error) {
	stanRevenues, err := s.GetAllStanRevenue(startDate, endDate)
	if err != nil {
		return nil, err
	}

	totalRevenue := 0.0
	totalOrders := 0
	for _, revenue := range stanRevenues {
		totalRevenue += revenue.TotalRevenue
		totalOrders += revenue.TotalOrders
	}

	return &RevenueReport{
		TotalRevenue: totalRevenue,
		TotalOrders:  totalOrders,
		StanRevenues: stanRevenues,
		StartDate:    startDate,
		EndDate:      endDate,
	}, nil
}

// GetGlobalDiscounts retrieves all global discounts
func (s *SuperadminService) GetGlobalDiscounts() ([]models.Diskon, error) {
	var diskon []models.Diskon
	err := s.db.Where("tipe_diskon = ? OR id_stan IS NULL", models.DiskonGlobal).
		Preload("Stan").
		Find(&diskon).Error
	return diskon, err
}

// CreateGlobalDiscount creates a new global discount
func (s *SuperadminService) CreateGlobalDiscount(diskon *models.Diskon) error {
	diskon.TipeDiskon = models.DiskonGlobal
	diskon.IDStan = nil // Global discounts don't have a stan
	return s.db.Create(diskon).Error
}

// UpdateGlobalDiscount updates a global discount
func (s *SuperadminService) UpdateGlobalDiscount(id uint, updates map[string]interface{}) error {
	// Ensure it's a global discount
	var diskon models.Diskon
	if err := s.db.First(&diskon, id).Error; err != nil {
		return err
	}

	if diskon.TipeDiskon != models.DiskonGlobal && diskon.IDStan != nil {
		return gorm.ErrRecordNotFound
	}

	return s.db.Model(&models.Diskon{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteGlobalDiscount deletes a global discount
func (s *SuperadminService) DeleteGlobalDiscount(id uint) error {
	// Ensure it's a global discount
	var diskon models.Diskon
	if err := s.db.First(&diskon, id).Error; err != nil {
		return err
	}

	if diskon.TipeDiskon != models.DiskonGlobal && diskon.IDStan != nil {
		return gorm.ErrRecordNotFound
	}

	return s.db.Delete(&models.Diskon{}, id).Error
}

// GetStanStatistics returns statistics for a specific stan
type StanStatistics struct {
	StanID          uint    `json:"stan_id"`
	NamaStan        string  `json:"nama_stan"`
	TotalMenu       int     `json:"total_menu"`
	AvailableMenu   int     `json:"available_menu"`
	TotalOrders     int     `json:"total_orders"`
	TotalRevenue    float64 `json:"total_revenue"`
	AverageOrder    float64 `json:"average_order"`
}

func (s *SuperadminService) GetStanStatistics(stanID uint, startDate, endDate time.Time) (*StanStatistics, error) {
	var stan models.Stan
	if err := s.db.Preload("Menu").First(&stan, stanID).Error; err != nil {
		return nil, err
	}

	totalMenu := len(stan.Menu)
	availableMenu := 0
	for _, menu := range stan.Menu {
		if menu.IsAvailable && menu.Stock > 0 {
			availableMenu++
		}
	}

	var result struct {
		TotalRevenue float64
		TotalOrders  int
	}

	query := s.db.Model(&models.DetailTransaksi{}).
		Joins("JOIN transaksis ON detail_transaksis.id_transaksi = transaksis.id").
		Where("transaksis.id_stan = ?", stanID)

	if !startDate.IsZero() && !endDate.IsZero() {
		query = query.Where("transaksis.tanggal BETWEEN ? AND ?", startDate, endDate)
	}

	query.Select("COALESCE(SUM(detail_transaksis.qty * detail_transaksis.harga_beli), 0) as total_revenue, COUNT(DISTINCT transaksis.id) as total_orders").
		Scan(&result)

	averageOrder := 0.0
	if result.TotalOrders > 0 {
		averageOrder = result.TotalRevenue / float64(result.TotalOrders)
	}

	return &StanStatistics{
		StanID:        stanID,
		NamaStan:      stan.NamaStan,
		TotalMenu:     totalMenu,
		AvailableMenu: availableMenu,
		TotalOrders:   result.TotalOrders,
		TotalRevenue:  result.TotalRevenue,
		AverageOrder:  averageOrder,
	}, nil
}

// GetAllStanStatistics returns statistics for all stans
func (s *SuperadminService) GetAllStanStatistics(startDate, endDate time.Time) ([]StanStatistics, error) {
	var stans []models.Stan
	if err := s.db.Find(&stans).Error; err != nil {
		return nil, err
	}

	var statistics []StanStatistics
	for _, stan := range stans {
		stat, err := s.GetStanStatistics(stan.ID, startDate, endDate)
		if err != nil {
			continue
		}
		statistics = append(statistics, *stat)
	}

	return statistics, nil
}
