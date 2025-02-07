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
	return &user, nil
}

func AddUser(username, password, email, role string, db *gorm.DB) error {
	user := models.User{Username: username, Password: password, Email: email, Role: role}
	if err := db.Save(user).Error; err != nil {
		return fmt.Errorf("add user error: %v", err)
	}
	return nil
}
