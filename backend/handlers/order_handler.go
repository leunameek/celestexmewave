package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/leunameek/celestexmewave/internal/services"
	"github.com/leunameek/celestexmewave/models"
)

// CreateOrderRequest represents a request to create an order
type CreateOrderRequest struct {
	SessionID string `json:"session_id"`
	Email     string `json:"email"`
}

// ProcessPaymentRequest represents a payment request
type ProcessPaymentRequest struct {
	CardNumber  string `json:"card_number" binding:"required"`
	CardHolder  string `json:"card_holder" binding:"required"`
	ExpiryMonth int    `json:"expiry_month" binding:"required"`
	ExpiryYear  int    `json:"expiry_year" binding:"required"`
	CVV         string `json:"cvv" binding:"required"`
}

// CreateOrder creates an order from cart
func CreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	var userID *uuid.UUID
	var sessionID *string

	// Check if user is authenticated
	if userIDStr, exists := c.Get("user_id"); exists {
		if parsed, err := uuid.Parse(userIDStr.(string)); err == nil {
			userID = &parsed
		}
	}

	// Use session ID from request
	if req.SessionID != "" {
		sessionID = &req.SessionID
	}

	if userID == nil && sessionID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id or session_id required"})
		return
	}

	// Get or create cart
	cart, err := services.GetOrCreateCart(userID, sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create order from cart
	order, err := services.CreateOrderFromCart(cart.ID, userID, sessionID, req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch order items
	items, _ := services.GetCartItems(cart.ID)

	var formattedItems []gin.H
	for _, item := range items {
		formattedItems = append(formattedItems, gin.H{
			"product_id":   item.ProductID,
			"product_name": item.Product.Name,
			"quantity":     item.Quantity,
			"size":         item.Size,
			"unit_price":   item.Product.Price,
		})
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":             order.ID,
		"total_amount":   order.TotalAmount,
		"status":         order.Status,
		"payment_status": order.PaymentStatus,
		"items":          formattedItems,
		"created_at":     order.CreatedAt,
	})
}

// GetOrder retrieves an order by ID
func GetOrder(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	order, err := services.GetOrder(orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var formattedItems []gin.H
	for _, item := range order.OrderItems {
		formattedItems = append(formattedItems, gin.H{
			"product_id":   item.ProductID,
			"product_name": item.Product.Name,
			"quantity":     item.Quantity,
			"size":         item.Size,
			"unit_price":   item.UnitPrice,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"id":             order.ID,
		"total_amount":   order.TotalAmount,
		"status":         order.Status,
		"payment_status": order.PaymentStatus,
		"items":          formattedItems,
		"created_at":     order.CreatedAt,
		"updated_at":     order.UpdatedAt,
	})
}

// GetOrders retrieves all orders for user or session
func GetOrders(c *gin.Context) {
	page := 1
	limit := 10

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

	var orders []models.Order
	var total int64
	var err error

	// Check if user is authenticated
	if userIDStr, exists := c.Get("user_id"); exists {
		if userID, err := uuid.Parse(userIDStr.(string)); err == nil {
			orders, total, err = services.GetOrdersByUser(userID, page, limit)
		}
	} else if sessionID := c.Query("session_id"); sessionID != "" {
		orders, total, err = services.GetOrdersBySession(sessionID, page, limit)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id or session_id required"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var formattedOrders []gin.H
	for _, order := range orders {
		var items []gin.H
		for _, item := range order.OrderItems {
			items = append(items, gin.H{
				"product_id":   item.ProductID,
				"product_name": item.Product.Name,
				"quantity":     item.Quantity,
				"size":         item.Size,
				"unit_price":   item.UnitPrice,
			})
		}

		formattedOrders = append(formattedOrders, gin.H{
			"id":             order.ID,
			"total_amount":   order.TotalAmount,
			"status":         order.Status,
			"payment_status": order.PaymentStatus,
			"items":          items,
			"created_at":     order.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"total":  total,
		"page":   page,
		"limit":  limit,
		"orders": formattedOrders,
	})
}

// ProcessPayment processes payment for an order
func ProcessPayment(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	var req ProcessPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Process payment
	paymentReq := services.PaymentRequest{
		CardNumber:  req.CardNumber,
		CardHolder:  req.CardHolder,
		ExpiryMonth: req.ExpiryMonth,
		ExpiryYear:  req.ExpiryYear,
		CVV:         req.CVV,
	}

	response, err := services.ProcessMockPayment(orderID, paymentReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if response.PaymentStatus == "failed" {
		c.JSON(http.StatusBadRequest, gin.H{
			"order_id":       response.OrderID,
			"payment_status": response.PaymentStatus,
			"message":        response.Message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"order_id":                response.OrderID,
		"payment_status":          response.PaymentStatus,
		"transaction_id":          response.TransactionID,
		"message":                 response.Message,
		"confirmation_email_sent": response.ConfirmationSent,
	})
}

// GetConfirmation retrieves order confirmation
func GetConfirmation(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	order, err := services.GetOrder(orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var formattedItems []gin.H
	for _, item := range order.OrderItems {
		formattedItems = append(formattedItems, gin.H{
			"product_name": item.Product.Name,
			"quantity":     item.Quantity,
			"size":         item.Size,
			"unit_price":   item.UnitPrice,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"order_id":       order.ID,
		"order_date":     order.CreatedAt,
		"total_amount":   order.TotalAmount,
		"items":          formattedItems,
		"status":         order.Status,
		"payment_status": order.PaymentStatus,
		"message":        "Thank you for your purchase!",
	})
}
