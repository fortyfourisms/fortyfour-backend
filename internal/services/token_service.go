package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/utils"
	"fortyfour-backend/pkg/cache"
	"net/http"
	"time"
)

type TokenService struct {
	redis        cache.RedisInterface
	JWTSecret    string
	isProduction bool
	domain       string
}

func NewTokenService(redis cache.RedisInterface, jwtSecret string, isProduction bool, domain string) *TokenService {
	return &TokenService{
		redis:        redis,
		JWTSecret:    jwtSecret,
		isProduction: isProduction,
		domain:       domain,
	}
}

// GenerateTokenPair creates access & refresh tokens
func (s *TokenService) GenerateTokenPair(userID, username, role string) (*models.TokenPair, error) {
	accessToken, expiresAt, err := utils.GenerateAccessToken(userID, username, role, s.JWTSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	tokenData := models.RefreshTokenData{
		UserID:    userID,
		Username:  username,
		CreatedAt: time.Now(),
		Role:      role,
	}

	tokenDataJSON, err := json.Marshal(tokenData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal refresh token data: %w", err)
	}

	if err := s.redis.Set(
		fmt.Sprintf("refresh_token:%s", refreshToken),
		string(tokenDataJSON),
		7*24*time.Hour,
	); err != nil {
		return nil, err
	}

	return &models.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// SetAuthCookies sets secure HTTP-only cookies for authentication
func (s *TokenService) SetAuthCookies(w http.ResponseWriter, tokens *models.TokenPair) {
	// Set access token cookie (shorter expiry, 15 minutes)
	accessTokenCookie := &http.Cookie{
		Name:     "access_token",
		Value:    tokens.AccessToken,
		Path:     "/",
		Domain:   s.domain,
		MaxAge:   15 * 60, // 15 minutes
		HttpOnly: true,
		Secure:   s.isProduction, // Only send over HTTPS in production
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, accessTokenCookie)

	// Set refresh token cookie (longer expiry, 7 days)
	refreshTokenCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		Path:     "/api/auth/refresh", // Only send to refresh endpoint
		Domain:   s.domain,
		MaxAge:   7 * 24 * 60 * 60, // 7 days
		HttpOnly: true,
		Secure:   s.isProduction,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, refreshTokenCookie)
}

// GetAccessTokenFromCookie extracts access token from cookie
func (s *TokenService) GetAccessTokenFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("access_token")
	if err != nil {
		return "", errors.New("access token cookie not found")
	}
	return cookie.Value, nil
}

// GetRefreshTokenFromCookie extracts refresh token from cookie
func (s *TokenService) GetRefreshTokenFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		return "", errors.New("refresh token cookie not found")
	}
	return cookie.Value, nil
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

	// Revoke the used refresh token (Refresh Token Rotation)
	if err := s.redis.Delete(key); err != nil {
		// We proceed even if delete fails, though ideally this should be alerted
	}

	return s.GenerateTokenPair(tokenData.UserID, tokenData.Username, tokenData.Role)
}

// RevokeRefreshToken deletes a single refresh token
func (s *TokenService) RevokeRefreshToken(refreshToken string) error {
	key := fmt.Sprintf("refresh_token:%s", refreshToken)
	return s.redis.Delete(key)
}

// RevokeAllUserTokens deletes all refresh tokens for a user
func (s *TokenService) RevokeAllUserTokens(userID string) error {
	pattern := "refresh_token:*"
	keys, err := s.redis.Scan(pattern)
	if err != nil {
		return fmt.Errorf("failed to get keys: %w", err)
	}

	for _, key := range keys {
		data, err := s.redis.Get(key)
		if err != nil {
			continue
		}

		var tokenData models.RefreshTokenData
		if err := json.Unmarshal([]byte(data), &tokenData); err != nil {
			continue
		}

		if tokenData.UserID == userID {
			if err := s.redis.Delete(key); err != nil {
			}
		}
	}

	return nil
}

// ClearAuthCookies removes authentication cookies (for logout)
func (s *TokenService) ClearAuthCookies(w http.ResponseWriter) {
	// Clear access token cookie
	accessTokenCookie := &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		Domain:   s.domain,
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   s.isProduction,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, accessTokenCookie)

	// Clear refresh token cookie
	refreshTokenCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/api/auth/refresh",
		Domain:   s.domain,
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   s.isProduction,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, refreshTokenCookie)
}

// ValidateAndRefreshIfNeeded checks access token and refreshes if expired
func (s *TokenService) ValidateAndRefreshIfNeeded(w http.ResponseWriter, r *http.Request) (*models.TokenClaims, error) {
	// Try to get and validate access token
	accessToken, err := s.GetAccessTokenFromCookie(r)
	if err == nil {
		claims, err := utils.ValidateAccessToken(accessToken, s.JWTSecret)
		if err == nil {
			return claims, nil
		}
	}

	// Access token invalid/expired, try refresh token
	refreshToken, err := s.GetRefreshTokenFromCookie(r)
	if err != nil {
		return nil, errors.New("authentication required")
	}

	// Generate new token pair
	newTokens, err := s.RefreshAccessToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// Set new cookies
	s.SetAuthCookies(w, newTokens)

	// Validate and return new claims
	claims, err := utils.ValidateAccessToken(newTokens.AccessToken, s.JWTSecret)
	if err != nil {
		return nil, err
	}

	return claims, nil
}
