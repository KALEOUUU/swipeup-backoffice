package handlers

import (
	"strings"

	"swipeup-be/internal/models"
	"swipeup-be/internal/services"

	"github.com/gin-gonic/gin"
)

type TransaksiHandler struct {
	service      *services.TransaksiService
	stanService  *services.StanService
	siswaService *services.SiswaService
	menuService  *services.MenuService
}

func NewTransaksiHandler(service *services.TransaksiService) *TransaksiHandler {
	return &TransaksiHandler{service: service}
}

func NewTransaksiHandlerWithDeps(
	service *services.TransaksiService,
	stanService *services.StanService,
	siswaService *services.SiswaService,
	menuService *services.MenuService,
) *TransaksiHandler {
	return &TransaksiHandler{
		service:      service,
		stanService:  stanService,
		siswaService: siswaService,
		menuService:  menuService,
	}
}

type CreateTransaksiRequest struct {
	IDStan  uint                          `json:"id_stan" binding:"required"`
	IDSiswa uint                          `json:"id_siswa" binding:"required"`
	Details []models.DetailTransaksi      `json:"details" binding:"required,min=1"`
}

func (h *TransaksiHandler) Create(c *gin.Context) {
	var req CreateTransaksiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, "Invalid request body", err)
		return
	}

	transaksi := &models.Transaksi{
		IDStan:  req.IDStan,
		IDSiswa: req.IDSiswa,
		Status:  models.StatusBelumDikonfirm,
	}

	if err := h.service.CreateWithDetails(transaksi, req.Details); err != nil {
		// Check for specific FK constraint errors
		errMsg := err.Error()
		if strings.Contains(errMsg, "fk_stans_transaksi") {
			BadRequestResponse(c, "Stan with given id_stan does not exist", nil)
			return
		}
		if strings.Contains(errMsg, "fk_siswas_transaksi") || strings.Contains(errMsg, "fk_transaksis_siswa") {
			BadRequestResponse(c, "Siswa with given id_siswa does not exist", nil)
			return
		}
		if strings.Contains(errMsg, "fk_menus_detail") || strings.Contains(errMsg, "fk_detail_transaksis_menu") {
			BadRequestResponse(c, "One or more menu items in details do not exist", nil)
			return
		}
		InternalErrorResponse(c, "Failed to create transaction", err)
		return
	}

	result, _ := h.service.GetWithFullDetails(transaksi.ID)
	CreatedResponse(c, "Transaction created successfully", result)
}

func (h *TransaksiHandler) GetAll(c *gin.Context) {
	transaksi, err := h.service.FindAll("Stan", "Siswa", "DetailTransaksi")
	if err != nil {
		InternalErrorResponse(c, "Failed to get transactions", err)
		return
	}

	SuccessResponse(c, "Transactions retrieved successfully", transaksi)
}

func (h *TransaksiHandler) GetByID(c *gin.Context) {
	id, err := GetIDParam(c)
	if err != nil {
		BadRequestResponse(c, "Invalid ID", err)
		return
	}

	transaksi, err := h.service.GetWithFullDetails(id)
	if err != nil {
		NotFoundResponse(c, "Transaction not found")
		return
	}

	SuccessResponse(c, "Transaction retrieved successfully", transaksi)
}

func (h *TransaksiHandler) GetBySiswaID(c *gin.Context) {
	siswaID, err := GetQueryParamUint(c, "siswa_id")
	if err != nil || siswaID == 0 {
		BadRequestResponse(c, "Invalid siswa_id parameter", err)
		return
	}

	transaksi, err := h.service.GetBySiswaID(siswaID)
	if err != nil {
		InternalErrorResponse(c, "Failed to get transactions", err)
		return
	}

	SuccessResponse(c, "Transactions retrieved successfully", transaksi)
}

func (h *TransaksiHandler) GetByStanID(c *gin.Context) {
	stanID, err := GetQueryParamUint(c, "stan_id")
	if err != nil || stanID == 0 {
		BadRequestResponse(c, "Invalid stan_id parameter", err)
		return
	}

	transaksi, err := h.service.GetByStanID(stanID)
	if err != nil {
		InternalErrorResponse(c, "Failed to get transactions", err)
		return
	}

	SuccessResponse(c, "Transactions retrieved successfully", transaksi)
}

type UpdateStatusRequest struct {
	Status models.StatusTransaksi `json:"status" binding:"required"`
}

func (h *TransaksiHandler) UpdateStatus(c *gin.Context) {
	id, err := GetIDParam(c)
	if err != nil {
		BadRequestResponse(c, "Invalid ID", err)
		return
	}

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, "Invalid request body", err)
		return
	}

	if err := h.service.UpdateStatus(id, req.Status); err != nil {
		InternalErrorResponse(c, "Failed to update transaction status", err)
		return
	}

	transaksi, _ := h.service.GetWithFullDetails(id)
	SuccessResponse(c, "Transaction status updated successfully", transaksi)
}
