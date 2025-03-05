package main

import (
	"fmt"
	"os"
	"user-service/internal/handlers"
	"user-service/internal/models"
	"user-service/internal/repository"
	"user-service/internal/services"
	"user-service/utils"

	"github.com/gofiber/fiber/v2"
)

var (
	API_KEY   string
	had_admin bool = false
)

func init() {
	API_KEY = os.Getenv("API_KEY")
	err := services.AddUser(models.UserCreationParams{Username: "admin", Password: "admin", Role: "Admin", Email: "thaiphienn@gmail.com"}, had_admin, "", repository.DB)
	if err != nil {
		fmt.Println("Init admin account failed:", err)
		had_admin = false
		return
	}
	had_admin = false

	utils.SetupTimeZone()
}

func main() {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "ok",
			"service": "user-service",
		})
	})

	user := app.Group("/user", Middleware(API_KEY))

	// Routes
	user.Post("/get", handlers.GetUserWithUsernamePasswordHandler(repository.DB))
	user.Post("/add", handlers.AddUser(repository.DB, had_admin))
	user.Get("/get-public-user", handlers.GetPublicUser(repository.DB))
	user.Get("/get-all-user", handlers.GetAllUser(repository.DB))
	user.Get("", handlers.GetUserwithUsername(repository.DB))
	user.Put("/update", handlers.UpdateUser(repository.DB))
	user.Patch("/update/status", handlers.UpdateUserStatus(repository.DB))
	user.Delete("/delete", handlers.DeleteUserHandler(repository.DB))
	user.Post("/delete/cancel", handlers.CancelDeleteUserHandler(repository.DB))
	user.Put("/users/:username", handlers.AdminUpdateUserHandler(repository.DB))
	user.Put("/verify", handlers.VerifyUserHandler(repository.DB))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8085" // default port
	}

	if err := app.Listen(":" + port); err != nil {
		panic(fmt.Sprintf("Failed to start server: %v", err))
	}
}
