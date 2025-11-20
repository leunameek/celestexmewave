package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/leunameek/celestexmewave/internal/database"
	"github.com/leunameek/celestexmewave/internal/utils"
	"github.com/leunameek/celestexmewave/models"
)

// RegisterUser registers a new user
func RegisterUser(email, phone, firstName, lastName, password string) (*models.User, error) {
	// Validate input
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

	// Check if user already exists
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

	// Hash password
	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
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

	// Send registration email if email is provided
	if email != "" {
		_ = utils.SendRegistrationEmail(email, firstName)
	}

	return user, nil
}

// LoginUser authenticates a user and returns JWT tokens
func LoginUser(emailOrPhone, password string) (*models.User, string, string, error) {
	// Find user by email or phone
	var user models.User
	if err := database.DB.Where("email = ? OR phone = ?", emailOrPhone, emailOrPhone).First(&user).Error; err != nil {
		return nil, "", "", fmt.Errorf("invalid credentials")
	}

	// Verify password
	if !utils.VerifyPassword(user.PasswordHash, password) {
		return nil, "", "", fmt.Errorf("invalid credentials")
	}

	// Generate tokens
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

// RefreshAccessToken generates a new access token from a refresh token
func RefreshAccessToken(refreshToken string) (string, error) {
	userID, err := utils.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token: %w", err)
	}

	// Get user
	var user models.User
	if err := database.DB.First(&user, "id = ?", userID).Error; err != nil {
		return "", fmt.Errorf("user not found: %w", err)
	}

	// Generate new access token
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

// RequestPasswordReset creates a password reset code
func RequestPasswordReset(emailOrPhone string) (string, error) {
	// Find user
	var user models.User
	if err := database.DB.Where("email = ? OR phone = ?", emailOrPhone, emailOrPhone).First(&user).Error; err != nil {
		return "", fmt.Errorf("user not found")
	}

	// Generate reset code
	resetCode := utils.GenerateResetCode()

	// Create password reset record
	passwordReset := &models.PasswordReset{
		ID:        uuid.New(),
		UserID:    user.ID,
		ResetCode: resetCode,
		ExpiresAt: time.Now().Add(1 * time.Hour),
		Used:      false,
	}

	if err := database.DB.Create(passwordReset).Error; err != nil {
		return "", fmt.Errorf("failed to create password reset: %w", err)
	}

	// Send reset code via email
	if user.Email != nil {
		_ = utils.SendPasswordResetEmail(*user.Email, resetCode)
	}

	return resetCode, nil
}

// VerifyResetCode verifies a reset code and updates the password
func VerifyResetCode(emailOrPhone, resetCode, newPassword string) error {
	// Find user
	var user models.User
	if err := database.DB.Where("email = ? OR phone = ?", emailOrPhone, emailOrPhone).First(&user).Error; err != nil {
		return fmt.Errorf("user not found")
	}

	// Find and validate reset code
	var passwordReset models.PasswordReset
	if err := database.DB.Where("user_id = ? AND reset_code = ?", user.ID, resetCode).First(&passwordReset).Error; err != nil {
		return fmt.Errorf("invalid reset code")
	}

	// Check if code is expired
	if time.Now().After(passwordReset.ExpiresAt) {
		return fmt.Errorf("reset code expired")
	}

	// Check if code was already used
	if passwordReset.Used {
		return fmt.Errorf("reset code already used")
	}

	// Validate new password
	if !utils.ValidatePassword(newPassword) {
		return fmt.Errorf("password must be at least 8 characters")
	}

	// Hash new password
	passwordHash, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update user password
	if err := database.DB.Model(&user).Update("password_hash", passwordHash).Error; err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Mark reset code as used
	if err := database.DB.Model(&passwordReset).Update("used", true).Error; err != nil {
		return fmt.Errorf("failed to mark reset code as used: %w", err)
	}

	return nil
}
