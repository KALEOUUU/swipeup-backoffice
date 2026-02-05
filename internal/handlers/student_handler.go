package handlers

import (
	"swipeup-be/internal/models"
	"swipeup-be/internal/services"

	"github.com/gin-gonic/gin"
)

type StudentHandler struct {
	studentService *services.StudentService
}

func NewStudentHandler(studentService *services.StudentService) *StudentHandler {
	return &StudentHandler{
		studentService: studentService,
	}
}

// GetSiswaProfile retrieves the siswa profile for the authenticated student
func (h *StudentHandler) GetSiswaProfile(c *gin.Context) {
	userID, exists := GetUserIDFromContext(c)
	if !exists {
		BadRequestResponse(c, "User not authenticated", nil)
		return
	}

	siswa, err := h.studentService.GetSiswaByUserID(userID)
	if err != nil {
		NotFoundResponse(c, "Siswa profile not found")
		return
	}

	SuccessResponse(c, "Siswa profile retrieved successfully", siswa)
}

// UpdateSiswaProfile updates the siswa profile for the authenticated student
func (h *StudentHandler) UpdateSiswaProfile(c *gin.Context) {
	userID, exists := GetUserIDFromContext(c)
	if !exists {
		BadRequestResponse(c, "User not authenticated", nil)
		return
	}

	// Bind request body to a map to only update provided fields
	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		BadRequestResponse(c, "Invalid request body", err)
		return
	}

	// Build updates map with only allowed fields
	updates := make(map[string]interface{})
	if namaSiswa, ok := updateData["nama_siswa"].(string); ok {
		updates["nama_siswa"] = namaSiswa
	}
	if kelas, ok := updateData["kelas"].(string); ok {
		updates["kelas"] = kelas
	}
	if alamat, ok := updateData["alamat"].(string); ok {
		updates["alamat"] = alamat
	}
	if telp, ok := updateData["telp"].(string); ok {
		updates["telp"] = telp
	}

	if len(updates) == 0 {
		BadRequestResponse(c, "No valid fields to update", nil)
		return
	}

	if err := h.studentService.UpdateSiswaProfile(userID, updates); err != nil {
		InternalErrorResponse(c, "Failed to update siswa profile", err)
		return
	}

	// Get updated siswa
	siswa, _ := h.studentService.GetSiswaByUserID(userID)
	SuccessResponse(c, "Siswa profile updated successfully", siswa)
}

// AddToCart adds an item to the student's cart
func (h *StudentHandler) AddToCart(c *gin.Context) {
	userID, exists := GetUserIDFromContext(c)
	if !exists {
		BadRequestResponse(c, "User not authenticated", nil)
		return
	}

	// Get siswa ID from user
	siswa, err := h.studentService.GetSiswaByUserID(userID)
	if err != nil {
		NotFoundResponse(c, "Siswa profile not found")
		return
	}

	var cart models.Cart
	if err := c.ShouldBindJSON(&cart); err != nil {
		BadRequestResponse(c, "Invalid request body", err)
		return
	}

	if err := h.studentService.AddToCart(siswa.ID, &cart); err != nil {
		InternalErrorResponse(c, "Failed to add item to cart", err)
		return
	}

	SuccessResponse(c, "Item added to cart successfully", cart)
}

// GetCart retrieves the student's cart
func (h *StudentHandler) GetCart(c *gin.Context) {
	userID, exists := GetUserIDFromContext(c)
	if !exists {
		BadRequestResponse(c, "User not authenticated", nil)
		return
	}

	// Get siswa ID from user
	siswa, err := h.studentService.GetSiswaByUserID(userID)
	if err != nil {
		NotFoundResponse(c, "Siswa profile not found")
		return
	}

	carts, totalItems, totalPrice, err := h.studentService.GetCart(siswa.ID)
	if err != nil {
		InternalErrorResponse(c, "Failed to get cart", err)
		return
	}

	response := gin.H{
		"items":       carts,
		"total_items": totalItems,
		"total_price": totalPrice,
	}

	SuccessResponse(c, "Cart retrieved successfully", response)
}

