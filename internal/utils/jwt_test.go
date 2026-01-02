// internal/utils/jwt_test.go
package utils

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

const testSecret = "2EBCA63B41C6C772F11CE5EBF2501553"

func TestGenerateAccessToken(t *testing.T) {
	tests := []struct {
		name     string
		userID   string
		username string
		role     string
		secret   string
		wantErr  bool
	}{
		{
			name:     "Valid token generation",
			userID:   "user123",
			username: "testuser",
			role:     "admin",
			secret:   testSecret,
			wantErr:  false,
		},
		{
			name:     "Empty user ID",
			userID:   "",
			username: "testuser",
			role:     "user",
			secret:   testSecret,
			wantErr:  false, // Should still generate token
		},
		{
			name:     "Empty secret",
			userID:   "user123",
			username: "testuser",
			role:     "user",
			secret:   "",
			wantErr:  false, // Will generate but won't be secure
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenString, expiresAt, err := GenerateAccessToken(
				tt.userID,
				tt.username,
				tt.role,
				tt.secret,
			)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, tokenString)
			assert.True(t, expiresAt.After(time.Now()))
			assert.True(t, expiresAt.Before(time.Now().Add(2*time.Hour)))

			// Verify token structure (should have 3 parts separated by dots)
			parts := strings.Split(tokenString, ".")
			assert.Equal(t, 3, len(parts), "JWT should have 3 parts")
		})
	}
}

func TestGenerateAccessToken_Claims(t *testing.T) {
	userID := "user123"
	username := "testuser"
	role := "admin"

	tokenString, _, err := GenerateAccessToken(userID, username, role, testSecret)
	assert.NoError(t, err)

	// Parse token to verify claims
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(testSecret), nil
	})

	assert.NoError(t, err)
	assert.True(t, token.Valid)

	claims, ok := token.Claims.(*Claims)
	assert.True(t, ok)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, username, claims.Username)
	assert.Equal(t, role, claims.Role)
	assert.NotNil(t, claims.ExpiresAt)
	assert.NotNil(t, claims.IssuedAt)
}

func TestGenerateRefreshToken(t *testing.T) {
	t.Run("Generate valid refresh token", func(t *testing.T) {
		token, err := GenerateRefreshToken()

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.Equal(t, 64, len(token), "Hex encoded 32 bytes should be 64 characters")
	})

	t.Run("Generate unique tokens", func(t *testing.T) {
		token1, err1 := GenerateRefreshToken()
		token2, err2 := GenerateRefreshToken()

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NotEqual(t, token1, token2, "Tokens should be unique")
	})

	t.Run("Token is valid hex", func(t *testing.T) {
		token, err := GenerateRefreshToken()

		assert.NoError(t, err)

		// Try to decode hex - should not error
		for _, c := range token {
			assert.True(t,
				(c >= '0' && c <= '9') || (c >= 'a' && c <= 'f'),
				"Token should only contain hex characters",
			)
		}
	})
}

func TestVerifyToken_ValidToken(t *testing.T) {
	userID := "user123"
	username := "testuser"
	role := "admin"

	// Generate a token
	tokenString, _, err := GenerateAccessToken(userID, username, role, testSecret)
	assert.NoError(t, err)

	// Verify the token
	claims, err := VerifyToken(tokenString, testSecret)

	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, username, claims.Username)
	assert.Equal(t, role, claims.Role)
}

func TestVerifyToken_InvalidToken(t *testing.T) {
	tests := []struct {
		name        string
		token       string
		secret      string
		expectError bool
	}{
		{
			name:        "Invalid token format",
			token:       "invalid.token.format",
			secret:      testSecret,
			expectError: true,
		},
		{
			name:        "Empty token",
			token:       "",
			secret:      testSecret,
			expectError: true,
		},
		{
			name:        "Wrong secret",
			token:       generateValidToken(t, testSecret),
			secret:      "wrong-secret",
			expectError: true,
		},
		{
			name:        "Malformed token",
			token:       "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.malformed",
			secret:      testSecret,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := VerifyToken(tt.token, tt.secret)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
			}
		})
	}
}

func TestVerifyToken_ExpiredToken(t *testing.T) {
	// Create an expired token
	expiresAt := time.Now().Add(-1 * time.Hour) // Already expired

	claims := Claims{
		UserID:   "user123",
		Username: "testuser",
		Role:     "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(testSecret))
	assert.NoError(t, err)

	// Try to verify expired token
	verifiedClaims, err := VerifyToken(tokenString, testSecret)

	assert.Error(t, err)
	assert.Nil(t, verifiedClaims)
	assert.Contains(t, err.Error(), "token is expired")
}

func TestVerifyToken_WrongSigningMethod(t *testing.T) {
	// Create token with different signing method (RS256 instead of HS256)
	// This should fail because our VerifyToken expects HS256

	claims := Claims{
		UserID:   "user123",
		Username: "testuser",
		Role:     "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Create token with none signing method (will be rejected)
	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	tokenString, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	assert.NoError(t, err)

	// Try to verify
	verifiedClaims, err := VerifyToken(tokenString, testSecret)

	assert.Error(t, err)
	assert.Nil(t, verifiedClaims)
}

func TestTokenExpiration(t *testing.T) {
	userID := "user123"
	username := "testuser"
	role := "user"

	_, expiresAt, err := GenerateAccessToken(userID, username, role, testSecret)
	assert.NoError(t, err)

	// Verify expiration is approximately 1 hour from now
	expectedExpiry := time.Now().Add(1 * time.Hour)
	timeDiff := expiresAt.Sub(expectedExpiry)

	// Allow 5 second tolerance for test execution time
	assert.True(t, timeDiff < 5*time.Second && timeDiff > -5*time.Second,
		"Expiration should be approximately 1 hour from now")
}

func TestTokenRoleField(t *testing.T) {
	roles := []string{"admin", "user", "moderator", ""}

	for _, role := range roles {
		t.Run("Role: "+role, func(t *testing.T) {
			tokenString, _, err := GenerateAccessToken("user123", "testuser", role, testSecret)
			assert.NoError(t, err)

			claims, err := VerifyToken(tokenString, testSecret)
			assert.NoError(t, err)
			assert.Equal(t, role, claims.Role)
		})
	}
}

// Helper function to generate a valid token for testing
func generateValidToken(t *testing.T, secret string) string {
	token, _, err := GenerateAccessToken("user123", "testuser", "admin", secret)
	assert.NoError(t, err)
	return token
}

// Benchmark tests
func BenchmarkGenerateAccessToken(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateAccessToken("user123", "testuser", "admin", testSecret)
	}
}

func BenchmarkGenerateRefreshToken(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateRefreshToken()
	}
}

func BenchmarkVerifyToken(b *testing.B) {
	token, _, _ := GenerateAccessToken("user123", "testuser", "admin", testSecret)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		VerifyToken(token, testSecret)
	}
}
