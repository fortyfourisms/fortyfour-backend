package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthenticate(t *testing.T) {
	internalKey := "secret-test"
	middleware := NewAuthMiddleware(internalKey)

	tests := []struct {
		name           string
		headers        map[string]string
		expectedStatus int
		expectedBody   string
		checkContext   bool
	}{
		{
			name: "Success - All headers present",
			headers: map[string]string{
				"X-Internal-Key":  internalKey,
				"X-User-ID":       "user-123",
				"X-Username":      "johndoe",
				"X-User-Role":     "admin",
				"X-Perusahaan-ID": "comp-456",
			},
			expectedStatus: http.StatusOK,
			checkContext:   true,
		},
		{
			name: "Success - Optional headers empty",
			headers: map[string]string{
				"X-Internal-Key": internalKey,
				"X-User-ID":      "user-123",
				"X-User-Role":    "user",
			},
			expectedStatus: http.StatusOK,
			checkContext:   true,
		},
		{
			name: "Failure - Invalid internal key",
			headers: map[string]string{
				"X-Internal-Key": "wrong-key",
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Unauthorized: Direct access not allowed or invalid internal key",
		},
		{
			name: "Failure - Missing User ID",
			headers: map[string]string{
				"X-Internal-Key": internalKey,
				"X-User-Role":    "admin",
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Unauthorized: User identification missing in gateway headers",
		},
		{
			name: "Failure - Missing User Role",
			headers: map[string]string{
				"X-Internal-Key": internalKey,
				"X-User-ID":      "user-123",
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Unauthorized: User identification missing in gateway headers",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup request
			req := httptest.NewRequest(http.MethodGet, "http://example.com/test", nil)
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}

			// Setup response recorder
			rr := httptest.NewRecorder()

			// Next handler to check context
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.checkContext {
					assert.Equal(t, tt.headers["X-User-ID"], r.Context().Value(UserIDKey))
					assert.Equal(t, tt.headers["X-Username"], r.Context().Value(Username))
					assert.Equal(t, tt.headers["X-User-Role"], r.Context().Value(Role))
					assert.Equal(t, tt.headers["X-Perusahaan-ID"], r.Context().Value(PerusahaanIDKey))
				}
				w.WriteHeader(http.StatusOK)
			})

			// Run middleware
			handler := middleware.Authenticate(nextHandler)
			handler.ServeHTTP(rr, req)

			// Assert status code
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// Assert body if present
			if tt.expectedBody != "" {
				var resp map[string]string
				err := json.NewDecoder(rr.Body).Decode(&resp)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, resp["message"])
			}
		})
	}
}
