package services

import (
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/leunameek/celestexmewave/internal/database"
	"github.com/leunameek/celestexmewave/models"
)

// SeedDatabase populates the database with initial data if it's empty
func SeedDatabase() error {
	log.Println("Seeding/Syncing database...")

	// Create Stores if they don't exist
	celesteID := uuid.New()
	mewaveID := uuid.New()

	stores := []models.Store{
		{
			Name:        "Celeste",
			Description: "Women's fashion brand",
		},
		{
			Name:        "Mewave",
			Description: "Streetwear brand",
		},
	}

	for _, store := range stores {
		var existingStore models.Store
		if err := database.DB.Where("name = ?", store.Name).First(&existingStore).Error; err != nil {
			// Si no existe, la creamos
			store.ID = uuid.New()
			if store.Name == "Celeste" {
				celesteID = store.ID
			} else {
				mewaveID = store.ID
			}
			if err := database.DB.Create(&store).Error; err != nil {
				return err
			}
		} else {
			// Si existe, usamos el ID que ya tiene
			if store.Name == "Celeste" {
				celesteID = existingStore.ID
			} else {
				mewaveID = existingStore.ID
			}
		}
	}
	log.Println("✓ Stores synced")

	// Vemos la ruta de assets (asumimos backend/)
	assetsPath := "../assets/products"
	
	// Si no existe, intentamos otra ruta
	if _, err := os.Stat(assetsPath); os.IsNotExist(err) {
		if _, err := os.Stat("assets/products"); err == nil {
			assetsPath = "assets/products"
		}
	}

	// Cargamos productos de Celeste
	celesteJSON := filepath.Join(assetsPath, "celeste.json")
	if err := LoadProductsFromJSON(celesteID, "Celeste", celesteJSON); err != nil {
		log.Printf("Warning: Failed to load Celeste products: %v", err)
	} else {
		log.Println("✓ Celeste products loaded")
	}

	// Cargamos productos de Mewave
	mewaveJSON := filepath.Join(assetsPath, "mewave.json")
	if err := LoadProductsFromJSON(mewaveID, "Mewave", mewaveJSON); err != nil {
		log.Printf("Warning: Failed to load Mewave products: %v", err)
	} else {
		log.Println("✓ Mewave products loaded")
	}

	return nil
}
