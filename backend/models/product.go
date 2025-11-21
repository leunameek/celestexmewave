package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Product struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	StoreID        uuid.UUID      `gorm:"type:uuid;not null;index" json:"store_id"`
	Name           string         `gorm:"type:varchar(255);not null" json:"name"`
	Description    string         `gorm:"type:text" json:"description"`
	Category       string         `gorm:"type:varchar(100)" json:"category"`
	Price          float64        `gorm:"type:decimal(10,2);not null" json:"price"`
	AvailableUnits int            `gorm:"type:integer;default:0" json:"available_units"`
	ImagePath      string         `gorm:"type:varchar(255)" json:"image_path"`
	Sizes          datatypes.JSON `gorm:"type:jsonb" json:"sizes"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`

	Store      Store       `gorm:"foreignKey:StoreID" json:"-"`
	CartItems  []CartItem  `gorm:"foreignKey:ProductID" json:"-"`
	OrderItems []OrderItem `gorm:"foreignKey:ProductID" json:"-"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

func (p *Product) SetSizes(sizes []string) error {
	data, err := json.Marshal(sizes)
	if err != nil {
		return err
	}
	p.Sizes = data
	return nil
}

func (p *Product) GetSizes() ([]string, error) {
	var sizes []string
	err := json.Unmarshal(p.Sizes, &sizes)
	return sizes, err
}