// UpdateCartItem updates a cart item
func (h *StudentHandler) UpdateCartItem(c *gin.Context) {
	cartID, err := GetIDParam(c)
	if err != nil {
		BadRequestResponse(c, "Invalid cart ID", err)
		return
	}

	var req struct {
		Qty int `json:"qty" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, "Invalid request body", err)
		return
	}

	if err := h.studentService.UpdateCartItem(cartID, req.Qty); err != nil {
		InternalErrorResponse(c, "Failed to update cart item", err)
		return
	}

	SuccessResponse(c, "Cart item updated successfully", nil)
}

// RemoveFromCart removes an item from the cart
func (h *StudentHandler) RemoveFromCart(c *gin.Context) {
	cartID, err := GetIDParam(c)
	if err != nil {
		BadRequestResponse(c, "Invalid cart ID", err)
		return
	}

	if err := h.studentService.RemoveFromCart(cartID); err != nil {
		InternalErrorResponse(c, "Failed to remove item from cart", err)
		return
	}

	SuccessResponse(c, "Item removed from cart successfully", nil)
}

// ClearCart clears all items from the cart
func (h *StudentHandler) ClearCart(c *gin.Context) {
	userID, exists := GetUserIDFromContext(c)
	if !exists {
		BadRequestResponse(c, "User not authenticated", nil)
		return
	}

	// Get siswa ID from user
	siswa, err := h.studentService.GetSiswaByUserID(userID)
	if err != nil {
		NotFoundResponse(c, "Siswa profile not found")
		return
	}

	if err := h.studentService.ClearCart(siswa.ID); err != nil {
		InternalErrorResponse(c, "Failed to clear cart", err)
		return
	}

	SuccessResponse(c, "Cart cleared successfully", nil)
}

// CheckoutCart converts cart items to a transaction
func (h *StudentHandler) CheckoutCart(c *gin.Context) {
	userID, exists := GetUserIDFromContext(c)
	if !exists {
		BadRequestResponse(c, "User not authenticated", nil)
		return
	}

	// Get siswa ID from user
	siswa, err := h.studentService.GetSiswaByUserID(userID)
	if err != nil {
		NotFoundResponse(c, "Siswa profile not found")
		return
	}

	var req struct {
		StanID uint `json:"stan_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, "Invalid request body", err)
		return
	}

	transaksi, details, err := h.studentService.CheckoutCart(siswa.ID, req.StanID)
	if err != nil {
		if err.Error() == "record not found" {
			BadRequestResponse(c, "Cart is empty", nil)
		} else {
			InternalErrorResponse(c, "Failed to checkout cart", err)
		}
		return
	}

	response := gin.H{
		"transaksi": transaksi,
		"details":   details,
		"message":   "Checkout successful",
	}

	CreatedResponse(c, "Checkout successful", response)
}

// GetTransactions retrieves all transactions for the student
func (h *StudentHandler) GetTransactions(c *gin.Context) {
	userID, exists := GetUserIDFromContext(c)
	if !exists {
		BadRequestResponse(c, "User not authenticated", nil)
		return
	}

	// Get siswa ID from user
	siswa, err := h.studentService.GetSiswaByUserID(userID)
	if err != nil {
		NotFoundResponse(c, "Siswa profile not found")
		return
	}

	transaksi, err := h.studentService.GetTransactions(siswa.ID)
	if err != nil {
		InternalErrorResponse(c, "Failed to get transactions", err)
		return
	}

	SuccessResponse(c, "Transactions retrieved successfully", transaksi)
}

// GetTransactionByID retrieves a specific transaction for the student
func (h *StudentHandler) GetTransactionByID(c *gin.Context) {
	userID, exists := GetUserIDFromContext(c)
	if !exists {
		BadRequestResponse(c, "User not authenticated", nil)
		return
	}

	// Get siswa ID from user
	siswa, err := h.studentService.GetSiswaByUserID(userID)
	if err != nil {
		NotFoundResponse(c, "Siswa profile not found")
		return
	}

	transaksiID, err := GetIDParam(c)
	if err != nil {
		BadRequestResponse(c, "Invalid transaction ID", err)
		return
	}

	transaksi, err := h.studentService.GetTransactionByID(siswa.ID, transaksiID)
	if err != nil {
		if err.Error() == "record not found" {
			NotFoundResponse(c, "Transaction not found or you don't have permission")
		} else {
			InternalErrorResponse(c, "Failed to get transaction", err)
		}
		return
	}

	SuccessResponse(c, "Transaction retrieved successfully", transaksi)
}
