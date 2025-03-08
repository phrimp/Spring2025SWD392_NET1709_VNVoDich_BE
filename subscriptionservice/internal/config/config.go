package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server            ServerConfig
	APIKey            string
	PaymentServiceURL string
	UserServiceURL    string
}

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// New creates a new Config instance with values from environment variables
func New() *Config {
	return &Config{
		Server:            loadServerConfig(),
		APIKey:            os.Getenv("API_KEY"),
		PaymentServiceURL: os.Getenv("PAYMENT_SERVICE_URL"),
		UserServiceURL:    os.Getenv("USER_SERVICE_URL"),
	}
}

func loadServerConfig() ServerConfig {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8086" // default port for subscription service
	}

	readTimeout, err := strconv.Atoi(os.Getenv("SERVER_READ_TIMEOUT"))
	if err != nil {
		readTimeout = 10 // default 10 seconds
	}

	writeTimeout, err := strconv.Atoi(os.Getenv("SERVER_WRITE_TIMEOUT"))
	if err != nil {
		writeTimeout = 10 // default 10 seconds
	}

	return ServerConfig{
		Port:         port,
		ReadTimeout:  time.Duration(readTimeout) * time.Second,
		WriteTimeout: time.Duration(writeTimeout) * time.Second,
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
