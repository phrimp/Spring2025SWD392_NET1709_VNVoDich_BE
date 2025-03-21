package routes

import (
	"subscription/internal/handlers"
	"subscription/internal/services"

	"github.com/gofiber/fiber/v2"
)

// SetupWebhookRoutes configures the routes for payment webhooks
func SetupWebhookRoutes(api fiber.Router, subscriptionService services.SubscriptionService) {
	// Initialize the handler
	subscriptionHandler := handlers.NewSubscriptionHandler(subscriptionService)

	// Payment webhook routes
	webhooks := api.Group("/webhooks")
	webhooks.Post("/payment", subscriptionHandler.HandlePaymentWebhook)
}
