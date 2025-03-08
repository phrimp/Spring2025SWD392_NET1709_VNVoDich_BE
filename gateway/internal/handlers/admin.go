package handlers

import (
	"fmt"
	"gateway/internal/config"
	"gateway/internal/middleware"
	"gateway/internal/routes"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type AdminServiceHandler struct {
	adminServiceURL string
	userServiceURL  string
}

func NewAdminService(config *config.Config) *AdminServiceHandler {
	return &AdminServiceHandler{
		adminServiceURL: config.AdminServiceURL,
		userServiceURL:  config.UserServiceURL,
	}
}

func (a *AdminServiceHandler) HandleAllGetUser() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)

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

		query_url := fmt.Sprintf("?page=%d&limit=%d", page, limit)

		if role != "" {
			query_url += fmt.Sprintf("&role=%s", role)
		}
		if status != "" {
			query_url += fmt.Sprintf("&status=%s", status)
		}
		if search != "" {
			query_url += fmt.Sprintf("&search=%s", search)
		}
		query_url += fmt.Sprintf("&sort=%s&sort_dir=%s", sort, sortDir)
		if dateFrom != "" {
			query_url += fmt.Sprintf("&created_from=%s", dateFrom)
		}
		if dateTo != "" {
			query_url += fmt.Sprintf("&created_to=%s", dateTo)
		}
		if isVerified != nil {
			query_url += fmt.Sprintf("&is_verified=%t", *isVerified)
		}
		return routes.GetAllUser(req, resp, c, a.adminServiceURL+"/api/users"+query_url)
	}
}

func (a *AdminServiceHandler) HandleAdminUpdateUser() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)
		query := fmt.Sprintf("?username=%s", c.Query("username"))

		return routes.AdminUpdateUser(req, resp, c, a.userServiceURL+"/user/admin/update"+query)
	}
}

func (a *AdminServiceHandler) HandleAdminGetUSerDetail() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)
		query := fmt.Sprintf("?username=%s", c.Query("username"))

		return routes.AdminGetUserDetail(req, resp, c, a.userServiceURL+"/user"+query)
	}
}

func (a *AdminServiceHandler) HandleUpdateUserStatus() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)

		username := c.Params("username")
		status := c.Query("status")

		if username == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Username is required",
			})
		}

		if status == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Status is required",
			})
		}

		query := fmt.Sprintf("?username=%s&status=%s", username, status)
		return routes.UpdateUserStatus(req, resp, c, a.userServiceURL+"/user/update/status"+query)
	}
}

func (a *AdminServiceHandler) HandleDeleteUser() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)

		claims, ok := c.Locals("user").(*middleware.Claims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "cannot find username in token claim"})
		}
		current_username := claims.Username

		username := c.Params("username")

		if username == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Username is required",
			})
		}

		if username == current_username {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Can not delete yourself using this route",
			})
		}

		query := fmt.Sprintf("?username=%s", username)
		return routes.AdminDeleteUser(req, resp, c, a.userServiceURL+"/user/admin/delete"+query)
	}
}

func (a *AdminServiceHandler) HandleAssignRole() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)

		username := c.Params("username")

		if username == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Username is required",
			})
		}

		role := c.Query("role")

		if role == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Role is required",
			})
		}

		// Forward the request to the user service
		return routes.AdminAssignRole(req, resp, c, a.userServiceURL+"/user/admin/role?username="+username, role)
	}
}
