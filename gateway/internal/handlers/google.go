package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"gateway/internal/middleware"
	"gateway/internal/routes"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type GoogleHandler struct {
	googleServiceURL string
}

func NewGoogleHandler(googleServiceURL string) *GoogleHandler {
	return &GoogleHandler{
		googleServiceURL: googleServiceURL,
	}
}

func (h *GoogleHandler) HandleLogin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Generate state in gateway
		state := generateRandomState()
		state = state[:len(state)-1]
		fmt.Println("state", state)

		c.Cookie(&fiber.Cookie{
			Name:     "oauth_state",
			Value:    state,
			Expires:  time.Now().Add(time.Minute * 5),
			HTTPOnly: true,
			Secure:   true,
			SameSite: "Lax",
			Path:     "/",
		})

		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)

		query_url := fmt.Sprintf("?state=%s", state)
		return routes.GoogleLoginRoute(req, resp, c, h.googleServiceURL+"/api/auth/google/login"+query_url)
	}
}

func (h *GoogleHandler) HandleCallback() fiber.Handler {
	return func(c *fiber.Ctx) error {
		query := string(c.Request().URI().QueryString())
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)
		return routes.GoogleLoginRoute(req, resp, c, h.googleServiceURL+"/api/auth/google/callback?"+query)
	}
}

type OTP struct {
	code         string
	expired_time int64
}

var verification_code map[string]OTP

func (h *GoogleHandler) HandleSendVerificationEmail() fiber.Handler {
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
		otp := OTP{}
		if cur_OTP, ok := verification_code[current_email]; !ok {
			otp = OTP{code: generateRandomOTP(), expired_time: time.Now().Unix()}
		} else {
			if cur_OTP.expired_time < time.Now().Unix() {
				otp = OTP{code: generateRandomOTP(), expired_time: time.Now().Unix()}
			} else {
				otp = cur_OTP
			}
		}

		query_url := fmt.Sprintf("?to=%s&body=%s", current_email, otp.code)
		return routes.SendVerificationEmail(req, resp, c, h.googleServiceURL+"/api/email/send/verify/email"+query_url)
	}
}

func generateRandomState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func generateRandomOTP() string {
	b := make([]byte, 8)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
