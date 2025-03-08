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
	config       *config.Config
	app          *fiber.App
	auth         *handlers.AuthHandler
	google       *handlers.GoogleHandler
	user         *handlers.UserServiceHandler
	node         *handlers.NodeServiceHandler
	admin        *handlers.AdminServiceHandler
	payment      *handlers.PaymentHandler
	subscription *handlers.SubscriptionHandler
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
		config:       config,
		app:          app,
		auth:         handlers.NewAuthHandler(config.AuthServiceURL),
		google:       handlers.NewGoogleHandler(config),
		user:         handlers.NewUserService(config.UserServiceURL),
		node:         handlers.NewNodeServiceHandler(config.NodeServiceURL),
		admin:        handlers.NewAdminService(config),
		payment:      handlers.NewPaymentHandler(config.PaymentServiceURL),
		subscription: handlers.NewSubscriptionHandler(config.SubscriptionURL),
	}

	gateway.setupRoutes()
	return gateway
}

func (g *Gateway) setupRoutes() {
	// Public routes
	g.app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "ok",
			"service": "gateway",
		})
	})
	g.app.Post("/auth/login", g.auth.HandleLogin())
	g.app.Post("/auth/register", g.auth.HandleRegister())
	g.app.Get("/google/auth/login", g.google.HandleLogin())
	g.app.Get("/google/auth/login/callback", g.google.HandleCallback())

	g.app.Get("/public/user/:username", g.user.HandleGetPublicUserProfile())
	g.app.Get("/public/course/all", g.node.HandleGetAllCourse())
	g.app.Get("/public/course/:id", g.node.HandleGetACourse())
	g.app.Get("/payment/success", g.payment.HandleCompletePayPalPayment())
	g.app.Get("/payment/cancel", g.payment.HandleCancelPayPalPayment())

	g.app.Get("/subscription/plans", g.subscription.HandleGetPlans())
	g.app.Get("/subscription/plans/:id", g.subscription.HandleGetPlan())

	// Protected routes
	api := g.app.Group("/api")
	api.Use(middleware.JWTMiddleware(g.config.JWTSecret))
	api.Get("/get/me", g.user.HandleGetMe())
	api.Put("/update/me", g.user.HandleUpdateMe())
	api.Delete("/delete/me", g.user.HandleDeleteMe())
	api.Post("/delete/me/cancel", g.user.HandleCancelDeleteMe())
	api.Post("verify-email/send", g.google.HandleSendVerificationEmail())
	api.Post("verify-email/verify", g.google.HandleVerifyEmail())
	api.Post("/payment/create", g.payment.HandleCreatePayment())

	api.Get("/subscription/tutor/:tutorId", g.subscription.HandleGetTutorSubscription())
	api.Post("/subscription", g.subscription.HandleCreateSubscription())
	api.Post("/subscription/confirm", g.subscription.HandleConfirmSubscription())
	api.Put("/subscription/:id/cancel", g.subscription.HandleCancelSubscription())
	api.Put("/subscription/:id/change-plan", g.subscription.HandleChangePlan())

	// Tutor routes
	tutor_api := api.Group("/tutor").Use(middleware.RequireRole("Tutor"))
	tutor_api.Get("/meet", g.google.HandleCreateMeetLink())
	//// Admin routes
	admin_api := api.Group("/admin")
	admin_api.Use(middleware.RequireRole("Admin"))
	admin_api.Get("/users", g.user.HandleAllGetUser())
	admin_api.Put("/user/update", g.admin.HandleAdminUpdateUser())
	admin_api.Get("/user", g.admin.HandleAdminGetUSerDetail())
	admin_api.Patch("/users/:username/status", g.admin.HandleUpdateUserStatus())
	admin_api.Delete("/users/:id", g.admin.HandleDeleteUser())
	admin_api.Post("/users/:username/roles", g.admin.HandleAssignRole())

	admin_api.Get("/subscriptions", g.subscription.HandleGetAllSubscriptions())
	admin_api.Put("/subscription/:id/status", g.subscription.HandleUpdateSubscriptionStatus())
	admin_api.Post("/subscription/plans", g.subscription.HandleAdminCreatePlan())
	admin_api.Put("/subscription/plans/:id", g.subscription.HandleAdminUpdatePlan())
	admin_api.Delete("/subscription/plans/:id", g.subscription.HandleAdminDeletePlan())
	//// Specific role-based routes
	//api.Get("/sensitive-data", middleware.RequireRole("admin", "data_analyst"), g.auth.HandleSensitiveData())
}

func (g *Gateway) Start(addr string) error {
	return g.app.Listen(addr)
}
