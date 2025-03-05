package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"gateway/internal/config"
	"gateway/internal/middleware"
	"gateway/internal/routes"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type GoogleHandler struct {
	googleServiceURL string
	userServiceURL   string
}

func NewGoogleHandler(config *config.Config) *GoogleHandler {
	return &GoogleHandler{
		googleServiceURL: config.GoogleServiceURL,
		userServiceURL:   config.UserServiceURL,
	}
}

func (h *GoogleHandler) HandleLogin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Generate state in gateway
		state := generateRandomState()
		state = state[:len(state)-1]

		if len(state) < 32 {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to generate secure state",
			})
		}

		c.Cookie(&fiber.Cookie{
			Name:     "oauth_state",
			Value:    state,
			Expires:  time.Now().Add(time.Minute * 5),
			HTTPOnly: true,
			Secure:   true,
			SameSite: "None",
			Path:     "/",
		})

		fmt.Println("Set cookie:", c.Response().Header.Peek("Set-Cookie"))

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
		fmt.Println("Incoming cookies:", string(c.Request().Header.Peek("Cookie")))

		stateCookie := c.Cookies("oauth_state")
		stateParam := c.Query("state")

		fmt.Println("State from cookie:", stateCookie)
		fmt.Println("State from param:", stateParam)

		if stateCookie == "" || stateParam == "" || stateCookie != stateParam {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":        "Invalid state parameter",
				"cookie_state": stateCookie,
				"param_state":  stateParam,
			})
		}

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

var verification_code map[string]*OTP

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
		code := generateRandomOTP()
		code = code[:len(code)-4]
		expired_time := time.Now().Add(10 * time.Minute).Unix()
		otp := &OTP{code: code, expired_time: expired_time}

		verification_code[current_email] = otp
		fmt.Println("VERIFICATION CODE DEBUG:", verification_code[current_email])

		query_url := fmt.Sprintf("?to=%s&body=%s", current_email, otp.code)
		return routes.SendVerificationEmail(req, resp, c, h.googleServiceURL+"/api/email/send/verify/email"+query_url)
	}
}

func (h *GoogleHandler) HandleVerifyEmail() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)
		code := c.Query("code")
		claims, ok := c.Locals("user").(*middleware.Claims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "cannot find username in token claim"})
		}
		current_email := claims.Email
		if verification_code[current_email].code != code {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "wrong otp code"})
		}
		delete(verification_code, current_email)

		query := fmt.Sprintf("?username=%s", claims.Username)
		return routes.VerifyEmail(req, resp, c, h.userServiceURL+"/user/verify"+query)
	}
}

func (h *GoogleHandler) HandleCreateMeetLink() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)
		query_url := fmt.Sprintf("?title=%s", c.Query("title"))
		return routes.CreateMeetLink(req, resp, c, h.googleServiceURL+"/api/meet/create"+query_url)
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
