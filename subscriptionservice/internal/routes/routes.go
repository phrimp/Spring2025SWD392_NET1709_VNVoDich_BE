package routes

import (
	"subscription/internal/services"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(api fiber.Router, planService services.PlanService, subscriptionService services.SubscriptionService) {
	// Register all route groups
	SetupPlanRoutes(api, planService)
	SetupSubscriptionRoutes(api, subscriptionService)
	SetupWebhookRoutes(api, subscriptionService)
}
