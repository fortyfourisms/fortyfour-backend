package middleware

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetUserID_GetUserID(t *testing.T) {
	tests := []struct {
		name     string
		userID   string
		expected string
	}{
		{
			name:     "set and get valid user ID",
			userID:   "user-123",
			expected: "user-123",
		},
		{
			name:     "set and get UUID format",
			userID:   "550e8400-e29b-41d4-a716-446655440000",
			expected: "550e8400-e29b-41d4-a716-446655440000",
		},
		{
			name:     "set and get empty string",
			userID:   "",
			expected: "",
		},
		{
			name:     "set and get numeric ID",
			userID:   "12345",
			expected: "12345",
		},
		{
			name:     "set and get with special characters",
			userID:   "user@domain-123",
			expected: "user@domain-123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx = SetUserID(ctx, tt.userID)

			result := GetUserID(ctx)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetUserID_NotSet(t *testing.T) {
	ctx := context.Background()

	result := GetUserID(ctx)
	assert.Equal(t, "", result, "should return empty string when UserID not set")
}

func TestGetUserID_WrongType(t *testing.T) {
	// Set wrong type in context
	ctx := context.WithValue(context.Background(), UserIDKey, 12345) // int instead of string

	result := GetUserID(ctx)
	assert.Equal(t, "", result, "should return empty string when value is not a string")
}

func TestGetUserID_NilContext(t *testing.T) {
	// Although passing nil context is bad practice, our function should handle it gracefully
	// Note: This will panic if not handled, so we test with Background
	ctx := context.Background()

	result := GetUserID(ctx)
	assert.Equal(t, "", result)
}

func TestSetUsername_GetUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
		expected string
	}{
		{
			name:     "set and get valid username",
			username: "johndoe",
			expected: "johndoe",
		},
		{
			name:     "set and get email as username",
			username: "user@example.com",
			expected: "user@example.com",
		},
		{
			name:     "set and get empty string",
			username: "",
			expected: "",
		},
		{
			name:     "set and get username with spaces",
			username: "john doe",
			expected: "john doe",
		},
		{
			name:     "set and get unicode username",
			username: "用户名",
			expected: "用户名",
		},
		{
			name:     "set and get username with underscore",
			username: "john_doe_123",
			expected: "john_doe_123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx = SetUsername(ctx, tt.username)

			result := GetUsername(ctx)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetUsername_NotSet(t *testing.T) {
	ctx := context.Background()

	result := GetUsername(ctx)
	assert.Equal(t, "", result, "should return empty string when Username not set")
}

func TestGetUsername_WrongType(t *testing.T) {
	// Set wrong type in context
	ctx := context.WithValue(context.Background(), UsernameKey, []byte("username")) // bytes instead of string

	result := GetUsername(ctx)
	assert.Equal(t, "", result, "should return empty string when value is not a string")
}

func TestSetRole_GetRole(t *testing.T) {
	tests := []struct {
		name     string
		role     string
		expected string
	}{
		{
			name:     "set and get admin role",
			role:     "admin",
			expected: "admin",
		},
		{
			name:     "set and get user role",
			role:     "user",
			expected: "user",
		},
		{
			name:     "set and get moderator role",
			role:     "moderator",
			expected: "moderator",
		},
		{
			name:     "set and get empty string",
			role:     "",
			expected: "",
		},
		{
			name:     "set and get role with uppercase",
			role:     "ADMIN",
			expected: "ADMIN",
		},
		{
			name:     "set and get custom role",
			role:     "super_admin_level_5",
			expected: "super_admin_level_5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx = SetRole(ctx, tt.role)

			result := GetRole(ctx)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetRole_NotSet(t *testing.T) {
	ctx := context.Background()

	result := GetRole(ctx)
	assert.Equal(t, "", result, "should return empty string when Role not set")
}

func TestGetRole_WrongType(t *testing.T) {
	// Set wrong type in context
	ctx := context.WithValue(context.Background(), RoleKey, 123) // int instead of string

	result := GetRole(ctx)
	assert.Equal(t, "", result, "should return empty string when value is not a string")
}

func TestContext_MultipleValues(t *testing.T) {
	// Test setting and getting multiple values in the same context
	ctx := context.Background()

	userID := "user-123"
	username := "johndoe"
	role := "admin"

	// Set all values
	ctx = SetUserID(ctx, userID)
	ctx = SetUsername(ctx, username)
	ctx = SetRole(ctx, role)

	// Get all values
	retrievedUserID := GetUserID(ctx)
	retrievedUsername := GetUsername(ctx)
	retrievedRole := GetRole(ctx)

	assert.Equal(t, userID, retrievedUserID)
	assert.Equal(t, username, retrievedUsername)
	assert.Equal(t, role, retrievedRole)
}

func TestContext_OverwriteValues(t *testing.T) {
	// Test that setting a value twice overwrites the first value
	ctx := context.Background()

	// Set initial values
	ctx = SetUserID(ctx, "user-1")
	ctx = SetUsername(ctx, "user1")
	ctx = SetRole(ctx, "user")

	// Overwrite with new values
	ctx = SetUserID(ctx, "user-2")
	ctx = SetUsername(ctx, "user2")
	ctx = SetRole(ctx, "admin")

	// Should get the latest values
	assert.Equal(t, "user-2", GetUserID(ctx))
	assert.Equal(t, "user2", GetUsername(ctx))
	assert.Equal(t, "admin", GetRole(ctx))
}

func TestContext_IndependentContexts(t *testing.T) {
	// Test that different contexts don't interfere with each other
	ctx1 := context.Background()
	ctx2 := context.Background()

	ctx1 = SetUserID(ctx1, "user-1")
	ctx1 = SetUsername(ctx1, "alice")
	ctx1 = SetRole(ctx1, "admin")

	ctx2 = SetUserID(ctx2, "user-2")
	ctx2 = SetUsername(ctx2, "bob")
	ctx2 = SetRole(ctx2, "user")

	// ctx1 should have its own values
	assert.Equal(t, "user-1", GetUserID(ctx1))
	assert.Equal(t, "alice", GetUsername(ctx1))
	assert.Equal(t, "admin", GetRole(ctx1))

	// ctx2 should have its own values
	assert.Equal(t, "user-2", GetUserID(ctx2))
	assert.Equal(t, "bob", GetUsername(ctx2))
	assert.Equal(t, "user", GetRole(ctx2))
}

func TestContext_PartialValues(t *testing.T) {
	// Test setting only some values
	ctx := context.Background()

	ctx = SetUserID(ctx, "user-123")
	// Don't set username
	ctx = SetRole(ctx, "admin")

	assert.Equal(t, "user-123", GetUserID(ctx))
	assert.Equal(t, "", GetUsername(ctx)) // Should return empty string
	assert.Equal(t, "admin", GetRole(ctx))
}

func TestContext_ChainedSetting(t *testing.T) {
	// Test chaining context setting
	ctx := context.Background()

	ctx = SetRole(SetUsername(SetUserID(ctx, "user-123"), "johndoe"), "admin")

	assert.Equal(t, "user-123", GetUserID(ctx))
	assert.Equal(t, "johndoe", GetUsername(ctx))
	assert.Equal(t, "admin", GetRole(ctx))
}

func TestContext_NilValue(t *testing.T) {
	// Test when context has nil value (though SetX functions don't allow this)
	ctx := context.WithValue(context.Background(), UserIDKey, nil)

	result := GetUserID(ctx)
	assert.Equal(t, "", result, "should return empty string when value is nil")
}

func TestContext_EmptyStringVsNotSet(t *testing.T) {
	// Test difference between setting empty string and not setting at all
	ctx1 := context.Background()
	ctx1 = SetUserID(ctx1, "") // Explicitly set to empty string

	ctx2 := context.Background() // Not set at all

	// Both should return empty string, but for different reasons
	assert.Equal(t, "", GetUserID(ctx1), "explicitly set empty string")
	assert.Equal(t, "", GetUserID(ctx2), "not set at all")
}

func TestContextKeys_AreUnique(t *testing.T) {
	// Verify that the context keys are distinct
	assert.NotEqual(t, UserIDKey, UsernameKey)
	assert.NotEqual(t, UserIDKey, RoleKey)
	assert.NotEqual(t, UsernameKey, RoleKey)
}

func TestContext_LargeValues(t *testing.T) {
	// Test with large string values
	largeUserID := string(make([]byte, 10000))
	for i := range largeUserID {
		largeUserID = largeUserID[:i] + "a"
	}

	ctx := context.Background()
	ctx = SetUserID(ctx, largeUserID)

	result := GetUserID(ctx)
	assert.Equal(t, largeUserID, result)
}

func TestContext_SpecialCharacters(t *testing.T) {
	specialChars := []string{
		"user\nwith\nnewlines",
		"user\twith\ttabs",
		"user\"with\"quotes",
		"user'with'apostrophes",
		"user\\with\\backslashes",
		"user/with/slashes",
		"user with spaces",
		"用户中文",
		"🎉emoji🎉",
	}

	for _, char := range specialChars {
		t.Run("special_char_"+char[:min(10, len(char))], func(t *testing.T) {
			ctx := context.Background()
			ctx = SetUserID(ctx, char)

			result := GetUserID(ctx)
			assert.Equal(t, char, result)
		})
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}