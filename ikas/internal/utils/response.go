package utils

import (
	"encoding/json"
	"net/http"
)

// JSONResponse is the standard API envelope for single-item or non-paginated responses.
type JSONResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Total   interface{} `json:"total,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// PaginationMeta holds pagination metadata, following the industry-standard
// pattern used by GitHub, Stripe, and Google APIs.
type PaginationMeta struct {
	Total      int `json:"total"`
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalPages int `json:"total_pages"`
}

// PaginatedJSONResponse is the standard API envelope for paginated list responses.
// Data contains the array of items directly (no double-nesting).
// Pagination metadata is isolated in the "pagination" object.
type PaginatedJSONResponse struct {
	Status     string         `json:"status"`
	Message    string         `json:"message"`
	Pagination PaginationMeta `json:"pagination"`
	Data       interface{}    `json:"data"`
}

func RespondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func RespondError(w http.ResponseWriter, status int, message string) {
	RespondJSON(w, status, JSONResponse{
		Status:  "error",
		Message: message,
	})
}
