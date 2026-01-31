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