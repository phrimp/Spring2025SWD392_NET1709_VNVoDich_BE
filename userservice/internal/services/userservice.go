package services

import (
	"errors"
	"fmt"
	"math"
	"strings"
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

func RegisterUserWithRole(params models.UserCreationParams, db *gorm.DB) error {
	// Start a database transaction
	return db.Transaction(func(tx *gorm.DB) error {
		user := models.User{
			Username: params.Username,
			Email:    params.Email,
			Role:     models.UserRole(params.Role),
			Status:   models.StatusActive,
			Password: params.Password,
		}

		// Set optional fields if provided
		if params.FullName != "" {
			user.FullName = &params.FullName
		}
		if params.Phone != "" {
			user.Phone = &params.Phone
		}

		// Validate user data
		if err := user.Validate(); err != nil {
			return fmt.Errorf("validation failed: %v", err)
		}

		// Create user in database
		if err := tx.Create(&user).Error; err != nil {
			return fmt.Errorf("failed to create user: %v", err)
		}

		switch user.Role {
		case models.RoleParent:
			return createParentRecord(tx, user.ID)
		case models.RoleTutor:
			return createTutorRecord(tx, user.ID)
		case models.RoleAdmin:
			// Admin role doesn't need additional records
			return nil
		default:
			return errors.New("unsupported user role")
		}
	})
}

func AddUser(params models.UserCreationParams, had_admin bool, db *gorm.DB) error {
	// Check for existing username
	fmt.Printf("Adding user with username: %s\n", params.Username)

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
	if params.Role == string(models.RoleAdmin) && !had_admin {
		user := models.User{
			Username: params.Username,
			Email:    params.Email,
			Role:     models.RoleAdmin,
			Status:   models.StatusActive,
			Password: params.Password,
		}

		if params.FullName != "" {
			user.FullName = &params.FullName
		}

		if err := db.Create(&user).Error; err != nil {
			return fmt.Errorf("failed to create admin user: %v", err)
		}
		return nil
	}

	return RegisterUserWithRole(params, db)
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

func GetAllUser(db *gorm.DB, page, limit int, filters map[string]interface{}) (*models.PaginatedResponse, error) {
	// Default values if not provided
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10 // Default limit with max cap at 100
	}

	var users []models.User
	query := db.Model(&models.User{})

	if filters != nil {
		// Filter by role
		if role, ok := filters["role"].(string); ok && role != "" {
			query = query.Where("role = ?", role)
		}

		// Filter by status
		if status, ok := filters["status"].(string); ok && status != "" {
			query = query.Where("status = ?", status)
		}

		// Filter by search term (username, email, or full_name)
		if search, ok := filters["search"].(string); ok && search != "" {
			searchTerm := "%" + search + "%"
			query = query.Where("username LIKE ? OR email LIKE ? OR full_name LIKE ?",
				searchTerm, searchTerm, searchTerm)
		}

		// Filter by verified status
		if verified, ok := filters["is_verified"].(bool); ok {
			query = query.Where("is_verified = ?", verified)
		}

		// Filter by created date range
		if from, ok := filters["created_from"].(string); ok && from != "" {
			query = query.Where("created_at >= ?", from)
		}

		if to, ok := filters["created_to"].(string); ok && to != "" {
			query = query.Where("created_at <= ?", to)
		}
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("error counting users: %w", err)
	}

	// Apply sorting - default to newest first
	sortBy := "created_at"
	sortDir := "DESC"

	if sort, ok := filters["sort"].(string); ok && sort != "" {
		sortBy = sort
	}

	if dir, ok := filters["sort_dir"].(string); ok &&
		(strings.ToUpper(dir) == "ASC" || strings.ToUpper(dir) == "DESC") {
		sortDir = strings.ToUpper(dir)
	}

	offset := (page - 1) * limit

	// Execute the query with pagination
	result := query.Order(fmt.Sprintf("%s %s", sortBy, sortDir)).
		Limit(limit).
		Offset(offset).
		Find(&users)

	if result.Error != nil {
		return nil, fmt.Errorf("error fetching users: %w", result.Error)
	}

	for i := range users {
		users[i].Password = ""
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	// Prepare paginated response
	response := &models.PaginatedResponse{
		Data: users,
		Pagination: models.Pagination{
			Total:      total,
			Page:       page,
			PageSize:   limit,
			TotalPages: totalPages,
		},
	}

	return response, nil
}

func DeleteUser(db *gorm.DB, email string) error {
	// Find user
	var user models.User
	if result := db.First(&user, email); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return fmt.Errorf("user not found: %s", result.Error)
		}
		return fmt.Errorf("failed to fetch user: %s", result.Error)
	}

	// Perform soft delete
	if result := db.Delete(&user); result.Error != nil {
		return fmt.Errorf("failed to delete user: %s", result.Error)
	}

	fmt.Println("User successfully deleted")
	return nil
}

func createParentRecord(tx *gorm.DB, userID uint) error {
	parent := models.Parent{
		ID:                   userID,
		PreferredLanguage:    "English", // Default language
		NotificationsEnabled: true,      // Default to enabled
	}

	if err := tx.Create(&parent).Error; err != nil {
		return fmt.Errorf("failed to create parent record: %v", err)
	}

	return nil
}

func createTutorRecord(tx *gorm.DB, userID uint) error {
	tutor := models.Tutor{
		ID:             userID,
		Bio:            "",    // Empty bio initially
		Qualifications: "",    // Empty qualifications initially
		TeachingStyle:  "",    // Empty teaching style initially
		IsAvailable:    false, // Not available by default
		DemoVideoURL:   "",    // No demo video initially
		Image:          "",    // No image initially
	}

	if err := tx.Create(&tutor).Error; err != nil {
		return fmt.Errorf("failed to create tutor record: %v", err)
	}

	// Create default availability for the tutor
	availability := models.Availability{
		TutorID: userID,
		TimeGap: 10, // Default 10-minute gap between sessions
	}

	if err := tx.Create(&availability).Error; err != nil {
		return fmt.Errorf("failed to create tutor availability: %v", err)
	}

	return nil
}
