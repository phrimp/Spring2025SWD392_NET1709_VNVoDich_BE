package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func handleLogin(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		fmt.Println("DADWAFAWFWAFAWFAWFWA")
		var req LoginRequest
		if err := c.BodyParser(&req); err != nil {
			fmt.Println("Invalid request:", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request",
			})
		}

		// Find user
		var user User
		if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid credentials",
			})
		}

		// Check password
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid credentials",
			})
		}

		// Generate JWT
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["userId"] = user.ID
		claims["username"] = user.Username
		claims["role"] = user.Role
		claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

		// Sign token
		tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not generate token",
			})
		}

		// Clear password before sending
		user.Password = ""

		return c.JSON(LoginResponse{
			Token: tokenString,
			User:  user,
		})
	}
}

func handleRegister(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var user User
		if err := c.BodyParser(&user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		// Check if username already exists
		var existingUser User
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
