package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValueOrNull(t *testing.T) {
	testCases := []struct {
		name     string
		input    *string
		expected interface{}
	}{
		{"nil pointer", nil, nil},
		{"valid string", stringPtr("test"), "test"},
		{"empty string", stringPtr(""), ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ValueOrNull(tc.input)
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestValueOrEmpty(t *testing.T) {
	testCases := []struct {
		name     string
		input    *string
		expected string
	}{
		{"nil pointer", nil, ""},
		{"valid string", stringPtr("test"), "test"},
		{"empty string", stringPtr(""), ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ValueOrEmpty(tc.input)
			if result != tc.expected {
				t.Errorf("expected '%s', got '%s'", tc.expected, result)
			}
		})
	}
}

func TestIntOrNull(t *testing.T) {
	testCases := []struct {
		name     string
		input    *int
		expected interface{}
	}{
		{"nil pointer", nil, nil},
		{"valid int", intPtr(42), 42},
		{"zero int", intPtr(0), 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := IntOrNull(tc.input)
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestBoolOrNull(t *testing.T) {
	testCases := []struct {
		name     string
		input    *bool
		expected interface{}
	}{
		{"nil pointer", nil, nil},
		{"true bool", boolPtr(true), true},
		{"false bool", boolPtr(false), false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := BoolOrNull(tc.input)
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestGenerateUUID(t *testing.T) {
	u1 := GenerateUUID()
	u2 := GenerateUUID()
	if u1 == "" || u2 == "" {
		t.Error("GenerateUUID returned empty string")
	}
	if u1 == u2 {
		t.Error("GenerateUUID returned non-unique IDs")
	}
}

func TestRoundToTwo(t *testing.T) {
	tests := []struct {
		input    float64
		expected float64
	}{
		{3.14159, 3.14},
		{3.145, 3.15},
		{3.1, 3.1},
	}
	for _, tc := range tests {
		if got := RoundToTwo(tc.input); got != tc.expected {
			t.Errorf("RoundToTwo(%f) = %f, want %f", tc.input, got, tc.expected)
		}
	}
}

func TestStringToInt(t *testing.T) {
	val, err := StringToInt("123")
	if err != nil || val != 123 {
		t.Errorf("StringToInt('123') = %v, %v", val, err)
	}

	_, err = StringToInt("abc")
	if err == nil {
		t.Error("StringToInt('abc') should have failed")
	}
}

func TestExtractIntID(t *testing.T) {
	val, err := ExtractIntID("/api/v1/users/123", "users")
	if err != nil || val != 123 {
		t.Errorf("ExtractIntID() = %v, %v", val, err)
	}

	val, err = ExtractIntID("/api/v1/users", "users")
	if err != nil || val != 0 {
		t.Errorf("ExtractIntID() with empty ID = %v, %v", val, err)
	}

	_, err = ExtractIntID("/api/v1/users/abc", "users")
	if err == nil {
		t.Error("ExtractIntID() with non-int ID should have failed")
	}
}

func TestAdaptHandler(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	})
	adapted := AdaptHandler(handler)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	adapted(w, r)
	if w.Code != http.StatusAccepted {
		t.Errorf("expected status 202, got %d", w.Code)
	}
}

func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func boolPtr(b bool) *bool {
	return &b
}
