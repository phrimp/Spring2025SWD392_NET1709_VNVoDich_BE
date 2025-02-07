package handlers

import (
	"gateway/internal/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type AuthHandler struct {
	authServiceURL string
}

func NewAuthHandler(authServiceURL string) *AuthHandler {
	return &AuthHandler{
		authServiceURL: authServiceURL,
	}
}

func (h *AuthHandler) HandleLogin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)
		return routes.LoginRoute(req, resp, c, h.authServiceURL+"/login")
	}
}

func (h *AuthHandler) HandleRegister() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)
		return routes.RegisterRoute(req, resp, c, h.authServiceURL+"/register")
	}
}
