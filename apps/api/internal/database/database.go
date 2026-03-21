package database

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log"

	"pocketpanel/api/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect(path string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Server{},
	)
}

// generateRandomPassword generates a secure random password (base64 encoded, 32 bytes)
func generateRandomPassword() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// CreateDefaultAdmin creates a default admin user if it doesn't exist
// Returns the generated password (only when newly created) or empty string if user already exists
func CreateDefaultAdmin(db *gorm.DB) (string, error) {
	// Check if admin user already exists
	var existingUser models.User
	result := db.Where("username = ?", "admin").First(&existingUser)
	if result.Error == nil {
		// Admin already exists
		return "", nil
	}
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// Some other error occurred
		return "", result.Error
	}

	// Generate random password
	plainPassword, err := generateRandomPassword()
	if err != nil {
		return "", fmt.Errorf("failed to generate password: %w", err)
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	// Create admin user
	admin := models.User{
		Username: "admin",
		Password: string(hashedPassword),
	}

	if err := db.Create(&admin).Error; err != nil {
		return "", fmt.Errorf("failed to create admin user: %w", err)
	}

	log.Println("Default admin account created successfully")
	return plainPassword, nil
}
