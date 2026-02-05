package handlers

import (
	"swipeup-be/internal/models"
	"swipeup-be/internal/services"
	"swipeup-be/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type StanAdminHandler struct {
	stanAdminService *services.StanAdminService
	menuService     *services.MenuService
}

func NewStanAdminHandler(
	stanAdminService *services.StanAdminService,
	menuService *services.MenuService,
) *StanAdminHandler {
	return &StanAdminHandler{
		stanAdminService: stanAdminService,
		menuService:     menuService,
	}
}

// GetStanProfile retrieves the stan profile for the authenticated admin
func (h *StanAdminHandler) GetStanProfile(c *gin.Context) {
	userID, exists := GetUserIDFromContext(c)
	if !exists {
		BadRequestResponse(c, "User not authenticated", nil)
		return
	}

	stan, err := h.stanAdminService.GetStanByUserID(userID)
	if err != nil {
		NotFoundResponse(c, "Stan not found")
		return
	}

	SuccessResponse(c, "Stan profile retrieved successfully", stan)
}

// UpdateStanProfile updates the stan profile for the authenticated admin
func (h *StanAdminHandler) UpdateStanProfile(c *gin.Context) {
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
	if namaStan, ok := updateData["nama_stan"].(string); ok {
		updates["nama_stan"] = namaStan
	}
	if namaPemilik, ok := updateData["nama_pemilik"].(string); ok {
		updates["nama_pemilik"] = namaPemilik
	}
	if telp, ok := updateData["telp"].(string); ok {
		updates["telp"] = telp
	}
	if foto, ok := updateData["foto"].(string); ok {
		// Handle base64 image
		if utils.IsBase64Image(foto) {
			imagePath, err := utils.SaveBase64Image(foto)
			if err != nil {
				BadRequestResponse(c, "Failed to process image", err)
				return
			}
			// Get existing stan to delete old image
			stan, _ := h.stanAdminService.GetStanByUserID(userID)
			if stan.Foto != "" {
				utils.DeleteImage(stan.Foto)
			}
			updates["foto"] = imagePath
		} else {
			updates["foto"] = foto
		}
	}

	if len(updates) == 0 {
		BadRequestResponse(c, "No valid fields to update", nil)
		return
	}

	if err := h.stanAdminService.UpdateStanProfile(userID, updates); err != nil {
		InternalErrorResponse(c, "Failed to update stan profile", err)
		return
	}

	// Get updated stan
	stan, _ := h.stanAdminService.GetStanByUserID(userID)
	SuccessResponse(c, "Stan profile updated successfully", stan)
}

