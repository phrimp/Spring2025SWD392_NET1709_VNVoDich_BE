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
	// Copy request body
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

	// Handle non-JSON responses (like during OAuth redirects)
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

func GoogleForwardResquest(req *fasthttp.Request, resp *fasthttp.Response, c *fiber.Ctx, url, method string, body []byte) error {
	// Copy request headers from original request
	c.Request().Header.VisitAll(func(key, value []byte) {
		req.Header.SetBytesV(string(key), value)
	})

	utils.BuildRequest(req, method, body, utils.API_KEY, url)

	if c.Request().URI().QueryString() != nil {
		req.URI().SetQueryStringBytes(c.Request().URI().QueryString())
	}

	// Forward request
	if err := fasthttp.Do(req, resp); err != nil {
		log.Printf("Error forwarding request to %s: %v\n", url, err)
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   "Failed to reach service",
			"details": err.Error(),
		})
	}

	// Handle redirects
	if resp.StatusCode() >= 300 && resp.StatusCode() < 400 {
		redirectURL := string(resp.Header.Peek("Location"))
		return c.Redirect(redirectURL, resp.StatusCode())
	}

	// Set response headers
	c.Set("Content-Type", string(resp.Header.Peek("Content-Type")))

	// Handle non-JSON responses (like during OAuth redirects)
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
