package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRespondJSON(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]string{"message": "test"}

	RespondJSON(w, http.StatusOK, data)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got '%s'", w.Header().Get("Content-Type"))
	}

	var response map[string]string
	json.NewDecoder(w.Body).Decode(&response)
	if response["message"] != "test" {
		t.Errorf("expected message 'test', got '%s'", response["message"])
	}
}

func TestRespondError(t *testing.T) {
	w := httptest.NewRecorder()

	RespondError(w, http.StatusBadRequest, "test error")

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got '%s'", w.Header().Get("Content-Type"))
	}

	var response map[string]string
	json.NewDecoder(w.Body).Decode(&response)
	if response["error"] != "test error" {
		t.Errorf("expected error 'test error', got '%s'", response["error"])
	}
}

func TestRespondError_DifferentStatusCodes(t *testing.T) {
	testCases := []struct {
		name           string
		statusCode     int
		message        string
		expectedStatus int
	}{
		{"Bad Request", http.StatusBadRequest, "bad request", http.StatusBadRequest},
		{"Unauthorized", http.StatusUnauthorized, "unauthorized", http.StatusUnauthorized},
		{"Forbidden", http.StatusForbidden, "forbidden", http.StatusForbidden},
		{"Not Found", http.StatusNotFound, "not found", http.StatusNotFound},
		{"Internal Server Error", http.StatusInternalServerError, "internal error", http.StatusInternalServerError},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			RespondError(w, tc.statusCode, tc.message)

			if w.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, w.Code)
			}

			var response map[string]string
			json.NewDecoder(w.Body).Decode(&response)
			if response["error"] != tc.message {
				t.Errorf("expected error '%s', got '%s'", tc.message, response["error"])
			}
		})
	}
}

