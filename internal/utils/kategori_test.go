package utils

import (
	"testing"
)

func TestGetKategoriTingkatKematangan(t *testing.T) {
	testCases := []struct {
		name     string
		nilai    float64
		expected string
	}{
		{"Level 1 - Awal (min)", 0.0, "Level 1 - Awal"},
		{"Level 1 - Awal (mid)", 0.75, "Level 1 - Awal"},
		{"Level 1 - Awal (max)", 1.5, "Level 1 - Awal"},
		{"Level 2 - Berulang (min)", 1.51, "Level 2 - Berulang"},
		{"Level 2 - Berulang (mid)", 2.0, "Level 2 - Berulang"},
		{"Level 2 - Berulang (max)", 2.5, "Level 2 - Berulang"},
		{"Level 3 - Terdefinisi (min)", 2.51, "Level 3 - Terdefinisi"},
		{"Level 3 - Terdefinisi (mid)", 3.0, "Level 3 - Terdefinisi"},
		{"Level 3 - Terdefinisi (max)", 3.5, "Level 3 - Terdefinisi"},
		{"Level 4 - Terkelola (min)", 3.51, "Level 4 - Terkelola"},
		{"Level 4 - Terkelola (mid)", 4.0, "Level 4 - Terkelola"},
		{"Level 4 - Terkelola (max)", 4.5, "Level 4 - Terkelola"},
		{"Level 5 - Inovatif (min)", 4.51, "Level 5 - Inovatif"},
		{"Level 5 - Inovatif (mid)", 4.75, "Level 5 - Inovatif"},
		{"Level 5 - Inovatif (max)", 5.0, "Level 5 - Inovatif"},
		{"Out of range (negative)", -1.0, "Tidak Terdefinisi"},
		{"Out of range (high)", 6.0, "Tidak Terdefinisi"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := GetKategoriTingkatKematangan(tc.nilai)
			if result != tc.expected {
				t.Errorf("expected '%s', got '%s'", tc.expected, result)
			}
		})
	}
}

