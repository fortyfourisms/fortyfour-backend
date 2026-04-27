package utils

import (
	"encoding/json"
	"net/http"
)

type JSONResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Total   interface{} `json:"total,omitempty"`
	Data    interface{} `json:"data,omitempty"`
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
