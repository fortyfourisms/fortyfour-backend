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

	log.Printf("User Question: %s", req.Message)

	// STEP 1: Generate SQL Query
	sqlQuery, err := h.service.GenerateSQLQuery(req.Message)
	if err != nil {
		log.Printf("SQL Generation Error: %v", err)
		sendSSEError(w, flusher, "Gagal membuat query database")
		return
	}
	log.Printf("Generated SQL: %s", sqlQuery)

	// STEP 2: Execute SQL Query
	results, err := h.service.ExecuteQuery(sqlQuery)
	if err != nil {
		log.Printf("Query Execution Error: %v", err)
		sendSSEError(w, flusher, "Gagal menjalankan query database")
		return
	}
	log.Printf("Query Results: %d rows", len(results))

	// STEP 3: Format hasil query
	answer, err := h.service.FormatQueryResults(req.Message, results)
	if err != nil {
		log.Printf("Formatting Error: %v", err)
		sendSSEError(w, flusher, "Gagal memformat jawaban")
		return
	}

	// STEP 4: Stream response word by word
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

	// Send done event
	sendSSEEvent(w, flusher, map[string]interface{}{
		"type":    "done",
		"content": fullAnswer,
		"done":    true,
	})

	// STEP 5: Save history
	_ = h.service.Repo().Save(req.SessionID, req.Message, fullAnswer)
	log.Printf("Saved to history for session: %s", req.SessionID)
}

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
