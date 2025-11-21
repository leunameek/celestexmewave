package utils

import (
	"fmt"
	"net/smtp"
	"strconv"
	"strings"

	"github.com/leunameek/celestexmewave/internal/config"
)

func SendEmail(to, subject, body string) error {
	cfg := config.Get()

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
		return nil
	}

	return nil
}

func formatPriceColombian(price float64) string {
	intPrice := int(price)
	str := strconv.Itoa(intPrice)

	var parts []string
	for len(str) > 3 {
		parts = append([]string{str[len(str)-3:]}, parts...)
		str = str[:len(str)-3]
	}
	if str != "" {
		parts = append([]string{str}, parts...)
	}
	result := strings.Join(parts, ".")

	return result
}

func SendPasswordResetEmail(to, resetCode string) error {
	subject := "Código de Restablecimiento de Contraseña - CelestexMewave"
	body := fmt.Sprintf(`
Hola,

Has solicitado restablecer la contraseña de tu cuenta en CelestexMewave.

Tu código de restablecimiento es: %s

Este código expirará en 1 hora.

Si no solicitaste esto, por favor ignora este correo.

Saludos,
Equipo de CelestexMewave
`, resetCode)

	return SendEmail(to, subject, body)
}

func SendOrderConfirmationEmail(to, orderID string, totalAmount float64, items []map[string]interface{}) error {
	subject := "Confirmación de Pedido - CelestexMewave"

	itemsHTML := ""
	for _, item := range items {
		price := item["unit_price"].(float64)
		itemsHTML += fmt.Sprintf("- %v x %v (Talla: %v) - $%s\n",
			item["product_name"],
			item["quantity"],
			item["size"],
			formatPriceColombian(price),
		)
	}

	body := fmt.Sprintf(`
Hola,

¡Gracias por tu compra!

ID del Pedido: %s
Monto Total: $%s

Productos:
%s

Tu pedido ha sido confirmado y será procesado en breve.

Saludos,
Equipo de CelestexMewave
`, orderID, formatPriceColombian(totalAmount), itemsHTML)

	return SendEmail(to, subject, body)
}

func SendRegistrationEmail(to, firstName string) error {
	subject := "Bienvenido a CelestexMewave"
	body := fmt.Sprintf(`
Hola %s,

¡Bienvenido a CelestexMewave!

Tu cuenta ha sido creada exitosamente. Ya puedes iniciar sesión y comenzar a comprar.

Saludos,
Equipo de CelestexMewave
`, firstName)

	return SendEmail(to, subject, body)
}
