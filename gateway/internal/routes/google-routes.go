package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func GoogleLoginRoute(req *fasthttp.Request, resp *fasthttp.Response, c *fiber.Ctx, url string) error {
	return GoogleForwardResquest(req, resp, c, url, "GET", c.Body())
}
