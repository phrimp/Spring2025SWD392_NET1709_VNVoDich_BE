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
	user   *handlers.UserServiceHandler
	node   *handlers.NodeServiceHandler
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
		user:   handlers.NewUserService(config.UserServiceURL),
		node:   handlers.NewNodeServiceHandler(config.NodeServiceURL),
	}

	gateway.setupRoutes()
	return gateway
}

func (g *Gateway) setupRoutes() {
	// Public routes
	g.app.Post("/auth/login", g.auth.HandleLogin())
	g.app.Post("/auth/register", g.auth.HandleRegister())
	g.app.Get("/google/auth/login", g.google.HandleLogin())
	g.app.Get("/public/user/:username", g.user.HandleGetUserwithUsername())
	g.app.Get("/public/course/all", g.node.HandleGetAllCourse())
	g.app.Get("/public/course/:id", g.node.HandleGetACourse())

	// Protected routes
	api := g.app.Group("/api")
	api.Use(middleware.JWTMiddleware(g.config.JWTSecret))
	api.Get("/get/me", g.user.HandleGetMe())

	// User routes (accessible by all authenticated users)
	// api.Get("/profile", g.auth.HandleGetProfile())

	//// Admin routes
	//admin := api.Group("/admin")
	//admin.Use(middleware.RequireAdmin())
	//admin.Get("/users", g.auth.HandleListUsers())
	//admin.Delete("/users/:id", g.auth.HandleDeleteUser())

	admin_api := api.Group("/admin")
	admin_api.Use(middleware.RequireRole("admin"))
	admin_api.Get("/get-all/user", g.user.HandleAllGetUser())
	//// Specific role-based routes
	//api.Get("/sensitive-data", middleware.RequireRole("admin", "data_analyst"), g.auth.HandleSensitiveData())
}

func (g *Gateway) Start(addr string) error {
	return g.app.Listen(addr)
}
