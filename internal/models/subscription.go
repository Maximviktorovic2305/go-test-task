package models

import (
	"time"

	"gorm.io/gorm"
)

type Subscription struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	ServiceName string     `gorm:"not null" json:"service_name"`
	Price       int        `gorm:"not null" json:"price"`
	UserID      string     `gorm:"type:uuid;not null" json:"user_id"`
	StartDate   time.Time  `gorm:"not null" json:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
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
