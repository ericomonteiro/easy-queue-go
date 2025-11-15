package models

import (
	"time"

	"github.com/google/uuid"
)

// UserRole representa os tipos de usuário no sistema
type UserRole string

const (
	RoleBusinessOwner UserRole = "BO"
	RoleCustomer      UserRole = "CU"
)

// User representa um usuário no sistema
type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Não expor no JSON
	Phone        string    `json:"phone"`
	Role         UserRole  `json:"role"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CreateUserRequest representa a requisição para criar um usuário
type CreateUserRequest struct {
	Email    string   `json:"email" binding:"required,email"`
	Password string   `json:"password" binding:"required,min=6"`
	Phone    string   `json:"phone" binding:"required"`
	Role     UserRole `json:"role" binding:"required,oneof=BO CU"`
}

// UserResponse representa a resposta com dados do usuário
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Role      UserRole  `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToResponse converte um User para UserResponse
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Phone:     u.Phone,
		Role:      u.Role,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
