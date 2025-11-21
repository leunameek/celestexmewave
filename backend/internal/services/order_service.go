package services

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/leunameek/celestexmewave/internal/database"
	"github.com/leunameek/celestexmewave/models"
)

// ShippingInfo es la info de envio sin misterio
type ShippingInfo struct {
	Name       string
	Phone      string
	Email      string
	City       string
	Address    string
	Address2   string
	PostalCode string
	Notes      string
}

// CreateOrderFromCart arma el pedido con lo que haya en el carrito
func CreateOrderFromCart(cartID uuid.UUID, userID *uuid.UUID, sessionID *string, shipping ShippingInfo) (*models.Order, error) {
	// Sacamos los items del carro
	cartItems, err := GetCartItems(cartID)
	if err != nil {
		return nil, err
	}

	if len(cartItems) == 0 {
		return nil, fmt.Errorf("cart is empty")
	}

	// Calculamos total
	total := 0.0
	for _, item := range cartItems {
		total += item.Product.Price * float64(item.Quantity)
	}

	// Creamos el pedido
	order := &models.Order{
		ID:                 uuid.New(),
		UserID:             userID,
		SessionID:          sessionID,
		TotalAmount:        total,
		Status:             "pending",
		PaymentStatus:      "pending",
		ShippingName:       shipping.Name,
		ShippingPhone:      shipping.Phone,
		ShippingEmail:      shipping.Email,
		ShippingCity:       shipping.City,
		ShippingAddress:    shipping.Address,
		ShippingAddress2:   shipping.Address2,
		ShippingPostalCode: shipping.PostalCode,
		ShippingNotes:      shipping.Notes,
	}

	if err := database.DB.Create(order).Error; err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Creamos los items del pedido
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

	// Limpiamos el carrito
	if err := ClearCart(cartID); err != nil {
		return nil, err
	}

	return order, nil
}

// GetOrder trae pedido por id
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

// GetOrdersByUser trae pedidos de un user
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

	// Conteo total
	if err := database.DB.Where("user_id = ?", userID).Model(&models.Order{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count orders: %w", err)
	}

	// Pedidos paginados
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

// GetOrdersBySession trae pedidos por session
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

	// Conteo total
	if err := database.DB.Where("session_id = ?", sessionID).Model(&models.Order{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count orders: %w", err)
	}

	// Pedidos paginados
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

// UpdateOrderStatus cambia el status del pedido
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

// UpdatePaymentStatus cambia el estado del pago
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
