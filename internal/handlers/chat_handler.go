package handlers

import (
	"encoding/json"
	"fmt"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/repository"
	"fortyfour-backend/internal/services"
	"log"
	"net/http"
	"strings"
)

type ChatHandler struct {
	service *services.ChatService
}

func NewChatHandler(s *services.ChatService) *ChatHandler {
	return &ChatHandler{service: s}
}

// Stream handles chat with SSE streaming
func (h *ChatHandler) Stream(w http.ResponseWriter, r *http.Request) {
	// SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE not supported", http.StatusInternalServerError)
		return
	}

	var req dto.ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendSSEError(w, flusher, "Invalid request body")
		return
	}

	//DETECT INTENT
	intent := services.DetectIntent(req.Message)

	//GET DATA FROM DB
	dataContext := h.service.BuildDataContext(intent)

	//GET CHAT HISTORY
	history, _ := h.service.Repo().GetHistory(req.SessionID)

	//BUILD PROMPT (ANTI HALU)
	systemPrompt := `
Kamu adalah chatbot CS internal.
ATURAN:
- Jawaban hanya boleh berdasarkan DATA SISTEM
- Jangan mengarang
- Jangan gunakan pengetahuan umum
- Jika data tidak ada, jawab: "Data tidak tersedia di sistem"
`

	prompt := systemPrompt + "\n\n"
	prompt += "DATA SISTEM:\n" + dataContext + "\n\n"

	for _, h := range history {
		prompt += "User: " + h.User + "\n"
		prompt += "Bot: " + h.Bot + "\n"
	}

	prompt += "User: " + req.Message + "\n"

	log.Printf("Generating response for session: %s", req.SessionID)

	//GENERATE AI RESPONSE
	answer, err := h.service.GetGemini().Generate(prompt)
	if err != nil {
		log.Printf("Gemini error: %v", err)
		sendSSEError(w, flusher, "Failed to generate response")
		return
	}

	//SSE STREAMING (TETAP SAMA)
	words := strings.Fields(answer)
	fullAnswer := ""

	for i, word := range words {
		if i > 0 {
			fullAnswer += " "
		}
		fullAnswer += word

		sendSSEEvent(w, flusher, map[string]interface{}{
			"type":    "chunk",
			"content": word + " ",
			"done":    false,
		})
	}

	sendSSEEvent(w, flusher, map[string]interface{}{
		"type":    "done",
		"content": fullAnswer,
		"done":    true,
	})

	//SAVE HISTORY
	_ = h.service.Repo().Save(req.SessionID, req.Message, fullAnswer)
}

// DeleteSession deletes a chat session
func (h *ChatHandler) DeleteSession(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SessionID string `json:"session_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if repo, ok := h.service.Repo().(*repository.InMemoryChatRepo); ok {
		if err := repo.DeleteSession(req.SessionID); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": fmt.Sprintf("Session %s deleted", req.SessionID),
		})
		return
	}

	http.Error(w, "unsupported repository type", http.StatusInternalServerError)
}

// Helper functions
func sendSSEEvent(w http.ResponseWriter, flusher http.Flusher, data interface{}) {
	jsonData, _ := json.Marshal(data)
	fmt.Fprintf(w, "data: %s\n\n", jsonData)
	flusher.Flush()
}

func sendSSEError(w http.ResponseWriter, flusher http.Flusher, message string) {
	event := map[string]interface{}{
		"type":    "error",
		"content": message,
		"done":    true,
	}
	sendSSEEvent(w, flusher, event)
}
