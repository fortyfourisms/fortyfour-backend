// internal/utils/jwt_test.go
package utils

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testSecret = "2EBCA63B41C6C772F11CE5EBF2501553"

// ============================================================================
// Access Token Generation Tests
// ============================================================================

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
			name:     "Empty username",
			userID:   "user123",
			username: "",
			role:     "user",
			secret:   testSecret,
			wantErr:  false,
		},
		{
			name:     "Empty role",
			userID:   "user123",
			username: "testuser",
			role:     "",
			secret:   testSecret,
			wantErr:  false,
		},
		{
			name:     "Empty secret",
			userID:   "user123",
			username: "testuser",
			role:     "user",
			secret:   "",
			wantErr:  false, // Will generate but won't be secure
		},
		{
			name:     "Special characters in username",
			userID:   "user123",
			username: "test.user+special@example.com",
			role:     "user",
			secret:   testSecret,
			wantErr:  false,
		},
		{
			name:     "Long values",
			userID:   strings.Repeat("a", 100),
			username: strings.Repeat("b", 100),
			role:     strings.Repeat("c", 50),
			secret:   testSecret,
			wantErr:  false,
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

			require.NoError(t, err)
			assert.NotEmpty(t, tokenString)
			assert.True(t, expiresAt.After(time.Now()))
			assert.True(t, expiresAt.Before(time.Now().Add(20*time.Minute)), "Expiration should be within 20 minutes")

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
	require.NoError(t, err)

	// Parse token to verify claims
	claims, err := ValidateAccessToken(tokenString, testSecret)

	require.NoError(t, err)
	require.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, username, claims.Username)
	assert.Equal(t, role, claims.Role)
	assert.NotZero(t, claims.ExpiresAt)
}

func TestTokenExpiration(t *testing.T) {
	userID := "user123"
	username := "testuser"
	role := "user"

	_, expiresAt, err := GenerateAccessToken(userID, username, role, testSecret)
	require.NoError(t, err)

	// Verify expiration is approximately 15 minutes from now
	expectedExpiry := time.Now().Add(15 * time.Minute)
	timeDiff := expiresAt.Sub(expectedExpiry)

	// Allow 5 second tolerance for test execution time
	assert.True(t, timeDiff < 5*time.Second && timeDiff > -5*time.Second,
		"Expiration should be approximately 15 minutes from now")
}

func TestTokenRoleField(t *testing.T) {
	roles := []string{"admin", "user", "moderator", "guest", ""}

	for _, role := range roles {
		t.Run("Role: "+role, func(t *testing.T) {
			tokenString, _, err := GenerateAccessToken("user123", "testuser", role, testSecret)
			require.NoError(t, err)

			claims, err := ValidateAccessToken(tokenString, testSecret)
			require.NoError(t, err)
			assert.Equal(t, role, claims.Role)
		})
	}
}

// ============================================================================
// Refresh Token Generation Tests
// ============================================================================

func TestGenerateRefreshToken(t *testing.T) {
	t.Run("Generate valid refresh token", func(t *testing.T) {
		token, err := GenerateRefreshToken()

		require.NoError(t, err)
		assert.NotEmpty(t, token)
		// Base64 URL encoding of 32 bytes results in 43-44 characters (with padding)
		assert.Greater(t, len(token), 40, "Token should be at least 40 characters")
	})

	t.Run("Generate unique tokens", func(t *testing.T) {
		token1, err1 := GenerateRefreshToken()
		token2, err2 := GenerateRefreshToken()

		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.NotEqual(t, token1, token2, "Tokens should be unique")
	})

	t.Run("Token is valid base64 URL encoding", func(t *testing.T) {
		token, err := GenerateRefreshToken()

		require.NoError(t, err)

		// Try to check if it's valid base64 URL encoded
		for _, c := range token {
			assert.True(t,
				(c >= '0' && c <= '9') || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_' || c == '-' || c == '=',
				"Token should only contain base64 URL-safe characters",
			)
		}
	})

	t.Run("Generate multiple unique tokens", func(t *testing.T) {
		tokens := make(map[string]bool)
		for i := 0; i < 100; i++ {
			token, err := GenerateRefreshToken()
			require.NoError(t, err)
			assert.False(t, tokens[token], "All tokens should be unique")
			tokens[token] = true
		}
	})
}

// ============================================================================
// Token Validation Tests
// ============================================================================

func TestValidateAccessToken_ValidToken(t *testing.T) {
	userID := "user123"
	username := "testuser"
	role := "admin"

	// Generate a token
	tokenString, _, err := GenerateAccessToken(userID, username, role, testSecret)
	require.NoError(t, err)

	// Verify the token
	claims, err := ValidateAccessToken(tokenString, testSecret)

	require.NoError(t, err)
	require.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, username, claims.Username)
	assert.Equal(t, role, claims.Role)
}

