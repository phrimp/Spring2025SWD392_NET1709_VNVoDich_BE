package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	AuthServiceURL    string
	NodeServiceURL    string
	GoogleServiceURL  string
	UserServiceURL    string
	AdminServiceURL   string
	PaymentServiceURL string
	SubscriptionURL   string
	JWTSecret         string
	ServerCfg         ServerConfig
}

func New() *Config {
	return &Config{
		AuthServiceURL:    os.Getenv("AUTH_SERVICE_URL"),
		NodeServiceURL:    os.Getenv("NODE_SERVICE_URL"),
		GoogleServiceURL:  os.Getenv("GOOGLE_SERVICE_URL"),
		UserServiceURL:    os.Getenv("USER_SERVICE_URL"),
		AdminServiceURL:   os.Getenv("ADMIN_SERVICE_URL"),
		PaymentServiceURL: os.Getenv("PAYMENT_SERVICE_URL"),
		SubscriptionURL:   os.Getenv("SUBSCRIPTION_SERVICE_URL"),
		JWTSecret:         os.Getenv("JWT_SECRET"),
		ServerCfg:         loadServerConfig(),
	}
}

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func loadServerConfig() ServerConfig {
	port := os.Getenv("PORT")
	if port == "" {
		fmt.Println("PORT env is empty")
		port = "8080"
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
