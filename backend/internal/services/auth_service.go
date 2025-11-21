package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/leunameek/celestexmewave/internal/database"
	"github.com/leunameek/celestexmewave/internal/utils"
	"github.com/leunameek/celestexmewave/models"
)

// ErrInvalidCredentials sale cuando las credenciales estan mal
var ErrInvalidCredentials = errors.New("invalid credentials")

// RegisterUser mete un user nuevo en la DB
func RegisterUser(email, phone, firstName, lastName, password string) (*models.User, error) {
	// Validamos los datos
	if email == "" && phone == "" {
		return nil, fmt.Errorf("email or phone is required")
	}

	if !utils.ValidateName(firstName) || !utils.ValidateName(lastName) {
		return nil, fmt.Errorf("invalid first name or last name")
	}

	if !utils.ValidatePassword(password) {
		return nil, fmt.Errorf("password must be at least 8 characters")
	}

	if email != "" && !utils.ValidateEmail(email) {
		return nil, fmt.Errorf("invalid email format")
	}

	if phone != "" && !utils.ValidatePhone(phone) {
		return nil, fmt.Errorf("invalid phone format")
	}

	// Revisamos si ya existe
	var existingUser models.User
	if email != "" {
		if err := database.DB.Where("email = ?", email).First(&existingUser).Error; err == nil {
			return nil, fmt.Errorf("email already registered")
		}
	}

	if phone != "" {
		if err := database.DB.Where("phone = ?", phone).First(&existingUser).Error; err == nil {
			return nil, fmt.Errorf("phone already registered")
		}
	}

	// Hasheamos la clave
	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Creamos el user
	user := &models.User{
		ID:           uuid.New(),
		FirstName:    firstName,
		LastName:     lastName,
		PasswordHash: passwordHash,
		IsRegistered: true,
	}

	if email != "" {
		user.Email = &email
	}

	if phone != "" {
		user.Phone = &phone
	}

	if err := database.DB.Create(user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Enviamos correo de bienvenida si hay email
	if email != "" {
		_ = utils.SendRegistrationEmail(email, firstName)
	}

	return user, nil
}

// LoginUser valida user y devuelve tokens
func LoginUser(emailOrPhone, password string) (*models.User, string, string, error) {
	// Buscamos user por email o telefono
	var user models.User
	if err := database.DB.Where("email = ? OR phone = ?", emailOrPhone, emailOrPhone).First(&user).Error; err != nil {
		return nil, "", "", ErrInvalidCredentials
	}

	// Verificamos la clave
	if !utils.VerifyPassword(user.PasswordHash, password) {
		return nil, "", "", ErrInvalidCredentials
	}

	// Sacamos tokens
	accessToken, err := utils.GenerateAccessToken(user.ID, emailOrPhone, user.FirstName, user.LastName)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &user, accessToken, refreshToken, nil
}

// RefreshAccessToken saca un access token nuevo desde refresh
func RefreshAccessToken(refreshToken string) (string, error) {
	userID, err := utils.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token: %w", err)
	}

	// Buscamos el user
	var user models.User
	if err := database.DB.First(&user, "id = ?", userID).Error; err != nil {
		return "", fmt.Errorf("user not found: %w", err)
	}

	// Generamos token nuevo
	identifier := ""
	if user.Email != nil {
		identifier = *user.Email
	} else if user.Phone != nil {
		identifier = *user.Phone
	}
	accessToken, err := utils.GenerateAccessToken(user.ID, identifier, user.FirstName, user.LastName)
	if err != nil {
		return "", fmt.Errorf("failed to generate access token: %w", err)
	}

	return accessToken, nil
}

// RequestPasswordReset crea codigo de reset
func RequestPasswordReset(emailOrPhone string) (string, error) {
	// Buscamos el user
	var user models.User
	if err := database.DB.Where("email = ? OR phone = ?", emailOrPhone, emailOrPhone).First(&user).Error; err != nil {
		return "", fmt.Errorf("user not found")
	}

	// Intentamos crear codigo, reintentamos si choca
	var resetCode string
	const maxAttempts = 5
	for i := 0; i < maxAttempts; i++ {
		resetCode = utils.GenerateResetCode()
		passwordReset := &models.PasswordReset{
			ID:        uuid.New(),
			UserID:    user.ID,
			ResetCode: resetCode,
			ExpiresAt: time.Now().Add(1 * time.Hour),
			Used:      false,
		}

		if err := database.DB.Create(passwordReset).Error; err != nil {
			// Si choca por duplicado, seguimos intentando
			if strings.Contains(err.Error(), "idx_password_resets_reset_code") {
				continue
			}
			return "", fmt.Errorf("failed to crear el código de recuperación, intenta de nuevo")
		}

		// Mandamos el codigo por correo si hay
		if user.Email != nil {
			_ = utils.SendPasswordResetEmail(*user.Email, resetCode)
		}

		return resetCode, nil
	}

	return "", fmt.Errorf("no pudimos generar un código de recuperación, por favor intenta nuevamente")
}

// VerifyResetCode valida el codigo y cambia la clave
func VerifyResetCode(emailOrPhone, resetCode, newPassword string) error {
	// Buscamos user
	var user models.User
	if err := database.DB.Where("email = ? OR phone = ?", emailOrPhone, emailOrPhone).First(&user).Error; err != nil {
		return fmt.Errorf("user not found")
	}

	// Buscamos y validamos el codigo
	var passwordReset models.PasswordReset
	if err := database.DB.Where("user_id = ? AND reset_code = ?", user.ID, resetCode).First(&passwordReset).Error; err != nil {
		return fmt.Errorf("invalid reset code")
	}

	// Revisamos expiracion
	if time.Now().After(passwordReset.ExpiresAt) {
		return fmt.Errorf("reset code expired")
	}

	// Revisamos si ya fue usado
	if passwordReset.Used {
		return fmt.Errorf("reset code already used")
	}

	// Revisamos la clave nueva
	if !utils.ValidatePassword(newPassword) {
		return fmt.Errorf("password must be at least 8 characters")
	}

	// Hasheamos la clave nueva
	passwordHash, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Guardamos la nueva clave
	if err := database.DB.Model(&user).Update("password_hash", passwordHash).Error; err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Marcamos codigo como usado
	if err := database.DB.Model(&passwordReset).Update("used", true).Error; err != nil {
		return fmt.Errorf("failed to mark reset code as used")
	}

	return nil
}
