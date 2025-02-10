package handlers

import (
	"authservice/internal/repository"
	"authservice/internal/services"
	"authservice/utils"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID   uint   `json:"userId"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

func HandleLogin(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req repository.LoginRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request",
			})
		}

		// Forward to user service
		user, err := services.GetUserFromUserService(utils.SERVICES_ROUTES.UserService, req.Username, req.Password)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid credentials",
			})
		}

		// Generate JWT
		claims := Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				Issuer:    "auth-service",
			},
			UserID:   user.ID,
			Username: user.Username,
			Role:     user.Role,
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not generate token",
			})
		}

		return c.JSON(repository.LoginResponse{
			Token: tokenString,
			User:  *user,
		})
	}
}

func HandleRegister(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var user repository.User
		if err := c.BodyParser(&user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		// Check if username already exists
		var existingUser repository.User
		if err := db.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Username already exists",
			})
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not hash password",
			})
		}
		user.Password = string(hashedPassword)

		// Set default role if not provided
		if user.Role == "" {
			user.Role = "user"
		}

		// Create user
		if err := db.Create(&user).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not create user",
			})
		}

		// Clear password before sending response
		user.Password = ""

		return c.Status(fiber.StatusCreated).JSON(user)
	}
}
