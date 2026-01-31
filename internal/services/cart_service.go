package services

import (
	"swipeup-be/internal/models"

	"gorm.io/gorm"
)

type CartService struct {
	*BaseService[models.Cart]
	db *gorm.DB
}

func NewCartService(db *gorm.DB) *CartService {
	return &CartService{
		BaseService: NewBaseService[models.Cart](db),
		db:          db,
	}
}

// AddToCart adds item to cart or updates quantity if already exists
func (s *CartService) AddToCart(cart *models.Cart) error {
	var existingCart models.Cart
	err := s.db.Where("id_siswa = ? AND id_menu = ?", cart.IDSiswa, cart.IDMenu).First(&existingCart).Error

	if err == gorm.ErrRecordNotFound {
		// Create new cart item
		return s.db.Create(cart).Error
	} else if err != nil {
		return err
	} else {
		// Update existing cart item quantity
		existingCart.Qty += cart.Qty
		return s.db.Save(&existingCart).Error
	}
}

// GetCartBySiswaID gets all cart items for a siswa
func (s *CartService) GetCartBySiswaID(siswaID uint) ([]models.Cart, error) {
	var carts []models.Cart
	err := s.db.Where("id_siswa = ?", siswaID).
		Preload("Menu.Stan").
		Preload("Siswa").
		Find(&carts).Error
	return carts, err
}

// UpdateCartItem updates quantity of a cart item
func (s *CartService) UpdateCartItem(cartID uint, qty int) error {
	return s.db.Model(&models.Cart{}).Where("id = ?", cartID).Update("qty", qty).Error
}

// RemoveFromCart removes item from cart
func (s *CartService) RemoveFromCart(cartID uint) error {
	return s.db.Delete(&models.Cart{}, cartID).Error
}

// ClearCart removes all items from cart for a siswa
func (s *CartService) ClearCart(siswaID uint) error {
	return s.db.Where("id_siswa = ?", siswaID).Delete(&models.Cart{}).Error
}

// GetCartTotal calculates total items and price for a siswa's cart
func (s *CartService) GetCartTotal(siswaID uint) (int, float64, error) {
	var result struct {
		TotalItems int
		TotalPrice float64
	}

	err := s.db.Table("carts").
		Select("COALESCE(SUM(carts.qty), 0) as total_items, COALESCE(SUM(carts.qty * menus.harga), 0) as total_price").
		Joins("JOIN menus ON carts.id_menu = menus.id").
		Where("carts.id_siswa = ?", siswaID).
		Scan(&result).Error

	return result.TotalItems, result.TotalPrice, err
}

// CheckoutCart converts cart items to transaction details and clears cart
func (s *CartService) CheckoutCart(siswaID uint, stanID uint) ([]models.DetailTransaksi, error) {
	var carts []models.Cart
	err := s.db.Where("id_siswa = ?", siswaID).
		Preload("Menu").
		Find(&carts).Error
	if err != nil {
		return nil, err
	}

	var details []models.DetailTransaksi
	for _, cart := range carts {
		details = append(details, models.DetailTransaksi{
			IDMenu:    cart.IDMenu,
			Qty:       cart.Qty,
			HargaBeli: cart.Menu.Harga,
		})
	}

	// Clear cart after checkout
	if err := s.ClearCart(siswaID); err != nil {
		return nil, err
	}

	return details, nil
}