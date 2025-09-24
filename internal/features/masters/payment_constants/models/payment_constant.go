package models

import (
	"time"

	"github.com/google/uuid"
)

// PaymentConstant represents payment constants in the system
type PaymentConstant struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"size:100;not null;unique"`
	Value       int       `json:"value" gorm:"not null"`
	Description *string   `json:"description" gorm:"size:255"`
	IsActive    bool      `json:"is_active" gorm:"not null;default:true"`
	CreatedAt   time.Time `json:"created_at" gorm:"not null;type:timestamptz"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"not null;type:timestamptz"`
}

// TableName returns the table name for the PaymentConstant model
func (PaymentConstant) TableName() string {
	return "payment_constants"
}
