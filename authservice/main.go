package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

func init() {
	godotenv.Load(".env")
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

	// Database connection
	db := initDB()

	// Health check for docker compose
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Routes
	app.Post("/login", handleLogin(db))
	app.Post("/register", handleRegister(db))

	port := os.Getenv("AUTH_SERVICE_PORT")
	if port == "" {
		port = "8081" // default port
	}

	if err := app.Listen(":" + port); err != nil {
		panic(fmt.Sprintf("Failed to start server: %v", err))
	}
}

func initDB() *gorm.DB {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	if dbUser == "" || dbPass == "" || dbHost == "" || dbName == "" {
		fmt.Println("Warning: Some database environment variables are missing. Using default values.")
		// Fallback to default values if env vars are not set
		dbUser = "user"
		dbPass = "password"
		dbHost = "mysql"
		dbName = "auth_db"
	}
	if dbPort == "" {
		dbPort = "3306" // default MySQL port
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	// Auto migrate schemas
	if err := db.AutoMigrate(&User{}); err != nil {
		panic(fmt.Sprintf("Failed to migrate database: %v", err))
	}

	return db
}
