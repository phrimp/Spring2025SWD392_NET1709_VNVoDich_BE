package handlers

import (
	"fmt"
	"gateway/internal/config"
	"gateway/internal/routes"

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
		query_url := fmt.Sprintf("?page=%s&limit=%s", c.Query("page"), c.Query("limit"))
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

		username := c.Params("username")

		if username == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Username is required",
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

		// Get role from request body
		var requestBody struct {
			Role string `json:"role"`
		}

		if err := c.BodyParser(&requestBody); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		if requestBody.Role == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Role is required",
			})
		}

		// Forward the request to the user service
		return routes.AdminAssignRole(req, resp, c, a.userServiceURL+"/user/admin/role?username="+username, requestBody.Role)
	}
}
