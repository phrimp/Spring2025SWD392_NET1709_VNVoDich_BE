package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"google-service/internal/config"
	"google-service/internal/services"

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
	state := c.Query("state")
	if state == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "State parameter required",
		})
	}

	url := h.oauthService.GetAuthURL(state)
	return c.Redirect(url)
}

func (h *GoogleHandler) HandleGoogleCallback(c *fiber.Ctx) error {
	state := c.Cookies("oauth_state")
	if state != c.Query("state") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid state",
		})
	}

	code := c.Query("code")
	token, err := h.oauthService.Exchange(c.Context(), code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to exchange token",
		})
	}

	userInfo, err := h.oauthService.GetUserInfo(token)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get user info",
		})
	}
	fmt.Println(userInfo)

	//// Initialize email service with the token
	//emailService, err := services.NewEmailService(c.Context(), token, h.config)
	//if err != nil {
	//	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	//		"error": "Failed to initialize email service",
	//	})
	//}

	// h.emailService = emailService

	return c.JSON(fiber.Map{
		"token": token.AccessToken,
		"user":  userInfo,
	})
}

func generateRandomState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
