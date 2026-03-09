package models

import "time"

type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type RefreshTokenData struct {
	UserID       string    `json:"user_id"`
	Username     string    `json:"username"`
	CreatedAt    time.Time `json:"created_at"`
	Role         string    `json:"role"`
	IDPerusahaan string    `json:"id_perusahaan"`
}

type TokenClaims struct {
	UserID       string `json:"user_id"`
	Username     string `json:"username"`
	Role         string `json:"role"`
	IDPerusahaan string `json:"id_perusahaan"`
	ExpiresAt    int64  `json:"exp"`
}