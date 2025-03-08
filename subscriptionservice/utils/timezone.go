package utils

import (
	"log"
	"os"
	"time"
)

func SetupTimeZone() {
	// Set default timezone to value in TZ environment variable
	tz := os.Getenv("TZ")
	if tz == "" {
		tz = "Asia/Ho_Chi_Minh" // Default timezone
	}

	loc, err := time.LoadLocation(tz)
	if err != nil {
		log.Printf("Failed to load %s location: %v", tz, err)
		return
	}

	// Set the default timezone
	time.Local = loc
	log.Printf("Default timezone set to: %s", time.Local.String())
}
