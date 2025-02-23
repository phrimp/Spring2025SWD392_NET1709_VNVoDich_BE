package services

import (
	"errors"
	"fmt"
	"log"
	"user-service/internal/models"

	"gorm.io/gorm"
)

var LoginAttempt map[string]int

// Login
func FindUserWithUsernamePassword(username, password string, db *gorm.DB) (*models.User, error) {
	var user models.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("invalid credentials")
		}
		return nil, fmt.Errorf("database error: %v", err)
	}
	if err := user.IsAccountValid(); err != nil {
		return nil, fmt.Errorf("account invalid: %v", err)
	}

	// Check password
	if err := user.CheckPassword(password); err != nil {
		if login_attempt, exist := LoginAttempt[user.Username]; exist {
			if login_attempt == 5 {
				user.Lock(10)
				return nil, fmt.Errorf("account is temporary locked due to exceeding number of logins allowed")
			}
			LoginAttempt[user.Username] = login_attempt + 1
		} else {
			LoginAttempt[user.Username] = 0
		}
		return nil, fmt.Errorf("invalid credentials")
	}

	// Clear password before returning
	user.Password = ""
	return &user, nil
}

func AddUser(params models.UserCreationParams, db *gorm.DB) error {
	// Check for existing username
	log.Printf("Adding user with username: %s", params.Username)

	var existingUser models.User
	if err := db.Where("username = ?", params.Username).First(&existingUser).Error; err == nil {
		return fmt.Errorf("username already exists")
	}

	// Check for existing email
	if err := db.Where("email = ?", params.Email).First(&existingUser).Error; err == nil {
		return fmt.Errorf("email already exists")
	}

	// Create new user instance
	user := models.User{
		Username: params.Username,
		Email:    params.Email,
		Role:     models.UserRole(params.Role),
		Status:   models.StatusActive, // Default status
	}

	if ok, err := user.EmailValidation(); !ok {
		return fmt.Errorf("email validation failed: %v", err)
	}

	// Set optional fields if provided
	if params.FullName != "" {
		user.FullName = &params.FullName
	}
	if params.Phone != "" {
		user.Phone = &params.Phone
	}

	user.Password = params.Password
	log.Printf("User struct created with password: %s", user.Password)

	// Set default role if not provided
	if user.Role == "" {
		user.Role = models.RoleParent
	}

	// Validate the entire user model
	if err := user.Validate(); err != nil {
		return fmt.Errorf("validation failed: %v", err)
	}

	// Create user in database within a transaction
	err := db.Transaction(func(tx *gorm.DB) error {
		log.Println("Starting transaction")
		if err := tx.Create(&user).Error; err != nil {
			log.Printf("Error creating user: %v", err)
			return fmt.Errorf("failed to create user: %v", err)
		}
		log.Printf("User created with final password: %s", user.Password)
		// more here

		return nil
	})
	if err != nil {
		return err
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
