package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"fortyfour-backend/internal/services"
)

func setupSSEServiceWithClients(t *testing.T, n int) *services.SSEService {
	svc := services.NewSSEService()
	// register n clients
	for i := 0; i < n; i++ {
		c := &services.Client{
			ID:      fmt.Sprintf("client-%d", i),
			UserID:  "user-1",
			Channel: make(chan services.SSEEvent, 10),
		}
		svc.RegisterClient(c)
	}

	// wait until service registers clients (with a small timeout)
	deadline := time.Now().Add(1 * time.Second)
	for time.Now().Before(deadline) {
		if svc.GetClientCount() >= n {
			return svc
		}
		time.Sleep(10 * time.Millisecond)
	}
	if svc.GetClientCount() < n {
		t.Fatalf("expected %d clients registered but found %d", n, svc.GetClientCount())
	}
	return svc
}

func TestSSEHandler_GetStats_NoClients(t *testing.T) {
	svc := services.NewSSEService()
	h := &SSEHandler{sseService: svc}

	req := httptest.NewRequest(http.MethodGet, "/sse/stats", nil)
	rr := httptest.NewRecorder()

	h.GetStats(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200 OK, got %d", rr.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json response: %v", err)
	}

	if v, ok := resp["connected_clients"]; !ok {
		t.Fatalf("response missing connected_clients field")
	} else {
		// numbers are decoded as float64 by default
		if int(v.(float64)) != 0 {
			t.Fatalf("expected 0 connected clients, got %v", v)
		}
	}
}

func TestSSEHandler_GetStats_WithClients(t *testing.T) {
	svc := setupSSEServiceWithClients(t, 1)
	h := &SSEHandler{sseService: svc}

	req := httptest.NewRequest(http.MethodGet, "/sse/stats", nil)
	rr := httptest.NewRecorder()

	h.GetStats(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200 OK, got %d", rr.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json response: %v", err)
	}

	if v, ok := resp["connected_clients"]; !ok {
		t.Fatalf("response missing connected_clients field")
	} else {
		if int(v.(float64)) != 1 {
			t.Fatalf("expected 1 connected client, got %v", v)
		}
	}
}
