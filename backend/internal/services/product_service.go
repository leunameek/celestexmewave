package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/google/uuid"
	"github.com/leunameek/celestexmewave/internal/database"
	"github.com/leunameek/celestexmewave/models"
)

// ProductJSON es la forma del producto en el JSON
type ProductJSON struct {
	Category       string   `json:"category"`
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	Price          float64  `json:"price"`
	Sizes          []string `json:"sizes"`
	AvailableUnits int      `json:"available_units"`
	Image          string   `json:"image"`
}

// LoadProductsFromJSON carga productos desde un archivo JSON, todo casero
func LoadProductsFromJSON(storeID uuid.UUID, storeName, jsonPath string) error {
	// Leemos el archivo JSON
	data, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("failed to read JSON file: %w", err)
	}

	var products []ProductJSON
	if err := json.Unmarshal(data, &products); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Insertamos o actualizamos productos
	for _, p := range products {
		sizesJSON, _ := json.Marshal(p.Sizes)
		
		// Revisamos si el producto ya existe
		var existingProduct models.Product
		err := database.DB.Where("store_id = ? AND name = ?", storeID, p.Name).First(&existingProduct).Error

		if err == nil {
			// Si existe, lo actualizamos
			existingProduct.Description = p.Description
			existingProduct.Price = p.Price
			existingProduct.AvailableUnits = p.AvailableUnits
			existingProduct.ImagePath = p.Image
			existingProduct.Sizes = sizesJSON
			
			if err := database.DB.Save(&existingProduct).Error; err != nil {
				fmt.Printf("Failed to update product %s: %v\n", p.Name, err)
			}
		} else {
			// Producto nuevo
			product := &models.Product{
				ID:             uuid.New(),
				StoreID:        storeID,
				Name:           p.Name,
				Description:    p.Description,
				Category:       p.Category,
				Price:          p.Price,
				AvailableUnits: p.AvailableUnits,
				ImagePath:      p.Image,
				Sizes:          sizesJSON,
			}

			if err := database.DB.Create(product).Error; err != nil {
				fmt.Printf("Failed to create product %s: %v\n", p.Name, err)
			}
		}
	}

	return nil
}

// GetAllProducts trae productos con filtros opcionales
func GetAllProducts(store, category string, minPrice, maxPrice float64, page, limit int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit
	query := database.DB

	// Filtros opcionales
	if store != "" {
		query = query.Joins("JOIN stores ON products.store_id = stores.id").
			Where("stores.name = ?", store)
	}

	if category != "" {
		query = query.Where("category = ?", category)
	}

	if minPrice > 0 {
		query = query.Where("price >= ?", minPrice)
	}

	if maxPrice > 0 {
		query = query.Where("price <= ?", maxPrice)
	}

	// Total de productos
	if err := query.Model(&models.Product{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count products: %w", err)
	}

	// Productos paginados
	if err := query.
		Preload("Store").
		Offset(offset).
		Limit(limit).
		Find(&products).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to fetch products: %w", err)
	}

	return products, total, nil
}

// GetProductByID trae producto por ID
func GetProductByID(productID uuid.UUID) (*models.Product, error) {
	var product models.Product
	if err := database.DB.Preload("Store").First(&product, "id = ?", productID).Error; err != nil {
		return nil, fmt.Errorf("product not found")
	}
	return &product, nil
}

// GetProductsByStore trae productos por tienda
func GetProductsByStore(storeID uuid.UUID, page, limit int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	// Total de productos
	if err := database.DB.Where("store_id = ?", storeID).Model(&models.Product{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count products: %w", err)
	}

	// Productos paginados
	if err := database.DB.
		Where("store_id = ?", storeID).
		Preload("Store").
		Offset(offset).
		Limit(limit).
		Find(&products).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to fetch products: %w", err)
	}

	return products, total, nil
}

// GetProductsByCategory trae productos por categoria
func GetProductsByCategory(category string, page, limit int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	// Total de productos
	if err := database.DB.Where("category = ?", category).Model(&models.Product{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count products: %w", err)
	}

	// Productos paginados
	if err := database.DB.
		Where("category = ?", category).
		Preload("Store").
		Offset(offset).
		Limit(limit).
		Find(&products).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to fetch products: %w", err)
	}

	return products, total, nil
}
