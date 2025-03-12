package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server           ServerConfig
	DatabaseConfig   DatabaseConfig
	APIKey           string
	ExternalServices ExternalServices
}

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type DatabaseConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
}

type ExternalServices struct {
	UserService   string
	NodeService   string
	GoogleService string
}

// New creates a new Config instance with values from environment variables
func New() *Config {
	return &Config{
		Server:           loadServerConfig(),
		DatabaseConfig:   loadDatabaseConfig(),
		APIKey:           os.Getenv("API_KEY"),
		ExternalServices: loadExternalServiceConfig(),
	}
}

func loadExternalServiceConfig() ExternalServices {
	user_service := os.Getenv("USER_SERVICE_URL")
	if user_service == "" {
		user_service = "user-service:8085"
	}
	node_service := os.Getenv("NODE_SERVICE_URL")
	if user_service == "" {
		user_service = "node-service:8000"
	}
	google_service := os.Getenv("GOOGLE_SERVICE_URL")
	if google_service == "" {
		google_service = "http://google-service:8084"
	}

	return ExternalServices{
		UserService:   user_service,
		NodeService:   node_service,
		GoogleService: google_service,
	}
}

func loadServerConfig() ServerConfig {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8083" // default port for admin service
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

func loadDatabaseConfig() DatabaseConfig {
	// Use the same environment variables as in your docker-compose.yml file
	return DatabaseConfig{
		Host:     getEnvOrDefault("DB_HOST", "mysql"),
		Port:     getEnvOrDefault("DB_PORT", "3306"),
		Username: getEnvOrDefault("DB_USER", "appuser"),
		Password: getEnvOrDefault("DB_PASSWORD", "apppassword"),
		DBName:   getEnvOrDefault("DB_NAME", "online_tutoring_platform"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