// UpdatePaymentSettings updates payment settings for the stan
func (h *StanAdminHandler) UpdatePaymentSettings(c *gin.Context) {
	userID, exists := GetUserIDFromContext(c)
	if !exists {
		BadRequestResponse(c, "User not authenticated", nil)
		return
	}

	var req struct {
		AcceptCash bool   `json:"accept_cash"`
		AcceptQris bool   `json:"accept_qris"`
		QrisImage  string `json:"qris_image"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, "Invalid request body", err)
		return
	}

	// Handle base64 image if provided
	qrisImage := req.QrisImage
	if qrisImage != "" && utils.IsBase64Image(qrisImage) {
		imagePath, err := utils.SaveBase64Image(qrisImage)
		if err != nil {
			BadRequestResponse(c, "Failed to process QRIS image", err)
			return
		}
		// Get existing stan to delete old image
		stan, _ := h.stanAdminService.GetStanByUserID(userID)
		if stan.QrisImage != "" {
			utils.DeleteImage(stan.QrisImage)
		}
		qrisImage = imagePath
	}

	if err := h.stanAdminService.UpdatePaymentSettings(userID, req.AcceptCash, req.AcceptQris, qrisImage); err != nil {
		InternalErrorResponse(c, "Failed to update payment settings", err)
		return
	}

	SuccessResponse(c, "Payment settings updated successfully", nil)
}

// CreateMenu creates a new menu item for the stan
func (h *StanAdminHandler) CreateMenu(c *gin.Context) {
	userID, exists := GetUserIDFromContext(c)
	if !exists {
		BadRequestResponse(c, "User not authenticated", nil)
		return
	}

	var menu models.Menu
	if err := c.ShouldBindJSON(&menu); err != nil {
		BadRequestResponse(c, "Invalid request body", err)
		return
	}

	// Handle base64 image if provided
	if menu.Foto != "" && utils.IsBase64Image(menu.Foto) {
		imagePath, err := utils.SaveBase64Image(menu.Foto)
		if err != nil {
			BadRequestResponse(c, "Failed to process image", err)
			return
		}
		menu.Foto = imagePath
	}

	if err := h.stanAdminService.CreateMenu(userID, &menu); err != nil {
		InternalErrorResponse(c, "Failed to create menu", err)
		return
	}

	CreatedResponse(c, "Menu created successfully", menu)
}

// UpdateMenu updates a menu item
func (h *StanAdminHandler) UpdateMenu(c *gin.Context) {
	userID, exists := GetUserIDFromContext(c)
	if !exists {
		BadRequestResponse(c, "User not authenticated", nil)
		return
	}

	menuID, err := GetIDParam(c)
	if err != nil {
		BadRequestResponse(c, "Invalid menu ID", err)
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
	if namaMakanan, ok := updateData["nama_makanan"].(string); ok {
		updates["nama_makanan"] = namaMakanan
	}
	if harga, ok := updateData["harga"].(float64); ok {
		updates["harga"] = harga
	}
	if jenis, ok := updateData["jenis"].(string); ok {
		updates["jenis"] = jenis
	}
	if foto, ok := updateData["foto"].(string); ok {
		// Handle base64 image
		if utils.IsBase64Image(foto) {
			imagePath, err := utils.SaveBase64Image(foto)
			if err != nil {
				BadRequestResponse(c, "Failed to process image", err)
				return
			}
			// Delete old image if exists
			menu, _ := h.menuService.FindByID(menuID)
			if menu.Foto != "" {
				utils.DeleteImage(menu.Foto)
			}
			updates["foto"] = imagePath
		} else {
			updates["foto"] = foto
		}
	}
	if deskripsi, ok := updateData["deskripsi"].(string); ok {
		updates["deskripsi"] = deskripsi
	}
	if stock, ok := updateData["stock"].(float64); ok {
		updates["stock"] = int(stock)
	}

	if len(updates) == 0 {
		BadRequestResponse(c, "No valid fields to update", nil)
		return
	}

	if err := h.stanAdminService.UpdateMenu(userID, menuID, updates); err != nil {
		if err.Error() == "record not found" {
			NotFoundResponse(c, "Menu not found or you don't have permission")
		} else {
			InternalErrorResponse(c, "Failed to update menu", err)
		}
		return
	}

	// Get updated menu
	updatedMenu, _ := h.menuService.FindByID(menuID, "Stan")
	SuccessResponse(c, "Menu updated successfully", updatedMenu)
}

// DeleteMenu deletes a menu item
func (h *StanAdminHandler) DeleteMenu(c *gin.Context) {
	userID, exists := GetUserIDFromContext(c)
	if !exists {
		BadRequestResponse(c, "User not authenticated", nil)
		return
	}

	menuID, err := GetIDParam(c)
	if err != nil {
		BadRequestResponse(c, "Invalid menu ID", err)
		return
	}

	// Get existing menu to cleanup image
	menu, err := h.menuService.FindByID(menuID)
	if err != nil {
		NotFoundResponse(c, "Menu not found")
		return
	}

	if err := h.stanAdminService.DeleteMenu(userID, menuID); err != nil {
		if err.Error() == "record not found" {
			NotFoundResponse(c, "Menu not found or you don't have permission")
		} else {
			InternalErrorResponse(c, "Failed to delete menu", err)
		}
		return
	}

	// Delete image file if exists
	if menu.Foto != "" {
		utils.DeleteImage(menu.Foto)
	}

	SuccessResponse(c, "Menu deleted successfully", nil)
}

// GetMenus retrieves all menu items for the stan
func (h *StanAdminHandler) GetMenus(c *gin.Context) {
	userID, exists := GetUserIDFromContext(c)
	if !exists {
		BadRequestResponse(c, "User not authenticated", nil)
		return
	}

	menus, err := h.stanAdminService.GetMenusByStan(userID)
	if err != nil {
		InternalErrorResponse(c, "Failed to get menus", err)
		return
	}

	SuccessResponse(c, "Menus retrieved successfully", menus)
}

// UpdateStock updates the stock of a menu item
func (h *StanAdminHandler) UpdateStock(c *gin.Context) {
	menuID, err := GetIDParam(c)
	if err != nil {
		BadRequestResponse(c, "Invalid menu ID", err)
		return
	}

	var req struct {
		Stock int `json:"stock" binding:"required,min=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, "Invalid request body", err)
		return
	}

	if err := h.menuService.UpdateStock(menuID, req.Stock); err != nil {
		InternalErrorResponse(c, "Failed to update stock", err)
		return
	}

	menu, _ := h.menuService.FindByID(menuID, "Stan")
	SuccessResponse(c, "Stock updated successfully", menu)
}

