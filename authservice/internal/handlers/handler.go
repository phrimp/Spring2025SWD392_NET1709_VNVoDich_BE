package handlers

import (
	"authservice/internal/services"
	"authservice/utils"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
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
		var req LoginRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request",
			})
		}

		// Forward to user service
		user, err := services.GetUserFromUserService(utils.SERVICES_ROUTES.UserService, req.Username, req.Password)
		if err != nil {
			fmt.Println("Login error:", err)
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

		return c.JSON(LoginResponse{
			Token: tokenString,
			User:  *user,
		})
	}
}

func HandleRegister(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req RegisterRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Request"})
		}

		err := services.AddUserUserService(utils.SERVICES_ROUTES.UserService, req.Username, req.Password, req.Email, req.Role)
		if err != nil {
			fmt.Println("Register Error:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"Status": "User Created"})
	}
}