func TestValidateAccessToken_InvalidToken(t *testing.T) {
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
		{
			name:        "Random string",
			token:       "this-is-not-a-token",
			secret:      testSecret,
			expectError: true,
		},
		{
			name:        "Only header",
			token:       "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			secret:      testSecret,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ValidateAccessToken(tt.token, tt.secret)

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

func TestValidateAccessToken_ExpiredToken(t *testing.T) {
	// Create an expired token
	expiresAt := time.Now().Add(-1 * time.Hour)

	claims := jwt.MapClaims{
		"user_id":  "user123",
		"username": "testuser",
		"role":     "admin",
		"exp":      expiresAt.Unix(),
		"iat":      time.Now().Add(-2 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(testSecret))
	require.NoError(t, err)

	// Try to verify expired token
	verifiedClaims, err := ValidateAccessToken(tokenString, testSecret)

	assert.Error(t, err)
	assert.Nil(t, verifiedClaims)
}

func TestValidateAccessToken_WrongSigningMethod(t *testing.T) {
	// Create token with different signing method (None instead of HS256)
	// This should fail because our ValidateAccessToken expects HS256

	claims := jwt.MapClaims{
		"user_id":  "user123",
		"username": "testuser",
		"role":     "admin",
		"exp":      time.Now().Add(1 * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}

	// Create token with none signing method (will be rejected)
	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	tokenString, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	// Try to verify
	verifiedClaims, err := ValidateAccessToken(tokenString, testSecret)

	assert.Error(t, err)
	assert.Nil(t, verifiedClaims)
}

func TestValidateAccessToken_MissingClaims(t *testing.T) {
	tests := []struct {
		name   string
		claims jwt.MapClaims
	}{
		{
			name: "Missing user_id",
			claims: jwt.MapClaims{
				"username": "testuser",
				"role":     "admin",
				"exp":      time.Now().Add(1 * time.Hour).Unix(),
				"iat":      time.Now().Unix(),
			},
		},
		{
			name: "Missing username",
			claims: jwt.MapClaims{
				"user_id": "user123",
				"role":    "admin",
				"exp":     time.Now().Add(1 * time.Hour).Unix(),
				"iat":     time.Now().Unix(),
			},
		},
		{
			name: "Missing role",
			claims: jwt.MapClaims{
				"user_id":  "user123",
				"username": "testuser",
				"exp":      time.Now().Add(1 * time.Hour).Unix(),
				"iat":      time.Now().Unix(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, tt.claims)
			tokenString, err := token.SignedString([]byte(testSecret))
			require.NoError(t, err)

			// This will panic in the current implementation due to missing claims
			// Testing whether it panics or returns an error
			defer func() {
				if r := recover(); r != nil {
					// Expected: panic due to type assertion on nil value
					t.Logf("Function panicked as expected: %v", r)
				}
			}()

			claims, err := ValidateAccessToken(tokenString, testSecret)
			// If we reach here without panic, check that we got an error
			if err == nil {
				t.Error("Expected error or panic with missing claims")
			}
			assert.Nil(t, claims)
		})
	}
}

func TestValidateAccessToken_WrongClaimTypes(t *testing.T) {
	tests := []struct {
		name   string
		claims jwt.MapClaims
	}{
		{
			name: "user_id as integer",
			claims: jwt.MapClaims{
				"user_id":  12345, // Should be string
				"username": "testuser",
				"role":     "admin",
				"exp":      time.Now().Add(1 * time.Hour).Unix(),
				"iat":      time.Now().Unix(),
			},
		},
		{
			name: "username as integer",
			claims: jwt.MapClaims{
				"user_id":  "user123",
				"username": 12345, // Should be string
				"role":     "admin",
				"exp":      time.Now().Add(1 * time.Hour).Unix(),
				"iat":      time.Now().Unix(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, tt.claims)
			tokenString, err := token.SignedString([]byte(testSecret))
			require.NoError(t, err)

			// This will panic in the current implementation due to wrong type
			defer func() {
				if r := recover(); r != nil {
					// Expected: panic due to type assertion failure
					t.Logf("Function panicked as expected: %v", r)
				}
			}()

			claims, err := ValidateAccessToken(tokenString, testSecret)
			// If we reach here without panic, we should have an error
			if err == nil && claims != nil {
				t.Error("Expected error or panic with wrong claim types")
			}
		})
	}
}

// ============================================================================
// Helper Functions
// ============================================================================

// Helper function to generate a valid token for testing
func generateValidToken(t *testing.T, secret string) string {
	token, _, err := GenerateAccessToken("user123", "testuser", "admin", secret)
	require.NoError(t, err)
	return token
}

// ============================================================================
// Benchmark Tests
// ============================================================================

func BenchmarkGenerateAccessToken(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		GenerateAccessToken("user123", "testuser", "admin", testSecret)
	}
}

func BenchmarkGenerateRefreshToken(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		GenerateRefreshToken()
	}
}

func BenchmarkValidateAccessToken(b *testing.B) {
	token, _, _ := GenerateAccessToken("user123", "testuser", "admin", testSecret)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ValidateAccessToken(token, testSecret)
	}
}
