package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func GetAllCourse(req *fasthttp.Request, resp *fasthttp.Response, c *fiber.Ctx, url string) error {
	return ForwardRequest(req, resp, c, url, "GET", c.Body())
}

func GetACourse(req *fasthttp.Request, resp *fasthttp.Response, c *fiber.Ctx, url string) error {
	return ForwardRequest(req, resp, c, url, "GET", c.Body())
}

func CreateCourse(req *fasthttp.Request, resp *fasthttp.Response, c *fiber.Ctx, url string) error {
	return ForwardRequest(req, resp, c, url, "POST", c.Body())
}

func AddLessonToCourse(req *fasthttp.Request, resp *fasthttp.Response, c *fiber.Ctx, url string) error {
	return ForwardRequest(req, resp, c, url, "PUT", c.Body())
}
