package utils

import (
	"encoding/json"
	"testing"
)

// ═══════════════════════════════════════════════════════════════════════════
// TEST: EventRegistrationQRPayload.JSON
// ═══════════════════════════════════════════════════════════════════════════

func TestEventRegistrationQRPayload_JSON(t *testing.T) {
	payload := EventRegistrationQRPayload{
		EventID:        1,
		EventTitle:     "Workshop Keamanan Siber",
		RegistrationID: 42,
		Nama:           "John Doe",
		Email:          "john@example.com",
		Perusahaan:     "PT Test",
		Jabatan:        "Manager",
		NoHP:           "08123456789",
		Sektor:         "Keuangan",
		QRToken:        "token-abc-123",
	}

	jsonStr, err := payload.JSON()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Parse back to verify
	var parsed EventRegistrationQRPayload
	if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
		t.Fatalf("result is not valid JSON: %v", err)
	}

	if parsed.EventID != 1 {
		t.Errorf("expected EventID 1, got %d", parsed.EventID)
	}
	if parsed.Nama != "John Doe" {
		t.Errorf("expected Nama 'John Doe', got '%s'", parsed.Nama)
	}
	if parsed.QRToken != "token-abc-123" {
		t.Errorf("expected QRToken 'token-abc-123', got '%s'", parsed.QRToken)
	}
}

func TestEventRegistrationQRPayload_JSON_Empty(t *testing.T) {
	payload := EventRegistrationQRPayload{}

	jsonStr, err := payload.JSON()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if jsonStr == "" {
		t.Error("expected non-empty JSON string")
	}
}

// ═══════════════════════════════════════════════════════════════════════════
// TEST: GenerateQRCodePNG
// ═══════════════════════════════════════════════════════════════════════════

func TestGenerateQRCodePNG_Success(t *testing.T) {
	data, err := GenerateQRCodePNG("Hello QR", 256)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty PNG data")
	}

	// PNG magic bytes: 0x89 0x50 0x4E 0x47
	if data[0] != 0x89 || data[1] != 0x50 || data[2] != 0x4E || data[3] != 0x47 {
		t.Error("output is not valid PNG (wrong magic bytes)")
	}
}

func TestGenerateQRCodePNG_EmptyPayload(t *testing.T) {
	_, err := GenerateQRCodePNG("", 256)
	if err == nil {
		t.Error("expected error for empty payload")
	}
}

func TestGenerateQRCodePNG_WhitespacePayload(t *testing.T) {
	_, err := GenerateQRCodePNG("   ", 256)
	if err == nil {
		t.Error("expected error for whitespace-only payload")
	}
}

func TestGenerateQRCodePNG_DefaultSize(t *testing.T) {
	// size <= 0 should use default 256
	data, err := GenerateQRCodePNG("test", 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty PNG data with default size")
	}
}

func TestGenerateQRCodePNG_NegativeSize(t *testing.T) {
	// Negative size should also use default
	data, err := GenerateQRCodePNG("test", -1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty PNG data with negative size")
	}
}

func TestGenerateQRCodePNG_LongPayload(t *testing.T) {
	longPayload := `{"event_id":12345,"event_title":"Workshop Keamanan Siber Nasional 2024","registration_id":999,"nama":"John Doe","email":"john.doe@example.com","perusahaan":"PT Teknologi Indonesia","jabatan":"Chief Information Security Officer","no_hp":"081234567890","sektor":"Keuangan","qr_token":"very-long-token-abc-def-123-456-789"}`

	data, err := GenerateQRCodePNG(longPayload, 512)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty PNG data for long payload")
	}
}

func TestGenerateQRCodePNG_JSONPayload(t *testing.T) {
	payload := EventRegistrationQRPayload{
		EventID:        1,
		RegistrationID: 42,
		Nama:           "Test User",
		QRToken:        "token-123",
	}

	jsonStr, err := payload.JSON()
	if err != nil {
		t.Fatalf("JSON() error: %v", err)
	}

	data, err := GenerateQRCodePNG(jsonStr, 256)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty PNG data")
	}
}
