package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"fortyfour-backend/internal/middleware"
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

// ============================================================
// TestSSEHandler_GetStats — tambahan
// ============================================================

func TestSSEHandler_GetStats_ContentTypeJSON(t *testing.T) {
	h := &SSEHandler{sseService: services.NewSSEService()}

	req := httptest.NewRequest(http.MethodGet, "/api/events/stats", nil)
	rr := httptest.NewRecorder()

	h.GetStats(rr, req)

	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}
}

func TestSSEHandler_GetStats_MultipleClients(t *testing.T) {
	svc := setupSSEServiceWithClients(t, 3)
	h := &SSEHandler{sseService: svc}

	req := httptest.NewRequest(http.MethodGet, "/api/events/stats", nil)
	rr := httptest.NewRecorder()

	h.GetStats(rr, req)

	var resp map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &resp)

	if int(resp["connected_clients"].(float64)) != 3 {
		t.Errorf("expected 3 connected clients, got %v", resp["connected_clients"])
	}
}

// ============================================================
// TestSSEHandler_HandleSSE
// HandleSSE adalah long-lived streaming — ditest dengan context
// cancellation untuk mensimulasikan client disconnect.
// ============================================================

func TestSSEHandler_HandleSSE_SetsSSEHeaders(t *testing.T) {
	svc := services.NewSSEService()
	h := NewSSEHandler(svc)

	ctx, cancel := context.WithCancel(context.Background())
	req := httptest.NewRequest(http.MethodGet, "/api/events", nil).WithContext(ctx)
	rr := httptest.NewRecorder()

	// Batalkan context segera setelah handler mulai berjalan
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	h.HandleSSE(rr, req)

	if ct := rr.Header().Get("Content-Type"); ct != "text/event-stream" {
		t.Errorf("expected Content-Type text/event-stream, got %s", ct)
	}
	if cc := rr.Header().Get("Cache-Control"); cc != "no-cache" {
		t.Errorf("expected Cache-Control no-cache, got %s", cc)
	}
	if conn := rr.Header().Get("Connection"); conn != "keep-alive" {
		t.Errorf("expected Connection keep-alive, got %s", conn)
	}
}

func TestSSEHandler_HandleSSE_SendsInitialConnectedEvent(t *testing.T) {
	svc := services.NewSSEService()
	h := NewSSEHandler(svc)

	ctx, cancel := context.WithCancel(context.Background())
	req := httptest.NewRequest(http.MethodGet, "/api/events", nil).WithContext(ctx)
	rr := httptest.NewRecorder()

	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	h.HandleSSE(rr, req)

	body := rr.Body.String()
	if !strings.Contains(body, "data: ") {
		t.Fatal("response harus mengandung SSE data line")
	}

	// Parse event pertama
	lines := strings.Split(body, "\n")
	var firstEventJSON string
	for _, line := range lines {
		if strings.HasPrefix(line, "data: ") {
			firstEventJSON = strings.TrimPrefix(line, "data: ")
			break
		}
	}

	var event map[string]interface{}
	if err := json.Unmarshal([]byte(firstEventJSON), &event); err != nil {
		t.Fatalf("event pertama bukan JSON valid: %v", err)
	}

	if event["type"] != "connected" {
		t.Errorf("expected type 'connected', got %v", event["type"])
	}
	if event["resource"] != "system" {
		t.Errorf("expected resource 'system', got %v", event["resource"])
	}
}

func TestSSEHandler_HandleSSE_ClientIDDiInitialData(t *testing.T) {
	svc := services.NewSSEService()
	h := NewSSEHandler(svc)

	ctx, cancel := context.WithCancel(context.Background())
	req := httptest.NewRequest(http.MethodGet, "/api/events", nil).WithContext(ctx)
	rr := httptest.NewRecorder()

	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	h.HandleSSE(rr, req)

	body := rr.Body.String()
	var firstJSON string
	for _, line := range strings.Split(body, "\n") {
		if strings.HasPrefix(line, "data: ") {
			firstJSON = strings.TrimPrefix(line, "data: ")
			break
		}
	}

	var event map[string]interface{}
	json.Unmarshal([]byte(firstJSON), &event)

	data, ok := event["data"].(map[string]interface{})
	if !ok {
		t.Fatal("event.data harus berupa object")
	}
	if data["client_id"] == "" || data["client_id"] == nil {
		t.Error("client_id harus ada dan tidak kosong di initial event")
	}
	if data["message"] != "Connected to SSE" {
		t.Errorf("expected message 'Connected to SSE', got %v", data["message"])
	}
}

