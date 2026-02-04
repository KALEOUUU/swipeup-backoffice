package services

import (
	"swipeup-be/internal/models"
	"time"

	"gorm.io/gorm"
)

type ActivityLogService struct {
	db *gorm.DB
}

func NewActivityLogService(db *gorm.DB) *ActivityLogService {
	return &ActivityLogService{db: db}
}

// LogActivity records a user activity
func (s *ActivityLogService) LogActivity(userID uint, action, description, ipAddress, userAgent string) error {
	activity := &models.ActivityLog{
		IDUser:      userID,
		Action:      action,
		Description: description,
		IPAddress:   ipAddress,
		UserAgent:   userAgent,
		CreatedAt:   time.Now(),
	}

	return s.db.Create(activity).Error
}

// GetActivityByID retrieves a single activity log by ID
func (s *ActivityLogService) GetActivityByID(id uint) (*models.ActivityLog, error) {
	var activity models.ActivityLog
	err := s.db.Preload("User").First(&activity, id).Error
	if err != nil {
		return nil, err
	}
	return &activity, nil
}

// GetUserActivities retrieves activities for a specific user
func (s *ActivityLogService) GetUserActivities(userID uint, limit, offset int) ([]models.ActivityLog, error) {
	var activities []models.ActivityLog
	query := s.db.Preload("User").Where("id_user = ?", userID).Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&activities).Error
	return activities, err
}

// GetAllActivities retrieves all activities with optional filtering
func (s *ActivityLogService) GetAllActivities(actionFilter string, limit, offset int) ([]models.ActivityLog, error) {
	var activities []models.ActivityLog
	query := s.db.Preload("User").Order("created_at DESC")

	if actionFilter != "" {
		query = query.Where("action = ?", actionFilter)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&activities).Error
	return activities, err
}

// GetActivitiesByDateRange retrieves activities within a date range
func (s *ActivityLogService) GetActivitiesByDateRange(startDate, endDate time.Time, limit, offset int) ([]models.ActivityLog, error) {
	var activities []models.ActivityLog
	query := s.db.Preload("User").
		Where("created_at >= ? AND created_at <= ?", startDate, endDate).
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&activities).Error
	return activities, err
}

// GetActivityStats returns activity statistics
func (s *ActivityLogService) GetActivityStats() (map[string]interface{}, error) {
	var stats struct {
		TotalActivities int64 `json:"total_activities"`
		UniqueUsers     int64 `json:"unique_users"`
		TodayActivities int64 `json:"today_activities"`
	}

	// Total activities
	s.db.Model(&models.ActivityLog{}).Count(&stats.TotalActivities)

	// Unique users
	s.db.Model(&models.ActivityLog{}).Distinct("id_user").Count(&stats.UniqueUsers)

	// Today's activities
	today := time.Now().Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)
	s.db.Model(&models.ActivityLog{}).Where("created_at >= ? AND created_at < ?", today, tomorrow).Count(&stats.TodayActivities)

	result := map[string]interface{}{
		"total_activities": stats.TotalActivities,
		"unique_users":     stats.UniqueUsers,
		"today_activities": stats.TodayActivities,
	}

	return result, nil
}

// CleanOldLogs removes activity logs older than specified days
func (s *ActivityLogService) CleanOldLogs(days int) error {
	cutoffDate := time.Now().AddDate(0, 0, -days)
	return s.db.Where("created_at < ?", cutoffDate).Delete(&models.ActivityLog{}).Error
}