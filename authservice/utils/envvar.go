package utils

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type ServicesRoute struct {
	UserService string
}

var (
	API_KEY         string
	SERVICES_ROUTES ServicesRoute
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("Warning: .env file not found, using environment variables")
	}
	API_KEY = os.Getenv("API_KEY")
	SERVICES_ROUTES.UserService = "http://user-service:" + os.Getenv("USER_SERVICE_PORT")
}
