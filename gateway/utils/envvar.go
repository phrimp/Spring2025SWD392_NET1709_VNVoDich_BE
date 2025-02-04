package utils

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var API_KEY string

func init() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("Warning: .env file not found, using environment variables")
	}
	API_KEY = os.Getenv("API_KEY")
}
