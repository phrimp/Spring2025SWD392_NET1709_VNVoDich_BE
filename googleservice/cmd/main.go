package main

import (
	"google-service/internal/config"
	"google-service/internal/handlers"
	"google-service/internal/middleware"
	"google-service/utils"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func init() {
	utils.SetupTimeZone()
}

func main() {
	cfg := config.Google_config

	// Initialize app
	app := fiber.New(fiber.Config{
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	})

	app.Use(cors.New())
	app.Use(logger.New())

	// Initialize handlers
	googleHandler := handlers.NewGoogleHandler(cfg.GoogleAuth)
	emailHandler := handlers.NewEmailHandler(cfg.Email)
	meetHandler := handlers.NewMeetHandler(cfg.GoogleAuth)
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "ok",
			"service": "google-service",
		})
	})
	// Routes
	api := app.Group("/api", middleware.Middleware(os.Getenv("API_KEY")))

	// Google OAuth routes
	auth := api.Group("/auth")
	auth.Get("/google/login", googleHandler.HandleGoogleLogin)
	auth.Get("/google/callback", googleHandler.HandleGoogleCallback)

	// Email routes
	email := api.Group("/email")
	email.Post("/send", emailHandler.HandleSendPlainEmail)
	email.Post("/send/verify/email", emailHandler.HandleVerifyEmail)

	// Meet routes
	meet := api.Group("/meet")
	meet.Post("/create", meetHandler.CreateMeetWithEmail)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8084"
	}
	log.Fatal(app.Listen(":" + port))
}