// AdjustStock adjusts stock by a delta value
func (h *StanAdminHandler) AdjustStock(c *gin.Context) {
	menuID, err := GetIDParam(c)
	if err != nil {
		BadRequestResponse(c, "Invalid menu ID", err)
		return
	}

	var req struct {
		Delta int `json:"delta" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, "Invalid request body", err)
		return
	}

	if err := h.menuService.AdjustStock(menuID, req.Delta); err != nil {
		InternalErrorResponse(c, "Failed to adjust stock", err)
		return
	}

	menu, _ := h.menuService.FindByID(menuID, "Stan")
	SuccessResponse(c, "Stock adjusted successfully", menu)
}

// CreateStanDiscount creates a new stan-level discount
func (h *StanAdminHandler) CreateStanDiscount(c *gin.Context) {
	userID, exists := GetUserIDFromContext(c)
	if !exists {
		BadRequestResponse(c, "User not authenticated", nil)
		return
	}

	var req struct {
		NamaDiskon       string  `json:"nama_diskon" binding:"required"`
		PersentaseDiskon float64 `json:"persentase_diskon" binding:"required,min=0,max=100"`
		TanggalAwal      string  `json:"tanggal_awal" binding:"required"`
		TanggalAkhir     string  `json:"tanggal_akhir" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, "Invalid request body", err)
		return
	}

	// Parse dates
	tanggalAwal, err := time.Parse(time.RFC3339, req.TanggalAwal)
	if err != nil {
		BadRequestResponse(c, "Invalid tanggal_awal format", err)
		return
	}
	tanggalAkhir, err := time.Parse(time.RFC3339, req.TanggalAkhir)
	if err != nil {
		BadRequestResponse(c, "Invalid tanggal_akhir format", err)
		return
	}

	diskon := models.Diskon{
		NamaDiskon:       req.NamaDiskon,
		PersentaseDiskon: req.PersentaseDiskon,
		TanggalAwal:      tanggalAwal,
		TanggalAkhir:     tanggalAkhir,
	}

	if err := h.stanAdminService.CreateStanDiscount(userID, &diskon); err != nil {
		InternalErrorResponse(c, "Failed to create discount", err)
		return
	}

	CreatedResponse(c, "Discount created successfully", diskon)
}

