package services

import (
	"encoding/json"
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

func RegisterUserWithRole(params models.UserCreationParams, google_access_token string, db *gorm.DB) error {
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
		if google_access_token != "" {
			user.GoogleToken = google_access_token
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

func AddUser(params models.UserCreationParams, had_admin bool, google_access_token string, db *gorm.DB) error {
	// Check for existing username
	fmt.Printf("Adding user with username: %s\n", params.Username)

	var existingUser models.User
	if err := db.Where("username = ?", params.Username).First(&existingUser).Error; err == nil {
		if google_access_token != "" {
			err := UpdateGoogleToken(params.Username, google_access_token, db)
			if err != nil {
				return fmt.Errorf("update google token failed: %s", err.Error())
			}
			return fmt.Errorf("user is already exists, update google token")
		}
		return fmt.Errorf("username already exists")
	}

	// Check for existing email
	if err := db.Where("email = ?", params.Email).First(&existingUser).Error; err == nil {
		if google_access_token != "" {
			err := UpdateGoogleToken(existingUser.Username, google_access_token, db)
			if err != nil {
				return fmt.Errorf("update google token failed: %s", err.Error())
			}
			return fmt.Errorf("user email is already exist, update google token")
		}
		return fmt.Errorf("email already exists")
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

	return RegisterUserWithRole(params, google_access_token, db)
}

func FindUserWithUsername(username string, db *gorm.DB) (*map[string]interface{}, error) {
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

	// Remove sensitive data
	user.Password = ""

	response := make(map[string]interface{})

	userData, err := json.Marshal(user)
	if err != nil {
		return nil, fmt.Errorf("error serializing user data: %w", err)
	}

	if err := json.Unmarshal(userData, &response); err != nil {
		return nil, fmt.Errorf("error building response: %w", err)
	}

	switch user.Role {
	case models.RoleTutor:
		var tutor models.Tutor
		if err := db.Where("id = ?", user.ID).First(&tutor).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("error fetching tutor data: %w", err)
			}
			// If no tutor record is found, continue without it
		} else {
			tutorData, err := json.Marshal(tutor)
			if err != nil {
				return nil, fmt.Errorf("error serializing tutor data: %w", err)
			}

			var tutorMap map[string]interface{}
			if err := json.Unmarshal(tutorData, &tutorMap); err != nil {
				return nil, fmt.Errorf("error building tutor response: %w", err)
			}

			response["tutor"] = tutorMap

			// Fetch tutor specialties
			var specialties []models.TutorSpecialty
			if err := db.Where("tutor_id = ?", user.ID).Find(&specialties).Error; err != nil {
				// Just log the error but continue
				fmt.Printf("Error fetching tutor specialties: %v\n", err)
			} else if len(specialties) > 0 {
				specialtiesData, _ := json.Marshal(specialties)
				var specialtiesMap []map[string]interface{}
				json.Unmarshal(specialtiesData, &specialtiesMap)

				response["tutor"].(map[string]interface{})["specialties"] = specialtiesMap
			}

			// Fetch availability
			var availability models.Availability
			if err := db.Where("tutor_id = ?", user.ID).First(&availability).Error; err != nil {
				fmt.Printf("Error fetching tutor availability: %v\n", err)
			} else {
				// Fetch days of availability
				var days []models.DayAvailability
				if err := db.Where("availability_id = ?", availability.ID).Find(&days).Error; err == nil && len(days) > 0 {
					daysData, _ := json.Marshal(days)
					var daysMap []map[string]interface{}
					json.Unmarshal(daysData, &daysMap)

					availabilityData, _ := json.Marshal(availability)
					var availabilityMap map[string]interface{}
					json.Unmarshal(availabilityData, &availabilityMap)

					availabilityMap["days"] = daysMap
					response["tutor"].(map[string]interface{})["availability"] = availabilityMap
				}
			}
		}

	case models.RoleParent:
		var parent models.Parent
		if err := db.Where("id = ?", user.ID).First(&parent).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("error fetching parent data: %w", err)
			}
			// If no parent record is found, continue without it
		} else {
			parentData, err := json.Marshal(parent)
			if err != nil {
				return nil, fmt.Errorf("error serializing parent data: %w", err)
			}

			var parentMap map[string]interface{}
			if err := json.Unmarshal(parentData, &parentMap); err != nil {
				return nil, fmt.Errorf("error building parent response: %w", err)
			}

			response["parent"] = parentMap

			// Fetch children
			var children []models.Children
			if err := db.Where("parent_id = ?", user.ID).Find(&children).Error; err != nil {
				fmt.Printf("Error fetching children: %v\n", err)
			} else if len(children) > 0 {
				childrenData, _ := json.Marshal(children)
				var childrenMap []map[string]interface{}
				json.Unmarshal(childrenData, &childrenMap)

				// Remove sensitive data like passwords
				for i := range childrenMap {
					delete(childrenMap[i], "password")
				}

				response["parent"].(map[string]interface{})["children"] = childrenMap
			}
		}
	}

	return &response, nil
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

func UpdateUser(username string, params models.UserUpdateParams, db *gorm.DB) (*models.User, error) {
	var user models.User
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}

	result := db.Model(&models.User{}).
		Where("username = ?", username).
		First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("error finding user: %w", result.Error)
	}
	// Update fields if provided
	updates := make(map[string]interface{})

	if params.Email != "" {
		// Validate email
		emailValidation, _ := (&models.User{Email: params.Email}).EmailValidation()
		if !emailValidation {
			return nil, fmt.Errorf("invalid email format")
		}

		// Check if email already exists for another user
		var existingUser models.User
		if err := db.Where("email = ?", params.Email).First(&existingUser).Error; err == nil {
			return nil, fmt.Errorf("email already in use by another account")
		}

		updates["email"] = params.Email
	}

	if params.FullName != "" {
		fullName := params.FullName
		updates["full_name"] = &fullName
	}

	if params.Phone != "" {
		phone := params.Phone
		updates["phone"] = &phone
	}

	if params.Picture != "" {
		updates["picture"] = params.Picture
	}

	if params.Status != "" && user.Role == models.RoleAdmin {
		// Validate status
		switch models.UserStatus(params.Status) {
		case models.StatusActive, models.StatusBanned, models.StatusDeleted:
			updates["status"] = params.Status
		default:
			return nil, fmt.Errorf("invalid status value")
		}
	}
	// Apply updates
	if len(updates) > 0 {
		if err := db.Model(&user).Updates(updates).Error; err != nil {
			return nil, fmt.Errorf("failed to update user: %w", err)
		}

		// Fetch updated user
		if err := db.First(&user, user.ID).Error; err != nil {
			return nil, fmt.Errorf("error retrieving updated user: %w", err)
		}
	}
	user.Password = ""
	return &user, nil
}

