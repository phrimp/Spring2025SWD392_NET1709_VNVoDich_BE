package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
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
		err := routes.GoogleLoginRoute(req, resp, c, h.googleServiceURL+"/api/auth/google/callback?"+query)
		if err != nil {
			fmt.Println("GOOGLE SERVICE CALL BACK FAILED:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to get user information",
			})
		}

		var responseData map[string]interface{}
		if err := json.Unmarshal(resp.Body(), &responseData); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to parse authentication response",
			})
		}

		token, hasToken := responseData["token"].(string)
		userData, hasUserData := responseData["user"].(map[string]interface{})

		if !hasToken || !hasUserData {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Invalid authentication response format",
			})
		}

		c.Cookie(&fiber.Cookie{
			Name:     "authToken",
			Value:    token,
			Path:     "/",
			Expires:  time.Now().Add(24 * time.Hour), // 24 hour expiration
			HTTPOnly: true,
			Secure:   true,
			SameSite: "Lax",
		})

		userDataJSON, _ := json.Marshal(userData)
		c.Cookie(&fiber.Cookie{
			Name:     "user_info",
			Value:    string(userDataJSON),
			Path:     "/",
			Expires:  time.Now().Add(24 * time.Hour),
			HTTPOnly: false, // Not HTTP-only so JavaScript can read it
			Secure:   true,
			SameSite: "Lax",
		})

		return c.Redirect("http://localhost:3000", fiber.StatusTemporaryRedirect)
	}
}

type OTP struct {
	code         string
	expired_time int64
}

var verification_code map[string]*OTP = make(map[string]*OTP)

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
		claims, ok := c.Locals("user").(*middleware.Claims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "cannot find user in token claim"})
		}

		query_url := fmt.Sprintf("?title=%s&email=%s", c.Query("title"), claims.Email)
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
