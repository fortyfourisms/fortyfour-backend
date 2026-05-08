package utils

import (
	"net/url"
	"testing"
)

// ═══════════════════════════════════════════════════════════════════════════
// TEST: ExtractIntID
// ═══════════════════════════════════════════════════════════════════════════

func TestExtractIntID_ValidPath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		resource string
		expected int
	}{
		{"simple path", "/api/aktivitas/123", "aktivitas", 123},
		{"path with trailing slash", "/api/aktivitas/456/", "aktivitas", 456},
		{"nested resource", "/api/v2/aktivitas/789", "aktivitas", 789},
		{"first ID only", "/api/aktivitas/10/detail", "aktivitas", 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := ExtractIntID(tt.path, tt.resource)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if id != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, id)
			}
		})
	}
}

func TestExtractIntID_ResourceNotFound(t *testing.T) {
	id, err := ExtractIntID("/api/berita/1", "aktivitas")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if id != 0 {
		t.Errorf("expected 0, got %d", id)
	}
}

func TestExtractIntID_NonNumericID(t *testing.T) {
	_, err := ExtractIntID("/api/aktivitas/abc", "aktivitas")
	if err == nil {
		t.Error("expected error for non-numeric ID")
	}
}

func TestExtractIntID_ResourceAtEnd(t *testing.T) {
	// Resource is last segment — no ID after it
	id, err := ExtractIntID("/api/aktivitas", "aktivitas")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if id != 0 {
		t.Errorf("expected 0, got %d", id)
	}
}

func TestExtractIntID_WithQueryParams(t *testing.T) {
	id, err := ExtractIntID("/api/aktivitas/42?page=1", "aktivitas")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if id != 42 {
		t.Errorf("expected 42, got %d", id)
	}
}

func TestExtractIntID_EmptyPath(t *testing.T) {
	id, err := ExtractIntID("", "aktivitas")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if id != 0 {
		t.Errorf("expected 0, got %d", id)
	}
}

// ═══════════════════════════════════════════════════════════════════════════
// TEST: ExtractStringID
// ═══════════════════════════════════════════════════════════════════════════

func TestExtractStringID_ValidPath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		resource string
		expected string
	}{
		{"UUID path", "/api/perusahaan/uuid-123-abc", "perusahaan", "uuid-123-abc"},
		{"simple path", "/api/users/admin-1", "users", "admin-1"},
		{"nested", "/api/v2/kelas/kls-1/materi", "kelas", "kls-1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := ExtractStringID(tt.path, tt.resource)
			if id != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, id)
			}
		})
	}
}

func TestExtractStringID_ResourceNotFound(t *testing.T) {
	id := ExtractStringID("/api/berita/1", "perusahaan")
	if id != "" {
		t.Errorf("expected empty string, got '%s'", id)
	}
}

func TestExtractStringID_ResourceAtEnd(t *testing.T) {
	id := ExtractStringID("/api/perusahaan", "perusahaan")
	if id != "" {
		t.Errorf("expected empty string, got '%s'", id)
	}
}

func TestExtractStringID_WithQueryParams(t *testing.T) {
	id := ExtractStringID("/api/perusahaan/uuid-456?detail=true", "perusahaan")
	if id != "uuid-456" {
		t.Errorf("expected 'uuid-456', got '%s'", id)
	}
}

func TestExtractStringID_EmptyPath(t *testing.T) {
	id := ExtractStringID("", "perusahaan")
	if id != "" {
		t.Errorf("expected empty string, got '%s'", id)
	}
}

// ═══════════════════════════════════════════════════════════════════════════
// TEST: ParseQueryParams
// ═══════════════════════════════════════════════════════════════════════════

func TestParseQueryParams_WithParams(t *testing.T) {
	u, _ := url.Parse("http://localhost/api/data?page=1&limit=10&sort=name")
	params := ParseQueryParams(u)

	if params.Get("page") != "1" {
		t.Errorf("expected page=1, got '%s'", params.Get("page"))
	}
	if params.Get("limit") != "10" {
		t.Errorf("expected limit=10, got '%s'", params.Get("limit"))
	}
	if params.Get("sort") != "name" {
		t.Errorf("expected sort=name, got '%s'", params.Get("sort"))
	}
}

func TestParseQueryParams_NoParams(t *testing.T) {
	u, _ := url.Parse("http://localhost/api/data")
	params := ParseQueryParams(u)

	if len(params) != 0 {
		t.Errorf("expected 0 params, got %d", len(params))
	}
}

func TestParseQueryParams_EmptyValue(t *testing.T) {
	u, _ := url.Parse("http://localhost/api/data?key=")
	params := ParseQueryParams(u)

	if params.Get("key") != "" {
		t.Errorf("expected empty value, got '%s'", params.Get("key"))
	}
}

func TestParseQueryParams_MissingKey(t *testing.T) {
	u, _ := url.Parse("http://localhost/api/data?page=1")
	params := ParseQueryParams(u)

	if params.Get("nonexistent") != "" {
		t.Errorf("expected empty for missing key, got '%s'", params.Get("nonexistent"))
	}
}
