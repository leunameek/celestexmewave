package services

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/leunameek/celestexmewave/internal/database"
	"github.com/leunameek/celestexmewave/internal/utils"
	"github.com/leunameek/celestexmewave/models"
)

// GetUserProfile retrieves a user's profile
func GetUserProfile(userID uuid.UUID) (*models.User, error) {
	var user models.User
	if err := database.DB.First(&user, "id = ?", userID).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return &user, nil
}

// UpdateUserProfile updates a user's profile
func UpdateUserProfile(userID uuid.UUID, firstName, lastName, phone string) (*models.User, error) {
	var user models.User
	if err := database.DB.First(&user, "id = ?", userID).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Validate input
	if firstName != "" && !utils.ValidateName(firstName) {
		return nil, fmt.Errorf("invalid first name")
	}

	if lastName != "" && !utils.ValidateName(lastName) {
		return nil, fmt.Errorf("invalid last name")
	}

	if phone != "" && !utils.ValidatePhone(phone) {
		return nil, fmt.Errorf("invalid phone format")
	}

	// Check if phone is already in use
	if phone != "" && (user.Phone == nil || phone != *user.Phone) {
		var existingUser models.User
		if err := database.DB.Where("phone = ?", phone).First(&existingUser).Error; err == nil {
			return nil, fmt.Errorf("phone already in use")
		}
	}

	// Update fields
	updates := map[string]interface{}{}
	if firstName != "" {
		updates["first_name"] = firstName
	}
	if lastName != "" {
		updates["last_name"] = lastName
	}
	if phone != "" {
		updates["phone"] = phone
	}

	if err := database.DB.Model(&user).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	return &user, nil
}

// GetUserOrders retrieves a user's orders
func GetUserOrders(userID uuid.UUID, page, limit int) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	// Get total count
	if err := database.DB.Where("user_id = ?", userID).Model(&models.Order{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count orders: %w", err)
	}

	// Get paginated orders
	if err := database.DB.
		Where("user_id = ?", userID).
		Preload("OrderItems").
		Preload("OrderItems.Product").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&orders).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to fetch orders: %w", err)
	}

	return orders, total, nil
}

// ChangeUserPassword changes a user's password
func ChangeUserPassword(userID uuid.UUID, currentPassword, newPassword string) error {
	var user models.User
	if err := database.DB.First(&user, "id = ?", userID).Error; err != nil {
		return fmt.Errorf("user not found")
	}

	// Verify current password
	if !utils.VerifyPassword(user.PasswordHash, currentPassword) {
		return fmt.Errorf("current password is incorrect")
	}

	// Validate new password
	if !utils.ValidatePassword(newPassword) {
		return fmt.Errorf("new password must be at least 8 characters long")
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password
	if err := database.DB.Model(&user).Update("password_hash", hashedPassword).Error; err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

// DeleteUser deletes a user account
func DeleteUser(userID uuid.UUID) error {
	// Delete user (assuming cascade deletes related data)
	if err := database.DB.Delete(&models.User{}, "id = ?", userID).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
