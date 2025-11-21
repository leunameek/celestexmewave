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
	// Cargamos la config toda chill
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	log.Println("✓ Configuration loaded")

	// Montamos la base de datos sin drama
	if err := database.Initialize(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	log.Println("✓ Database initialized")

	// Semillita de datos si toca
	if err := services.SeedDatabase(); err != nil {
		log.Printf("Warning: Failed to seed database: %v", err)
	}

	// Armamos el router bacan
	r := router.SetupRouter()
	log.Println("✓ Router configured")

	// Levantamos el server en una go routine pa no bloquear
	go func() {
		addr := fmt.Sprintf("%s:%d", cfg.ServerHost, cfg.ServerPort)
		log.Printf("✓ Starting server on %s\n", addr)
		if err := r.Run(addr); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Esperamos la senal de cierre (ctrl+c vibes)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down server...")
	if err := database.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	}
	log.Println("✓ Server stopped")
}
