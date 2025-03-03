package handlers

import (
	"fmt"
	"gateway/internal/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type PaymentHandler struct {
	paymentServiceURL string
}

func NewPaymentHandler(paymentHanderURL string) *PaymentHandler {
	return &PaymentHandler{
		paymentServiceURL: paymentHanderURL,
	}
}

func (p *PaymentHandler) HandleCreatePayment() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)
		query := fmt.Sprintf("?amount=%s&description=%s&orderId=%s", c.Query("amount"), c.Query("description"), c.Query("orderId"))
		return routes.CreatePayment(req, resp, c, p.paymentServiceURL+"/api/payment/paypal/create"+query)
	}
}

func (h *PaymentHandler) HandleCompletePayPalPayment() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)
		query := fmt.Sprintf("?paymentId=%s&PayerID=%s&orderId=%s", c.Query("paymentId"), c.Query("PayerID"), c.Query("orderId"))
		return routes.CompletePayment(req, resp, c, h.paymentServiceURL+"/api/payment/paypal/success"+query)
	}
}

func (h *PaymentHandler) HandleCancelPayPalPayment() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)
		query := fmt.Sprintf("?orderId=%s", c.Query("orderId"))
		return routes.CompletePayment(req, resp, c, h.paymentServiceURL+"/api/payment/paypal/success"+query)
	}
}
