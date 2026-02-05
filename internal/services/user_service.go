package services

import (
	"swipeup-be/internal/models"

	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := s.db.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (s *UserService) CreateUser(user *models.User) error {
	return s.db.Create(user).Error
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	var users []models.User
	if err := s.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserService) GetAllUsersPaginated(limit, offset int) ([]models.User, int64, error) {
	var users []models.User
	var count int64
	err := s.db.Model(&models.User{}).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	err = s.db.Limit(limit).Offset(offset).Find(&users).Error
	return users, count, err
}

func (s *UserService) GetUserByID(id string) (*models.User, error) {
	var user models.User
	err := s.db.First(&user, id).Error
	return &user, err
}

// UpdateUser updates a user's information
func (s *UserService) UpdateUser(id uint, updates map[string]interface{}) error {
	return s.db.Model(&models.User{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteUser soft deletes a user
func (s *UserService) DeleteUser(id uint) error {
	return s.db.Delete(&models.User{}, id).Error
}