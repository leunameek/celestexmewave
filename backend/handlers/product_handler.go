package handlers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/leunameek/celestexmewave/internal/config"
	"github.com/leunameek/celestexmewave/internal/services"
)

// Trae productos con filtros opcionales, sin tanto show
func GetAllProducts(c *gin.Context) {
	store := c.Query("store")
	category := c.Query("category")
	minPrice := 0.0
	maxPrice := 0.0
	page := 1
	limit := 20

	if mp := c.Query("min_price"); mp != "" {
		if parsed, err := strconv.ParseFloat(mp, 64); err == nil {
			minPrice = parsed
		}
	}

	if mp := c.Query("max_price"); mp != "" {
		if parsed, err := strconv.ParseFloat(mp, 64); err == nil {
			maxPrice = parsed
		}
	}

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	products, total, err := services.GetAllProducts(store, category, minPrice, maxPrice, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Formateamos la respuesta pa que llegue chevere
	var formattedProducts []gin.H
	for _, product := range products {
		sizes, _ := product.GetSizes()
		formattedProducts = append(formattedProducts, gin.H{
			"id":              product.ID,
			"store_id":        product.StoreID,
			"store_name":      product.Store.Name,
			"name":            product.Name,
			"description":     product.Description,
			"category":        product.Category,
			"price":           product.Price,
			"available_units": product.AvailableUnits,
			"image_url":       imageURL(product.ImagePath),
			"sizes":           sizes,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"total":    total,
		"page":     page,
		"limit":    limit,
		"products": formattedProducts,
	})
}

// GetProductByID retrieves a product by ID
func GetProductByID(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	product, err := services.GetProductByID(productID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	sizes, _ := product.GetSizes()
	c.JSON(http.StatusOK, gin.H{
		"id":              product.ID,
		"store_id":        product.StoreID,
		"store_name":      product.Store.Name,
		"name":            product.Name,
		"description":     product.Description,
		"category":        product.Category,
		"price":           product.Price,
		"available_units": product.AvailableUnits,
		"image_url":       imageURL(product.ImagePath),
		"sizes":           sizes,
		"created_at":      product.CreatedAt,
	})
}

// GetProductsByStore retrieves products from a specific store
func GetProductsByStore(c *gin.Context) {
	storeID, err := uuid.Parse(c.Param("store_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid store id"})
		return
	}

	page := 1
	limit := 20

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	products, total, err := services.GetProductsByStore(storeID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var formattedProducts []gin.H
	for _, product := range products {
		sizes, _ := product.GetSizes()
		formattedProducts = append(formattedProducts, gin.H{
			"id":              product.ID,
			"store_id":        product.StoreID,
			"store_name":      product.Store.Name,
			"name":            product.Name,
			"description":     product.Description,
			"category":        product.Category,
			"price":           product.Price,
			"available_units": product.AvailableUnits,
			"image_url":       imageURL(product.ImagePath),
			"sizes":           sizes,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"total":    total,
		"page":     page,
		"limit":    limit,
		"products": formattedProducts,
	})
}

// GetProductsByCategory retrieves products by category
func GetProductsByCategory(c *gin.Context) {
	category := c.Param("category")
	page := 1
	limit := 20

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	products, total, err := services.GetProductsByCategory(category, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var formattedProducts []gin.H
	for _, product := range products {
		sizes, _ := product.GetSizes()
		formattedProducts = append(formattedProducts, gin.H{
			"id":              product.ID,
			"store_id":        product.StoreID,
			"store_name":      product.Store.Name,
			"name":            product.Name,
			"description":     product.Description,
			"category":        product.Category,
			"price":           product.Price,
			"available_units": product.AvailableUnits,
			"image_url":       imageURL(product.ImagePath),
			"sizes":           sizes,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"total":    total,
		"page":     page,
		"limit":    limit,
		"products": formattedProducts,
	})
}

// ServeImage serves a product image
func ServeImage(c *gin.Context) {
	filename := c.Param("filename")
	// Remove leading slash if present (Gin's *param includes it)
	if len(filename) > 0 && filename[0] == '/' {
		filename = filename[1:]
	}
	cfg := config.Get()

	// Prevent directory traversal
	if filename == ".." || filename == "." {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid filename"})
		return
	}

	// Log raw request
	log.Printf("DEBUG: ServeImage requested filename: '%s'", filename)

	// Construct full path
	// Note: filename here is already "cleaned" by the router/URL construction if it came from our API
	// But let's see what we actually get.
	fullPath := filepath.Join(cfg.UploadDir, filename)

	// Get absolute path for debugging
	absPath, _ := filepath.Abs(fullPath)
	log.Printf("DEBUG: ServeImage looking for file at: '%s' (Absolute: '%s')", fullPath, absPath)

	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		log.Printf("ERROR: ServeImage File NOT FOUND at '%s'", fullPath)
		// List directory contents to see what's actually there
		dir := filepath.Dir(fullPath)
		files, _ := os.ReadDir(dir)
		log.Printf("DEBUG: Listing directory '%s':", dir)
		for _, f := range files {
			log.Printf("  - %s", f.Name())
		}

		c.JSON(http.StatusNotFound, gin.H{"error": "image not found"})
		return
	}

	c.File(fullPath)
}
