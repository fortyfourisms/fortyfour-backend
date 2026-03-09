package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fortyfour-backend/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateAccessToken(userID, username, role, secret, idPerusahaan string) (string, time.Time, error) {
	expiresAt := time.Now().Add(15 * time.Minute)

	claims := jwt.MapClaims{
		"user_id":       userID,
		"username":      username,
		"role":          role,
		"id_perusahaan": idPerusahaan,
		"exp":           expiresAt.Unix(),
		"iat":           time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// ValidateAccessToken validates a JWT token and returns the claims.
// This version includes safe type assertions to prevent panics from malformed tokens.
func ValidateAccessToken(tokenString, secret string) (*models.TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// Safe type assertions with validation to prevent panics
	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, errors.New("invalid user_id claim")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return nil, errors.New("invalid username claim")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return nil, errors.New("invalid role claim")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, errors.New("invalid exp claim")
	}

	// id_perusahaan opsional, default ke "" jika tidak ada
	idPerusahaan, _ := claims["id_perusahaan"].(string)

	return &models.TokenClaims{
		UserID:       userID,
		Username:     username,
		Role:         role,
		IDPerusahaan: idPerusahaan,
		ExpiresAt:    int64(exp),
	}, nil
}