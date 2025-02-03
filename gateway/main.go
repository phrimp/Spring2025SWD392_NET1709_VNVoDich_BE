package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

var API_KEY string

type Config struct {
	AuthServiceURL string
	NodeServiceURL string
	JWTSecret      string
}

type Gateway struct {
	config Config
	app    *fiber.App
}

func init() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("Warning: .env file not found, using environment variables")
	}
	API_KEY = os.Getenv("API_KEY")
}

func NewGateway(config Config) *Gateway {
	app := fiber.New(fiber.Config{})

	// Middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET, POST, PUT, DELETE",
		AllowHeaders: "*",
	}))
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path}\n",
	}))

	return &Gateway{
		config: config,
		app:    app,
	}
}

func main() {
	config := Config{
		AuthServiceURL: os.Getenv("AUTH_SERVICE_URL"),
		NodeServiceURL: os.Getenv("NODE_SERVICE_URL"),
		JWTSecret:      os.Getenv("JWT_SECRET"),
	}

	gateway := NewGateway(config)

	// Routes without login
	gateway.app.Post("/auth/login", gateway.handleLogin())
	gateway.app.Post("/auth/register", gateway.handleRegister())

	api := gateway.app.Group("/api")
	api.Use(gateway.jwtMiddleware())

	// Protected routes
	api.Get("/users", gateway.handleGetUsers())

	port := os.Getenv("GATEWAY_PORT")
	if port == "" {
		port = "8080" // default port
	}

	if err := gateway.app.Listen(":" + port); err != nil {
		panic(fmt.Sprintf("Failed to start server: %v", err))
	}
}
