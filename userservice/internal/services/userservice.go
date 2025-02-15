package services

import (
	"errors"
	"fmt"
	"user-service/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func FindUserWithUsernamePassword(username, password string, db *gorm.DB) (*models.User, error) {
	var user models.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("invalid credentials")
		}
		return nil, fmt.Errorf("database error: %v", err)
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Clear password before returning
	user.Password = ""
	return &user, nil
}

func AddUser(username, password, email, role, full_name string, db *gorm.DB) error {
	// Check if user exists
	var existingUser models.User
	if err := db.Where("username = ?", username).First(&existingUser).Error; err == nil {
		return fmt.Errorf("username already exists")
	}

	// Hash password before saving
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	user := models.User{
		Username: username,
		Password: string(hashedPassword),
		Email:    email,
		Role:     role,
		Full_name: full_name,
	}
	// Set default if no provide
	if role == "" {
		user.Role = "user"
	}

	if err := db.Create(&user).Error; err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	return nil
}

func FindUserWithUsername(username string, db *gorm.DB) (*models.User, error) {
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}

	var user models.User
	result := db.Model(&models.User{}).
		Where("username = ?", username).
		First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("error finding user: %w", result.Error)
	}
	user.Password = ""

	return &user, nil
}

func GetAllUser(db *gorm.DB, page, limit int) ([]models.User, error) {
	var users []models.User

	offset := (page - 1) * limit

	result := db.Model(&models.User{}).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&users)

	if result.Error != nil {
		return nil, fmt.Errorf("error fetching users: %w", result.Error)
	}

	return users, nil
}
