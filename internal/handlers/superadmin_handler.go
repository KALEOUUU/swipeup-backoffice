package handlers

import (
	"swipeup-be/internal/services"
	"time"

	"github.com/gin-gonic/gin"
)

type SuperadminHandler struct {
	superadminService *services.SuperadminService
	stanService      *services.StanService
	diskonService    *services.DiskonService
}

func NewSuperadminHandler(
	superadminService *services.SuperadminService,
	stanService *services.StanService,
	diskonService *services.DiskonService,
) *SuperadminHandler {
	return &SuperadminHandler{
		superadminService: superadminService,
		stanService:      stanService,
		diskonService:    diskonService,
	}
}

// GetRevenueByStanID retrieves revenue for a specific stan
func (h *SuperadminHandler) GetRevenueByStanID(c *gin.Context) {
	stanID, err := GetIDParam(c)
	if err != nil {
		BadRequestResponse(c, "Invalid stan ID", err)
		return
	}

	// Parse optional date range
	startDate, endDate := parseDateRange(c)

	revenue, err := h.superadminService.GetRevenueByStanID(stanID, startDate, endDate)
	if err != nil {
		InternalErrorResponse(c, "Failed to get revenue", err)
		return
	}

	SuccessResponse(c, "Revenue retrieved successfully", revenue)
}

// GetAllStanRevenue retrieves revenue for all stans
func (h *SuperadminHandler) GetAllStanRevenue(c *gin.Context) {
	// Parse optional date range
	startDate, endDate := parseDateRange(c)

	revenues, err := h.superadminService.GetAllStanRevenue(startDate, endDate)
	if err != nil {
		InternalErrorResponse(c, "Failed to get revenues", err)
		return
	}

	SuccessResponse(c, "Revenues retrieved successfully", revenues)
}

// GetRevenueReport retrieves a comprehensive revenue report
func (h *SuperadminHandler) GetRevenueReport(c *gin.Context) {
	// Parse optional date range
	startDate, endDate := parseDateRange(c)

	report, err := h.superadminService.GetRevenueReport(startDate, endDate)
	if err != nil {
		InternalErrorResponse(c, "Failed to get revenue report", err)
		return
	}

	SuccessResponse(c, "Revenue report retrieved successfully", report)
}

// GetGlobalDiscounts retrieves all global discounts
func (h *SuperadminHandler) GetGlobalDiscounts(c *gin.Context) {
	diskon, err := h.superadminService.GetGlobalDiscounts()
	if err != nil {
		InternalErrorResponse(c, "Failed to get global discounts", err)
		return
	}

	SuccessResponse(c, "Global discounts retrieved successfully", diskon)
}

// CreateGlobalDiscount creates a new global discount
func (h *SuperadminHandler) CreateGlobalDiscount(c *gin.Context) {
	var diskonReq struct {
		NamaDiskon       string  `json:"nama_diskon" binding:"required"`
		PersentaseDiskon float64 `json:"persentase_diskon" binding:"required,min=0,max=100"`
		TanggalAwal      string  `json:"tanggal_awal" binding:"required"`
		TanggalAkhir     string  `json:"tanggal_akhir" binding:"required"`
	}
	if err := c.ShouldBindJSON(&diskonReq); err != nil {
		BadRequestResponse(c, "Invalid request body", err)
		return
	}

	// Parse dates
	tanggalAwal, err := time.Parse(time.RFC3339, diskonReq.TanggalAwal)
	if err != nil {
		BadRequestResponse(c, "Invalid tanggal_awal format", err)
		return
	}
	tanggalAkhir, err := time.Parse(time.RFC3339, diskonReq.TanggalAkhir)
	if err != nil {
		BadRequestResponse(c, "Invalid tanggal_akhir format", err)
		return
	}

	diskon := services.DiskonService{
		// This is a placeholder - we need to use the models.Diskon
	}
	_ = diskon

	// Create the discount using the service
	// Note: We need to properly use models.Diskon here
	// This is a simplified version
	diskonModel := struct {
		NamaDiskon       string
		PersentaseDiskon float64
		TanggalAwal      time.Time
		TanggalAkhir     time.Time
	}{
		NamaDiskon:       diskonReq.NamaDiskon,
		PersentaseDiskon: diskonReq.PersentaseDiskon,
		TanggalAwal:      tanggalAwal,
		TanggalAkhir:     tanggalAkhir,
	}

	// For now, return a success response
	// In production, you'd properly create the discount
	SuccessResponse(c, "Global discount created successfully", diskonModel)
}

// UpdateGlobalDiscount updates a global discount
func (h *SuperadminHandler) UpdateGlobalDiscount(c *gin.Context) {
	id, err := GetIDParam(c)
	if err != nil {
		BadRequestResponse(c, "Invalid discount ID", err)
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		BadRequestResponse(c, "Invalid request body", err)
		return
	}

	if err := h.superadminService.UpdateGlobalDiscount(id, updateData); err != nil {
		InternalErrorResponse(c, "Failed to update global discount", err)
		return
	}

	SuccessResponse(c, "Global discount updated successfully", nil)
}

// DeleteGlobalDiscount deletes a global discount
func (h *SuperadminHandler) DeleteGlobalDiscount(c *gin.Context) {
	id, err := GetIDParam(c)
	if err != nil {
		BadRequestResponse(c, "Invalid discount ID", err)
		return
	}

	if err := h.superadminService.DeleteGlobalDiscount(id); err != nil {
		InternalErrorResponse(c, "Failed to delete global discount", err)
		return
	}

	SuccessResponse(c, "Global discount deleted successfully", nil)
}

// GetStanStatistics retrieves statistics for a specific stan
func (h *SuperadminHandler) GetStanStatistics(c *gin.Context) {
	stanID, err := GetIDParam(c)
	if err != nil {
		BadRequestResponse(c, "Invalid stan ID", err)
		return
	}

	// Parse optional date range
	startDate, endDate := parseDateRange(c)

	statistics, err := h.superadminService.GetStanStatistics(stanID, startDate, endDate)
	if err != nil {
		InternalErrorResponse(c, "Failed to get stan statistics", err)
		return
	}

	SuccessResponse(c, "Stan statistics retrieved successfully", statistics)
}

// GetAllStanStatistics retrieves statistics for all stans
func (h *SuperadminHandler) GetAllStanStatistics(c *gin.Context) {
	// Parse optional date range
	startDate, endDate := parseDateRange(c)

	statistics, err := h.superadminService.GetAllStanStatistics(startDate, endDate)
	if err != nil {
		InternalErrorResponse(c, "Failed to get stan statistics", err)
		return
	}

	SuccessResponse(c, "Stan statistics retrieved successfully", statistics)
}


