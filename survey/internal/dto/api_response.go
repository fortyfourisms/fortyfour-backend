package dto

// APIResponse adalah standar response JSON API
type APIResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message,omitempty"`
	Data    interface{}       `json:"data,omitempty"`
	Error   string            `json:"error,omitempty"`
	Errors  map[string]string `json:"errors,omitempty"`
}
