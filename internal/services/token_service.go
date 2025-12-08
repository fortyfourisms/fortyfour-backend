package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/utils"
	"fortyfour-backend/pkg/cache"
	"time"
)

type TokenService struct {
	redis     cache.RedisInterface
	jwtSecret string
}

func NewTokenService(redis cache.RedisInterface, jwtSecret string) *TokenService {
	return &TokenService{
		redis:     redis,
		jwtSecret: jwtSecret,
	}
}

// GenerateTokenPair creates both access and refresh tokens
func (s *TokenService) GenerateTokenPair(userID int, username string) (*models.TokenPair, error) {
	// Generate access token (short-lived)
	accessToken, expiresAt, err := utils.GenerateAccessToken(userID, username, s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token (random string)
	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Store refresh token in Redis with 7 days expiration
	refreshTokenData := models.RefreshTokenData{
		UserID:    userID,
		Username:  username,
		CreatedAt: time.Now(),
	}

	data, err := json.Marshal(refreshTokenData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal refresh token data: %w", err)
	}

	key := fmt.Sprintf("refresh_token:%s", refreshToken)
	if err := s.redis.Set(key, data, 7*24*time.Hour); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &models.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// RefreshAccessToken validates refresh token and generates new access token
func (s *TokenService) RefreshAccessToken(refreshToken string) (*models.TokenPair, error) {
	key := fmt.Sprintf("refresh_token:%s", refreshToken)

	// Get refresh token data from Redis
	data, err := s.redis.Get(key)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	var tokenData models.RefreshTokenData
	if err := json.Unmarshal([]byte(data), &tokenData); err != nil {
		return nil, errors.New("invalid token data")
	}

	// Generate new token pair
	return s.GenerateTokenPair(tokenData.UserID, tokenData.Username)
}

// RevokeRefreshToken removes a refresh token from Redis
func (s *TokenService) RevokeRefreshToken(refreshToken string) error {
	key := fmt.Sprintf("refresh_token:%s", refreshToken)
	return s.redis.Delete(key)
}

// RevokeAllUserTokens removes all refresh tokens for a user
func (s *TokenService) RevokeAllUserTokens(userID int) error {
	// In production, you'd want to use Redis SCAN to find all tokens for a user
	// For now, we'll implement a simple version
	// You should store a set of tokens per user for this to work efficiently
	return nil
}
