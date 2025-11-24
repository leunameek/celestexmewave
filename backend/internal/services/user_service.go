package services

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/leunameek/celestexmewave/internal/database"
	"github.com/leunameek/celestexmewave/internal/utils"
	"github.com/leunameek/celestexmewave/models"
	"gorm.io/gorm"
)

// GetUserProfile trae el perfil del user
func GetUserProfile(userID uuid.UUID) (*models.User, error) {
	var user models.User
	if err := database.DB.First(&user, "id = ?", userID).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return &user, nil
}

// UpdateUserProfile actualiza el perfil del user
func UpdateUserProfile(userID uuid.UUID, firstName, lastName, phone string) (*models.User, error) {
	var user models.User
	if err := database.DB.First(&user, "id = ?", userID).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Validamos los datos
	if firstName != "" && !utils.ValidateName(firstName) {
		return nil, fmt.Errorf("invalid first name")
	}

	if lastName != "" && !utils.ValidateName(lastName) {
		return nil, fmt.Errorf("invalid last name")
	}

	if phone != "" && !utils.ValidatePhone(phone) {
		return nil, fmt.Errorf("invalid phone format")
	}

	// Revisamos si el telefono ya esta en uso
	if phone != "" && (user.Phone == nil || phone != *user.Phone) {
		var existingUser models.User
		if err := database.DB.Where("phone = ?", phone).First(&existingUser).Error; err == nil {
			return nil, fmt.Errorf("phone already in use")
		}
	}

	// Campos a actualizar
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

// GetUserOrders trae pedidos del user
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

	// Conteo total
	if err := database.DB.Where("user_id = ?", userID).Model(&models.Order{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count orders: %w", err)
	}

	// Pedidos paginados
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

// ChangeUserPassword cambia la clave del user
func ChangeUserPassword(userID uuid.UUID, currentPassword, newPassword string) error {
	var user models.User
	if err := database.DB.First(&user, "id = ?", userID).Error; err != nil {
		return fmt.Errorf("user not found")
	}

	// Revisamos la clave actual
	if !utils.VerifyPassword(user.PasswordHash, currentPassword) {
		return fmt.Errorf("tu contraseña actual no coincide con la ingresada")
	}

	// Validamos la clave nueva
	if !utils.ValidatePassword(newPassword) {
		return fmt.Errorf("new password must be at least 8 characters long")
	}

	// Hacemos hash de la clave nueva
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Guardamos la nueva clave
	if err := database.DB.Model(&user).Update("password_hash", hashedPassword).Error; err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

// DeleteUser borra la cuenta del user y todo lo relacionado (cascada manual)
func DeleteUser(userID uuid.UUID) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Borrar items de ordenes del usuario (via orders)
		// Primero buscamos los IDs de las ordenes
		var orderIDs []uuid.UUID
		if err := tx.Model(&models.Order{}).Where("user_id = ?", userID).Pluck("id", &orderIDs).Error; err != nil {
			return err
		}

		if len(orderIDs) > 0 {
			if err := tx.Where("order_id IN ?", orderIDs).Delete(&models.OrderItem{}).Error; err != nil {
				return err
			}
		}

		// 2. Borrar ordenes
		if err := tx.Where("user_id = ?", userID).Delete(&models.Order{}).Error; err != nil {
			return err
		}

		// 3. Borrar items del carrito (via carts)
		var cartIDs []uuid.UUID
		if err := tx.Model(&models.Cart{}).Where("user_id = ?", userID).Pluck("id", &cartIDs).Error; err != nil {
			return err
		}

		if len(cartIDs) > 0 {
			if err := tx.Where("cart_id IN ?", cartIDs).Delete(&models.CartItem{}).Error; err != nil {
				return err
			}
		}

		// 4. Borrar carritos
		if err := tx.Where("user_id = ?", userID).Delete(&models.Cart{}).Error; err != nil {
			return err
		}

		// 5. Borrar resets de password
		if err := tx.Where("user_id = ?", userID).Delete(&models.PasswordReset{}).Error; err != nil {
			return err
		}

		// 6. Finalmente borrar el usuario
		if err := tx.Delete(&models.User{}, "id = ?", userID).Error; err != nil {
			return err
		}

		return nil
	})
}
