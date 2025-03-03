package handlers

import (
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
		return routes.CreatePayment(req, resp, c, p.paymentServiceURL+"/api/payment/paypal/create")
	}
}

func (h *PaymentHandler) HandleCompletePayPalPayment() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)
		return routes.CompletePayment(req, resp, c, h.paymentServiceURL+"/api/payment/paypal/success")
	}
}
