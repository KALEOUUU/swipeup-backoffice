package handlers

import (
	"swipeup-be/internal/models"
	"swipeup-be/internal/services"
	"swipeup-be/pkg/utils"

	"github.com/gin-gonic/gin"
)

type MenuHandler struct {
	service     *services.MenuService
	authService *services.AuthService
	stanService *services.StanService
}

func NewMenuHandler(service *services.MenuService) *MenuHandler {
	return &MenuHandler{service: service}
}

func NewMenuHandlerWithDeps(service *services.MenuService, authService *services.AuthService, stanService *services.StanService) *MenuHandler {
	return &MenuHandler{
		service:     service,
		authService: authService,
		stanService: stanService,
	}
}

func (h *MenuHandler) Create(c *gin.Context) {
	// Get user ID from context (authenticated user)
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

	// Get user information
	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		InternalErrorResponse(c, "Failed to get user information", err)
		return
	}

	// Check if user is admin_stan
	if user.Role != "admin_stan" {
		BadRequestResponse(c, "Only stan admins can create menu items", nil)
		return
	}

	// Get stan for this user
	stan, err := h.stanService.GetByUserID(userID)
	if err != nil {
		InternalErrorResponse(c, "Failed to get stan information", err)
		return
	}

	menu.IDStan = stan.ID

	// Handle base64 image if provided
	if menu.Foto != "" && utils.IsBase64Image(menu.Foto) {
		imagePath, err := utils.SaveBase64Image(menu.Foto)
		if err != nil {
			BadRequestResponse(c, "Failed to process image", err)
			return
		}
		menu.Foto = imagePath
	}

	if err := h.service.Create(&menu); err != nil {
		InternalErrorResponse(c, "Failed to create menu", err)
		return
	}

	CreatedResponse(c, "Menu created successfully", menu)
}

func (h *MenuHandler) GetAll(c *gin.Context) {
	page, limit, offset := ParsePaginationParams(c)
	menus, total, err := h.service.FindAllPaginated(limit, offset, "Stan")
	if err != nil {
		InternalErrorResponse(c, "Failed to get menus", err)
		return
	}

	PaginatedSuccessResponse(c, "Menus retrieved successfully", menus, page, limit, int(total))
}

func (h *MenuHandler) GetByID(c *gin.Context) {
	id, err := GetIDParam(c)
	if err != nil {
		BadRequestResponse(c, "Invalid ID", err)
		return
	}

	menu, err := h.service.GetWithDiskon(id)
	if err != nil {
		NotFoundResponse(c, "Menu not found")
		return
	}

	SuccessResponse(c, "Menu retrieved successfully", menu)
}

func (h *MenuHandler) Update(c *gin.Context) {
	id, err := GetIDParam(c)
	if err != nil {
		BadRequestResponse(c, "Invalid ID", err)
		return
	}

	// Check if menu exists
	existingMenu, err := h.service.FindByID(id)
	if err != nil {
		NotFoundResponse(c, "Menu not found")
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
			if existingMenu.Foto != "" {
				utils.DeleteImage(existingMenu.Foto)
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

	if err := h.service.UpdateFields(id, updates); err != nil {
		InternalErrorResponse(c, "Failed to update menu", err)
		return
	}

	// Get updated menu to return
	updatedMenu, _ := h.service.FindByID(id, "Stan")
	SuccessResponse(c, "Menu updated successfully", updatedMenu)
}

func (h *MenuHandler) Delete(c *gin.Context) {
	id, err := GetIDParam(c)
	if err != nil {
		BadRequestResponse(c, "Invalid ID", err)
		return
	}

	// Get existing menu to cleanup image
	menu, err := h.service.FindByID(id)
	if err != nil {
		NotFoundResponse(c, "Menu not found")
		return
	}

	if err := h.service.Delete(id); err != nil {
		InternalErrorResponse(c, "Failed to delete menu", err)
		return
	}

	// Delete image file if exists
	if menu.Foto != "" {
		utils.DeleteImage(menu.Foto)
	}

	SuccessResponse(c, "Menu deleted successfully", nil)
}

func (h *MenuHandler) GetByStanID(c *gin.Context) {
	stanID, err := GetQueryParamUint(c, "stan_id")
	if err != nil || stanID == 0 {
		BadRequestResponse(c, "Invalid stan_id parameter", err)
		return
	}

	menus, err := h.service.GetByStanID(stanID)
	if err != nil {
		InternalErrorResponse(c, "Failed to get menus", err)
		return
	}

	SuccessResponse(c, "Menus retrieved successfully", menus)
}

func (h *MenuHandler) SearchByName(c *gin.Context) {
	name := GetQueryParam(c, "name")
	if name == "" {
		BadRequestResponse(c, "Name parameter is required", nil)
		return
	}

	menus, err := h.service.SearchByName(name)
	if err != nil {
		InternalErrorResponse(c, "Failed to search menus", err)
		return
	}

	SuccessResponse(c, "Menus retrieved successfully", menus)
}

// UpdateStock updates the stock of a menu item (inventory management)
func (h *MenuHandler) UpdateStock(c *gin.Context) {
	id, err := GetIDParam(c)
	if err != nil {
		BadRequestResponse(c, "Invalid ID", err)
		return
	}

	var req struct {
		Stock int `json:"stock" binding:"required,min=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, "Invalid request body", err)
		return
	}

	if err := h.service.UpdateStock(id, req.Stock); err != nil {
		InternalErrorResponse(c, "Failed to update stock", err)
		return
	}

	menu, _ := h.service.FindByID(id, "Stan")
	SuccessResponse(c, "Stock updated successfully", menu)
}

// AdjustStock adjusts stock by a delta value (positive to add, negative to reduce)
func (h *MenuHandler) AdjustStock(c *gin.Context) {
	id, err := GetIDParam(c)
	if err != nil {
		BadRequestResponse(c, "Invalid ID", err)
		return
	}

	var req struct {
		Delta int `json:"delta" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, "Invalid request body", err)
		return
	}

	if err := h.service.AdjustStock(id, req.Delta); err != nil {
		InternalErrorResponse(c, "Failed to adjust stock", err)
		return
	}

	menu, _ := h.service.FindByID(id, "Stan")
	SuccessResponse(c, "Stock adjusted successfully", menu)
}

// GetAvailableByStanID gets only available (in-stock) menu items by stan
func (h *MenuHandler) GetAvailableByStanID(c *gin.Context) {
	stanID, err := GetQueryParamUint(c, "stan_id")
	if err != nil || stanID == 0 {
		BadRequestResponse(c, "Invalid stan_id parameter", err)
		return
	}

	menus, err := h.service.GetAvailableMenuByStanID(stanID)
	if err != nil {
		InternalErrorResponse(c, "Failed to get available menus", err)
		return
	}

	SuccessResponse(c, "Available menus retrieved successfully", menus)
}
