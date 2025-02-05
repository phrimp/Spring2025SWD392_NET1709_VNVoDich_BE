package server

import (
	"gateway/internal/config"
	"gateway/internal/handlers"
	"gateway/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type Gateway struct {
	config *config.Config
	app    *fiber.App
	auth   *handlers.AuthHandler
	google *handlers.GoogleHandler
}

func NewGateway(config *config.Config) *Gateway {
	app := fiber.New(fiber.Config{
		ReadTimeout:  config.ServerCfg.ReadTimeout,
		WriteTimeout: config.ServerCfg.WriteTimeout,
	})

	// Middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET, POST, PUT, DELETE",
		AllowHeaders: "*",
	}))

	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path}\n",
	}))

	gateway := &Gateway{
		config: config,
		app:    app,
		auth:   handlers.NewAuthHandler(config.AuthServiceURL),
		google: handlers.NewGoogleHandler(config.GoogleServiceURL),
	}

	gateway.setupRoutes()
	return gateway
}

func (g *Gateway) setupRoutes() {
	// Public routes
	g.app.Post("/auth/login", g.auth.HandleLogin())
	g.app.Post("/auth/register", g.auth.HandleRegister())
	g.app.Post("/google/auth/login", g.google.HandleLogin())

	// Protected routes
	api := g.app.Group("/api")
	api.Use(middleware.JWTMiddleware(g.config.JWTSecret))

	// User routes (accessible by all authenticated users)
	// api.Get("/profile", g.auth.HandleGetProfile())

	//// Admin routes
	//admin := api.Group("/admin")
	//admin.Use(middleware.RequireAdmin())
	//admin.Get("/users", g.auth.HandleListUsers())
	//admin.Delete("/users/:id", g.auth.HandleDeleteUser())

	//// Manager routes
	//manager := api.Group("/manager")
	//manager.Use(middleware.RequireRole("admin", "manager"))
	//manager.Get("/reports", g.auth.HandleReports())

	//// Specific role-based routes
	//api.Get("/sensitive-data", middleware.RequireRole("admin", "data_analyst"), g.auth.HandleSensitiveData())
}

func (g *Gateway) Start(addr string) error {
	return g.app.Listen(addr)
}
