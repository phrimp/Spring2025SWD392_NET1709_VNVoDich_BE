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

		return routes.AdminUpdateUser(req, resp, c, a.userServiceURL+"/admin/users/"+query)
	}
}
