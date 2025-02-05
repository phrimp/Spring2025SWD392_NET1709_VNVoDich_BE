package handlers

import (
	"user-service/internal/services"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type RequestParam struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func GetUserWithUsernamePasswordHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req RequestParam
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request",
			})
		}
		user, err := services.FindUserWithUsernamePassword(req.Username, req.Password, db)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err,
			})
		}
		return c.JSON(user)
	}
}
