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

func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func boolPtr(b bool) *bool {
	return &b
}

// ================================================================
// AdaptHandler
// ================================================================

func TestAdaptHandler_WrapsHandlerAndServes(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("created"))
	})

	adapted := AdaptHandler(inner)

	req := httptest.NewRequest(http.MethodPost, "/resource", nil)
	w := httptest.NewRecorder()
	adapted(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}
	if w.Body.String() != "created" {
		t.Errorf("expected body 'created', got '%s'", w.Body.String())
	}
}

func TestAdaptHandler_ReturnTypeIsHandlerFunc(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	result := AdaptHandler(inner)

	// Pastikan hasil adalah http.HandlerFunc yang bisa dipanggil langsung
	var _ http.HandlerFunc = result
}

func TestAdaptHandler_PassesRequestThrough(t *testing.T) {
	var gotMethod, gotPath string

	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
	})

	adapted := AdaptHandler(inner)

	req := httptest.NewRequest(http.MethodDelete, "/api/items/99", nil)
	w := httptest.NewRecorder()
	adapted(w, req)

	if gotMethod != http.MethodDelete {
		t.Errorf("expected method DELETE, got '%s'", gotMethod)
	}
	if gotPath != "/api/items/99" {
		t.Errorf("expected path '/api/items/99', got '%s'", gotPath)
	}
}
