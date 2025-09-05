package models

import (
	"time"
)

// swagger:model
type Subscription struct {
	// The unique identifier of the subscription
	// Read Only: true
	// Example: 1
	ID uint `gorm:"primaryKey" json:"id"`

	// The name of the service
	// Required: true
	// Example: Netflix
	ServiceName string `gorm:"not null" json:"service_name"`

	// The price of the subscription in rubles
	// Required: true
	// Minimum: 0
	// Example: 990
	Price int `gorm:"not null" json:"price"`

	// The UUID of the user
	// Required: true
	// Example: 550e8400-e29b-41d4-a716-446655440000
	UserID string `gorm:"type:uuid;not null" json:"user_id"`

	// The start date of the subscription
	// Required: true
	// Example: 2023-01-01T00:00:00Z
	StartDate time.Time `gorm:"not null" json:"start_date"`

	// The end date of the subscription
	// Example: 2023-12-31T00:00:00Z
	EndDate *time.Time `json:"end_date,omitempty"`

	// The creation timestamp
	// Read Only: true
	// Example: 2023-01-01T00:00:00Z
	CreatedAt time.Time `json:"created_at"`

	// The last update timestamp
	// Read Only: true
	// Example: 2023-01-01T00:00:00Z
	UpdatedAt time.Time `json:"updated_at"`
}