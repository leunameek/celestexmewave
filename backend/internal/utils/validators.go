package utils

import (
	"crypto/rand"
	"regexp"
	"strings"
	"time"
)

func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(pattern, email)
	return match
}

func ValidatePhone(phone string) bool {
	cleaned := strings.ReplaceAll(phone, "-", "")
	cleaned = strings.ReplaceAll(cleaned, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "(", "")
	cleaned = strings.ReplaceAll(cleaned, ")", "")
	cleaned = strings.ReplaceAll(cleaned, "+", "")

	return len(cleaned) >= 10 && len(cleaned) <= 15 && isNumeric(cleaned)
}

func ValidatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	return true
}

func ValidateName(name string) bool {
	if len(strings.TrimSpace(name)) == 0 {
		return false
	}
	pattern := `^[a-zA-ZáéíóúÁÉÍÓÚñÑ\s]+$`
	match, _ := regexp.MatchString(pattern, name)
	return match
}

func ValidateCardNumber(cardNumber string) bool {
	cleaned := strings.ReplaceAll(cardNumber, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")

	if len(cleaned) != 16 || !isNumeric(cleaned) {
		return false
	}

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

func ValidateCVV(cvv string) bool {
	if len(cvv) < 3 || len(cvv) > 4 {
		return false
	}
	return isNumeric(cvv)
}

func ValidateExpiryDate(month, year int) bool {
	if month < 1 || month > 12 {
		return false
	}
	return year >= 2024
}

func isNumeric(s string) bool {
	for _, char := range s {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}

func GenerateResetCode() string {
	const charset = "0123456789"
	code := ""
	for i := 0; i < 6; i++ {
		code += string(charset[randomInt(0, len(charset))])
	}
	return code
}

func randomInt(min, max int) int {
	if max <= min {
		return min
	}
	return min + (int(randomByte()) % (max - min))
}

func randomByte() byte {
	var b [1]byte
	if _, err := rand.Read(b[:]); err == nil {
		return b[0]
	}
	return byte(time.Now().UnixNano())
}
