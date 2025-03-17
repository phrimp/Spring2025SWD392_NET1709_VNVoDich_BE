package main

import (
	"log"
	"os"
	"subscription/internal/config"
	"subscription/internal/middleware"
	"subscription/internal/repository"
	"subscription/internal/routes"
	"subscription/internal/services"
	"subscription/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func init() {
	utils.SetupTimeZone()
}

func main() {
	cfg := config.New()

	// Initialize database connection
	db := repository.DB

	// Setup repositories
	planRepo := repository.NewPlanRepository(db)
	subscriptionRepo := repository.NewSubscriptionRepository(db)

	// Setup services
	paymentService := services.NewPaymentService(cfg.PaymentServiceURL, cfg.APIKey)
	planService := services.NewPlanService(planRepo)
	subscriptionService := services.NewSubscriptionService(subscriptionRepo, planRepo, paymentService)

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

	// Middleware
	app.Use(cors.New())
	app.Use(logger.New())

	// Setup API router group with middleware
	api := app.Group("/api", middleware.Middleware(cfg.APIKey))

	// Setup all routes
	routes.SetupRoutes(api, planService, subscriptionService)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8086" // Default port for subscription service
	}

	log.Printf("Starting subscription service on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
