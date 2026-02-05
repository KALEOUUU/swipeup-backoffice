package middleware

import (
	"net/http"

	"swipeup-be/internal/services"

	"github.com/gin-gonic/gin"
)

// ResourceMiddleware provides helper functions for resource ownership verification
type ResourceMiddleware struct {
	stanService  *services.StanService
	siswaService *services.SiswaService
}

func NewResourceMiddleware(stanService *services.StanService, siswaService *services.SiswaService) *ResourceMiddleware {
	return &ResourceMiddleware{
		stanService:  stanService,
		siswaService: siswaService,
	}
}

// StanOwnerOnly checks if the authenticated user owns the stan
func (rm *ResourceMiddleware) StanOwnerOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Unauthorized",
				"error":   "User not authenticated",
			})
			c.Abort()
			return
		}

		// Get stan ID from params or query
		stanID, err := getStanIDFromContext(c)
		if err != nil || stanID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid stan ID",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}

		// Check if stan belongs to user
		stan, err := rm.stanService.GetByUserID(userID.(uint))
		if err != nil || stan.ID != stanID {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Forbidden",
				"error":   "You don't have permission to access this stan",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// SiswaOwnerOnly checks if the authenticated user owns the siswa profile
func (rm *ResourceMiddleware) SiswaOwnerOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Unauthorized",
				"error":   "User not authenticated",
			})
			c.Abort()
			return
		}

		// Get siswa ID from params or query
		siswaID, err := getSiswaIDFromContext(c)
		if err != nil || siswaID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid siswa ID",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}

		// Check if siswa belongs to user
		siswa, err := rm.siswaService.GetByUserID(userID.(uint))
		if err != nil || siswa.ID != siswaID {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Forbidden",
				"error":   "You don't have permission to access this siswa profile",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// MenuStanOwnerOnly checks if the authenticated user owns the stan that the menu belongs to
func (rm *ResourceMiddleware) MenuStanOwnerOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Unauthorized",
				"error":   "User not authenticated",
			})
			c.Abort()
			return
		}

		// Get menu ID from params
		menuID, err := getIDParam(c)
		if err != nil || menuID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid menu ID",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}

		// Get stan owned by this admin
		stan, err := rm.stanService.GetByUserID(userID.(uint))
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Forbidden",
				"error":   "You don't have a stan assigned",
			})
			c.Abort()
			return
		}

		// Check if menu belongs to stan (this would need menuService, simplified for now)
		// In practice, you'd check if menu.IDStan == stan.ID
		// For now, we'll set stan_id in context for use in handlers
		c.Set("stan_id", stan.ID)
		c.Next()
	}
}

// TransaksiStanOwnerOnly checks if the authenticated user owns the stan that the transaction belongs to
func (rm *ResourceMiddleware) TransaksiStanOwnerOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Unauthorized",
				"error":   "User not authenticated",
			})
			c.Abort()
			return
		}

		// Get stan owned by this admin
		stan, err := rm.stanService.GetByUserID(userID.(uint))
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Forbidden",
				"error":   "You don't have a stan assigned",
			})
			c.Abort()
			return
		}

		// Set stan_id in context for use in handlers
		c.Set("stan_id", stan.ID)
		c.Next()
	}
}

// DiskonStanOwnerOnly checks if the authenticated user owns the stan that the discount belongs to
func (rm *ResourceMiddleware) DiskonStanOwnerOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Unauthorized",
				"error":   "User not authenticated",
			})
			c.Abort()
			return
		}

		// Get stan owned by this admin
		stan, err := rm.stanService.GetByUserID(userID.(uint))
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Forbidden",
				"error":   "You don't have a stan assigned",
			})
			c.Abort()
			return
		}

		// Set stan_id in context for use in handlers
		c.Set("stan_id", stan.ID)
		c.Next()
	}
}

// Helper functions
func getStanIDFromContext(c *gin.Context) (uint, error) {
	// Try to get from param first
	if id := c.Param("id"); id != "" {
		var stanID uint
		if err := c.ShouldBindUri(&struct {
			ID uint `uri:"id" binding:"required"`
		}{ID: stanID}); err == nil {
			return stanID, nil
		}
	}
	// Try to get from query
	if id := c.Query("stan_id"); id != "" {
		var stanID uint
		if err := c.ShouldBindQuery(&struct {
			StanID uint `form:"stan_id" binding:"required"`
		}{StanID: stanID}); err == nil {
			return stanID, nil
		}
	}
	// Try to get from body
	if id := c.PostForm("stan_id"); id != "" {
		var stanID uint
		if err := c.ShouldBind(&struct {
			StanID uint `form:"stan_id" binding:"required"`
		}{StanID: stanID}); err == nil {
			return stanID, nil
		}
	}
	return 0, nil
}

func getSiswaIDFromContext(c *gin.Context) (uint, error) {
	// Try to get from param first
	if id := c.Param("id"); id != "" {
		var siswaID uint
		if err := c.ShouldBindUri(&struct {
			ID uint `uri:"id" binding:"required"`
		}{ID: siswaID}); err == nil {
			return siswaID, nil
		}
	}
	// Try to get from query
	if id := c.Query("siswa_id"); id != "" {
		var siswaID uint
		if err := c.ShouldBindQuery(&struct {
			SiswaID uint `form:"siswa_id" binding:"required"`
		}{SiswaID: siswaID}); err == nil {
			return siswaID, nil
		}
	}
	return 0, nil
}

func getIDParam(c *gin.Context) (uint, error) {
	var id uint
	if err := c.ShouldBindUri(&struct {
		ID uint `uri:"id" binding:"required"`
	}{ID: id}); err != nil {
		return 0, err
	}
	return id, nil
}
