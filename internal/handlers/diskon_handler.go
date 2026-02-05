package handlers

import (
	"fmt"
	"swipeup-be/internal/models"
	"swipeup-be/internal/services"
	"time"

	"github.com/gin-gonic/gin"
)

type DiskonHandler struct {
	service     *services.DiskonService
	stanService *services.StanService
	authService *services.AuthService
}

func NewDiskonHandler(service *services.DiskonService, stanService *services.StanService, authService *services.AuthService) *DiskonHandler {
	return &DiskonHandler{
		service:     service,
		stanService: stanService,
		authService: authService,
	}
}

// checkDiskonPermission checks if user has permission to modify the diskon
func (h *DiskonHandler) checkDiskonPermission(c *gin.Context, diskon *models.Diskon) bool {
	userRole, _ := c.Get("role")
	role := userRole.(string)

	if role == "superadmin" {
		return true // superadmin can modify all
	}

	if role == "admin_stan" {
		userID, _ := c.Get("user_id")
		uid := userID.(uint)

		// Get stan owned by this admin
		stan, err := h.stanService.GetByUserID(uid)
		if err != nil {
			return false
		}

		// Admin can only modify diskon for their own stan or global (but global should be superadmin only)
		// For stan and menu diskon, check if id_stan matches
		if diskon.TipeDiskon == models.DiskonStan || diskon.TipeDiskon == models.DiskonMenu {
			return diskon.IDStan != nil && *diskon.IDStan == stan.ID
		}
	}

	return false
}

type CreateDiskonRequest struct {
	NamaDiskon       string  `json:"nama_diskon" binding:"required"`
	PersentaseDiskon float64 `json:"persentase_diskon" binding:"required,min=0,max=100"`
	TanggalAwal      string  `json:"tanggal_awal" binding:"required"`
	TanggalAkhir     string  `json:"tanggal_akhir" binding:"required"`
	TipeDiskon       string  `json:"tipe_diskon" binding:"required,oneof=global stan menu"`
	IDStan           *uint   `json:"id_stan,omitempty"`
	IDMenu           []uint  `json:"id_menu,omitempty"`
}

func (h *DiskonHandler) Create(c *gin.Context) {
	var req CreateDiskonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, "Invalid request body", err)
		return
	}

	// Check permission for creation
	userRole, _ := c.Get("role")
	role := userRole.(string)

	if role == "admin_stan" {
		if req.TipeDiskon == "global" {
			ErrorResponse(c, 403, "Admin stan cannot create global diskon", nil)
			return
		}
		if req.TipeDiskon == "stan" || req.TipeDiskon == "menu" {
			if req.IDStan == nil {
				BadRequestResponse(c, "id_stan is required", nil)
				return
			}
			userID, _ := c.Get("user_id")
			uid := userID.(uint)
			stan, err := h.stanService.GetByUserID(uid)
			if err != nil || stan.ID != *req.IDStan {
				ErrorResponse(c, 403, "You can only create diskon for your own stan", err)
				return
			}
		}
	}

	// Validate logic
	if req.TipeDiskon == "stan" || req.TipeDiskon == "menu" {
		if req.IDStan == nil {
			BadRequestResponse(c, "id_stan is required for stan or menu discount", nil)
			return
		}
		// Check if stan exists
		_, err := h.stanService.FindByID(*req.IDStan)
		if err != nil {
			NotFoundResponse(c, "Stan not found")
			return
		}
	}
	if req.TipeDiskon == "menu" && len(req.IDMenu) == 0 {
		BadRequestResponse(c, "id_menu is required for menu discount", nil)
		return
	}
	if req.TipeDiskon == "menu" {
		// Check if all menus exist and belong to the stan
		menuService := services.NewMenuService(h.service.GetDB()) // Assuming we can create it here
		for _, menuID := range req.IDMenu {
			menu, err := menuService.FindByID(menuID)
			if err != nil {
				NotFoundResponse(c, fmt.Sprintf("Menu with ID %d not found", menuID))
				return
			}
			if menu.IDStan != *req.IDStan {
				ErrorResponse(c, 400, fmt.Sprintf("Menu with ID %d does not belong to the specified stan", menuID), nil)
				return
			}
		}
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

	// Create diskon
	diskon := models.Diskon{
		NamaDiskon:       req.NamaDiskon,
		PersentaseDiskon: req.PersentaseDiskon,
		TanggalAwal:      tanggalAwal,
		TanggalAkhir:     tanggalAkhir,
		TipeDiskon:       models.TipeDiskon(req.TipeDiskon),
		IDStan:           req.IDStan,
	}

	if err := h.service.Create(&diskon); err != nil {
		InternalErrorResponse(c, "Failed to create diskon", err)
		return
	}

	// Assign to menus if menu discount
	if req.TipeDiskon == "menu" {
		for _, menuID := range req.IDMenu {
			if err := h.service.AssignToMenu(diskon.ID, menuID); err != nil {
				// Rollback? For simplicity, just log error
				InternalErrorResponse(c, "Failed to assign diskon to menu", err)
				return
			}
		}
	}

	CreatedResponse(c, "Diskon created successfully", diskon)
}

func (h *DiskonHandler) GetAll(c *gin.Context) {
	page, limit, offset := ParsePaginationParams(c)
	diskon, total, err := h.service.FindAllPaginated(limit, offset)
	if err != nil {
		InternalErrorResponse(c, "Failed to get diskon", err)
		return
	}

	PaginatedSuccessResponse(c, "Diskon retrieved successfully", diskon, page, limit, int(total))
}

