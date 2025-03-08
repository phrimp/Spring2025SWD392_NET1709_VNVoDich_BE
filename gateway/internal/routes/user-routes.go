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

func GetMe(req *fasthttp.Request, resp *fasthttp.Response, c *fiber.Ctx, url string, dataTransformer func(originalData interface{}) (interface{}, error)) error {
	return CustomForwardRequest(req, resp, c, url, "GET", c.Body(), dataTransformer)
}

func DeleteMe(req *fasthttp.Request, resp *fasthttp.Response, c *fiber.Ctx, url string) error {
	return ForwardRequest(req, resp, c, url, "DELETE", c.Body())
}

func CancelDeleteMe(req *fasthttp.Request, resp *fasthttp.Response, c *fiber.Ctx, url string) error {
	return ForwardRequest(req, resp, c, url, "POST", c.Body())
}

func UpdateMe(req *fasthttp.Request, resp *fasthttp.Response, c *fiber.Ctx, url string) error {
	return ForwardRequest(req, resp, c, url, "PUT", c.Body())
}

func UpdateMePassword(req *fasthttp.Request, resp *fasthttp.Response, c *fiber.Ctx, url string) error {
	return ForwardRequest(req, resp, c, url, "PATCH", c.Body())
}
