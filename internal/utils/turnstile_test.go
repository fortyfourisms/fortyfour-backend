package utils

import (
	"testing"
)

// ═══════════════════════════════════════════════════════════════════════════
// TEST: NewTurnstileValidator
// ═══════════════════════════════════════════════════════════════════════════

func TestNewTurnstileValidator(t *testing.T) {
	v := NewTurnstileValidator("test-secret")

	if v.SecretKey != "test-secret" {
		t.Errorf("expected SecretKey 'test-secret', got '%s'", v.SecretKey)
	}
	if v.Client == nil {
		t.Error("expected Client to be initialized")
	}
}

func TestNewTurnstileValidator_EmptyKey(t *testing.T) {
	v := NewTurnstileValidator("")

	if v.SecretKey != "" {
		t.Errorf("expected empty SecretKey, got '%s'", v.SecretKey)
	}
	if v.Client == nil {
		t.Error("expected Client to be initialized even with empty key")
	}
}

// ═══════════════════════════════════════════════════════════════════════════
// TEST: Validate — token kosong
// ═══════════════════════════════════════════════════════════════════════════

func TestTurnstileValidator_Validate_EmptyToken(t *testing.T) {
	v := NewTurnstileValidator("test-secret")

	result, err := v.Validate("", "127.0.0.1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Success {
		t.Error("expected Success=false for empty token")
	}
	if len(result.ErrorCodes) == 0 {
		t.Error("expected error codes for empty token")
	}
	if result.ErrorCodes[0] != "missing-input-response" {
		t.Errorf("expected 'missing-input-response', got '%s'", result.ErrorCodes[0])
	}
}
