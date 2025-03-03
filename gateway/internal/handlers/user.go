package handlers

import (
	"fmt"
	"gateway/internal/middleware"
	"gateway/internal/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type UserServiceHandler struct {
	userServiceURL string
}

func NewUserService(userServiceURL string) *UserServiceHandler {
	return &UserServiceHandler{
		userServiceURL: userServiceURL,
	}
}

func (h *UserServiceHandler) HandleGetUserwithUsername() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)
		return routes.GetUserwithUsername(req, resp, c, h.userServiceURL+"/user/get-public-user")
	}
}

func (h *UserServiceHandler) HandleAllGetUser() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)
		query_url := fmt.Sprintf("?page=%s&limit=%s", c.Query("page"), c.Query("limit"))
		return routes.GetAllUser(req, resp, c, h.userServiceURL+"/user/get-all-user"+query_url)
	}
}

func (h *UserServiceHandler) HandleGetMe() fiber.Handler {
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
		query_url := fmt.Sprintf("?username=%s", current_username)
		return routes.GetMe(req, resp, c, h.userServiceURL+"/user"+query_url)
	}
}

func (h *UserServiceHandler) HandleDeleteMe() fiber.Handler {
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
		query_url := fmt.Sprintf("?username=%s", current_username)
		return routes.DeleteMe(req, resp, c, h.userServiceURL+"/delete"+query_url)
	}
}

func (h *UserServiceHandler) HandleCancelDeleteMe() fiber.Handler {
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
		query_url := fmt.Sprintf("?username=%s", current_username)
		return routes.CancelDeleteMe(req, resp, c, h.userServiceURL+"/delete/cancel"+query_url)
	}
}

func (h *UserServiceHandler) HandleUpdateMe() fiber.Handler {
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
		query_url := fmt.Sprintf("?username=%s", current_username)
		return routes.UpdateMe(req, resp, c, h.userServiceURL+"/user/update"+query_url)
	}
}
