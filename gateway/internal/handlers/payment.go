package handlers

import (
	"encoding/json"
	"fmt"
	"gateway/internal/routes"
	"net/url"
	"os"
	"time"

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
		err := routes.CompletePayment(req, resp, c, endpoint)
		if err != nil {
			fmt.Println("PAYMENT COMPLETEMENT FAILED:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "payment is success but failed to complete payment in final steps",
			})
		}

		var responseData map[string]interface{}
		if err := json.Unmarshal(resp.Body(), &responseData); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to parse payment complete response",
			})
		}

		c.Cookie(&fiber.Cookie{
			Name:     "status",
			Value:    responseData["status"].(string),
			Path:     "/",
			Expires:  time.Now().Add(24 * time.Hour),
			HTTPOnly: false,
			Secure:   true,
			SameSite: "Lax",
		})

		c.Cookie(&fiber.Cookie{
			Name:     "paymentID",
			Value:    responseData["paymentId"].(string),
			Path:     "/",
			Expires:  time.Now().Add(24 * time.Hour),
			HTTPOnly: false,
			Secure:   true,
			SameSite: "Lax",
		})

		return c.Redirect(os.Getenv("REDIRECT_URL"), fiber.StatusTemporaryRedirect)
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
