package services

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/leunameek/celestexmewave/internal/database"
	"github.com/leunameek/celestexmewave/models"
)

// GetOrCreateCart gets or creates a cart for a user or session
func GetOrCreateCart(userID *uuid.UUID, sessionID *string) (*models.Cart, error) {
	var cart models.Cart

	// Try to find existing cart
	query := database.DB
	if userID != nil {
		query = query.Where("user_id = ?", userID)
	} else if sessionID != nil {
		query = query.Where("session_id = ?", sessionID)
	} else {
		return nil, fmt.Errorf("user_id or session_id required")
	}

	err := query.First(&cart).Error
	if err == nil {
		return &cart, nil
	}

	// Create new cart
	cart = models.Cart{
		ID:        uuid.New(),
		UserID:    userID,
		SessionID: sessionID,
	}

	if err := database.DB.Create(&cart).Error; err != nil {
		return nil, fmt.Errorf("failed to create cart: %w", err)
	}

	return &cart, nil
}

// AddItemToCart adds an item to the cart
func AddItemToCart(cartID, productID uuid.UUID, quantity int, size string) (*models.CartItem, error) {
	// Check if product exists and has stock
	var product models.Product
	if err := database.DB.First(&product, "id = ?", productID).Error; err != nil {
		return nil, fmt.Errorf("product not found")
	}

	if product.AvailableUnits < quantity {
		return nil, fmt.Errorf("insufficient stock")
	}

	// Check if item already in cart
	var existingItem models.CartItem
	if err := database.DB.Where("cart_id = ? AND product_id = ? AND size = ?", cartID, productID, size).First(&existingItem).Error; err == nil {
		// Update quantity
		existingItem.Quantity += quantity
		if err := database.DB.Save(&existingItem).Error; err != nil {
			return nil, fmt.Errorf("failed to update cart item: %w", err)
		}
		return &existingItem, nil
	}

	// Create new cart item
	cartItem := &models.CartItem{
		ID:        uuid.New(),
		CartID:    cartID,
		ProductID: productID,
		Quantity:  quantity,
		Size:      size,
	}

	if err := database.DB.Create(cartItem).Error; err != nil {
		return nil, fmt.Errorf("failed to add item to cart: %w", err)
	}

	return cartItem, nil
}

// UpdateCartItem updates a cart item
func UpdateCartItem(itemID uuid.UUID, quantity int, size string) (*models.CartItem, error) {
	var item models.CartItem
	if err := database.DB.First(&item, "id = ?", itemID).Error; err != nil {
		return nil, fmt.Errorf("cart item not found")
	}

	// Check stock
	var product models.Product
	if err := database.DB.First(&product, "id = ?", item.ProductID).Error; err != nil {
		return nil, fmt.Errorf("product not found")
	}

	if product.AvailableUnits < quantity {
		return nil, fmt.Errorf("insufficient stock")
	}

	item.Quantity = quantity
	if size != "" {
		item.Size = size
	}

	if err := database.DB.Save(&item).Error; err != nil {
		return nil, fmt.Errorf("failed to update cart item: %w", err)
	}

	return &item, nil
}

// RemoveCartItem removes an item from the cart
func RemoveCartItem(itemID uuid.UUID) error {
	if err := database.DB.Delete(&models.CartItem{}, "id = ?", itemID).Error; err != nil {
		return fmt.Errorf("failed to remove cart item: %w", err)
	}
	return nil
}

// ClearCart removes all items from a cart
func ClearCart(cartID uuid.UUID) error {
	if err := database.DB.Delete(&models.CartItem{}, "cart_id = ?", cartID).Error; err != nil {
		return fmt.Errorf("failed to clear cart: %w", err)
	}
	return nil
}

// GetCartItems retrieves all items in a cart
func GetCartItems(cartID uuid.UUID) ([]models.CartItem, error) {
	var items []models.CartItem
	if err := database.DB.
		Where("cart_id = ?", cartID).
		Preload("Product").
		Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch cart items: %w", err)
	}
	return items, nil
}

// GetCartTotal calculates the total price of a cart
func GetCartTotal(cartID uuid.UUID) (float64, error) {
	var total float64
	if err := database.DB.
		Table("cart_items").
		Joins("JOIN products ON cart_items.product_id = products.id").
		Where("cart_items.cart_id = ?", cartID).
		Select("COALESCE(SUM(cart_items.quantity * products.price), 0)").
		Row().
		Scan(&total); err != nil {
		return 0, fmt.Errorf("failed to calculate total: %w", err)
	}
	return total, nil
}
