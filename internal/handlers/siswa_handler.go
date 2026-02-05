package handlers

import (
	"swipeup-be/internal/models"
	"swipeup-be/internal/services"
	"swipeup-be/pkg/utils"

	"github.com/gin-gonic/gin"
)

type SiswaHandler struct {
	service *services.SiswaService
}

func NewSiswaHandler(service *services.SiswaService) *SiswaHandler {
	return &SiswaHandler{service: service}
}

func (h *SiswaHandler) Create(c *gin.Context) {
	var siswa models.Siswa
	if err := c.ShouldBindJSON(&siswa); err != nil {
		BadRequestResponse(c, "Invalid request body", err)
		return
	}
// Handle base64 image if provided
	if siswa.Foto != "" && utils.IsBase64Image(siswa.Foto) {
		imagePath, err := utils.SaveBase64Image(siswa.Foto)
		if err != nil {
			BadRequestResponse(c, "Failed to process image", err)
			return
		}
		siswa.Foto = imagePath
	}

	
	if err := h.service.Create(&siswa); err != nil {
		InternalErrorResponse(c, "Failed to create siswa", err)
		return
	}

	CreatedResponse(c, "Siswa created successfully", siswa)
}

func (h *SiswaHandler) GetAll(c *gin.Context) {
	page, limit, offset := ParsePaginationParams(c)
	siswa, total, err := h.service.FindAllPaginated(limit, offset, "User")
	if err != nil {
		InternalErrorResponse(c, "Failed to get siswa", err)
		return
	}

	PaginatedSuccessResponse(c, "Siswa retrieved successfully", siswa, page, limit, int(total))
}

func (h *SiswaHandler) GetByID(c *gin.Context) {
	id, err := GetIDParam(c)
	if err != nil {
		BadRequestResponse(c, "Invalid ID", err)
		return
	}

	siswa, err := h.service.FindByID(id, "User")
	if err != nil {
		NotFoundResponse(c, "Siswa not found")
		return
	}

	SuccessResponse(c, "Siswa retrieved successfully", siswa)
}

func (h *SiswaHandler) Update(c *gin.Context) {
	id, err := GetIDParam(c)
	if err != nil {
		BadRequestResponse(c, "Invalid ID", err)
		return
	}

	// Check if siswa exists
	existingSiswa, err := h.service.FindByID(id)
	if err != nil {
		NotFoundResponse(c, "Siswa not found")
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
	if nama, ok := updateData["nama_siswa"].(string); ok {
		updates["nama_siswa"] = nama
	}
	if alamat, ok := updateData["alamat"].(string); ok {
		updates["alamat"] = alamat
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
			if existingSiswa.Foto != "" {
				utils.DeleteImage(existingSiswa.Foto)
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
		InternalErrorResponse(c, "Failed to update siswa", err)
		return
	}

	// Get updated siswa to return
	updatedSiswa, _ := h.service.FindByID(id, "User")
	SuccessResponse(c, "Siswa updated successfully", updatedSiswa)
}

func (h *SiswaHandler) Delete(c *gin.Context) {
	id, err := GetIDParam(c)
	if err != nil {
		BadRequestResponse(c, "Invalid ID", err)
		return
	}

	// Get siswa to delete image file
	siswa, err := h.service.FindByID(id)
	if err != nil {
		NotFoundResponse(c, "Siswa not found")
		return
	}

	// Delete from database
	if err := h.service.Delete(id); err != nil {
		InternalErrorResponse(c, "Failed to delete siswa", err)
		return
	}

	// Delete image file if exists
	if siswa.Foto != "" {
		utils.DeleteImage(siswa.Foto)
	}

	SuccessResponse(c, "Siswa deleted successfully", nil)
}

func (h *SiswaHandler) GetByUserID(c *gin.Context) {
	userID, err := GetQueryParamUint(c, "user_id")
	if err != nil || userID == 0 {
		BadRequestResponse(c, "Invalid user_id parameter", err)
		return
	}

	siswa, err := h.service.GetByUserID(userID)
	if err != nil {
		NotFoundResponse(c, "Siswa not found")
		return
	}

	SuccessResponse(c, "Siswa retrieved successfully", siswa)
}
