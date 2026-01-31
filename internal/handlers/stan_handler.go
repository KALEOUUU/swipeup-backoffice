package handlers

import (
	"swipeup-be/internal/models"
	"swipeup-be/internal/services"
	"swipeup-be/pkg/utils"

	"github.com/gin-gonic/gin"
)

type StanHandler struct {
	service *services.StanService
}

func NewStanHandler(service *services.StanService) *StanHandler {
	return &StanHandler{service: service}
}

func (h *StanHandler) Create(c *gin.Context) {
	var stan models.Stan
	if err := c.ShouldBindJSON(&stan); err != nil {
		BadRequestResponse(c, "Invalid request body", err)
		return
	}

	// Handle base64 image if provided
	if stan.Foto != "" && utils.IsBase64Image(stan.Foto) {
		imagePath, err := utils.SaveBase64Image(stan.Foto)
		if err != nil {
			BadRequestResponse(c, "Failed to process image", err)
			return
		}
		stan.Foto = imagePath
	}

	if err := h.service.Create(&stan); err != nil {
		InternalErrorResponse(c, "Failed to create stan", err)
		return
	}

	CreatedResponse(c, "Stan created successfully", stan)
}

func (h *StanHandler) GetAll(c *gin.Context) {
	stan, err := h.service.FindAll("User")
	if err != nil {
		InternalErrorResponse(c, "Failed to get stan", err)
		return
	}

	SuccessResponse(c, "Stan retrieved successfully", stan)
}

func (h *StanHandler) GetByID(c *gin.Context) {
	id, err := GetIDParam(c)
	if err != nil {
		BadRequestResponse(c, "Invalid ID", err)
		return
	}

	stan, err := h.service.GetWithMenu(id)
	if err != nil {
		NotFoundResponse(c, "Stan not found")
		return
	}

	SuccessResponse(c, "Stan retrieved successfully", stan)
}

func (h *StanHandler) Update(c *gin.Context) {
	id, err := GetIDParam(c)
	if err != nil {
		BadRequestResponse(c, "Invalid ID", err)
		return
	}

	// Check if stan exists
	existingStan, err := h.service.FindByID(id)
	if err != nil {
		NotFoundResponse(c, "Stan not found")
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
			// Delete old image if exists
			if existingStan.Foto != "" {
				utils.DeleteImage(existingStan.Foto)
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

	if err := h.service.UpdateFields(id, updates); err != nil {
		InternalErrorResponse(c, "Failed to update stan", err)
		return
	}

	// Get updated stan to return
	updatedStan, _ := h.service.FindByID(id, "User")
	SuccessResponse(c, "Stan updated successfully", updatedStan)
}

func (h *StanHandler) Delete(c *gin.Context) {
	id, err := GetIDParam(c)
	if err != nil {
		BadRequestResponse(c, "Invalid ID", err)
		return
	}

	// Get existing stan to cleanup image
	stan, err := h.service.FindByID(id)
	if err != nil {
		NotFoundResponse(c, "Stan not found")
		return
	}

	if err := h.service.Delete(id); err != nil {
		InternalErrorResponse(c, "Failed to delete stan", err)
		return
	}

	// Delete image file if exists
	if stan.Foto != "" {
		utils.DeleteImage(stan.Foto)
	}

	SuccessResponse(c, "Stan deleted successfully", nil)
}

func (h *StanHandler) GetByUserID(c *gin.Context) {
	userID, err := GetQueryParamUint(c, "user_id")
	if err != nil || userID == 0 {
		BadRequestResponse(c, "Invalid user_id parameter", err)
		return
	}

	stan, err := h.service.GetByUserID(userID)
	if err != nil {
		NotFoundResponse(c, "Stan not found")
		return
	}

	SuccessResponse(c, "Stan retrieved successfully", stan)
}
