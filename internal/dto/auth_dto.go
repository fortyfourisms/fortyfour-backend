package dto

import "fortyfour-backend/internal/models"

type RegisterRequest struct {
	Username  string  `json:"username" validate:"required,min=3,max=50"`
	Password  string  `json:"password" validate:"required,min=6"`
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
