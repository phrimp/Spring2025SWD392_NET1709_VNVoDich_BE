package handlers

import (
	"fmt"
	"gateway/internal/routes"
	"net/url"

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

		params := url.Values{}
		params.Add("amount", c.Query("amount"))
		params.Add("description", c.Query("description"))
		params.Add("orderId", c.Query("orderId"))

		endpoint := fmt.Sprintf("%s/api/payment/paypal/create?%s", p.paymentServiceURL, params.Encode())
		fmt.Println(endpoint)

		return routes.CreatePayment(req, resp, c, endpoint)
	}
}

func (h *PaymentHandler) HandleCompletePayPalPayment() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)

		params := url.Values{}
		params.Add("paymentId", c.Query("paymentId"))
		params.Add("PayerID", c.Query("PayerID"))
		params.Add("orderId", c.Query("orderId"))

		endpoint := fmt.Sprintf("%s/api/payment/paypal/success?%s", h.paymentServiceURL, params.Encode())
		return routes.CompletePayment(req, resp, c, endpoint)
	}
}

func (h *PaymentHandler) HandleCancelPayPalPayment() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)

		params := url.Values{}
		params.Add("orderId", c.Query("orderId"))

		endpoint := fmt.Sprintf("%s/api/payment/paypal/cancel?%s", h.paymentServiceURL, params.Encode())
		return routes.CancelPayment(req, resp, c, endpoint)
	}
}
