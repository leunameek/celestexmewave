package utils

import (
	"regexp"
	"strings"
)

// ValidateEmail validates email format
func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(pattern, email)
	return match
}

// ValidatePhone validates phone format (basic validation)
func ValidatePhone(phone string) bool {
	// Remove common separators
	cleaned := strings.ReplaceAll(phone, "-", "")
	cleaned = strings.ReplaceAll(cleaned, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "(", "")
	cleaned = strings.ReplaceAll(cleaned, ")", "")
	cleaned = strings.ReplaceAll(cleaned, "+", "")

	// Check if it's between 10-15 digits
	return len(cleaned) >= 10 && len(cleaned) <= 15 && isNumeric(cleaned)
}

// ValidatePassword validates password strength
func ValidatePassword(password string) bool {
	// Minimum 8 characters
	if len(password) < 8 {
		return false
	}
	return true
}

// ValidateName validates name format
func ValidateName(name string) bool {
	// Name should not be empty and should contain only letters and spaces
	if len(strings.TrimSpace(name)) == 0 {
		return false
	}
	pattern := `^[a-zA-ZáéíóúÁÉÍÓÚñÑ\s]+$`
	match, _ := regexp.MatchString(pattern, name)
	return match
}

// ValidateCardNumber validates credit card number (basic Luhn algorithm)
func ValidateCardNumber(cardNumber string) bool {
	// Remove spaces and dashes
	cleaned := strings.ReplaceAll(cardNumber, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")

	// Check if it's 16 digits
	if len(cleaned) != 16 || !isNumeric(cleaned) {
		return false
	}

	// Basic Luhn algorithm
	sum := 0
	for i, digit := range cleaned {
		d := int(digit - '0')
		if (len(cleaned)-i)%2 == 0 {
			d *= 2
			if d > 9 {
				d -= 9
			}
		}
		sum += d
	}
	return sum%10 == 0
}

// ValidateCVV validates CVV format
func ValidateCVV(cvv string) bool {
	// CVV should be 3-4 digits
	if len(cvv) < 3 || len(cvv) > 4 {
		return false
	}
	return isNumeric(cvv)
}

// ValidateExpiryDate validates card expiry date
func ValidateExpiryDate(month, year int) bool {
	// Month should be 1-12
	if month < 1 || month > 12 {
		return false
	}
	// Year should be current year or later
	// This is a simple check, in production you'd compare with current date
	return year >= 2024
}

// Helper function to check if string contains only digits
func isNumeric(s string) bool {
	for _, char := range s {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}

// GenerateResetCode generates a random 6-digit reset code
func GenerateResetCode() string {
	const charset = "0123456789"
	code := ""
	for i := 0; i < 6; i++ {
		code += string(charset[randomInt(0, len(charset))])
	}
	return code
}

// randomInt generates a random integer between min and max
func randomInt(min, max int) int {
	// Simple pseudo-random for demo purposes
	// In production, use crypto/rand
	return min + (int(randomByte()) % (max - min))
}

// randomByte generates a random byte
func randomByte() byte {
	// Simple implementation for demo
	// In production, use crypto/rand
	return byte(42) // Placeholder
}
