package utils

import (
	"log"
	"os"
	"time"
)

type ServicesRoute struct {
	UserService string
}

var (
	API_KEY         string
	SERVICES_ROUTES ServicesRoute
)

func init() {
	API_KEY = os.Getenv("API_KEY")
	SERVICES_ROUTES.UserService = "http://user-service:" + os.Getenv("USER_SERVICE_PORT")
}

func SetupTimeZone() {
	// Set default timezone to Asia/Ho_Chi_Minh
	loc, err := time.LoadLocation(os.Getenv("TZ"))
	if err != nil {
		log.Printf("Failed to load %s location: %v", loc.String(), err)
		return
	}

	// Set the default timezone
	time.Local = loc
	log.Printf("Default timezone set to: %s", time.Local.String())
}
