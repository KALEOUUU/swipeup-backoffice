package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
)

// parseDateRange parses start_date and end_date from query params
func parseDateRange(c *gin.Context) (startDate, endDate time.Time) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr != "" {
		startDate, _ = time.Parse(time.RFC3339, startDateStr)
	}
	if endDateStr != "" {
		endDate, _ = time.Parse(time.RFC3339, endDateStr)
	}

	return startDate, endDate
}
