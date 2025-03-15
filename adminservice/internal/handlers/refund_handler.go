package handlers

import (
	"adminservice/internal/models"
	"adminservice/internal/services"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type RefundHandler struct {
	refundService services.RefundService
}

func NewRefundHandler(refundService services.RefundService) *RefundHandler {
	return &RefundHandler{
		refundService: refundService,
	}
}

// HandleCreateRefundRequest handles the creation of a refund request
func (h *RefundHandler) HandleCreateRefundRequest() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Parse request body
		userid_str := c.Query("userid")
		userid, _ := strconv.Atoi(userid_str)

		username := c.Query("username")
		email := c.Query("email")
		var input models.RefundRequestInput
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body: " + err.Error(),
			})
		}

		// Basic validation
		if input.OrderID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Order ID is required",
			})
		}

		if input.Amount <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Amount must be greater than 0",
			})
		}

		if input.CardNumber == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Card number is required",
			})
		}

		// Create the refund request
		refund, err := h.refundService.CreateRefundRequest(
			uint(userid),
			username, email,
			input,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create refund request: " + err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "Refund request created successfully",
			"data":    refund,
		})
	}
}

// HandleGetRefundRequest handles getting a specific refund request
func (h *RefundHandler) HandleGetRefundRequest() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.ParseUint(c.Params("id"), 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid refund request ID",
			})
		}

		refund, err := h.refundService.GetRefundRequestByID(uint(id))
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Refund request not found: " + err.Error(),
			})
		}

		return c.JSON(refund)
	}
}

// HandleGetAllRefundRequests handles getting all refund requests with filtering and pagination
func (h *RefundHandler) HandleGetAllRefundRequests() fiber.Handler {
	return func(c *fiber.Ctx) error {
		page, _ := strconv.Atoi(c.Query("page", "1"))
		limit, _ := strconv.Atoi(c.Query("limit", "10"))

		filters := make(map[string]interface{})

		if status := c.Query("status"); status != "" {
			filters["status"] = models.RefundStatus(status)
		}

		if userID := c.Query("user_id"); userID != "" {
			if id, err := strconv.ParseUint(userID, 10, 32); err == nil {
				filters["user_id"] = uint(id)
			}
		}

		if orderID := c.Query("order_id"); orderID != "" {
			filters["order_id"] = orderID
		}

		if username := c.Query("username"); username != "" {
			filters["username"] = username
		}

		if email := c.Query("email"); email != "" {
			filters["email"] = email
		}

		if fromDate := c.Query("from_date"); fromDate != "" {
			if date, err := time.Parse("2006-01-02", fromDate); err == nil {
				filters["from_date"] = date
			}
		}

		if toDate := c.Query("to_date"); toDate != "" {
			if date, err := time.Parse("2006-01-02", toDate); err == nil {
				// Add a day to include the entire end date
				filters["to_date"] = date.Add(24 * time.Hour)
			}
		}

		response, err := h.refundService.GetAllRefundRequests(page, limit, filters)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to get refund requests: " + err.Error(),
			})
		}

		return c.JSON(response)
	}
}

// HandleProcessRefundRequest handles approving or rejecting a refund request
func (h *RefundHandler) HandleProcessRefundRequest() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Query("role")
		// Verify the user is an admin
		if role != "Admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Only administrators can process refund requests",
			})
		}
		adminid_str := c.Query("adminid")
		adminid, _ := strconv.Atoi(adminid_str)

		// Get refund ID from path
		id, err := strconv.ParseUint(c.Params("id"), 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid refund request ID",
			})
		}

		// Parse request body
		var input models.RefundProcessInput
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body: " + err.Error(),
			})
		}

		// Validate status
		if input.Status != models.RefundStatusApproved && input.Status != models.RefundStatusRejected {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Status must be either 'approved' or 'rejected'",
			})
		}

		// Process the refund request
		if err := h.refundService.ProcessRefundRequest(uint(id), uint(adminid), input); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to process refund request: " + err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"message": fmt.Sprintf("Refund request %s successfully", input.Status),
		})
	}
}

// HandleGetRefundStatistics handles getting statistics about refund requests
func (h *RefundHandler) HandleGetRefundStatistics() fiber.Handler {
	return func(c *fiber.Ctx) error {
		stats, err := h.refundService.GetRefundStatistics()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to get refund statistics: " + err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"data": stats,
		})
	}
}
