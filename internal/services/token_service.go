// services/token_service.go
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

// GenerateTokenPair creates access & refresh tokens
func (s *TokenService) GenerateTokenPair(userID string, username string) (*models.TokenPair, error) {
	accessToken, expiresAt, err := utils.GenerateAccessToken(userID, username, s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

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

// RefreshAccessToken validates a refresh token and issues new access token
func (s *TokenService) RefreshAccessToken(refreshToken string) (*models.TokenPair, error) {
	key := fmt.Sprintf("refresh_token:%s", refreshToken)

	data, err := s.redis.Get(key)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	var tokenData models.RefreshTokenData
	if err := json.Unmarshal([]byte(data), &tokenData); err != nil {
		return nil, errors.New("invalid token data")
	}

	return s.GenerateTokenPair(tokenData.UserID, tokenData.Username)
}

// RevokeRefreshToken deletes a single refresh token
func (s *TokenService) RevokeRefreshToken(refreshToken string) error {
	key := fmt.Sprintf("refresh_token:%s", refreshToken)
	return s.redis.Delete(key)
}

// RevokeAllUserTokens deletes all refresh tokens for a user
func (s *TokenService) RevokeAllUserTokens(userID string) error {
	// Implement Redis SCAN or key set per user in production
	return nil
}