func UpdateUserStatus(username string, status string, db *gorm.DB) (*models.User, error) {
	// Find user
	var user models.User
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}

	result := db.Model(&models.User{}).
		Where("username = ?", username).
		First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("error finding user: %w", result.Error)
	}

	switch models.UserStatus(status) {
	case models.StatusActive:
		user.Status = models.StatusActive
	case models.StatusBanned:
		user.Status = models.StatusBanned
	default:
		return nil, fmt.Errorf("invalid status value: must be 'active', 'banned', or 'pending'")
	}

	// Save updated status
	if err := db.Save(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to update user status: %w", err)
	}

	user.Password = ""

	return &user, nil
}

func UpdateGoogleToken(username string, googleToken string, db *gorm.DB) error {
	if username == "" {
		return errors.New("username cannot be empty")
	}

	if googleToken == "" {
		return errors.New("google token cannot be empty")
	}

	// Find the user
	var user models.User
	result := db.Model(&models.User{}).
		Where("username = ?", username).
		First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return fmt.Errorf("error finding user: %w", result.Error)
	}

	// Update Google token
	if err := db.Model(&user).Update("google_token", googleToken).Error; err != nil {
		return fmt.Errorf("failed to update google token: %w", err)
	}

	return nil
}

func SoftDeleteUser(username string, db *gorm.DB) error {
	// Find user
	var user models.User
	if username == "" {
		return errors.New("username cannot be empty")
	}

	result := db.Model(&models.User{}).
		Where("username = ?", username).
		First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return fmt.Errorf("error finding user: %w", result.Error)
	}

	// Update user status to "Deleted" instead of actually deleting
	if err := db.Model(&user).Update("status", models.StatusDeleted).Error; err != nil {
		return fmt.Errorf("failed to update user status: %w", err)
	}

	return nil
}

func CancelDeleteUser(username string, db *gorm.DB) error {
	// Find user
	var user models.User
	if username == "" {
		return errors.New("username cannot be empty")
	}

	result := db.Model(&models.User{}).
		Where("username = ?", username).
		First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return fmt.Errorf("error finding user: %w", result.Error)
	}

	// Check if the user status is actually "Deleted"
	if user.Status != models.StatusDeleted {
		return errors.New("user is not in deleted status")
	}

	// Update user status back to "Active"
	if err := db.Model(&user).Update("status", models.StatusActive).Error; err != nil {
		return fmt.Errorf("failed to update user status: %w", err)
	}

	return nil
}
