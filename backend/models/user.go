package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Email        *string   `gorm:"type:varchar(255);uniqueIndex" json:"email"`
	Phone        *string   `gorm:"type:varchar(20);uniqueIndex" json:"phone"`
	FirstName    string    `gorm:"type:varchar(255);not null" json:"first_name"`
	LastName     string    `gorm:"type:varchar(255);not null" json:"last_name"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"`
	IsRegistered bool      `gorm:"default:true" json:"is_registered"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	Carts          []Cart          `gorm:"foreignKey:UserID" json:"-"`
	Orders         []Order         `gorm:"foreignKey:UserID" json:"-"`
	PasswordResets []PasswordReset `gorm:"foreignKey:UserID" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

type PasswordReset struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	ResetCode string    `gorm:"type:varchar(10);uniqueIndex;not null" json:"reset_code"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	Used      bool      `gorm:"default:false" json:"used"`
	CreatedAt time.Time `json:"created_at"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}

func (pr *PasswordReset) BeforeCreate(tx *gorm.DB) error {
	if pr.ID == uuid.Nil {
		pr.ID = uuid.New()
	}
	return nil
}
