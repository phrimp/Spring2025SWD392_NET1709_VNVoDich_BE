package routes

import (
	"subscription/internal/handlers"
	"subscription/internal/services"

	"github.com/gofiber/fiber/v2"
)

// SetupPlanRoutes configures the routes for plan management
func SetupPlanRoutes(api fiber.Router, planService services.PlanService) {
	// Initialize the handler
	planHandler := handlers.NewPlanHandler(planService)

	// Public plan routes
	plans := api.Group("/plans")
	plans.Get("/", planHandler.HandleGetAllPlans)
	plans.Get("/:id", planHandler.HandleGetPlan)

	// Admin routes for plan management
	planAdmin := api.Group("/admin/plans")
	planAdmin.Post("/", planHandler.HandleCreatePlan)
	planAdmin.Put("/:id", planHandler.HandleUpdatePlan)
	planAdmin.Delete("/:id", planHandler.HandleDeletePlan)
}
