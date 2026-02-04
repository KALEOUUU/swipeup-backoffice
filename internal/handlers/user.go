package handlers

import (
	"swipeup-be/internal/models"
	"swipeup-be/internal/services"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		BadRequestResponse(c, "Invalid request body", err)
		return
	}
	if err := h.userService.CreateUser(&user); err != nil {
		InternalErrorResponse(c, "Failed to create user", err)
		return
	}
	CreatedResponse(c, "User created successfully", user)
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		InternalErrorResponse(c, "Failed to get users", err)
		return
	}
	SuccessResponse(c, "Users retrieved successfully", users)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := h.userService.GetUserByID(id)
	if err != nil {
		NotFoundResponse(c, "User not found")
		return
	}
	SuccessResponse(c, "User retrieved successfully", user)
}

// UpdateUser updates a user's information (superadmin only)
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := GetIDParam(c)
	if err != nil {
		BadRequestResponse(c, "Invalid ID", err)
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		BadRequestResponse(c, "Invalid request body", err)
		return
	}

	// Build updates map with only allowed fields
	updates := make(map[string]interface{})
	if username, ok := updateData["username"].(string); ok {
		updates["username"] = username
	}
	if role, ok := updateData["role"].(string); ok {
		// Validate role
		if role == "superadmin" || role == "admin_stan" || role == "siswa" {
			updates["role"] = role
		}
	}

	if len(updates) == 0 {
		BadRequestResponse(c, "No valid fields to update", nil)
		return
	}

	if err := h.userService.UpdateUser(id, updates); err != nil {
		InternalErrorResponse(c, "Failed to update user", err)
		return
	}

	user, _ := h.userService.GetUserByID(c.Param("id"))
	SuccessResponse(c, "User updated successfully", user)
}

// DeleteUser soft deletes a user (superadmin only)
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := GetIDParam(c)
	if err != nil {
		BadRequestResponse(c, "Invalid ID", err)
		return
	}

	if err := h.userService.DeleteUser(id); err != nil {
		InternalErrorResponse(c, "Failed to delete user", err)
		return
	}

	SuccessResponse(c, "User deleted successfully", nil)
}