// CreateMenuDiscount creates a new menu-level discount
func (h *StanAdminHandler) CreateMenuDiscount(c *gin.Context) {
	userID, exists := GetUserIDFromContext(c)
	if !exists {
		BadRequestResponse(c, "User not authenticated", nil)
		return
	}

	var req struct {
		NamaDiskon       string  `json:"nama_diskon" binding:"required"`
		PersentaseDiskon float64 `json:"persentase_diskon" binding:"required,min=0,max=100"`
		TanggalAwal      string  `json:"tanggal_awal" binding:"required"`
		TanggalAkhir     string  `json:"tanggal_akhir" binding:"required"`
		MenuIDs          []uint  `json:"menu_ids" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, "Invalid request body", err)
		return
	}

	// Parse dates
	tanggalAwal, err := time.Parse(time.RFC3339, req.TanggalAwal)
	if err != nil {
		BadRequestResponse(c, "Invalid tanggal_awal format", err)
		return
	}
	tanggalAkhir, err := time.Parse(time.RFC3339, req.TanggalAkhir)
	if err != nil {
		BadRequestResponse(c, "Invalid tanggal_akhir format", err)
		return
	}

	diskon := models.Diskon{
		NamaDiskon:       req.NamaDiskon,
		PersentaseDiskon: req.PersentaseDiskon,
		TanggalAwal:      tanggalAwal,
		TanggalAkhir:     tanggalAkhir,
	}

	if err := h.stanAdminService.CreateMenuDiscount(userID, &diskon, req.MenuIDs); err != nil {
		InternalErrorResponse(c, "Failed to create discount", err)
		return
	}

	CreatedResponse(c, "Discount created successfully", diskon)
}

// UpdateDiscount updates a discount
func (h *StanAdminHandler) UpdateDiscount(c *gin.Context) {
	userID, exists := GetUserIDFromContext(c)
	if !exists {
		BadRequestResponse(c, "User not authenticated", nil)
		return
	}

	diskonID, err := GetIDParam(c)
	if err != nil {
		BadRequestResponse(c, "Invalid discount ID", err)
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		BadRequestResponse(c, "Invalid request body", err)
		return
	}

	// Build updates map
	updates := make(map[string]interface{})
	if namaDiskon, ok := updateData["nama_diskon"].(string); ok {
		updates["nama_diskon"] = namaDiskon
	}
	if persentase, ok := updateData["persentase_diskon"].(float64); ok {
		updates["persentase_diskon"] = persentase
	}
	if tanggalAwal, ok := updateData["tanggal_awal"].(string); ok {
		updates["tanggal_awal"] = tanggalAwal
	}
	if tanggalAkhir, ok := updateData["tanggal_akhir"].(string); ok {
		updates["tanggal_akhir"] = tanggalAkhir
	}

	if len(updates) == 0 {
		BadRequestResponse(c, "No valid fields to update", nil)
		return
	}

	if err := h.stanAdminService.UpdateDiscount(userID, diskonID, updates); err != nil {
		if err.Error() == "record not found" {
			NotFoundResponse(c, "Discount not found or you don't have permission")
		} else {
			InternalErrorResponse(c, "Failed to update discount", err)
		}
		return
	}

	SuccessResponse(c, "Discount updated successfully", nil)
}

// DeleteDiscount deletes a discount
func (h *StanAdminHandler) DeleteDiscount(c *gin.Context) {
	userID, exists := GetUserIDFromContext(c)
	if !exists {
		BadRequestResponse(c, "User not authenticated", nil)
		return
	}

	diskonID, err := GetIDParam(c)
	if err != nil {
		BadRequestResponse(c, "Invalid discount ID", err)
		return
	}

	if err := h.stanAdminService.DeleteDiscount(userID, diskonID); err != nil {
		if err.Error() == "record not found" {
			NotFoundResponse(c, "Discount not found or you don't have permission")
		} else {
			InternalErrorResponse(c, "Failed to delete discount", err)
		}
		return
	}

	SuccessResponse(c, "Discount deleted successfully", nil)
}

// GetDiscounts retrieves all discounts for the stan
func (h *StanAdminHandler) GetDiscounts(c *gin.Context) {
	userID, exists := GetUserIDFromContext(c)
	if !exists {
		BadRequestResponse(c, "User not authenticated", nil)
		return
	}

	diskon, err := h.stanAdminService.GetDiscountsByStan(userID)
	if err != nil {
		InternalErrorResponse(c, "Failed to get discounts", err)
		return
	}

	SuccessResponse(c, "Discounts retrieved successfully", diskon)
}

// GetActiveDiscounts retrieves active discounts for the stan
func (h *StanAdminHandler) GetActiveDiscounts(c *gin.Context) {
	userID, exists := GetUserIDFromContext(c)
	if !exists {
		BadRequestResponse(c, "User not authenticated", nil)
		return
	}

	diskon, err := h.stanAdminService.GetActiveDiscountsByStan(userID)
	if err != nil {
		InternalErrorResponse(c, "Failed to get active discounts", err)
		return
	}

	SuccessResponse(c, "Active discounts retrieved successfully", diskon)
}

// GetTransactions retrieves all transactions for the stan
func (h *StanAdminHandler) GetTransactions(c *gin.Context) {
	userID, exists := GetUserIDFromContext(c)
	if !exists {
		BadRequestResponse(c, "User not authenticated", nil)
		return
	}

	transaksi, err := h.stanAdminService.GetTransactionsByStan(userID)
	if err != nil {
		InternalErrorResponse(c, "Failed to get transactions", err)
		return
	}

	SuccessResponse(c, "Transactions retrieved successfully", transaksi)
}

// GetTransactionsByDateRange retrieves transactions for the stan within a date range
func (h *StanAdminHandler) GetTransactionsByDateRange(c *gin.Context) {
	userID, exists := GetUserIDFromContext(c)
	if !exists {
		BadRequestResponse(c, "User not authenticated", nil)
		return
	}

	startDate, endDate := parseDateRange(c)

	transaksi, err := h.stanAdminService.GetTransactionsByStanAndDateRange(userID, startDate, endDate)
	if err != nil {
		InternalErrorResponse(c, "Failed to get transactions", err)
		return
	}

	SuccessResponse(c, "Transactions retrieved successfully", transaksi)
}

// UpdateTransactionStatus updates the status of a transaction
func (h *StanAdminHandler) UpdateTransactionStatus(c *gin.Context) {
	userID, exists := GetUserIDFromContext(c)
	if !exists {
		BadRequestResponse(c, "User not authenticated", nil)
		return
	}

	transaksiID, err := GetIDParam(c)
	if err != nil {
		BadRequestResponse(c, "Invalid transaction ID", err)
		return
	}

	var req struct {
		Status models.StatusTransaksi `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, "Invalid request body", err)
		return
	}

	if err := h.stanAdminService.UpdateTransactionStatus(userID, transaksiID, req.Status); err != nil {
		if err.Error() == "record not found" {
			NotFoundResponse(c, "Transaction not found or you don't have permission")
		} else {
			InternalErrorResponse(c, "Failed to update transaction status", err)
		}
		return
	}

	SuccessResponse(c, "Transaction status updated successfully", nil)
}

// GetRevenue retrieves revenue for the stan
func (h *StanAdminHandler) GetRevenue(c *gin.Context) {
	userID, exists := GetUserIDFromContext(c)
	if !exists {
		BadRequestResponse(c, "User not authenticated", nil)
		return
	}

	startDate, endDate := parseDateRange(c)

	totalRevenue, totalOrders, err := h.stanAdminService.GetStanRevenue(userID, startDate, endDate)
	if err != nil {
		InternalErrorResponse(c, "Failed to get revenue", err)
		return
	}

	response := gin.H{
		"total_revenue": totalRevenue,
		"total_orders":  totalOrders,
		"start_date":   startDate,
		"end_date":     endDate,
	}

	SuccessResponse(c, "Revenue retrieved successfully", response)
}


