package services

import (
	"github.com/google/uuid"
	"github.com/leunameek/celestexmewave/internal/utils"
)

// PaymentRequest represents a payment request
type PaymentRequest struct {
	CardNumber  string
	CardHolder  string
	ExpiryMonth int
	ExpiryYear  int
	CVV         string
}

// PaymentResponse represents a payment response
type PaymentResponse struct {
	OrderID          string
	PaymentStatus    string
	TransactionID    string
	Message          string
	ConfirmationSent bool
}

// ProcessMockPayment processes a mock payment
func ProcessMockPayment(orderID uuid.UUID, payment PaymentRequest) (*PaymentResponse, error) {
	// Validate card details
	if !utils.ValidateCardNumber(payment.CardNumber) {
		return &PaymentResponse{
			OrderID:       orderID.String(),
			PaymentStatus: "failed",
			Message:       "invalid card number",
		}, nil
	}

	if !utils.ValidateCVV(payment.CVV) {
		return &PaymentResponse{
			OrderID:       orderID.String(),
			PaymentStatus: "failed",
			Message:       "invalid CVV",
		}, nil
	}

	if !utils.ValidateExpiryDate(payment.ExpiryMonth, payment.ExpiryYear) {
		return &PaymentResponse{
			OrderID:       orderID.String(),
			PaymentStatus: "failed",
			Message:       "card expired",
		}, nil
	}

	// Mock payment processing - always succeeds for demo
	transactionID := "TXN_" + orderID.String()[:8]

	// Update order payment status
	_, err := UpdatePaymentStatus(orderID, "completed")
	if err != nil {
		return &PaymentResponse{
			OrderID:       orderID.String(),
			PaymentStatus: "failed",
			Message:       "failed to update order",
		}, nil
	}

	// Update order status to confirmed
	_, err = UpdateOrderStatus(orderID, "confirmed")
	if err != nil {
		return &PaymentResponse{
			OrderID:       orderID.String(),
			PaymentStatus: "failed",
			Message:       "failed to confirm order",
		}, nil
	}

	return &PaymentResponse{
		OrderID:          orderID.String(),
		PaymentStatus:    "completed",
		TransactionID:    transactionID,
		Message:          "payment processed successfully",
		ConfirmationSent: true,
	}, nil
}

// ValidateCardNumber validates a card number using Luhn algorithm
func ValidateCardNumber(cardNumber string) bool {
	return utils.ValidateCardNumber(cardNumber)
}

// ValidateCVV validates a CVV
func ValidateCVV(cvv string) bool {
	return utils.ValidateCVV(cvv)
}

// ValidateExpiryDate validates an expiry date
func ValidateExpiryDate(month, year int) bool {
	return utils.ValidateExpiryDate(month, year)
}
