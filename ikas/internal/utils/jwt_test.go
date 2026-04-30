package utils

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestJWTFlow(t *testing.T) {
	secret := "test-secret"
	userID := "user-123"
	username := "testuser"
	role := "admin"

	// 1. Test GenerateAccessToken
	token, expiresAt, err := GenerateAccessToken(userID, username, role, secret)
	if err != nil {
		t.Fatalf("GenerateAccessToken failed: %v", err)
	}
	if token == "" {
		t.Fatal("Token is empty")
	}
	if expiresAt.Before(time.Now()) {
		t.Fatal("Token already expired")
	}

	// 2. Test VerifyToken (Valid Case)
	claims, err := VerifyToken(token, secret)
	if err != nil {
		t.Fatalf("VerifyToken failed: %v", err)
	}
	if claims.UserID != userID {
		t.Errorf("Expected UserID %v, got %v", userID, claims.UserID)
	}
	if claims.Username != username {
		t.Errorf("Expected Username %v, got %v", username, claims.Username)
	}
	if claims.Role != role {
		t.Errorf("Expected Role %v, got %v", role, claims.Role)
	}

	// 3. Test VerifyToken (Invalid Secret)
	_, err = VerifyToken(token, "wrong-secret")
	if err == nil {
		t.Error("VerifyToken should have failed with wrong secret")
	}

	// 4. Test VerifyToken (Invalid Token)
	_, err = VerifyToken("invalid-token", secret)
	if err == nil {
		t.Error("VerifyToken should have failed with invalid token")
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	token1, err := GenerateRefreshToken()
	if err != nil {
		t.Fatalf("GenerateRefreshToken failed: %v", err)
	}
	if len(token1) != 64 { // 32 bytes hex encoded
		t.Errorf("Expected token length 64, got %d", len(token1))
	}

	token2, err := GenerateRefreshToken()
	if err != nil {
		t.Fatalf("GenerateRefreshToken failed: %v", err)
	}
	if token1 == token2 {
		t.Error("GenerateRefreshToken should produce unique tokens")
	}
}

func TestVerifyTokenEdgeCases(t *testing.T) {
	secret := "test-secret"

	// Test expired token
	claims := Claims{
		UserID: "123",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(secret))

	_, err := VerifyToken(tokenString, secret)
	if err == nil {
		t.Error("VerifyToken should have failed for expired token")
	}

}
