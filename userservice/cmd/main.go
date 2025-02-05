package main

import (
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

	auth := app.Group("", Middleware(API_KEY))

	// Health check for docker compose
	auth.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Routes

	port := os.Getenv("USER_SERVICE_PORT")
	if port == "" {
		port = "8085" // default port
	}

	if err := app.Listen(":" + port); err != nil {
		panic(fmt.Sprintf("Failed to start server: %v", err))
	}
}
