package routes

import (
	"encoding/json"
	"fmt"
	"gateway/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func ForwardRequest(req *fasthttp.Request, resp *fasthttp.Response, c *fiber.Ctx, url, method string, body []byte) error {
	// Copy request body
	utils.BuildRequest(req, method, body, utils.API_KEY, url)

	// Forward request
	if err := fasthttp.Do(req, resp); err != nil {
		fmt.Printf("Error forwarding request: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Service unavailable",
		})
	}

	c.Set("Content-Type", "application/json")
	var data interface{}
	var respErr interface{}

	err := json.Unmarshal(resp.Body(), &data)
	if err != nil {
		fmt.Println("Error Converting Data to Response:", err)
		data = nil
		respErr = err
		return c.Status(resp.StatusCode()).JSON(fiber.Map{
			"data":  data,
			"error": respErr,
		})
	}

	// Check if data is a map
	if mapData, ok := data.(map[string]interface{}); ok {
		if mapData["error"] != nil {
			respErr = mapData["error"]
			data = nil
		}
	}
	return c.Status(resp.StatusCode()).JSON(fiber.Map{"data": data, "error": respErr})
}
