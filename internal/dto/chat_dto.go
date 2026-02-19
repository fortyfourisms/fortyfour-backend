package dto

type ChatRequest struct {
	SessionID string `json:"session_id"`
	Message   string `json:"message"`
}

type ChatHistory struct {
	User string `json:"user"`
	Bot  string `json:"bot"`
}
