package handlers

import (
	"fmt"
	"strconv"
	"subscription/internal/models"
	"subscription/internal/services"

	"github.com/gofiber/fiber/v2"
)

type PlanHandler struct {
	planService services.PlanService
}

func NewPlanHandler(planService services.PlanService) *PlanHandler {
	return &PlanHandler{
		planService: planService,
	}
}

func (h *PlanHandler) HandleGetAllPlans(c *fiber.Ctx) error {
	activeOnly := c.Query("active_only", "true") == "true"

	plans, err := h.planService.GetAllPlans(activeOnly)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get subscription plans: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": plans,
	})
}

func (h *PlanHandler) HandleGetPlan(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid plan ID",
		})
	}

	plan, err := h.planService.GetPlanByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Subscription plan not found: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": plan,
	})
}

func (h *PlanHandler) HandleCreatePlan(c *fiber.Ctx) error {
	var req models.PlanRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body: " + err.Error(),
		})
	}

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Plan name is required",
		})
	}
	fmt.Println(req)

	if req.PriceMonthly <= 0 || req.PriceAnnually <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Plan prices must be greater than 0",
		})
	}

	if req.MaxCourses <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Maximum courses must be greater than 0",
		})
	}

	if req.CommissionRate < 0 || req.CommissionRate > 100 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Commission rate must be between 0 and 100",
		})
	}

	// Create the plan
	plan, err := h.planService.CreatePlan(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create subscription plan: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data":    plan,
		"message": "Subscription plan created successfully",
	})
}

func (h *PlanHandler) HandleUpdatePlan(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid plan ID",
		})
	}

	var req models.PlanRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body: " + err.Error(),
		})
	}

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Plan name is required",
		})
	}

	if req.PriceMonthly <= 0 || req.PriceAnnually <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Plan prices must be greater than 0",
		})
	}

	if req.MaxCourses <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Maximum courses must be greater than 0",
		})
	}

	if req.CommissionRate < 0 || req.CommissionRate > 100 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Commission rate must be between 0 and 100",
		})
	}

	plan, err := h.planService.UpdatePlan(uint(id), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update subscription plan: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data":    plan,
		"message": "Subscription plan updated successfully",
	})
}

func (h *PlanHandler) HandleDeletePlan(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid plan ID",
		})
	}

	if err := h.planService.DeletePlan(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete subscription plan: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Subscription plan deleted successfully",
	})
}
