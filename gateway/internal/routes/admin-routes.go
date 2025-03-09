package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func GetAllUser(req *fasthttp.Request, resp *fasthttp.Response, c *fiber.Ctx, url string) error {
	return ForwardRequest(req, resp, c, url, "GET", c.Body())
}

func AdminUpdateUser(req *fasthttp.Request, resp *fasthttp.Response, c *fiber.Ctx, url string) error {
	return ForwardRequest(req, resp, c, url, "PUT", c.Body())
}

func AdminGetUserDetail(req *fasthttp.Request, resp *fasthttp.Response, c *fiber.Ctx, url string) error {
	return ForwardRequest(req, resp, c, url, "GET", c.Body())
}

func UpdateUserStatus(req *fasthttp.Request, resp *fasthttp.Response, c *fiber.Ctx, url string) error {
	return ForwardRequest(req, resp, c, url, "PATCH", c.Body())
}

func AdminDeleteUser(req *fasthttp.Request, resp *fasthttp.Response, c *fiber.Ctx, url string) error {
	return ForwardRequest(req, resp, c, url, "DELETE", c.Body())
}

func AdminAssignRole(req *fasthttp.Request, resp *fasthttp.Response, c *fiber.Ctx, url string, role string) error {
	// Create request body with the role
	body := []byte(fmt.Sprintf(`{"role":"%s"}`, role))
	return ForwardRequest(req, resp, c, url, "POST", body)
}

func AdminBlockJWT(req *fasthttp.Request, resp *fasthttp.Response, c *fiber.Ctx, url string) error {
	return ForwardRequest(req, resp, c, url, "POST", c.Body())
}
