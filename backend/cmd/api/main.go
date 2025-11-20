package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/leunameek/celestexmewave/internal/config"
	"github.com/leunameek/celestexmewave/internal/database"
	"github.com/leunameek/celestexmewave/internal/services"
	"github.com/leunameek/celestexmewave/router"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	log.Println("✓ Configuration loaded")

	// Initialize database
	if err := database.Initialize(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	log.Println("✓ Database initialized")

	// Seed database
	if err := services.SeedDatabase(); err != nil {
		log.Printf("Warning: Failed to seed database: %v", err)
	}

	// Setup router
	r := router.SetupRouter()
	log.Println("✓ Router configured")

	// Start server in a goroutine
	go func() {
		addr := fmt.Sprintf("%s:%d", cfg.ServerHost, cfg.ServerPort)
		log.Printf("✓ Starting server on %s\n", addr)
		if err := r.Run(addr); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down server...")
	if err := database.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	}
	log.Println("✓ Server stopped")
}
