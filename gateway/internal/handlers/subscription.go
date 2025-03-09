package handlers

import (
	"fmt"
	"gateway/internal/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type SubscriptionHandler struct {
	subscriptionServiceURL string
}

func NewSubscriptionHandler(subscriptionServiceURL string) *SubscriptionHandler {
	return &SubscriptionHandler{
		subscriptionServiceURL: subscriptionServiceURL,
	}
}

func (h *SubscriptionHandler) HandleGetPlans() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)

		query := c.Request().URI().QueryString()
		url := fmt.Sprintf("%s/api/plans", h.subscriptionServiceURL)
		if len(query) > 0 {
			url = fmt.Sprintf("%s?%s", url, string(query))
		}

		return routes.ForwardRequest(req, resp, c, url, "GET", nil)
	}
}

func (h *SubscriptionHandler) HandleGetPlan() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)

		id := c.Params("id")
		url := fmt.Sprintf("%s/api/plans/%s", h.subscriptionServiceURL, id)

		return routes.ForwardRequest(req, resp, c, url, "GET", nil)
	}
}

func (h *SubscriptionHandler) HandleCreateSubscription() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)

		url := fmt.Sprintf("%s/api/subscriptions", h.subscriptionServiceURL)

		return routes.ForwardRequest(req, resp, c, url, "POST", c.Body())
	}
}

func (h *SubscriptionHandler) HandleConfirmSubscription() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)

		url := fmt.Sprintf("%s/api/subscriptions/confirm", h.subscriptionServiceURL)

		return routes.ForwardRequest(req, resp, c, url, "POST", c.Body())
	}
}

func (h *SubscriptionHandler) HandleGetTutorSubscription() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)

		tutorId := c.Params("tutorId")
		url := fmt.Sprintf("%s/api/subscriptions/tutor/%s", h.subscriptionServiceURL, tutorId)

		return routes.ForwardRequest(req, resp, c, url, "GET", nil)
	}
}

func (h *SubscriptionHandler) HandleCancelSubscription() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)

		id := c.Params("id")
		url := fmt.Sprintf("%s/api/subscriptions/%s/cancel", h.subscriptionServiceURL, id)

		return routes.ForwardRequest(req, resp, c, url, "PUT", nil)
	}
}

func (h *SubscriptionHandler) HandleChangePlan() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)

		id := c.Params("id")
		url := fmt.Sprintf("%s/api/subscriptions/%s/change-plan", h.subscriptionServiceURL, id)

		return routes.ForwardRequest(req, resp, c, url, "PUT", c.Body())
	}
}

// Admin handlers

func (h *SubscriptionHandler) HandleGetAllSubscriptions() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)

		query := c.Request().URI().QueryString()
		url := fmt.Sprintf("%s/api/admin/subscriptions", h.subscriptionServiceURL)
		if len(query) > 0 {
			url = fmt.Sprintf("%s?%s", url, string(query))
		}

		return routes.ForwardRequest(req, resp, c, url, "GET", nil)
	}
}

func (h *SubscriptionHandler) HandleUpdateSubscriptionStatus() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)

		id := c.Params("id")
		url := fmt.Sprintf("%s/api/admin/subscriptions/%s/status", h.subscriptionServiceURL, id)

		return routes.ForwardRequest(req, resp, c, url, "PUT", c.Body())
	}
}

func (h *SubscriptionHandler) HandleAdminCreatePlan() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)

		url := fmt.Sprintf("%s/api/admin/plans", h.subscriptionServiceURL)

		return routes.ForwardRequest(req, resp, c, url, "POST", c.Body())
	}
}

func (h *SubscriptionHandler) HandleAdminUpdatePlan() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)

		id := c.Params("id")
		url := fmt.Sprintf("%s/api/admin/plans/%s", h.subscriptionServiceURL, id)

		return routes.ForwardRequest(req, resp, c, url, "PUT", c.Body())
	}
}

func (h *SubscriptionHandler) HandleAdminDeletePlan() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)

		id := c.Params("id")
		url := fmt.Sprintf("%s/api/admin/plans/%s", h.subscriptionServiceURL, id)

		return routes.ForwardRequest(req, resp, c, url, "DELETE", nil)
	}
}

// Payment webhook handler
func (h *SubscriptionHandler) HandlePaymentWebhook() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)

		url := fmt.Sprintf("%s/api/webhooks/payment", h.subscriptionServiceURL)

		return routes.ForwardRequest(req, resp, c, url, "POST", c.Body())
	}
}
