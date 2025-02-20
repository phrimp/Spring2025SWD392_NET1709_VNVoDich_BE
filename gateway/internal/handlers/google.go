package handlers

import (
	"crypto/rand"
	"encoding/base64"
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
		return routes.GoogleLoginRoute(req, resp, c, h.googleServiceURL+"/api/auth/google/login?state="+state)
	}
}

func (h *GoogleHandler) HandleCallback() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)
		return routes.GoogleLoginRoute(req, resp, c, h.googleServiceURL+"/api/auth/google/callback")
	}
}

func generateRandomState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
