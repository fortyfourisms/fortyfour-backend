package dto

import "fortyfour-backend/internal/models"

type RegisterRequest struct {
	Username  string  `json:"username" validate:"required,min=3,max=50"`
	Password  string  `json:"password" validate:"required,min=8"`
	Email     string  `json:"email" validate:"required,email"`
	RoleID    *string `json:"role_id,omitempty"`
	IDJabatan *string `json:"id_jabatan,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type AuthResponse struct {
	User         models.User `json:"user"`
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	ExpiresAt    string      `json:"expires_at"`
}

type ErrorResponse struct {
	Message string `json:"message" example:"invalid credentials"`
}

type MessageResponse struct {
	Message string `json:"message" example:"Logged out successfully"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	ExpiresAt    string `json:"expires_at" example:"2025-12-15T15:04:05+07:00"`
}
