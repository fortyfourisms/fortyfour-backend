package utils

import (
	"testing"
	"time"
)

func TestGenerateToken_Success(t *testing.T) {
	userID := "1"
	username := "testuser"
	secret := "test-secret"
	roles := "user roles"
	expiresAt := time.Now().Add(15 * time.Minute)

	token, expiry, err := GenerateAccessToken(userID, username, secret, roles)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if expiry != expiresAt {
		t.Error("expiry time is not suitable")
	}

	if token == "" {
		t.Error("expected token to be generated")
	}
}

func TestVerifyToken_Success(t *testing.T) {
	userID := "1"
	username := "testuser"
	secret := "test-secret"
	roles := "user roles"

	token, _, _ := GenerateAccessToken(userID, username, secret, roles)

	claims, err := VerifyToken(token, secret)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("expected UserID %d, got %d", userID, claims.UserID)
	}

	if claims.Username != username {
		t.Errorf("expected Username '%s', got '%s'", username, claims.Username)
	}
}

func TestVerifyToken_InvalidToken(t *testing.T) {
	secret := "test-secret"
	invalidToken := "invalid.token.here"

	_, err := VerifyToken(invalidToken, secret)

	if err == nil {
		t.Fatal("expected error for invalid token")
	}
}

func TestVerifyToken_WrongSecret(t *testing.T) {
	userID := "1"
	username := "testuser"
	secret := "test-secret"
	roles := "user roles"
	wrongSecret := "wrong-secret"

	token, _, _ := GenerateAccessToken(userID, username, secret, roles)

	_, err := VerifyToken(token, wrongSecret)

	if err == nil {
		t.Fatal("expected error for wrong secret")
	}
}

// func TestVerifyToken_ExpiredToken(t *testing.T) {
// 	userID := 1
// 	username := "testuser"
// 	secret := "test-secret"
// 	expiry := -1 * time.Hour

// 	token, _ := GenerateAccessToken(userID, username, secret)

// 	time.Sleep(10 * time.Millisecond)

// 	_, err := VerifyToken(token, secret)

// 	if err == nil {
// 		t.Fatal("expected error for expired token")
// 	}
// }
