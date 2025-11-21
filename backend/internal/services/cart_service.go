package services

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/leunameek/celestexmewave/internal/database"
	"github.com/leunameek/celestexmewave/models"
)

// GetOrCreateCart trae o crea un carrito para user o sesion, sin misterio
func GetOrCreateCart(userID *uuid.UUID, sessionID *string) (*models.Cart, error) {
	var cart models.Cart

	// Intentamos encontrar el carrito ya creado
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

	// Crear carrito nuevo
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

// AddItemToCart mete un item al carrito
func AddItemToCart(cartID, productID uuid.UUID, quantity int, size string) (*models.CartItem, error) {
	// Revisamos que el producto exista y tenga stock
	var product models.Product
	if err := database.DB.First(&product, "id = ?", productID).Error; err != nil {
		return nil, fmt.Errorf("product not found")
	}

	if product.AvailableUnits < quantity {
		return nil, fmt.Errorf("insufficient stock")
	}

	// Si ya esta en el carrito, solo sumamos cantidad
	var existingItem models.CartItem
	if err := database.DB.Where("cart_id = ? AND product_id = ? AND size = ?", cartID, productID, size).First(&existingItem).Error; err == nil {
		// Update quantity
		existingItem.Quantity += quantity
		if err := database.DB.Save(&existingItem).Error; err != nil {
			return nil, fmt.Errorf("failed to update cart item: %w", err)
		}
		return &existingItem, nil
	}

	// Item nuevo en el carrito
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

// UpdateCartItem actualiza cantidad o talla
func UpdateCartItem(itemID uuid.UUID, quantity int, size string) (*models.CartItem, error) {
	var item models.CartItem
	if err := database.DB.First(&item, "id = ?", itemID).Error; err != nil {
		return nil, fmt.Errorf("cart item not found")
	}

	// Revisamos stock
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

// RemoveCartItem quita un item del carrito
func RemoveCartItem(itemID uuid.UUID) error {
	if err := database.DB.Delete(&models.CartItem{}, "id = ?", itemID).Error; err != nil {
		return fmt.Errorf("failed to remove cart item: %w", err)
	}
	return nil
}

// ClearCart borra todo el carrito
func ClearCart(cartID uuid.UUID) error {
	if err := database.DB.Delete(&models.CartItem{}, "cart_id = ?", cartID).Error; err != nil {
		return fmt.Errorf("failed to clear cart: %w", err)
	}
	return nil
}

// GetCartItems trae todos los items del carrito
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

// GetCartTotal calcula el total del carrito
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
