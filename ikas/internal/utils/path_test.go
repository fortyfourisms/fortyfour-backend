package utils

import (
	"testing"
)

func TestExtractID(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		resource string
		want     string
	}{
		{
			name:     "Valid path and resource",
			path:     "/api/v1/users/123",
			resource: "users",
			want:     "123",
		},
		{
			name:     "Resource not in path",
			path:     "/api/v1/orders/123",
			resource: "users",
			want:     "",
		},
		{
			name:     "Path exactly resource",
			path:     "/users",
			resource: "users",
			want:     "",
		},
		{
			name:     "Path starts with resource",
			path:     "users/123",
			resource: "users",
			want:     "",
		},
		{
			name:     "Resource with trailing slash",
			path:     "/api/v1/users/",
			resource: "users",
			want:     "",
		},
		{
			name:     "Multiple segments after resource",
			path:     "/api/v1/users/123/profile",
			resource: "users",
			want:     "123/profile",
		},
		{
			name:     "Empty path",
			path:     "",
			resource: "users",
			want:     "",
		},
		{
			name:     "Empty resource",
			path:     "/api/v1/users/123",
			resource: "",
			want:     "api/v1/users/123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractID(tt.path, tt.resource); got != tt.want {
				t.Errorf("ExtractID() = %v, want %v", got, tt.want)
			}
		})
	}
}
