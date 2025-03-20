package routes

import (
	"subscription/internal/handlers"
	"subscription/internal/services"

	"github.com/gofiber/fiber/v2"
)

// SetupSubscriptionRoutes configures the routes for subscription management
func SetupSubscriptionRoutes(api fiber.Router, subscriptionService services.SubscriptionService) {
	// Initialize the handler
	subscriptionHandler := handlers.NewSubscriptionHandler(subscriptionService)

	// User subscription routes
	subscriptions := api.Group("/subscriptions")
	subscriptions.Post("/", subscriptionHandler.HandleCreateSubscription)
	subscriptions.Get("/tutor/:tutorId", subscriptionHandler.HandleGetTutorSubscription)
	subscriptions.Put("/:id/cancel", subscriptionHandler.HandleCancelSubscription)
	subscriptions.Put("/:id/change-plan", subscriptionHandler.HandleChangePlan)
	subscriptions.Post("/confirm", subscriptionHandler.HandleConfirmSubscription)

	// Admin subscription routes
	subscriptionAdmin := api.Group("/admin/subscriptions")
	subscriptionAdmin.Get("/", subscriptionHandler.HandleGetAllSubscriptions)
	subscriptionAdmin.Put("/:id/status", subscriptionHandler.HandleUpdateSubscriptionStatus)
}
