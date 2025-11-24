package services

import (
	"github.com/google/uuid"
	"github.com/leunameek/celestexmewave/internal/database"
	"github.com/leunameek/celestexmewave/internal/utils"
	"github.com/leunameek/celestexmewave/models"
)

// PaymentRequest es la peti de pago fake
type PaymentRequest struct {
	CardNumber  string
	CardHolder  string
	ExpiryMonth int
	ExpiryYear  int
	CVV         string
}

// PaymentResponse es la respuesta del pago mock
type PaymentResponse struct {
	OrderID          string
	PaymentStatus    string
	TransactionID    string
	Message          string
	ConfirmationSent bool
}

// ProcessMockPayment procesa el pago de mentiras
func ProcessMockPayment(orderID uuid.UUID, payment PaymentRequest) (*PaymentResponse, error) {
	// Pago mock: siempre pasa en modo demo, sin validar nada de la tarjeta
	transactionID := "TXN_" + orderID.String()[:8]

	// Actualizamos status de pago
	_, err := UpdatePaymentStatus(orderID, "completed")
	if err != nil {
		return &PaymentResponse{
			OrderID:       orderID.String(),
			PaymentStatus: "failed",
			Message:       "failed to update order",
		}, nil
	}

	// Cambiamos el status del pedido a confirmado
	_, err = UpdateOrderStatus(orderID, "confirmed")
	if err != nil {
		return &PaymentResponse{
			OrderID:       orderID.String(),
			PaymentStatus: "failed",
			Message:       "failed to confirm order",
		}, nil
	}

	// Cargamos pedido e items para mandar correo de confirmacion
	var order models.Order
	if err := database.DB.Preload("OrderItems.Product").First(&order, "id = ?", orderID).Error; err == nil && order.ShippingEmail != "" {
		var items []map[string]interface{}
		for _, item := range order.OrderItems {
			items = append(items, map[string]interface{}{
				"product_name": item.Product.Name,
				"quantity":     item.Quantity,
				"size":         item.Size,
				"unit_price":   item.UnitPrice,
			})
		}
		_ = utils.SendOrderConfirmationEmail(order.ShippingEmail, order.ID.String(), order.TotalAmount, items)
	}

	return &PaymentResponse{
		OrderID:          orderID.String(),
		PaymentStatus:    "completed",
		TransactionID:    transactionID,
		Message:          "payment processed successfully",
		ConfirmationSent: true,
	}, nil
}

// ValidateCardNumber usa Luhn y devuelve si pasa
func ValidateCardNumber(cardNumber string) bool {
	return utils.ValidateCardNumber(cardNumber)
}

// ValidateCVV revisa el CVV
func ValidateCVV(cvv string) bool {
	return utils.ValidateCVV(cvv)
}

// ValidateExpiryDate mira si la fecha exp es valida
func ValidateExpiryDate(month, year int) bool {
	return utils.ValidateExpiryDate(month, year)
}
