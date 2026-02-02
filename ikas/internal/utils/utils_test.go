package utils

import (
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
