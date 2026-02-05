package services

import (
	"swipeup-be/internal/models"

	"gorm.io/gorm"
)

// StudentService provides student-specific operations
type StudentService struct {
	db            *gorm.DB
	siswaService  *SiswaService
	cartService   *CartService
	transaksiService *TransaksiService
}

func NewStudentService(
	db *gorm.DB,
	siswaService *SiswaService,
	cartService *CartService,
	transaksiService *TransaksiService,
) *StudentService {
	return &StudentService{
		db:            db,
		siswaService:  siswaService,
		cartService:   cartService,
		transaksiService: transaksiService,
	}
}

// GetSiswaByUserID retrieves the siswa profile for the authenticated user
func (s *StudentService) GetSiswaByUserID(userID uint) (*models.Siswa, error) {
	return s.siswaService.GetByUserID(userID)
}

// UpdateSiswaProfile updates the siswa profile for the authenticated user
func (s *StudentService) UpdateSiswaProfile(userID uint, updates map[string]interface{}) error {
	siswa, err := s.siswaService.GetByUserID(userID)
	if err != nil {
		return err
	}
	return s.siswaService.UpdateFields(siswa.ID, updates)
}

// AddToCart adds an item to the student's cart
func (s *StudentService) AddToCart(siswaID uint, cart *models.Cart) error {
	cart.IDSiswa = siswaID
	return s.cartService.AddToCart(cart)
}

// GetCart retrieves the student's cart
func (s *StudentService) GetCart(siswaID uint) ([]models.Cart, int, float64, error) {
	carts, err := s.cartService.GetCartBySiswaID(siswaID)
	if err != nil {
		return nil, 0, 0, err
	}

	totalItems, totalPrice, err := s.cartService.GetCartTotal(siswaID)
	if err != nil {
		return nil, 0, 0, err
	}

	return carts, totalItems, totalPrice, nil
}

// UpdateCartItem updates a cart item
func (s *StudentService) UpdateCartItem(cartID uint, qty int) error {
	return s.cartService.UpdateCartItem(cartID, qty)
}

// RemoveFromCart removes an item from the cart
func (s *StudentService) RemoveFromCart(cartID uint) error {
	return s.cartService.RemoveFromCart(cartID)
}

// ClearCart clears all items from the cart
func (s *StudentService) ClearCart(siswaID uint) error {
	return s.cartService.ClearCart(siswaID)
}

// CheckoutCart converts cart items to a transaction
func (s *StudentService) CheckoutCart(siswaID uint, stanID uint) (*models.Transaksi, []models.DetailTransaksi, error) {
	// Get cart details
	details, err := s.cartService.CheckoutCart(siswaID, stanID)
	if err != nil {
		return nil, nil, err
	}

	if len(details) == 0 {
		return nil, nil, gorm.ErrRecordNotFound
	}

	// Create transaction with cart details
	transaksi := &models.Transaksi{
		IDStan:  stanID,
		IDSiswa: siswaID,
		Status:  models.StatusBelumDikonfirm,
	}

	if err := s.transaksiService.CreateWithDetails(transaksi, details); err != nil {
		return nil, nil, err
	}

	// Get full transaction details
	fullTransaksi, err := s.transaksiService.GetWithFullDetails(transaksi.ID)
	if err != nil {
		return nil, nil, err
	}

	return fullTransaksi, details, nil
}

// GetTransactions retrieves all transactions for the student
func (s *StudentService) GetTransactions(siswaID uint) ([]models.Transaksi, error) {
	return s.transaksiService.GetBySiswaID(siswaID)
}

// GetTransactionByID retrieves a specific transaction for the student
func (s *StudentService) GetTransactionByID(siswaID uint, transaksiID uint) (*models.Transaksi, error) {
	// Verify transaction belongs to student
	var transaksi models.Transaksi
	if err := s.db.First(&transaksi, transaksiID).Error; err != nil {
		return nil, err
	}

	if transaksi.IDSiswa != siswaID {
		return nil, gorm.ErrRecordNotFound
	}

	return s.transaksiService.GetWithFullDetails(transaksiID)
}
