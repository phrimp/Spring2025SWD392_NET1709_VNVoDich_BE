package main

import (
	"authservice/internal/handlers"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

var API_KEY string

func init() {
	godotenv.Load(".env")
	API_KEY = os.Getenv("API_KEY")
}

func main() {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "ok",
			"service": "auth-service",
		})
	})

	auth := app.Group("", Middleware(API_KEY))

	// Routes
	auth.Post("/login", handlers.HandleLogin())
	auth.Post("/register", handlers.HandleRegister())

	port := os.Getenv("AUTH_SERVICE_PORT")
	if port == "" {
		port = "8081" // default port
	}

	if err := app.Listen(":" + port); err != nil {
		panic(fmt.Sprintf("Failed to start server: %v", err))
	}
}
