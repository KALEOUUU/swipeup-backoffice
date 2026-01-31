package handlers

import (
	"strconv"
	"time"

	"swipeup-be/internal/services"

	"github.com/gin-gonic/gin"
)

type ActivityLogHandler struct {
	service *services.ActivityLogService
}

func NewActivityLogHandler(service *services.ActivityLogService) *ActivityLogHandler {
	return &ActivityLogHandler{service: service}
}

// GetUserActivities retrieves activities for a specific user
func (h *ActivityLogHandler) GetUserActivities(c *gin.Context) {
	userID, err := GetQueryParamUint(c, "user_id")
	if err != nil || userID == 0 {
		BadRequestResponse(c, "user_id parameter is required", err)
		return
	}

	limit, offset := ParsePaginationParams(c)

	activities, err := h.service.GetUserActivities(userID, limit, offset)
	if err != nil {
		InternalErrorResponse(c, "Failed to retrieve activities", err)
		return
	}

	SuccessResponse(c, "Activities retrieved successfully", gin.H{
		"activities": activities,
		"count":      len(activities),
	})
}

// GetAllActivities retrieves all activities with optional filtering
func (h *ActivityLogHandler) GetAllActivities(c *gin.Context) {
	actionFilter := c.Query("action")
	limit, offset := ParsePaginationParams(c)

	activities, err := h.service.GetAllActivities(actionFilter, limit, offset)
	if err != nil {
		InternalErrorResponse(c, "Failed to retrieve activities", err)
		return
	}

	SuccessResponse(c, "Activities retrieved successfully", gin.H{
		"activities": activities,
		"count":      len(activities),
	})
}

// GetActivitiesByDateRange retrieves activities within a date range
func (h *ActivityLogHandler) GetActivitiesByDateRange(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		BadRequestResponse(c, "start_date and end_date parameters are required", nil)
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		BadRequestResponse(c, "Invalid start_date format (use YYYY-MM-DD)", err)
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		BadRequestResponse(c, "Invalid end_date format (use YYYY-MM-DD)", err)
		return
	}

	// Set end date to end of day
	endDate = endDate.Add(24*time.Hour - time.Second)

	limit, offset := ParsePaginationParams(c)

	activities, err := h.service.GetActivitiesByDateRange(startDate, endDate, limit, offset)
	if err != nil {
		InternalErrorResponse(c, "Failed to retrieve activities", err)
		return
	}

	SuccessResponse(c, "Activities retrieved successfully", gin.H{
		"activities": activities,
		"count":      len(activities),
	})
}

// GetActivityStats returns activity statistics
func (h *ActivityLogHandler) GetActivityStats(c *gin.Context) {
	stats, err := h.service.GetActivityStats()
	if err != nil {
		InternalErrorResponse(c, "Failed to retrieve activity stats", err)
		return
	}

	SuccessResponse(c, "Activity stats retrieved successfully", stats)
}

// CleanOldLogs removes old activity logs (admin only)
func (h *ActivityLogHandler) CleanOldLogs(c *gin.Context) {
	daysStr := c.DefaultQuery("days", "90") // Default 90 days

	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 {
		BadRequestResponse(c, "Invalid days parameter (must be positive integer)", err)
		return
	}

	err = h.service.CleanOldLogs(days)
	if err != nil {
		InternalErrorResponse(c, "Failed to clean old logs", err)
		return
	}

	SuccessResponse(c, "Old activity logs cleaned successfully", nil)
}