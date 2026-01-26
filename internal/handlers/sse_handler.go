package handlers

import (
	"encoding/json"
	"fortyfour-backend/internal/middleware"
	"fortyfour-backend/internal/services"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/rollbar/rollbar-go"
)

type SSEHandler struct {
	sseService *services.SSEService
}

func NewSSEHandler(sseService *services.SSEService) *SSEHandler {
	return &SSEHandler{
		sseService: sseService,
	}
}

// HandleSSE handles SSE connections
func (h *SSEHandler) HandleSSE(w http.ResponseWriter, r *http.Request) {
	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Create a new client
	clientID := uuid.New().String()

	// Get user ID from context (set by auth middleware)
	userID := ""
	if uid := r.Context().Value(middleware.UserIDKey); uid != nil {
		userID = uid.(string)
	}

	client := &services.Client{
		ID:      clientID,
		UserID:  userID,
		Channel: make(chan services.SSEEvent, 10),
	}

	// Register client
	h.sseService.RegisterClient(client)
	defer h.sseService.UnregisterClient(client)

	// Send initial connection message
	initialMsg := services.SSEEvent{
		Type:     "connected",
		Resource: "system",
		Data:     map[string]string{"message": "Connected to SSE", "client_id": clientID},
	}
	w.Write([]byte(services.FormatSSEMessage(initialMsg)))
	w.(http.Flusher).Flush()

	log.Printf("SSE: Client %s connected", clientID)

	// Listen for client disconnect
	notify := r.Context().Done()

	// Stream events
	for {
		select {
		case event := <-client.Channel:
			// Format and send event
			message := services.FormatSSEMessage(event)
			if message != "" {
				_, err := w.Write([]byte(message))
				if err != nil {
					log.Printf("SSE: Error writing to client %s: %v", clientID, err)
					rollbar.Error(err)
					return
				}
				w.(http.Flusher).Flush()
			}

		case <-notify:
			// Client disconnected
			log.Printf("SSE: Client %s disconnected", clientID)
			return
		}
	}
}

// GetStats returns SSE statistics
func (h *SSEHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	stats := map[string]interface{}{
		"connected_clients": h.sseService.GetClientCount(),
	}

	json.NewEncoder(w).Encode(stats)
}
