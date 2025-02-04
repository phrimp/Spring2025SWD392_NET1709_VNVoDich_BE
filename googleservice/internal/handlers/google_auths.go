package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"google-service/internal/config"
	"google-service/internal/services"
	"time"

	"github.com/gofiber/fiber/v2"
)

type GoogleHandler struct {
	oauthService *services.GoogleOAuthService
}

func NewGoogleHandler(config *config.GoogleOAuthConfig) *GoogleHandler {
	return &GoogleHandler{
		oauthService: services.NewGoogleOAuthService(config),
	}
}

func (h *GoogleHandler) HandleGoogleLogin(c *fiber.Ctx) error {
	state := generateRandomState()
	c.Cookie(&fiber.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Expires:  time.Now().Add(time.Minute * 5),
		HTTPOnly: true,
		Secure:   true,
	})

	url := h.oauthService.GetAuthURL(state)
	return c.Redirect(url)
}

func generateRandomState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
