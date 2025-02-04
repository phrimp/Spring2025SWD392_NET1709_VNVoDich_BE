package main

import (
	"fmt"
	"gateway/internal/config"
	"gateway/internal/server"
	"os"
)

func main() {
	cfg := config.New()
	gateway := server.NewGateway(cfg)

	port := os.Getenv("GATEWAY_PORT")
	if port == "" {
		port = "8080"
	}

	if err := gateway.Start(":" + port); err != nil {
		panic(fmt.Sprintf("Failed to start server: %v", err))
	}
}
