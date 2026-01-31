package handlers

import (
	"swipeup-be/internal/models"
	"swipeup-be/internal/services"

	"github.com/gin-gonic/gin"
)

type DiskonHandler struct {
	service *services.DiskonService
}

func NewDiskonHandler(service *services.DiskonService) *DiskonHandler {
	return &DiskonHandler{service: service}
}

func (h *DiskonHandler) Create(c *gin.Context) {
	var diskon models.Diskon
	if err := c.ShouldBindJSON(&diskon); err != nil {
		BadRequestResponse(c, "Invalid request body", err)
		return
	}

	if err := h.service.Create(&diskon); err != nil {
		InternalErrorResponse(c, "Failed to create diskon", err)
		return
	}

	CreatedResponse(c, "Diskon created successfully", diskon)
}

func (h *DiskonHandler) GetAll(c *gin.Context) {
	diskon, err := h.service.FindAll()
	if err != nil {
		InternalErrorResponse(c, "Failed to get diskon", err)
		return
	}

	SuccessResponse(c, "Diskon retrieved successfully", diskon)
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
	_, err = h.service.FindByID(id)
	if err != nil {
		NotFoundResponse(c, "Diskon not found")
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
