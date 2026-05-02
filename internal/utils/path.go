package utils

import (
	"net/url"
	"strconv"
	"strings"
)

// ExtractIntID mengambil ID bertipe integer dari URL path berdasarkan prefix resource.
// Contoh: /api/aktivitas/123 -> returns 123
func ExtractIntID(path string, resource string) (int, error) {
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if part == resource && i+1 < len(parts) {
			idStr := parts[i+1]
			// Hilangkan query params jika ada
			if idx := strings.Index(idStr, "?"); idx != -1 {
				idStr = idStr[:idx]
			}
			return strconv.Atoi(idStr)
		}
	}
	return 0, nil
}

// ExtractStringID mengambil ID bertipe string dari URL path.
func ExtractStringID(path string, resource string) string {
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if part == resource && i+1 < len(parts) {
			id := parts[i+1]
			if idx := strings.Index(id, "?"); idx != -1 {
				id = id[:idx]
			}
			return id
		}
	}
	return ""
}

// ParseQueryParams parses query parameters safely
func ParseQueryParams(rUrl *url.URL) url.Values {
	return rUrl.Query()
}
