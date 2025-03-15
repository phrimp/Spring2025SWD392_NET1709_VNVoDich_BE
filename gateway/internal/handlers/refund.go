package handlers

import (
	"fmt"
	"gateway/internal/middleware"
	"gateway/internal/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type RefundHandler struct {
	adminServiceURL string
}

func NewRefundHandler(adminServiceURL string) *RefundHandler {
	return &RefundHandler{
		adminServiceURL: adminServiceURL,
	}
}

func (h *RefundHandler) HandleCreateRefundRequest() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)

		claims, ok := c.Locals("user").(*middleware.Claims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "cannot find username in token claim"})
		}
		current_email := claims.Email
		current_id := claims.Id
		current_username := claims.Username
		url := fmt.Sprintf("%s/api/refunds?email=%s&id=%s&username=%s", h.adminServiceURL, current_email, current_id, current_username)
		return routes.ForwardRequest(req, resp, c, url, "POST", c.Body())
	}
}

func (h *RefundHandler) HandleGetRefundRequest() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)

		id := c.Params("id")
		url := fmt.Sprintf("%s/api/refunds/%s", h.adminServiceURL, id)
		return routes.ForwardRequest(req, resp, c, url, "GET", nil)
	}
}

func (h *RefundHandler) HandleGetAllRefundRequests() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)

		// Forward query parameters
		query := string(c.Request().URI().QueryString())
		var queryStr string
		if query != "" {
			queryStr = "?" + query
		}

		url := fmt.Sprintf("%s/api/admin/refunds%s", h.adminServiceURL, queryStr)
		return routes.ForwardRequest(req, resp, c, url, "GET", nil)
	}
}

func (h *RefundHandler) HandleGetRefundStatistics() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)

		url := fmt.Sprintf("%s/api/admin/refunds/statistics", h.adminServiceURL)
		return routes.ForwardRequest(req, resp, c, url, "GET", nil)
	}
}

func (h *RefundHandler) HandleProcessRefundRequest() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)

		claims, ok := c.Locals("user").(*middleware.Claims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "cannot find username in token claim"})
		}
		current_role := claims.Role
		current_id := claims.Id
		id := c.Params("id")
		url := fmt.Sprintf("%s/api/admin/refunds/%s/process?adminid=%s&role=%s", h.adminServiceURL, id, current_id, current_role)
		return routes.ForwardRequest(req, resp, c, url, "PUT", c.Body())
	}
}
