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
	cfg := config.GetConfig()

	// Initialize app
	app := fiber.New(fiber.Config{
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	})

	// Middleware
	app.Use(cors.New())
	app.Use(logger.New())

	// Initialize handlers
	googleHandler := handlers.NewGoogleHandler(cfg.GoogleAuth)
	emailHandler := handlers.NewEmailHandler(cfg.Email)

	// Routes
	api := app.Group("/api")

	// Google OAuth routes
	auth := api.Group("/auth")
	auth.Get("/google/login", googleHandler.HandleGoogleLogin)
	auth.Get("/google/callback", googleHandler.HandleGoogleCallback)

	// Email routes
	email := api.Group("/email")
	email.Post("/send", emailHandler.HandleSendPlainEmail)
	email.Post("/send/verify/email", emailHandler.HandleVerifyEmail)
	// email.Post("/send/verification")

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8084"
	}
	log.Fatal(app.Listen(":" + port))
}
