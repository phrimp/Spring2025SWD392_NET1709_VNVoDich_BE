// Create this file: adminservice/cmd/main.go if not already exists
package main

import (
	"adminservice/internal/config"
	"adminservice/internal/handlers"
	"adminservice/internal/middleware"
	"adminservice/internal/repository"
	"adminservice/utils"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func init() {
	utils.SetupTimeZone()
}

func main() {
	cfg := config.New()

	// Initialize database connection
	db, err := repository.InitDB(cfg.DatabaseConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize the app
	app := fiber.New(fiber.Config{
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
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

	adminHandler := handlers.NewAdminHandler(cfg)

	app.Use(cors.New())

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "ok",
			"service": "admin-service",
		})
	})

	// API routes with API key middleware
	api := app.Group("/api", middleware.Middleware(cfg.APIKey))
	api.Get("/", handlers.TestHandler(db))
	api.Get("/users", adminHandler.GetAllUsersHandler())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8083" // default port
	}

	if err := app.Listen(":" + port); err != nil {
		panic(fmt.Sprintf("Failed to start server: %v", err))
	}
}
