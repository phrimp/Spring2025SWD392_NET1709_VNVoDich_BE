package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

var Google_config *Config

func init() {
	Google_config = GetConfig()
}

type Config struct {
	Server           ServerConfig
	GoogleAuth       *GoogleOAuthConfig
	Email            *EmailConfig
	JWT              JWTConfig
	ServiceAccount   *ServiceAccountConfig
	API_KEY          string
	USER_SERVICE_URL string
}

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

type JWTConfig struct {
	Secret        string
	ExpiresIn     time.Duration
	RefreshSecret string
}

type ServiceAccountConfig struct {
	CredentialsJSON []byte
}

// New creates a new Config instance with values from environment variables
func New() *Config {
	return &Config{
		Server:           loadServerConfig(),
		GoogleAuth:       NewGoogleOAuthConfig(),
		Email:            NewEmailConfig(),
		JWT:              loadJWTConfig(),
		ServiceAccount:   NewServiceAccountConfig(),
		API_KEY:          os.Getenv("API_KEY"),
		USER_SERVICE_URL: os.Getenv("USER_SERVICE_URL"),
	}
}

func loadServerConfig() ServerConfig {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
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

func loadSMTPConfig() SMTPConfig {
	port := os.Getenv("SMTP_PORT")
	if port == "" {
		port = "587" // default SMTP port
	}

	return SMTPConfig{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     port,
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
		From:     os.Getenv("SMTP_FROM"),
	}
}

func loadJWTConfig() JWTConfig {
	expiresIn, err := strconv.Atoi(os.Getenv("JWT_EXPIRES_IN"))
	if err != nil {
		expiresIn = 24 // default 24 hours
	}

	return JWTConfig{
		Secret:        getRequiredEnv("JWT_SECRET"),
		ExpiresIn:     time.Duration(expiresIn) * time.Hour,
		RefreshSecret: getRequiredEnv("JWT_REFRESH_SECRET"),
	}
}

// getRequiredEnv gets an environment variable or panics if it's not set
func getRequiredEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic("Required environment variable not set: " + key)
	}
	return value
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	// Add validation logic here if needed
	return nil
}

// GetConfig is a helper function to get config instance
func GetConfig() *Config {
	config := New()
	if err := config.Validate(); err != nil {
		panic(err)
	}
	return config
}

func loadMeetCredentials() ([]byte, error) {
	// Load from environment variable
	if creds := os.Getenv("GOOGLE_SERVICE_ACCOUNT_CREDENTIALS"); creds != "" {
		return []byte(creds), nil
	}

	// Or load from file
	return os.ReadFile("./internal/config/credentials/service-account.json")
}

func NewServiceAccountConfig() *ServiceAccountConfig {
	credentials, err := loadMeetCredentials()
	if err != nil {
		fmt.Println("Create New Service Account failed:", err)
		return nil
	}
	return &ServiceAccountConfig{
		CredentialsJSON: credentials,
	}
}
