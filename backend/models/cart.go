package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Cart struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    *uuid.UUID `gorm:"type:uuid;index" json:"user_id"`
	SessionID *string    `gorm:"type:varchar(255);index" json:"session_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`

	// Relationships
	User      *User      `gorm:"foreignKey:UserID" json:"-"`
	CartItems []CartItem `gorm:"foreignKey:CartID" json:"-"`
}

// BeforeCreate hook to generate UUID
func (c *Cart) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

type CartItem struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	CartID    uuid.UUID `gorm:"type:uuid;not null;index" json:"cart_id"`
	ProductID uuid.UUID `gorm:"type:uuid;not null;index" json:"product_id"`
	Quantity  int       `gorm:"type:integer;not null;default:1" json:"quantity"`
	Size      string    `gorm:"type:varchar(10)" json:"size"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	Cart    Cart    `gorm:"foreignKey:CartID" json:"-"`
	Product Product `gorm:"foreignKey:ProductID" json:"-"`
}

// BeforeCreate hook to generate UUID
func (ci *CartItem) BeforeCreate(tx *gorm.DB) error {
	if ci.ID == uuid.Nil {
		ci.ID = uuid.New()
	}
	return nil
}
