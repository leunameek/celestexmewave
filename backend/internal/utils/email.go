package utils

import (
	"fmt"
	"net/smtp"

	"github.com/leunameek/celestexmewave/internal/config"
)

// SendEmail sends an email using SMTP
func SendEmail(to, subject, body string) error {
	cfg := config.Get()

	// Skip email sending if SMTP credentials are not configured
	if cfg.SMTPUser == "" || cfg.SMTPPass == "" {
		fmt.Printf("[EMAIL MOCK] To: %s\nSubject: %s\nBody:\n%s\n\n", to, subject, body)
		return nil
	}

	auth := smtp.PlainAuth("", cfg.SMTPUser, cfg.SMTPPass, cfg.SMTPHost)
	addr := fmt.Sprintf("%s:%d", cfg.SMTPHost, cfg.SMTPPort)

	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", cfg.SMTPFrom, to, subject, body)

	err := smtp.SendMail(addr, auth, cfg.SMTPFrom, []string{to}, []byte(message))
	if err != nil {
		fmt.Printf("[EMAIL ERROR] Failed to send email to %s: %v\n", to, err)
		// Don't return error to allow app to continue
		return nil
	}

	return nil
}

// SendPasswordResetEmail sends a password reset email
func SendPasswordResetEmail(to, resetCode string) error {
	subject := "Password Reset Code - CelestexMewave"
	body := fmt.Sprintf(`
Hello,

You requested a password reset for your CelestexMewave account.

Your reset code is: %s

This code will expire in 1 hour.

If you didn't request this, please ignore this email.

Best regards,
CelestexMewave Team
`, resetCode)

	return SendEmail(to, subject, body)
}

// SendOrderConfirmationEmail sends an order confirmation email
func SendOrderConfirmationEmail(to, orderID string, totalAmount float64, items []map[string]interface{}) error {
	subject := "Order Confirmation - CelestexMewave"

	itemsHTML := ""
	for _, item := range items {
		itemsHTML += fmt.Sprintf("- %v x %v (Size: %v) - $%.2f\n",
			item["product_name"],
			item["quantity"],
			item["size"],
			item["unit_price"],
		)
	}

	body := fmt.Sprintf(`
Hello,

Thank you for your purchase!

Order ID: %s
Total Amount: $%.2f

Items:
%s

Your order has been confirmed and will be processed shortly.

Best regards,
CelestexMewave Team
`, orderID, totalAmount, itemsHTML)

	return SendEmail(to, subject, body)
}

// SendRegistrationEmail sends a registration confirmation email
func SendRegistrationEmail(to, firstName string) error {
	subject := "Welcome to CelestexMewave"
	body := fmt.Sprintf(`
Hello %s,

Welcome to CelestexMewave!

Your account has been successfully created. You can now login and start shopping.

Best regards,
CelestexMewave Team
`, firstName)

	return SendEmail(to, subject, body)
}
