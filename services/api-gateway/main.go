package main

import (
	"log"

	"github.com/bitbiz/hias-core/configs"
)

func main() {
	config, localConfigLoaded, err := configs.LoadConfig("./configs")
	if err != nil && !localConfigLoaded {
		log.Println("No local config found, attempting SSM parameters...")
	}

	if !localConfigLoaded {
		if err := configs.LoadSSMParameters(&config); err != nil {
			log.Printf("Warning: Failed to load SSM parameters: %v", err)
		}
	}

	// Set defaults
	if config.HTTPServerAddress == "" {
		config.HTTPServerAddress = "0.0.0.0:8080"
	}
	if config.GRPCServerAddress == "" {
		config.GRPCServerAddress = "0.0.0.0:9090"
	}

	configs.AppConfig = config

	server, err := NewUnifiedServer(config)
	if err != nil {
		log.Fatalf("Cannot create server: %v", err)
	}

	if err := server.Start(); err != nil {
		log.Fatalf("Cannot start server: %v", err)
	}
}
