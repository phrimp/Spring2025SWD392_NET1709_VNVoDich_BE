package main

import (
	"google-service/internal/config"
	"google-service/internal/handlers"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Initialize configs
	// cfg := config.New()
	googleCfg := config.NewGoogleOAuthConfig()

	// Initialize app
	app := fiber.New()

	// Middleware
	app.Use(cors.New())
	app.Use(logger.New())

	// Initialize handlers
	googleHandler := handlers.NewGoogleHandler(googleCfg)

	// Routes
	api := app.Group("/api")

	// Google OAuth routes
	auth := api.Group("/auth")
	auth.Get("/google/login", googleHandler.HandleGoogleLogin)

	// Email routes
	// email := api.Group("/email")

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8084"
	}
	log.Fatal(app.Listen(":" + port))
}
