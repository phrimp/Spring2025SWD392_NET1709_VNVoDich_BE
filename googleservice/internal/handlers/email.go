package handlers

import (
	"google-service/internal/config"
	"google-service/internal/services"

	"github.com/gofiber/fiber/v2"
)

type EmailHandler struct {
	emailService *services.EmailService
}

func NewEmailHandler(config *config.EmailConfig) *EmailHandler {
	return &EmailHandler{
		emailService: services.NewEmailService(config),
	}
}

func (e *EmailHandler) HandleSendPlainEmail(c *fiber.Ctx) error {
	to := c.Query("to")
	title := c.Query("title")
	body := c.Query("body")
	err := e.emailService.SendEmail(title, body, []string{to})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Email Sent"})
}

func (e *EmailHandler) HandleVerifyEmail(c *fiber.Ctx) error {
	title := "Email Verification"
	to := c.Query("to")
	body := c.Query("body")
	err := e.emailService.SendEmail(title, body, []string{to})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Email Sent"})
}