func TestSSEHandler_HandleSSE_RegistersAndUnregistersClient(t *testing.T) {
	svc := services.NewSSEService()
	h := NewSSEHandler(svc)

	ctx, cancel := context.WithCancel(context.Background())
	req := httptest.NewRequest(http.MethodGet, "/api/events", nil).WithContext(ctx)
	rr := httptest.NewRecorder()

	done := make(chan struct{})
	go func() {
		h.HandleSSE(rr, req)
		close(done)
	}()

	// Tunggu sampai client terdaftar
	deadline := time.Now().Add(500 * time.Millisecond)
	for time.Now().Before(deadline) {
		if svc.GetClientCount() == 1 {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	if svc.GetClientCount() != 1 {
		t.Fatalf("expected 1 client setelah connect, got %d", svc.GetClientCount())
	}

	// Disconnect
	cancel()
	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("HandleSSE tidak selesai setelah context cancel")
	}

	// Tunggu sampai unregister diproses
	time.Sleep(80 * time.Millisecond)
	if svc.GetClientCount() != 0 {
		t.Errorf("expected 0 clients setelah disconnect, got %d", svc.GetClientCount())
	}
}

func TestSSEHandler_HandleSSE_UserIDFromContext(t *testing.T) {
	svc := services.NewSSEService()
	h := NewSSEHandler(svc)

	ctx, cancel := context.WithCancel(context.Background())
	ctx = context.WithValue(ctx, middleware.UserIDKey, "user-xyz")
	req := httptest.NewRequest(http.MethodGet, "/api/events", nil).WithContext(ctx)
	rr := httptest.NewRecorder()

	done := make(chan struct{})
	go func() {
		h.HandleSSE(rr, req)
		close(done)
	}()

	// Tunggu client terdaftar
	time.Sleep(80 * time.Millisecond)
	cancel()

	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("HandleSSE tidak selesai")
	}
}

func TestSSEHandler_HandleSSE_TanpaUserIDContext(t *testing.T) {
	svc := services.NewSSEService()
	h := NewSSEHandler(svc)

	ctx, cancel := context.WithCancel(context.Background())
	// Tidak inject UserID ke context
	req := httptest.NewRequest(http.MethodGet, "/api/events", nil).WithContext(ctx)
	rr := httptest.NewRecorder()

	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	// Tidak boleh panik meski tidak ada UserID
	h.HandleSSE(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
}

func TestSSEHandler_HandleSSE_StreamsEventFromBroadcast(t *testing.T) {
	svc := services.NewSSEService()
	h := NewSSEHandler(svc)

	ctx, cancel := context.WithCancel(context.Background())
	req := httptest.NewRequest(http.MethodGet, "/api/events", nil).WithContext(ctx)
	rr := httptest.NewRecorder()

	done := make(chan struct{})
	go func() {
		h.HandleSSE(rr, req)
		close(done)
	}()

	// Tunggu client terdaftar
	deadline := time.Now().Add(500 * time.Millisecond)
	for time.Now().Before(deadline) {
		if svc.GetClientCount() == 1 {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	// Broadcast event
	svc.Broadcast(services.SSEEvent{
		Type:     services.EventCreate,
		Resource: "perusahaan",
		Message:  "PT Baru berhasil ditambahkan",
	})

	// Beri waktu event dikirim
	time.Sleep(80 * time.Millisecond)
	cancel()

	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("HandleSSE tidak selesai")
	}

	// Verifikasi event terbroadcast muncul di response body
	body := rr.Body.String()
	if !strings.Contains(body, "perusahaan") {
		t.Errorf("expected broadcast event 'perusahaan' di body, got:\n%s", body)
	}
}