func (h *DiskonHandler) GetActive(c *gin.Context) {
	diskon, err := h.service.GetActiveDiskon()
	if err != nil {
		InternalErrorResponse(c, "Failed to get active diskon", err)
		return
	}

	SuccessResponse(c, "Active diskon retrieved successfully", diskon)
}

func (h *DiskonHandler) GetByID(c *gin.Context) {
	id, err := GetIDParam(c)
	if err != nil {
		BadRequestResponse(c, "Invalid ID", err)
		return
	}

	diskon, err := h.service.FindByID(id, "MenuDiskon", "MenuDiskon.Menu")
	if err != nil {
		NotFoundResponse(c, "Diskon not found")
		return
	}

	SuccessResponse(c, "Diskon retrieved successfully", diskon)
}

func (h *DiskonHandler) Update(c *gin.Context) {
	id, err := GetIDParam(c)
	if err != nil {
		BadRequestResponse(c, "Invalid ID", err)
		return
	}

	// Check if diskon exists
	diskon, err := h.service.FindByID(id)
	if err != nil {
		NotFoundResponse(c, "Diskon not found")
		return
	}

	// Check permission
	if !h.checkDiskonPermission(c, diskon) {
		ErrorResponse(c, 403, "You don't have permission to update this diskon", nil)
		return
	}

	// Bind request body to map
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
	if tipeDiskon, ok := updateData["tipe_diskon"].(string); ok {
		updates["tipe_diskon"] = tipeDiskon
	}
	if idStan, ok := updateData["id_stan"].(float64); ok {
		stanID := uint(idStan)
		updates["id_stan"] = &stanID
	}

	if len(updates) == 0 {
		BadRequestResponse(c, "No valid fields to update", nil)
		return
	}

	if err := h.service.UpdateFields(id, updates); err != nil {
		InternalErrorResponse(c, "Failed to update diskon", err)
		return
	}

	// Get updated diskon
	updatedDiskon, _ := h.service.FindByID(id, "Stan")
	SuccessResponse(c, "Diskon updated successfully", updatedDiskon)
}

func (h *DiskonHandler) Delete(c *gin.Context) {
	id, err := GetIDParam(c)
	if err != nil {
		BadRequestResponse(c, "Invalid ID", err)
		return
	}

	// Check if diskon exists
	diskon, err := h.service.FindByID(id)
	if err != nil {
		NotFoundResponse(c, "Diskon not found")
		return
	}

	// Check permission
	if !h.checkDiskonPermission(c, diskon) {
		ErrorResponse(c, 403, "You don't have permission to delete this diskon", nil)
		return
	}

	if err := h.service.Delete(id); err != nil {
		InternalErrorResponse(c, "Failed to delete diskon", err)
		return
	}

	SuccessResponse(c, "Diskon deleted successfully", nil)
}

func (h *DiskonHandler) GetByStan(c *gin.Context) {
	stanID, err := GetQueryParamUint(c, "stan_id")
	if err != nil || stanID == 0 {
		BadRequestResponse(c, "Invalid stan_id parameter", err)
		return
	}

	diskon, err := h.service.GetByStanID(stanID)
	if err != nil {
		InternalErrorResponse(c, "Failed to get diskon", err)
		return
	}

	SuccessResponse(c, "Diskon retrieved successfully", diskon)
}

func (h *DiskonHandler) GetGlobal(c *gin.Context) {
	diskon, err := h.service.GetGlobalDiskon()
	if err != nil {
		InternalErrorResponse(c, "Failed to get global diskon", err)
		return
	}

	SuccessResponse(c, "Global diskon retrieved successfully", diskon)
}

func (h *DiskonHandler) GetActiveByStanID(c *gin.Context) {
	stanID, err := GetQueryParamUint(c, "stan_id")
	if err != nil || stanID == 0 {
		BadRequestResponse(c, "Invalid stan_id parameter", err)
		return
	}

	diskon, err := h.service.GetActiveDiskonByStan(stanID)
	if err != nil {
		InternalErrorResponse(c, "Failed to get active diskon", err)
		return
	}

	SuccessResponse(c, "Active diskon retrieved successfully", diskon)
}

type AssignDiskonRequest struct {
	MenuID uint `json:"menu_id" binding:"required"`
}

func (h *DiskonHandler) AssignToMenu(c *gin.Context) {
	diskonID, err := GetIDParam(c)
	if err != nil {
		BadRequestResponse(c, "Invalid diskon ID", err)
		return
	}

	var req AssignDiskonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, "Invalid request body", err)
		return
	}

	if err := h.service.AssignToMenu(diskonID, req.MenuID); err != nil {
		InternalErrorResponse(c, "Failed to assign diskon to menu", err)
		return
	}

	SuccessResponse(c, "Diskon assigned to menu successfully", nil)
}

func (h *DiskonHandler) RemoveFromMenu(c *gin.Context) {
	diskonID, err := GetIDParam(c)
	if err != nil {
		BadRequestResponse(c, "Invalid diskon ID", err)
		return
	}

	menuID, err := GetQueryParamUint(c, "menu_id")
	if err != nil || menuID == 0 {
		BadRequestResponse(c, "Invalid menu_id parameter", err)
		return
	}

	if err := h.service.RemoveFromMenu(diskonID, menuID); err != nil {
		InternalErrorResponse(c, "Failed to remove diskon from menu", err)
		return
	}

	SuccessResponse(c, "Diskon removed from menu successfully", nil)
}
