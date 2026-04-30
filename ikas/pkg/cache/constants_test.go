package cache

import (
	"testing"
	"time"
)

func TestConstants(t *testing.T) {
	if DefaultCacheExpiration != 24*time.Hour {
		t.Errorf("Expected DefaultCacheExpiration to be 24h, got %v", DefaultCacheExpiration)
	}

	// Verify static keys
	if CacheKeyIkasRecords != "ikas:records" {
		t.Errorf("Unexpected CacheKeyIkasRecords: %s", CacheKeyIkasRecords)
	}
	if CacheKeyPertanyaanIdentifikasi != "ikas:questions:identifikasi" {
		t.Errorf("Unexpected CacheKeyPertanyaanIdentifikasi: %s", CacheKeyPertanyaanIdentifikasi)
	}
}

func TestKeyHelpers(t *testing.T) {
	ikasID := "test-id-123"

	tests := []struct {
		name     string
		got      string
		expected string
	}{
		{"Jawaban Identifikasi", GetJawabanIdentifikasiKey(ikasID), "ikas:jawaban:identifikasi:" + ikasID},
		{"Jawaban Proteksi", GetJawabanProteksiKey(ikasID), "ikas:jawaban:proteksi:" + ikasID},
		{"Jawaban Deteksi", GetJawabanDeteksiKey(ikasID), "ikas:jawaban:deteksi:" + ikasID},
		{"Jawaban Gulih", GetJawabanGulihKey(ikasID), "ikas:jawaban:gulih:" + ikasID},
		{"Domain Identifikasi", GetIdentifikasiKey(ikasID), "ikas:domain:identifikasi:" + ikasID},
		{"Domain Proteksi", GetProteksiKey(ikasID), "ikas:domain:proteksi:" + ikasID},
		{"Domain Deteksi", GetDeteksiKey(ikasID), "ikas:domain:deteksi:" + ikasID},
		{"Domain Gulih", GetGulihKey(ikasID), "ikas:domain:gulih:" + ikasID},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("%s: expected %s, got %s", tt.name, tt.expected, tt.got)
			}
		})
	}
}
