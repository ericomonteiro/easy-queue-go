package models

import (
	"time"

	"github.com/google/uuid"
)

// Business represents a business in the system
type Business struct {
	ID          uuid.UUID `json:"id"`
	OwnerID     uuid.UUID `json:"owner_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Address     string    `json:"address"`
	Phone       string    `json:"phone"`
	Email       string    `json:"email"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateBusinessRequest represents the request to create a business
type CreateBusinessRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=255"`
	Description string `json:"description" binding:"max=1000"`
	Address     string `json:"address" binding:"max=500"`
	Phone       string `json:"phone" binding:"required,min=10,max=50"`
	Email       string `json:"email" binding:"omitempty,email"`
}

// UpdateBusinessRequest represents the request to update a business
type UpdateBusinessRequest struct {
	Name        string `json:"name" binding:"omitempty,min=3,max=255"`
	Description string `json:"description" binding:"omitempty,max=1000"`
	Address     string `json:"address" binding:"omitempty,max=500"`
	Phone       string `json:"phone" binding:"omitempty,min=10,max=50"`
	Email       string `json:"email" binding:"omitempty,email"`
	IsActive    *bool  `json:"is_active" binding:"omitempty"`
}

// BusinessResponse represents the response with business data
type BusinessResponse struct {
	ID          uuid.UUID `json:"id"`
	OwnerID     uuid.UUID `json:"owner_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Address     string    `json:"address"`
	Phone       string    `json:"phone"`
	Email       string    `json:"email"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToResponse converts a Business to BusinessResponse
func (b *Business) ToResponse() *BusinessResponse {
	return &BusinessResponse{
		ID:          b.ID,
		OwnerID:     b.OwnerID,
		Name:        b.Name,
		Description: b.Description,
		Address:     b.Address,
		Phone:       b.Phone,
		Email:       b.Email,
		IsActive:    b.IsActive,
		CreatedAt:   b.CreatedAt,
		UpdatedAt:   b.UpdatedAt,
	}
}
