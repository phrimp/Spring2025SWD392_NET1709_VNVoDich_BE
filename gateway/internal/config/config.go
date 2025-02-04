package config

import "os"

type Config struct {
	AuthServiceURL   string
	NodeServiceURL   string
	GoogleServiceURL string
	JWTSecret        string
}

func New() *Config {
	return &Config{
		AuthServiceURL:   os.Getenv("AUTH_SERVICE_URL"),
		NodeServiceURL:   os.Getenv("NODE_SERVICE_URL"),
		GoogleServiceURL: os.Getenv("GOOGLE_SERVICE_URL"),
		JWTSecret:        os.Getenv("JWT_SECRET"),
	}
}
