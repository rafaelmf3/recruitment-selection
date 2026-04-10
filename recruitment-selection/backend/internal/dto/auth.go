package dto

import "recruitment-selection/internal/model"

// RegisterRequest is the payload for POST /api/v1/auth/register.
type RegisterRequest struct {
	Name     string          `json:"name"     binding:"required,min=2,max=255"`
	Email    string          `json:"email"    binding:"required,email"`
	Password string          `json:"password" binding:"required,min=8"`
	Role     model.UserRole  `json:"role"     binding:"required,oneof=recruiter candidate"`
}

// LoginRequest is the payload for POST /api/v1/auth/login.
type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UserResponse is the public representation of a user (no password hash).
type UserResponse struct {
	ID    string         `json:"id"`
	Name  string         `json:"name"`
	Email string         `json:"email"`
	Role  model.UserRole `json:"role"`
}

// LoginResponse is the body returned on a successful login.
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}
