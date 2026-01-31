package handlers

import (
	"swipeup-be/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService        *services.AuthService
	activityLogService *services.ActivityLogService
}

func NewAuthHandler(authService *services.AuthService, activityLogService *services.ActivityLogService) *AuthHandler {
	return &AuthHandler{
		authService:        authService,
		activityLogService: activityLogService,
	}
}

// Register creates a new user account
func (h *AuthHandler) Register(c *gin.Context) {
	var req services.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, "Invalid request body", err)
		return
	}

	user, err := h.authService.Register(req)
	if err != nil {
		BadRequestResponse(c, "Registration failed", err)
		return
	}

	SuccessResponse(c, "User registered successfully", user)
}

// Login authenticates user and returns JWT token
func (h *AuthHandler) Login(c *gin.Context) {
	var req services.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, "Invalid request body", err)
		return
	}

	authResponse, err := h.authService.Login(req)
	if err != nil {
		BadRequestResponse(c, "Login failed", err)
		return
	}

	// Log successful login activity
	ip, userAgent := GetClientInfo(c)
	h.activityLogService.LogActivity(authResponse.User.ID, "login", "User logged in successfully", ip, userAgent)

	SuccessResponse(c, "Login successful", authResponse)
}

// GetProfile returns current user profile (requires auth)
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := GetUserIDFromContext(c)
	if !exists {
		BadRequestResponse(c, "User not authenticated", nil)
		return
	}

	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		InternalErrorResponse(c, "Failed to get user profile", err)
		return
	}

	SuccessResponse(c, "Profile retrieved successfully", user)
}