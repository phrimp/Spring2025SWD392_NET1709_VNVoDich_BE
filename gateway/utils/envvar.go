package utils

import (
	"log"
	"os"
	"time"
)

var API_KEY string

func init() {
	API_KEY = os.Getenv("API_KEY")
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
