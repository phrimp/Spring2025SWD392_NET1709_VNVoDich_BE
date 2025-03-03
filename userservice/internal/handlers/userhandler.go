package handlers

import (
	"fmt"
	"strconv"
	"user-service/internal/models"
	"user-service/internal/services"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type RequestParam struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Full_name string `json:"fullname"`
	Phone     string `json:"phone"`
}

func GetUserWithUsernamePasswordHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req RequestParam
		if err := c.BodyParser(&req); err != nil {
			fmt.Println("Error Get User:", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request",
			})
		}
		user, err := services.FindUserWithUsernamePassword(req.Username, req.Password, db)
		if err != nil {
			fmt.Println("Error Get User:", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err,
			})
		}
		return c.JSON(user)
	}
}

func AddUser(db *gorm.DB, had_admin bool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req RequestParam
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request",
			})
		}

		// Convert RequestParam to UserCreationParams
		params := models.UserCreationParams{
			Username: req.Username,
			Password: req.Password,
			Email:    req.Email,
			Role:     req.Role,
			FullName: req.Full_name,
			Phone:    req.Phone,
		}

		// Use our updated AddUser function that handles role-specific records
		err := services.AddUser(params, had_admin, db)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "OK",
			"message": "User created successfully with " + req.Role + " role",
		})
	}
}

func GetPublicUser(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req RequestParam
		if err := c.BodyParser(&req); err != nil {
			fmt.Println("Error Get User:", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request",
			})
		}
		user, err := services.FindUserWithUsername(req.Username, db)
		if err != nil {
			fmt.Println("Error Get User:", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err,
			})
		}

		// Get role-specific information
		switch user.Role {
		case models.RoleTutor:
			// Get tutor details
			var tutor models.Tutor
			if err := db.Where("id = ?", user.ID).First(&tutor).Error; err != nil {
				fmt.Println("Error getting tutor details:", err)
				// Return basic user info if tutor details not found
				return c.JSON(user)
			}
			// Return combined user and tutor info
			return c.JSON(fiber.Map{
				"user":  user,
				"tutor": tutor,
			})
		case models.RoleParent:
			// Get parent details
			var parent models.Parent
			if err := db.Where("id = ?", user.ID).First(&parent).Error; err != nil {
				fmt.Println("Error getting parent details:", err)
				// Return basic user info if parent details not found
				return c.JSON(user)
			}
			// Return combined user and parent info
			return c.JSON(fiber.Map{
				"user":   user,
				"parent": parent,
			})
		default:
			// For other roles, just return the user info
			return c.JSON(user)
		}
	}
}

func GetAllUser(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		page, _ := strconv.Atoi(c.Query("page", "1"))
		limit, _ := strconv.Atoi(c.Query("limit", "10"))

		filters := make(map[string]interface{})

		if role := c.Query("role"); role != "" {
			filters["role"] = role
		}

		if status := c.Query("status"); status != "" {
			filters["status"] = status
		}
		if search := c.Query("search"); search != "" {
			filters["search"] = search
		}

		if verified := c.Query("is_verified"); verified != "" {
			isVerified, err := strconv.ParseBool(verified)
			if err == nil {
				filters["is_verified"] = isVerified
			}
		}

		if from := c.Query("created_from"); from != "" {
			filters["created_from"] = from
		}

		if to := c.Query("created_to"); to != "" {
			filters["created_to"] = to
		}

		if sort := c.Query("sort"); sort != "" {
			filters["sort"] = sort
		}

		if sortDir := c.Query("sort_dir"); sortDir != "" {
			filters["sort_dir"] = sortDir
		}

		// Get paginated users with filters
		response, err := services.GetAllUser(db, page, limit, filters)
		if err != nil {
			fmt.Println("Error getting users:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve users: " + err.Error(),
			})
		}

		return c.JSON(response)
	}
}

func GetUserwithUsername(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		username := c.Query("username")
		user, err := services.FindUserWithUsername(username, db)
		if err != nil {
			fmt.Println("Error Get User:", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err,
			})
		}
		return c.JSON(user)
	}
}

func UpdateUser(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		username := c.Query("username")

		var updateReq models.UserUpdateParams
		if err := c.BodyParser(&updateReq); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		updatedUser, err := services.UpdateUser(username, updateReq, db)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(updatedUser)
	}
}

func UpdateUserStatus(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		username := c.Query("username")
		status := c.Query("status")

		updatedUser, err := services.UpdateUserStatus(username, status, db)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(updatedUser)
	}
}
