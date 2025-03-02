package handlers

import (
	"adminservice/internal/config"
	"adminservice/internal/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AdminHandler struct {
	adminService *services.AdminService
}

func NewAdminHandler(cfg *config.Config) *AdminHandler {
	return &AdminHandler{
		adminService: services.NewAdminService(cfg),
	}
}

func TestHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return nil
	}
}

func (a AdminHandler) GetAllUsersHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		page, _ := strconv.Atoi(c.Query("page", "1"))
		limit, _ := strconv.Atoi(c.Query("limit", "10"))

		role := c.Query("role")
		status := c.Query("status")
		search := c.Query("search")
		sort := c.Query("sort", "created_at")
		sortDir := c.Query("sort_dir", "DESC")
		dateFrom := c.Query("created_from")
		dateTo := c.Query("created_to")

		var isVerified *bool
		if verified := c.Query("is_verified"); verified != "" {
			verifiedBool, err := strconv.ParseBool(verified)
			if err == nil {
				isVerified = &verifiedBool
			}
		}

		response, err := a.adminService.GetAllUsers(page, limit, role, status, search,
			isVerified, dateFrom, dateTo, sort, sortDir)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve users: " + err.Error(),
			})
		}

		return c.JSON(response)
	}
}
