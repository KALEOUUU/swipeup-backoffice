package services

import (
	"errors"
	"swipeup-be/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	db         *gorm.DB
	jwtSecret  []byte
	tokenExpiry time.Duration
}

type AuthClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required,oneof=superadmin admin_stan siswa"`
}

// RegisterAdminStanRequest is used by superadmin to register admin_stan
type RegisterAdminStanRequest struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required,min=6"`
	NamaStan    string `json:"nama_stan" binding:"required"`
	NamaPemilik string `json:"nama_pemilik" binding:"required"`
	Telp        string `json:"telp"`
}

type AuthResponse struct {
	User   models.User `json:"user"`
	Token  string      `json:"token"`
	StanID *uint       `json:"stan_id,omitempty"` // Only for admin_stan role
}

func NewAuthService(db *gorm.DB, jwtSecret string) *AuthService {
	return &AuthService{
		db:          db,
		jwtSecret:   []byte(jwtSecret),
		tokenExpiry: 24 * time.Hour, // 24 hours
	}
}

// HashPassword hashes a password using bcrypt
func (s *AuthService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword verifies a password against its hash
func (s *AuthService) CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateToken generates a JWT token for a user
func (s *AuthService) GenerateToken(user *models.User) (string, error) {
	claims := AuthClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

// ValidateToken validates a JWT token and returns claims
func (s *AuthService) ValidateToken(tokenString string) (*AuthClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*AuthClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// Register creates a new user account
func (s *AuthService) Register(req RegisterRequest) (*models.User, error) {
	// Check if username already exists
	var existingUser models.User
	if err := s.db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return nil, errors.New("username already exists")
	}

	// Hash password
	hashedPassword, err := s.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := models.User{
		Username: req.Username,
		Password: hashedPassword,
		Role:     models.UserRole(req.Role),
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}

	// Remove password from response
	user.Password = ""
	return &user, nil
}

// Login authenticates a user and returns token
func (s *AuthService) Login(req LoginRequest) (*AuthResponse, error) {
	// Find user by username
	var user models.User
	if err := s.db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("invalid username or password")
		}
		return nil, err
	}

	// Check password
	if !s.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("invalid username or password")
	}

	// Generate token
	token, err := s.GenerateToken(&user)
	if err != nil {
		return nil, err
	}

	// Remove password from response
	user.Password = ""

	// Get stan_id for admin_stan users
	var stanID *uint
	if user.Role == models.RoleAdminStan {
		var stan models.Stan
		if err := s.db.Where("id_user = ?", user.ID).First(&stan).Error; err == nil {
			stanID = &stan.ID
		}
	}

	return &AuthResponse{
		User:   user,
		Token:  token,
		StanID: stanID,
	}, nil
}

// RegisterAdminStan creates a new admin_stan account with stan (public access)
func (s *AuthService) RegisterAdminStan(req RegisterAdminStanRequest) (*models.User, error) {
	// Check if username already exists
	var existingUser models.User
	if err := s.db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return nil, errors.New("username already exists")
	}

	// Hash password
	hashedPassword, err := s.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user with role admin_stan
	user := models.User{
		Username: req.Username,
		Password: hashedPassword,
		Role:     models.RoleAdminStan,
	}

	// Start transaction
	tx := s.db.Begin()

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create stan for this admin
	stan := models.Stan{
		NamaStan:    req.NamaStan,
		NamaPemilik: req.NamaPemilik,
		Telp:        req.Telp,
		IDUser:      user.ID,
		AcceptCash:  true,
		AcceptQris:  false,
	}

	if err := tx.Create(&stan).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	// Remove password from response
	user.Password = ""
	return &user, nil
}

// GetUserByID gets user by ID (without password)
func (s *AuthService) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := s.db.Select("id, username, role, created_at, updated_at").Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}