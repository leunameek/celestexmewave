package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	UserID        *uuid.UUID `gorm:"type:uuid;index" json:"user_id"`
	SessionID     *string    `gorm:"type:varchar(255);index" json:"session_id"`
	TotalAmount   float64    `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	Status        string     `gorm:"type:varchar(50);default:'pending'" json:"status"`         // pending, confirmed, shipped, delivered
	PaymentStatus string     `gorm:"type:varchar(50);default:'pending'" json:"payment_status"` // pending, completed, failed
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`

	// Relationships
	User       *User       `gorm:"foreignKey:UserID" json:"-"`
	OrderItems []OrderItem `gorm:"foreignKey:OrderID" json:"-"`
}

// BeforeCreate hook to generate UUID
func (o *Order) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}

type OrderItem struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	OrderID   uuid.UUID `gorm:"type:uuid;not null;index" json:"order_id"`
	ProductID uuid.UUID `gorm:"type:uuid;not null;index" json:"product_id"`
	Quantity  int       `gorm:"type:integer;not null" json:"quantity"`
	Size      string    `gorm:"type:varchar(10)" json:"size"`
	UnitPrice float64   `gorm:"type:decimal(10,2);not null" json:"unit_price"`
	CreatedAt time.Time `json:"created_at"`

	// Relationships
	Order   Order   `gorm:"foreignKey:OrderID" json:"-"`
	Product Product `gorm:"foreignKey:ProductID" json:"-"`
}

// BeforeCreate hook to generate UUID
func (oi *OrderItem) BeforeCreate(tx *gorm.DB) error {
	if oi.ID == uuid.Nil {
		oi.ID = uuid.New()
	}
	return nil
}
