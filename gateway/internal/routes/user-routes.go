package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func GetUserwithUsername(req *fasthttp.Request, resp *fasthttp.Response, c *fiber.Ctx, url string) error {
	username := c.Params("username")
	body := fmt.Sprintf(`{"username":"%s"}`, username)
	return ForwardRequest(req, resp, c, url, "GET", []byte(body))
}

func GetMe(req *fasthttp.Request, resp *fasthttp.Response, c *fiber.Ctx, url string) error {
	return ForwardRequest(req, resp, c, url, "GET", c.Body())
}
