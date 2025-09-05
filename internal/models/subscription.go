package models

import (
	"time"

	"gorm.io/gorm"
)

// Subscription represents a user subscription to a service
//
// swagger:model
type Subscription struct {
	// The unique identifier of the subscription
	// Read Only: true
	ID uint `gorm:"primaryKey" json:"id"`

	// The name of the service
	// Required: true
	ServiceName string `gorm:"not null" json:"service_name"`

	// The price of the subscription in rubles
	// Required: true
	// Minimum: 0
	Price int `gorm:"not null" json:"price"`

	// The UUID of the user
	// Required: true
	UserID string `gorm:"type:uuid;not null" json:"user_id"`

	// The start date of the subscription
	// Required: true
	StartDate time.Time `gorm:"not null" json:"start_date"`

	// The end date of the subscription
	EndDate *time.Time `json:"end_date,omitempty"`

	// The creation timestamp
	// Read Only: true
	CreatedAt time.Time `json:"created_at"`

	// The last update timestamp
	// Read Only: true
	UpdatedAt time.Time `json:"updated_at"`
}

// BeforeCreate хук для установки CreatedAt и UpdatedAt
func (s *Subscription) BeforeCreate(tx *gorm.DB) error {
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate хук для установки UpdatedAt
func (s *Subscription) BeforeUpdate(tx *gorm.DB) error {
	s.UpdatedAt = time.Now()
	return nil
}
