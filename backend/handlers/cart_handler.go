package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/leunameek/celestexmewave/internal/services"
)

// Peticion pa meter item al carrito
type AddCartItemRequest struct {
	ProductID uuid.UUID `json:"product_id" binding:"required"`
	Quantity  int       `json:"quantity" binding:"required,min=1"`
	Size      string    `json:"size"`
	SessionID string    `json:"session_id"`
}

// Peticion pa actualizar item del carrito
type UpdateCartItemRequest struct {
	Quantity int    `json:"quantity" binding:"required,min=1"`
	Size     string `json:"size"`
}

// Traer el carrito del user o de la sesion
func GetCart(c *gin.Context) {
	var userID *uuid.UUID
	var sessionID *string

	// Miramos si esta autenticado
	if userIDStr, exists := c.Get("user_id"); exists {
		if parsed, err := uuid.Parse(userIDStr.(string)); err == nil {
			userID = &parsed
		}
	}

	// Session id desde query
	if sid := c.Query("session_id"); sid != "" {
		sessionID = &sid
	}

	if userID == nil && sessionID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id or session_id required"})
		return
	}

	cart, err := services.GetOrCreateCart(userID, sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	items, err := services.GetCartItems(cart.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	total, err := services.GetCartTotal(cart.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var formattedItems []gin.H
	for _, item := range items {
		formattedItems = append(formattedItems, gin.H{
			"id":           item.ID,
			"product_id":   item.ProductID,
			"product_name": item.Product.Name,
			"quantity":     item.Quantity,
			"size":         item.Size,
			"price":        item.Product.Price,
			"image_url":    "/api/products/images/" + cleanImagePath(item.Product.ImagePath),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          cart.ID,
		"total_items": len(items),
		"total_price": total,
		"items":       formattedItems,
	})
}

// Agregar item al carrito
func AddItem(c *gin.Context) {
	var req AddCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Log del request pa ver que llega
	log.Printf("AddItem request: %+v\n", req)

	var userID *uuid.UUID
	var sessionID *string

	// Miramos si esta autenticado
	if userIDStr, exists := c.Get("user_id"); exists {
		if parsed, err := uuid.Parse(userIDStr.(string)); err == nil {
			userID = &parsed
		}
	}

	// Usamos session que manden
	if req.SessionID != "" {
		sessionID = &req.SessionID
	}

	if userID == nil && sessionID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id or session_id required"})
		return
	}

	cart, err := services.GetOrCreateCart(userID, sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	item, err := services.AddItemToCart(cart.ID, req.ProductID, req.Quantity, req.Size)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         item.ID,
		"product_id": item.ProductID,
		"quantity":   item.Quantity,
		"size":       item.Size,
		"message":    "item added to cart",
	})
}

// Actualizar item del carrito
func UpdateItem(c *gin.Context) {
	itemID, err := uuid.Parse(c.Param("item_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item id"})
		return
	}

	var req UpdateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	item, err := services.UpdateCartItem(itemID, req.Quantity, req.Size)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       item.ID,
		"quantity": item.Quantity,
		"size":     item.Size,
		"message":  "cart item updated",
	})
}

// Quitar item del carrito
func RemoveItem(c *gin.Context) {
	itemID, err := uuid.Parse(c.Param("item_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item id"})
		return
	}

	if err := services.RemoveCartItem(itemID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "item removed from cart"})
}

// Vaciar el carrito completo
func ClearCart(c *gin.Context) {
	var userID *uuid.UUID
	var sessionID *string

	// Miramos si esta autenticado
	if userIDStr, exists := c.Get("user_id"); exists {
		if parsed, err := uuid.Parse(userIDStr.(string)); err == nil {
			userID = &parsed
		}
	}

	// Sacamos session id
	if sid := c.Query("session_id"); sid != "" {
		sessionID = &sid
	}

	if userID == nil && sessionID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id or session_id required"})
		return
	}

	cart, err := services.GetOrCreateCart(userID, sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := services.ClearCart(cart.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "cart cleared"})
}
