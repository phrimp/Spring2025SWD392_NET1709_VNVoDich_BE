package handlers

import (
	"fmt"
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
