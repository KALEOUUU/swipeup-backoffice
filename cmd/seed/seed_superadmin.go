package main

import (
	"fmt"
	"log"

	"swipeup-be/internal/config"
	"swipeup-be/internal/database"
	"swipeup-be/internal/models"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Setup database connection
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Check if superadmin already exists
	var existingUser models.User
	if err := db.Where("username = ?", "superadmin").First(&existingUser).Error; err == nil {
		fmt.Println("Superadmin user already exists!")
		return
	}

	// Hash password
	password := "superadmin123" // Change this immediately after first login!
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Failed to hash password:", err)
	}

	// Create superadmin user
	superadmin := models.User{
		Username: "superadmin",
		Password: string(hashedPassword),
		Role:     models.RoleSuperAdmin,
	}

	if err := db.Create(&superadmin).Error; err != nil {
		log.Fatal("Failed to create superadmin:", err)
	}

	fmt.Println("=========================================")
	fmt.Println("Superadmin user created successfully!")
	fmt.Println("=========================================")
	fmt.Println("Username: superadmin")
	fmt.Println("Password: superadmin123")
	fmt.Println("=========================================")
	fmt.Println("IMPORTANT: Change the password immediately after first login!")
	fmt.Println("=========================================")
}
