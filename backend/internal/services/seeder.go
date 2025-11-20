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
			// Store doesn't exist, create it
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
			// Store exists, use its ID
			if store.Name == "Celeste" {
				celesteID = existingStore.ID
			} else {
				mewaveID = existingStore.ID
			}
		}
	}
	log.Println("✓ Stores synced")

	// Determine assets path
	// Assuming running from backend/ directory
	assetsPath := "../assets/products"
	
	// Check if we are in root or backend
	if _, err := os.Stat(assetsPath); os.IsNotExist(err) {
		// Try absolute path or different relative path if needed
		// For now, let's try to find it relative to current working directory
		// If running from root: assets/products
		if _, err := os.Stat("assets/products"); err == nil {
			assetsPath = "assets/products"
		}
	}

	// Load Celeste Products
	celesteJSON := filepath.Join(assetsPath, "celeste.json")
	if err := LoadProductsFromJSON(celesteID, "Celeste", celesteJSON); err != nil {
		log.Printf("Warning: Failed to load Celeste products: %v", err)
	} else {
		log.Println("✓ Celeste products loaded")
	}

	// Load Mewave Products
	mewaveJSON := filepath.Join(assetsPath, "mewave.json")
	if err := LoadProductsFromJSON(mewaveID, "Mewave", mewaveJSON); err != nil {
		log.Printf("Warning: Failed to load Mewave products: %v", err)
	} else {
		log.Println("✓ Mewave products loaded")
	}

	return nil
}
