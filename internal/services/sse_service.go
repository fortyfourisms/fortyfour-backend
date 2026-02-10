package services

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
	"sync"
	"time"

)

// Event types
const (
	EventCreate = "create"
	EventUpdate = "update"
	EventDelete = "delete"
)

// SSEEvent represents a Server-Sent Event
type SSEEvent struct {
	Type      string      `json:"type"`
	Resource  string      `json:"resource"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	UserID    string      `json:"user_id"`
	Timestamp time.Time   `json:"timestamp"`
}

// Client represents an SSE client connection
type Client struct {
	ID       string
	UserID   string
	Channel  chan SSEEvent
	LastPing time.Time
}

// SSEService manages SSE connections and broadcasts
type SSEService struct {
	clients      map[string]*Client
	mu           sync.RWMutex
	broadcast    chan SSEEvent
	register     chan *Client
	unregister   chan *Client
	pingInterval time.Duration
}

// NewSSEService creates a new SSE service
func NewSSEService() *SSEService {
	service := &SSEService{
		clients:      make(map[string]*Client),
		broadcast:    make(chan SSEEvent, 100),
		register:     make(chan *Client),
		unregister:   make(chan *Client),
		pingInterval: 30 * time.Second,
	}

	go service.run()
	return service
}

func (s *SSEService) run() {
	ticker := time.NewTicker(s.pingInterval)
	defer ticker.Stop()

	for {
		select {
		case client := <-s.register:
			s.mu.Lock()
			s.clients[client.ID] = client
			s.mu.Unlock()
			log.Printf("SSE: Client %s registered. Total: %d", client.ID, len(s.clients))

		case client := <-s.unregister:
			s.mu.Lock()
			if _, ok := s.clients[client.ID]; ok {
				close(client.Channel)
				delete(s.clients, client.ID)
				log.Printf("SSE: Client %s unregistered. Total: %d", client.ID, len(s.clients))
			}
			s.mu.Unlock()

		case event := <-s.broadcast:
			s.mu.RLock()
			for _, client := range s.clients {
				select {
				case client.Channel <- event:
				default:
					log.Printf("SSE: Client %s channel full", client.ID)
				}
			}
			s.mu.RUnlock()

		case <-ticker.C:
			s.mu.RLock()
			for _, client := range s.clients {
				select {
				case client.Channel <- SSEEvent{Type: "ping", Timestamp: time.Now()}:
					client.LastPing = time.Now()
				default:
				}
			}
			s.mu.RUnlock()
		}
	}
}

func (s *SSEService) RegisterClient(client *Client) {
	s.register <- client
}

func (s *SSEService) UnregisterClient(client *Client) {
	s.unregister <- client
}

func (s *SSEService) Broadcast(event SSEEvent) {
	event.Timestamp = time.Now()
	s.broadcast <- event
}

// NotifyCreate sends a create event dengan pesan otomatis
func (s *SSEService) NotifyCreate(resource string, data interface{}, userID string) {
	message := s.generateMessage(resource, EventCreate, data)
	s.Broadcast(SSEEvent{
		Type:     EventCreate,
		Resource: resource,
		Message:  message,
		Data:     data,
		UserID:   userID,
	})
}

// NotifyUpdate sends an update event dengan pesan otomatis
func (s *SSEService) NotifyUpdate(resource string, data interface{}, userID string) {
	message := s.generateMessage(resource, EventUpdate, data)
	s.Broadcast(SSEEvent{
		Type:     EventUpdate,
		Resource: resource,
		Message:  message,
		Data:     data,
		UserID:   userID,
	})
}

// NotifyDelete sends a delete event dengan pesan otomatis
func (s *SSEService) NotifyDelete(resource string, id interface{}, userID string) {
	message := s.generateMessage(resource, EventDelete, id)
	s.Broadcast(SSEEvent{
		Type:     EventDelete,
		Resource: resource,
		Message:  message,
		Data:     map[string]interface{}{"id": id},
		UserID:   userID,
	})
}

// generateMessage membuat pesan notifikasi
func (s *SSEService) generateMessage(resource, eventType string, data interface{}) string {
	resourceName := s.formatResourceName(resource)
	name := s.extractName(data)

	switch eventType {
	case EventCreate:
		if name != "" {
			return fmt.Sprintf("%s baru %s berhasil ditambahkan", resourceName, name)
		}
		return fmt.Sprintf("%s baru berhasil ditambahkan", resourceName)

	case EventUpdate:
		if name != "" {
			return fmt.Sprintf("%s %s berhasil diperbarui", resourceName, name)
		}
		return fmt.Sprintf("%s berhasil diperbarui", resourceName)

	case EventDelete:
		return fmt.Sprintf("%s berhasil dihapus", resourceName)

	default:
		return fmt.Sprintf("Event %s pada %s", eventType, resourceName)
	}
}

// formatResourceName converts resource ke format yang readable
func (s *SSEService) formatResourceName(resource string) string {
	if len(resource) <= 4 {
		upper := strings.ToUpper(resource)
		if upper == resource {
			return resource
		}
	}

	if len(resource) > 0 {
		return strings.ToUpper(string(resource[0])) + resource[1:]
	}

	return resource
}

// extractName mencoba extract nama dari data
func (s *SSEService) extractName(data interface{}) string {
	if data == nil {
		return ""
	}

	if dataMap, ok := data.(map[string]interface{}); ok {
		return s.extractNameFromMap(dataMap)
	}

	return s.extractNameFromStruct(data)
}

func (s *SSEService) extractNameFromMap(dataMap map[string]interface{}) string {
	var candidates []struct {
		key   string
		value string
		score int
	}

	// Scan SEMUA field di map
	for key, val := range dataMap {
		// Skip non-string values
		strVal := ""
		switch v := val.(type) {
		case string:
			strVal = v
		case *string:
			if v != nil {
				strVal = *v
			}
		default:
			continue
		}

		// Skip empty values
		if strVal == "" {
			continue
		}

		lowerKey := strings.ToLower(key)
		score := 0

		// Scoring system untuk menentukan prioritas
		// Semakin tinggi score, semakin prioritas

		// Exact match "nama" atau "name" (highest priority)
		if lowerKey == "nama" || lowerKey == "name" {
			score = 100
		}

		// Starts with "nama_" atau "name_"
		if strings.HasPrefix(lowerKey, "nama_") || strings.HasPrefix(lowerKey, "name_") {
			score = 90
		}

		// Ends with "_nama" atau "_name"
		if strings.HasSuffix(lowerKey, "_nama") || strings.HasSuffix(lowerKey, "_name") {
			score = 80
		}

		// Contains "nama" atau "name"
		if strings.Contains(lowerKey, "nama") || strings.Contains(lowerKey, "name") {
			score = 70
		}

		// Contains "title", "label", "judul"
		if strings.Contains(lowerKey, "title") ||
			strings.Contains(lowerKey, "label") ||
			strings.Contains(lowerKey, "judul") {
			score = 60
		}

		// If found any match, add to candidates
		if score > 0 {
			candidates = append(candidates, struct {
				key   string
				value string
				score int
			}{key, strVal, score})
		}
	}

	// Return candidate dengan score tertinggi
	if len(candidates) > 0 {
		// Find highest score
		highest := candidates[0]
		for _, candidate := range candidates {
			if candidate.score > highest.score {
				highest = candidate
			}
		}
		return highest.value
	}

	return ""
}

// extractNameFromStruct - dengan reflection
func (s *SSEService) extractNameFromStruct(data interface{}) string {
	val := reflect.ValueOf(data)

	// Handle pointer
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return ""
		}
		val = val.Elem()
	}

	// Only process struct
	if val.Kind() != reflect.Struct {
		return ""
	}

	var candidates []struct {
		name  string
		value string
		score int
	}

	// Iterate SEMUA fields
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		fieldName := fieldType.Name

		// Skip unexported fields
		if !field.CanInterface() {
			continue
		}

		// Get string value
		strVal := ""
		switch field.Kind() {
		case reflect.String:
			strVal = field.String()
		case reflect.Ptr:
			if !field.IsNil() && field.Elem().Kind() == reflect.String {
				strVal = field.Elem().String()
			}
		default:
			continue
		}

		// Skip empty
		if strVal == "" {
			continue
		}

		lowerFieldName := strings.ToLower(fieldName)
		score := 0

		// Scoring
		if lowerFieldName == "nama" || lowerFieldName == "name" {
			score = 100
		} else if strings.HasPrefix(lowerFieldName, "nama") || strings.HasPrefix(lowerFieldName, "name") {
			score = 90
		} else if strings.HasSuffix(lowerFieldName, "nama") || strings.HasSuffix(lowerFieldName, "name") {
			score = 80
		} else if strings.Contains(lowerFieldName, "nama") || strings.Contains(lowerFieldName, "name") {
			score = 70
		} else if strings.Contains(lowerFieldName, "title") ||
			strings.Contains(lowerFieldName, "label") ||
			strings.Contains(lowerFieldName, "judul") {
			score = 60
		}

		if score > 0 {
			candidates = append(candidates, struct {
				name  string
				value string
				score int
			}{fieldName, strVal, score})
		}
	}

	// Return highest score
	if len(candidates) > 0 {
		highest := candidates[0]
		for _, candidate := range candidates {
			if candidate.score > highest.score {
				highest = candidate
			}
		}
		return highest.value
	}

	return ""
}

// GetClientCount returns the number of connected clients
func (s *SSEService) GetClientCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.clients)
}

// FormatSSEMessage formats an event as SSE protocol message
func FormatSSEMessage(event SSEEvent) string {
	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("Error marshaling SSE event: %v", err)
		return ""
	}
	return fmt.Sprintf("data: %s\n\n", data)
}
