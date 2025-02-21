package handlers

import (
	"fmt"
	"google-service/internal/config"
	"google-service/internal/services"

	"github.com/gofiber/fiber/v2"
)

type MeetHandler struct {
	MeetService *services.MeetService
}

func NewMeetHandler(config *config.ServiceAccountConfig) *MeetHandler {
	tmpMeetService, err := services.NewMeetService(*config)
	if err != nil {
		fmt.Println("Meet Service init failed:", err)
		return nil
	}

	return &MeetHandler{
		MeetService: tmpMeetService,
	}
}

func (m *MeetHandler) GenerateMeetLink(c *fiber.Ctx) error {
	title := c.Query("title")
	meetLink, err := m.MeetService.CreateMeetLink(title)
	if err != nil {
		fmt.Println("Failed to create meet link:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"link": meetLink})
}
