package routes

import (
	"encoding/json"
	"fmt"
	"gateway/utils"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func ForwardRequest(req *fasthttp.Request, resp *fasthttp.Response, c *fiber.Ctx, url, method string, body []byte) error {
	cookie := c.Request().Header.Peek("Cookie")
	if len(cookie) > 0 {
		fmt.Println("Detect Cookies\nOriginal Cookies:", string(cookie), "\nCopying")
		req.Header.SetBytesK([]byte("Cookie"), string(cookie))
	}
	utils.BuildRequest(req, method, body, utils.API_KEY, url)

	// Forward request
	if err := fasthttp.Do(req, resp); err != nil {
		fmt.Printf("Error forwarding request: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Service unavailable",
		})
	}

	// Handle redirects
	if resp.StatusCode() >= 300 && resp.StatusCode() < 400 {
		redirectURL := string(resp.Header.Peek("Location"))
		return c.Redirect(redirectURL, resp.StatusCode())
	}

	c.Set("Content-Type", "application/json")

	contentType := string(resp.Header.Peek("Content-Type"))
	if !strings.Contains(contentType, "application/json") {
		return c.Status(resp.StatusCode()).Send(resp.Body())
	}

	// Parse JSON response
	var responseData interface{}
	if err := json.Unmarshal(resp.Body(), &responseData); err != nil {
		log.Printf("Error parsing JSON response: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Invalid response format",
			"details": err.Error(),
		})
	}

	// Process response data
	var data interface{}
	var respErr interface{}

	switch v := responseData.(type) {
	case map[string]interface{}:
		// Handle error cases
		if errorMsg, exists := v["error"]; exists {
			respErr = errorMsg
		} else if msg, exists := v["message"]; exists {
			respErr = msg
		} else {
			data = v
		}
	default:
		data = v
	}

	return c.Status(resp.StatusCode()).JSON(fiber.Map{
		"data":    data,
		"message": respErr,
	})
}

func CustomForwardRequest(req *fasthttp.Request, resp *fasthttp.Response, c *fiber.Ctx, url string, method string, body []byte, dataTransformer func(originalData interface{}) (interface{}, error),
) error {
	cookie := c.Request().Header.Peek("Cookie")
	if len(cookie) > 0 {
		fmt.Println("Detect Cookies\nOriginal Cookies:", string(cookie), "\nCopying")
		req.Header.SetBytesK([]byte("Cookie"), string(cookie))
	}
	utils.BuildRequest(req, method, body, utils.API_KEY, url)

	// Forward request
	if err := fasthttp.Do(req, resp); err != nil {
		fmt.Printf("Error forwarding request: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Service unavailable",
		})
	}

	// Handle redirects
	if resp.StatusCode() >= 300 && resp.StatusCode() < 400 {
		redirectURL := string(resp.Header.Peek("Location"))
		return c.Redirect(redirectURL, resp.StatusCode())
	}

	c.Set("Content-Type", "application/json")

	contentType := string(resp.Header.Peek("Content-Type"))
	if !strings.Contains(contentType, "application/json") {
		return c.Status(resp.StatusCode()).Send(resp.Body())
	}

	// Parse JSON response
	var originalData interface{}
	if err := json.Unmarshal(resp.Body(), &originalData); err != nil {
		fmt.Printf("Error parsing JSON response: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Invalid response format",
			"details": err.Error(),
		})
	}

	// Process response with custom transformer
	// The transformer will return exactly what should be sent to the client
	responseData, err := dataTransformer(originalData)
	if err != nil {
		fmt.Printf("Error transforming data: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Error processing response",
			"details": err.Error(),
		})
	}

	// Return only what the transformer provided
	return c.Status(resp.StatusCode()).JSON(responseData)
}
