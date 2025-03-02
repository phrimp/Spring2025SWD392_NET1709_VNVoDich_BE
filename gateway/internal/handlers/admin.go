package handlers

import (
	"fmt"
	"gateway/internal/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type AdminServiceHandler struct {
	adminServiceURL string
}

func NewAdminService(adminServiceURL string) *AdminServiceHandler {
	return &AdminServiceHandler{
		adminServiceURL: adminServiceURL,
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
