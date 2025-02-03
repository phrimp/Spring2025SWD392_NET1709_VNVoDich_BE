package main

import (
	"fmt"
	"gateway/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/valyala/fasthttp"
)

func (g *Gateway) handleLogin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Forward request to auth service
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)

		// Copy request body
		utils.BuildRequest(req, "POST", c.Body(), API_KEY, g.config.AuthServiceURL+"/login")

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
}

func (g *Gateway) handleRegister() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)

		// Copy request body
		utils.BuildRequest(req, "POST", c.Body(), API_KEY, g.config.AuthServiceURL+"/register")

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
}

// JWT middleware
func (g *Gateway) jwtMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authorization token",
			})
		}

		// Validate token
		_, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return []byte(g.config.JWTSecret), nil
		})
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		return c.Next()
	}
}

func (g *Gateway) handleGetUsers() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)

		req.Header.SetMethod("GET")
		req.SetRequestURI(g.config.NodeServiceURL + "/user" + string(c.Request().URI().QueryString()))
		req.Header.Set("API_KEY", API_KEY)

		if err := fasthttp.Do(req, resp); err != nil {
			fmt.Printf("Error forwarding request: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Service unavailable",
			})
		}

		return c.Status(resp.StatusCode()).Send(resp.Body())
	}
}
