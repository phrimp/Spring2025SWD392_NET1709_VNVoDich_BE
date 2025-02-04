package routes

import (
	"fmt"
	"gateway/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func GoogleLoginRoute(req *fasthttp.Request, resp *fasthttp.Response, c *fiber.Ctx, url string) error {
	// Copy request body
	utils.BuildRequest(req, "POST", c.Body(), utils.API_KEY, url)

	// Forward request
	if err := fasthttp.Do(req, resp); err != nil {
		fmt.Printf("Error forwarding request: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Service unavailable",
		})
	}

	// Return response
	c.Set("Content-Type", "application/json")
	return c.Status(resp.StatusCode()).Send(resp.Body())
}
