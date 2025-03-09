package main

import (
	"log"
	"os"
	"subscription/internal/config"
	"subscription/internal/handlers"
	"subscription/internal/middleware"
	"subscription/internal/repository"
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

	// Setup handlers
	planHandler := handlers.NewPlanHandler(planService)
	subscriptionHandler := handlers.NewSubscriptionHandler(subscriptionService)

	// Middleware
	app.Use(cors.New())
	app.Use(logger.New())

	// Public routes

	// Protected routes with API key
	api := app.Group("/api", middleware.Middleware(cfg.APIKey))

	// Plans routes
	plans := api.Group("/plans")
	plans.Get("/", planHandler.HandleGetAllPlans)
	plans.Get("/:id", planHandler.HandleGetPlan)

	// Admin routes for plan management
	planAdmin := api.Group("/admin/plans")
	planAdmin.Post("/", planHandler.HandleCreatePlan)
	planAdmin.Put("/:id", planHandler.HandleUpdatePlan)
	planAdmin.Delete("/:id", planHandler.HandleDeletePlan)

	subscriptions := api.Group("/subscriptions")
	subscriptions.Post("/", subscriptionHandler.HandleCreateSubscription)
	subscriptions.Get("/tutor/:tutorId", subscriptionHandler.HandleGetTutorSubscription)
	subscriptions.Put("/:id/cancel", subscriptionHandler.HandleCancelSubscription)
	subscriptions.Put("/:id/change-plan", subscriptionHandler.HandleChangePlan)
	subscriptions.Post("/confirm", subscriptionHandler.HandleConfirmSubscription)

	subscriptionAdmin := api.Group("/admin/subscriptions")
	subscriptionAdmin.Get("/", subscriptionHandler.HandleGetAllSubscriptions)
	subscriptionAdmin.Put("/:id/status", subscriptionHandler.HandleUpdateSubscriptionStatus)

	// Payment webhooks
	webhooks := api.Group("/webhooks")
	webhooks.Post("/payment", subscriptionHandler.HandlePaymentWebhook)

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
