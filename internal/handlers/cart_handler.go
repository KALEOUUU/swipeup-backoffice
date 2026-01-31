package handlers

import (
	"swipeup-be/internal/models"
	"swipeup-be/internal/services"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	service            *services.CartService
	activityLogService *services.ActivityLogService
}

func NewCartHandler(service *services.CartService, activityLogService *services.ActivityLogService) *CartHandler {
	return &CartHandler{
		service:            service,
		activityLogService: activityLogService,
	}
}

// AddToCart adds item to cart
func (h *CartHandler) AddToCart(c *gin.Context) {
	var cart models.Cart
	if err := c.ShouldBindJSON(&cart); err != nil {
		BadRequestResponse(c, "Invalid request body", err)
		return
	}

	if err := h.service.AddToCart(&cart); err != nil {
		InternalErrorResponse(c, "Failed to add item to cart", err)
		return
	}

	// Log add to cart activity
	if userID, exists := GetUserIDFromContext(c); exists {
		ip, userAgent := GetClientInfo(c)
		h.activityLogService.LogActivity(userID, "add_to_cart", "Added item to cart", ip, userAgent)
	}

	SuccessResponse(c, "Item added to cart successfully", cart)
}

// GetCart gets all cart items for a siswa
func (h *CartHandler) GetCart(c *gin.Context) {
	siswaID, err := GetQueryParamUint(c, "siswa_id")
	if err != nil || siswaID == 0 {
		BadRequestResponse(c, "Invalid siswa_id parameter", err)
		return
	}

	carts, err := h.service.GetCartBySiswaID(siswaID)
	if err != nil {
		InternalErrorResponse(c, "Failed to get cart", err)
		return
	}

	// Calculate total
	totalItems, totalPrice, err := h.service.GetCartTotal(siswaID)
	if err != nil {
		InternalErrorResponse(c, "Failed to calculate cart total", err)
		return
	}

	response := gin.H{
		"items":       carts,
		"total_items": totalItems,
		"total_price": totalPrice,
	}

	SuccessResponse(c, "Cart retrieved successfully", response)
}

// UpdateCartItem updates cart item quantity
func (h *CartHandler) UpdateCartItem(c *gin.Context) {
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

	if err := h.service.UpdateCartItem(cartID, req.Qty); err != nil {
		InternalErrorResponse(c, "Failed to update cart item", err)
		return
	}

	SuccessResponse(c, "Cart item updated successfully", nil)
}

// RemoveFromCart removes item from cart
func (h *CartHandler) RemoveFromCart(c *gin.Context) {
	cartID, err := GetIDParam(c)
	if err != nil {
		BadRequestResponse(c, "Invalid cart ID", err)
		return
	}

	if err := h.service.RemoveFromCart(cartID); err != nil {
		InternalErrorResponse(c, "Failed to remove item from cart", err)
		return
	}

	SuccessResponse(c, "Item removed from cart successfully", nil)
}

// ClearCart clears all items from cart
func (h *CartHandler) ClearCart(c *gin.Context) {
	siswaID, err := GetQueryParamUint(c, "siswa_id")
	if err != nil || siswaID == 0 {
		BadRequestResponse(c, "Invalid siswa_id parameter", err)
		return
	}

	if err := h.service.ClearCart(siswaID); err != nil {
		InternalErrorResponse(c, "Failed to clear cart", err)
		return
	}

	SuccessResponse(c, "Cart cleared successfully", nil)
}

// CheckoutCart converts cart to transaction
func (h *CartHandler) CheckoutCart(c *gin.Context) {
	var req struct {
		SiswaID uint `json:"siswa_id" binding:"required"`
		StanID  uint `json:"stan_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, "Invalid request body", err)
		return
	}

	// Get cart details
	details, err := h.service.CheckoutCart(req.SiswaID, req.StanID)
	if err != nil {
		InternalErrorResponse(c, "Failed to checkout cart", err)
		return
	}

	if len(details) == 0 {
		BadRequestResponse(c, "Cart is empty", nil)
		return
	}

	// Create transaction with cart details
	transaksi := &models.Transaksi{
		IDStan:  req.StanID,
		IDSiswa: req.SiswaID,
		Status:  models.StatusBelumDikonfirm,
	}

	// Note: This would need transaksiService to create with details
	// For now, return the details that should be used for transaction creation
	response := gin.H{
		"transaksi": transaksi,
		"details":   details,
		"message":   "Cart converted to transaction details. Use these details to create transaction.",
	}

	// Log checkout activity
	if userID, exists := GetUserIDFromContext(c); exists {
		ip, userAgent := GetClientInfo(c)
		h.activityLogService.LogActivity(userID, "checkout", "Completed cart checkout", ip, userAgent)
	}

	SuccessResponse(c, "Cart checkout prepared successfully", response)
}