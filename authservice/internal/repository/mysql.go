package repository

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	DB = initMySQLDB()
}

func initMySQLDB() *gorm.DB {
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
