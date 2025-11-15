package models

import (
	"time"

	"github.com/google/uuid"
)

// UserRole represents the user types in the system
type UserRole string

const (
	RoleBusinessOwner UserRole = "BO"
	RoleCustomer      UserRole = "CU"
	RoleAdmin         UserRole = "AD"
)

// User represents a user in the system
type User struct {
	ID           uuid.UUID  `json:"id"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"` // Do not expose in JSON
	Phone        string     `json:"phone"`
	Roles        []UserRole `json:"roles"`
	IsActive     bool       `json:"is_active"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// CreateUserRequest represents the request to create a user
type CreateUserRequest struct {
	Email    string     `json:"email" binding:"required,email"`
	Password string     `json:"password" binding:"required,min=6"`
	Phone    string     `json:"phone" binding:"required"`
	Roles    []UserRole `json:"roles" binding:"required,min=1,dive,oneof=BO CU AD"`
}

// UserResponse represents the response with user data
type UserResponse struct {
	ID        uuid.UUID  `json:"id"`
	Email     string     `json:"email"`
	Phone     string     `json:"phone"`
	Roles     []UserRole `json:"roles"`
	IsActive  bool       `json:"is_active"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// HasRole checks if the user has a specific role
func (u *User) HasRole(role UserRole) bool {
	for _, r := range u.Roles {
		if r == role {
			return true
		}
	}
	return false
}

// ToResponse converts a User to UserResponse
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Phone:     u.Phone,
		Roles:     u.Roles,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
