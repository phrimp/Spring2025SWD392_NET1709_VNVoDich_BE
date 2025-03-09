package handlers

import (
	"encoding/json"
	"strconv"
	"subscription/internal/models"
	"subscription/internal/services"
	"time"

	"github.com/gofiber/fiber/v2"
)

type SubscriptionHandler struct {
	subscriptionService services.SubscriptionService
}

func NewSubscriptionHandler(subscriptionService services.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{
		subscriptionService: subscriptionService,
	}
}

func (h *SubscriptionHandler) HandleCreateSubscription(c *fiber.Ctx) error {
	var req models.SubscriptionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body: " + err.Error(),
		})
	}

	if req.PlanID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Plan ID is required",
		})
	}

	if req.BillingCycle != models.BillingMonthly && req.BillingCycle != models.BillingAnnually {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Billing cycle must be either 'monthly' or 'annually'",
		})
	}

	// Create the subscription
	subscription, err := h.subscriptionService.InitiateSubscription(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create subscription: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data":    subscription,
		"message": "Subscription initiated successfully. Please complete payment to activate.",
	})
}

func (h *SubscriptionHandler) HandleConfirmSubscription(c *fiber.Ctx) error {
	// Parse request body
	var req models.PaymentConfirmationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body: " + err.Error(),
		})
	}

	if req.OrderID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Order ID is required",
		})
	}

	// Confirm the subscription
	subscription, err := h.subscriptionService.ConfirmSubscription(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to confirm subscription: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data":    subscription,
		"message": "Subscription confirmed successfully",
	})
}

func (h *SubscriptionHandler) HandleGetTutorSubscription(c *fiber.Ctx) error {
	tutorIdParam := c.Params("tutorId")
	tutorId, err := strconv.ParseUint(tutorIdParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid tutor ID",
		})
	}

	// Get the subscription
	subscription, err := h.subscriptionService.GetTutorSubscription(uint(tutorId))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Subscription not found: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": subscription,
	})
}

func (h *SubscriptionHandler) HandleCancelSubscription(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid subscription ID",
		})
	}

	if err := h.subscriptionService.CancelSubscription(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to cancel subscription: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Subscription cancellation processed successfully",
	})
}

func (h *SubscriptionHandler) HandleChangePlan(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid subscription ID",
		})
	}

	var req models.ChangePlanRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body: " + err.Error(),
		})
	}

	if req.NewPlanID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "New plan ID is required",
		})
	}

	if req.BillingCycle != models.BillingMonthly && req.BillingCycle != models.BillingAnnually {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Billing cycle must be either 'monthly' or 'annually'",
		})
	}

	updatedSubscription, err := h.subscriptionService.ChangePlan(uint(id), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to change subscription plan: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data":    updatedSubscription,
		"message": "Subscription plan change initiated. Please complete payment to activate the new plan.",
	})
}

func (h *SubscriptionHandler) HandleGetAllSubscriptions(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("limit", "10"))

	filters := make(map[string]interface{})

	if status := c.Query("status"); status != "" {
		filters["status"] = models.SubscriptionStatus(status)
	}

	if tutorID := c.Query("tutor_id"); tutorID != "" {
		id, err := strconv.ParseUint(tutorID, 10, 32)
		if err == nil {
			filters["tutor_id"] = uint(id)
		}
	}

	if planID := c.Query("plan_id"); planID != "" {
		id, err := strconv.ParseUint(planID, 10, 32)
		if err == nil {
			filters["plan_id"] = uint(id)
		}
	}

	if fromDate := c.Query("from_date"); fromDate != "" {
		date, err := time.Parse("2006-01-02", fromDate)
		if err == nil {
			filters["from_date"] = date
		}
	}

	if toDate := c.Query("to_date"); toDate != "" {
		date, err := time.Parse("2006-01-02", toDate)
		if err == nil {
			filters["to_date"] = date
		}
	}

	// Get the subscriptions
	subscriptions, total, err := h.subscriptionService.GetAllSubscriptions(page, pageSize, filters)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get subscriptions: " + err.Error(),
		})
	}

	totalPages := (int(total) + pageSize - 1) / pageSize

	return c.JSON(fiber.Map{
		"data": subscriptions,
		"pagination": fiber.Map{
			"total":       total,
			"page":        page,
			"page_size":   pageSize,
			"total_pages": totalPages,
		},
	})
}

func (h *SubscriptionHandler) HandleUpdateSubscriptionStatus(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid subscription ID",
		})
	}

	var req models.SubscriptionStatusUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body: " + err.Error(),
		})
	}

	// Validate status
	validStatuses := map[models.SubscriptionStatus]bool{
		models.SubscriptionActive:     true,
		models.SubscriptionCanceled:   true,
		models.SubscriptionPastDue:    true,
		models.SubscriptionTrialing:   true,
		models.SubscriptionIncomplete: true,
	}

	if !validStatuses[req.Status] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid subscription status",
		})
	}

	if err := h.subscriptionService.UpdateSubscriptionStatus(uint(id), req.Status); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update subscription status: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Subscription status updated successfully",
	})
}

func (h *SubscriptionHandler) HandlePaymentWebhook(c *fiber.Ctx) error {
	var payload models.PaymentWebhookPayload
	if err := json.Unmarshal(c.Body(), &payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid webhook payload",
		})
	}

	if err := h.subscriptionService.ProcessPaymentWebhook(payload); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process payment webhook: " + err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusOK)
}
