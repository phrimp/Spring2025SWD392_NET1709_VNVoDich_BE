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
		google_access_token := c.Query("google_token")

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
		err := services.AddUser(params, had_admin, google_access_token, db)
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

		userData, err := services.FindUserWithUsername(req.Username, db)
		if err != nil {
			fmt.Println("Error Get User:", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(userData)
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
		userData, err := services.FindUserWithUsername(username, db)
		if err != nil {
			fmt.Println("Error Get User:", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(userData)
	}
}

func UpdateUser(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		username := c.Query("username")
		userid := c.Query("id")

		var updateReq models.UserUpdateParams
		if err := c.BodyParser(&updateReq); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		updatedUser, err := services.UpdateUser(username, userid, updateReq, db)
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

// In userservice/internal/handlers/userhandler.go

func DeleteUserHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		username := c.Query("username")
		if username == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Username cannot be empty",
			})
		}

		if err := services.SoftDeleteUser(username, db); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "User has been deleted",
		})
	}
}

func CancelDeleteUserHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		username := c.Query("username")
		if username == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Username cannot be empty",
			})
		}

		if err := services.CancelDeleteUser(username, db); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "User deletion has been canceled",
		})
	}
}

func AdminUpdateUserHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userid := c.Query("id")
		username := c.Query("name")

		var updateReq models.UserUpdateParams
		if err := c.BodyParser(&updateReq); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		updatedUser, err := services.UpdateUser(username, userid, updateReq, db)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"message": "User updated successfully",
			"user":    updatedUser,
		})
	}
}

func VerifyUserHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		username := c.Query("username")
		if username == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Username is required",
			})
		}
		err := services.VerifyUser(username, db)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "User verification status updated successfully",
		})
	}
}

func AdminDeleteUserHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		username := c.Query("username")
		userid := c.Query("id")

		// Call the service to permanently delete the user
		if err := services.HardDeleteUser(userid, username, db); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "User has been permanently deleted",
		})
	}
}

func AdminAssignRoleHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		username := c.Query("username")
		if username == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Username is required",
			})
		}

		// Parse the request body
		var req struct {
			Role string `json:"role"`
		}

		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		if req.Role == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Role is required",
			})
		}

		// Call the service to assign the role
		if err := services.AssignRole(username, req.Role, db); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "User role has been updated successfully",
		})
	}
}

func UpdatePassword(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		new_password := c.Query("new_password")
		if new_password == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "new password is required",
			})
		}

		cur_password := c.Query("cur_password")
		if cur_password == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "current password is required",
			})
		}

		username := c.Query("username")
		if username == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "username is required",
			})
		}

		if err := services.UpdatePassword(new_password, cur_password, username, db); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "User Password has been updated successfully",
		})
	}
}

func CheckUserStatusHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		username := c.Query("username")
		userID := c.Query("user_id")

		if username == "" && userID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Either username or user_id must be provided",
			})
		}

		var isActive bool
		var err error

		if username != "" {
			isActive, err = services.IsUserActive(username, db)
		} else {
			id, parseErr := strconv.ParseUint(userID, 10, 32)
			if parseErr != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Invalid user_id format",
				})
			}
			isActive, err = services.IsUserActiveByID(uint(id), db)
		}

		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"username":  username,
			"user_id":   userID,
			"is_active": isActive,
		})
	}
}
