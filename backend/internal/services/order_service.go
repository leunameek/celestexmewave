package services

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/leunameek/celestexmewave/internal/database"
	"github.com/leunameek/celestexmewave/internal/utils"
	"github.com/leunameek/celestexmewave/models"
)

// CreateOrderFromCart creates an order from cart items
func CreateOrderFromCart(cartID uuid.UUID, userID *uuid.UUID, sessionID *string, email string) (*models.Order, error) {
	// Get cart items
	cartItems, err := GetCartItems(cartID)
	if err != nil {
		return nil, err
	}

	if len(cartItems) == 0 {
		return nil, fmt.Errorf("cart is empty")
	}

	// Calculate total
	total := 0.0
	for _, item := range cartItems {
		total += item.Product.Price * float64(item.Quantity)
	}

	// Create order
	order := &models.Order{
		ID:            uuid.New(),
		UserID:        userID,
		SessionID:     sessionID,
		TotalAmount:   total,
		Status:        "pending",
		PaymentStatus: "pending",
	}

	if err := database.DB.Create(order).Error; err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Create order items
	for _, cartItem := range cartItems {
		orderItem := &models.OrderItem{
			ID:        uuid.New(),
			OrderID:   order.ID,
			ProductID: cartItem.ProductID,
			Quantity:  cartItem.Quantity,
			Size:      cartItem.Size,
			UnitPrice: cartItem.Product.Price,
		}

		if err := database.DB.Create(orderItem).Error; err != nil {
			return nil, fmt.Errorf("failed to create order item: %w", err)
		}
	}

	// Clear cart
	if err := ClearCart(cartID); err != nil {
		return nil, err
	}

	// Send confirmation email
	if email != "" {
		var items []map[string]interface{}
		for _, item := range cartItems {
			items = append(items, map[string]interface{}{
				"product_name": item.Product.Name,
				"quantity":     item.Quantity,
				"size":         item.Size,
				"unit_price":   item.Product.Price,
			})
		}
		_ = utils.SendOrderConfirmationEmail(email, order.ID.String(), total, items)
	}

	return order, nil
}

// GetOrder retrieves an order by ID
func GetOrder(orderID uuid.UUID) (*models.Order, error) {
	var order models.Order
	if err := database.DB.
		Preload("OrderItems").
		Preload("OrderItems.Product").
		First(&order, "id = ?", orderID).Error; err != nil {
		return nil, fmt.Errorf("order not found")
	}
	return &order, nil
}

// GetOrdersByUser retrieves all orders for a user
func GetOrdersByUser(userID uuid.UUID, page, limit int) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	// Get total count
	if err := database.DB.Where("user_id = ?", userID).Model(&models.Order{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count orders: %w", err)
	}

	// Get paginated orders
	if err := database.DB.
		Where("user_id = ?", userID).
		Preload("OrderItems").
		Preload("OrderItems.Product").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&orders).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to fetch orders: %w", err)
	}

	return orders, total, nil
}

// GetOrdersBySession retrieves all orders for a session
func GetOrdersBySession(sessionID string, page, limit int) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	// Get total count
	if err := database.DB.Where("session_id = ?", sessionID).Model(&models.Order{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count orders: %w", err)
	}

	// Get paginated orders
	if err := database.DB.
		Where("session_id = ?", sessionID).
		Preload("OrderItems").
		Preload("OrderItems.Product").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&orders).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to fetch orders: %w", err)
	}

	return orders, total, nil
}

// UpdateOrderStatus updates the status of an order
func UpdateOrderStatus(orderID uuid.UUID, status string) (*models.Order, error) {
	var order models.Order
	if err := database.DB.First(&order, "id = ?", orderID).Error; err != nil {
		return nil, fmt.Errorf("order not found")
	}

	if err := database.DB.Model(&order).Update("status", status).Error; err != nil {
		return nil, fmt.Errorf("failed to update order status: %w", err)
	}

	return &order, nil
}

// UpdatePaymentStatus updates the payment status of an order
func UpdatePaymentStatus(orderID uuid.UUID, paymentStatus string) (*models.Order, error) {
	var order models.Order
	if err := database.DB.First(&order, "id = ?", orderID).Error; err != nil {
		return nil, fmt.Errorf("order not found")
	}

	if err := database.DB.Model(&order).Update("payment_status", paymentStatus).Error; err != nil {
		return nil, fmt.Errorf("failed to update payment status: %w", err)
	}

	return &order, nil
}
