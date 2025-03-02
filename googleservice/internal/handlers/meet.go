package handlers

import (
	"fmt"
	"google-service/internal/config"
	"google-service/internal/services"

	"github.com/gofiber/fiber/v2"
)

type MeetHandler struct {
	MeetService  *services.MeetService
	OAuthService *services.GoogleOAuthService
}

func NewMeetHandler(config *config.GoogleOAuthConfig) *MeetHandler {
	oauthService := services.NewGoogleOAuthService(config)
	meetService, err := services.NewMeetService(config)
	if err != nil {
		fmt.Println("Meet Service init failed:", err)
		return nil
	}

	return &MeetHandler{
		MeetService:  meetService,
		OAuthService: oauthService,
	}
}

// For users who already have a valid token stored
func (m *MeetHandler) CreateMeetWithEmail(c *fiber.Ctx) error {
	title := c.Query("title")
	email := c.Query("email")

	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No user email provided",
		})
	}

	token, err := m.OAuthService.GetUserToken(email)
	if err != nil {
		fmt.Println("Retriving token from email failed:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Retriving token from email failed",
		})
	}

	meetLink, err := m.MeetService.CreateMeetLinkWithOAUTH(token, title)
	if err != nil {
		fmt.Println("Failed to create meet link:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to create meeting",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"link": meetLink,
	})
}
