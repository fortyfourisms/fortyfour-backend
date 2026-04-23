package dto

type APIResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message,omitempty"`
	Data    interface{}       `json:"data,omitempty"`
	Error   string            `json:"error,omitempty"`
	Errors  map[string]string `json:"errors,omitempty"`
}

type ErrorResponse struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
}