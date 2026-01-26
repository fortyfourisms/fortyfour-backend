package services

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

// helper client
func newTestClient(id string) *Client {
	return &Client{
		ID:      id,
		UserID:  "user-1",
		Channel: make(chan SSEEvent, 10),
	}
}

func TestRegisterAndUnregisterClient(t *testing.T) {
	sse := NewSSEService()

	client := newTestClient("client-1")
	sse.RegisterClient(client)

	time.Sleep(50 * time.Millisecond)

	if count := sse.GetClientCount(); count != 1 {
		t.Fatalf("expected 1 client, got %d", count)
	}

	sse.UnregisterClient(client)
	time.Sleep(50 * time.Millisecond)

	if count := sse.GetClientCount(); count != 0 {
		t.Fatalf("expected 0 client, got %d", count)
	}
}

func TestBroadcastEvent(t *testing.T) {
	sse := NewSSEService()
	client := newTestClient("client-1")
	sse.RegisterClient(client)
	time.Sleep(50 * time.Millisecond)

	event := SSEEvent{
		Type:     EventCreate,
		Resource: "user",
		Message:  "User baru berhasil ditambahkan",
		UserID:   "user-1",
	}

	sse.Broadcast(event)

	select {
	case received := <-client.Channel:
		if received.Type != EventCreate {
			t.Fatalf("expected event type %s, got %s", EventCreate, received.Type)
		}
		if received.Resource != "user" {
			t.Fatalf("expected resource user, got %s", received.Resource)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("did not receive broadcast event")
	}
}

func TestNotifyCreate_WithNameFromMap(t *testing.T) {
	sse := NewSSEService()
	client := newTestClient("client-1")
	sse.RegisterClient(client)
	time.Sleep(50 * time.Millisecond)

	data := map[string]interface{}{
		"nama": "Admin Sistem",
	}

	sse.NotifyCreate("role", data, "user-1")

	select {
	case evt := <-client.Channel:
		if !strings.Contains(evt.Message, "Admin Sistem") {
			t.Fatalf("expected message to contain name, got: %s", evt.Message)
		}
		if evt.Type != EventCreate {
			t.Fatalf("expected create event")
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("event not received")
	}
}

func TestNotifyUpdate_WithNameFromStruct(t *testing.T) {
	type Role struct {
		Name string
	}

	sse := NewSSEService()
	client := newTestClient("client-1")
	sse.RegisterClient(client)
	time.Sleep(50 * time.Millisecond)

	data := Role{Name: "Super Admin"}
	sse.NotifyUpdate("role", data, "user-1")

	select {
	case evt := <-client.Channel:
		if !strings.Contains(evt.Message, "Super Admin") {
			t.Fatalf("expected message to contain struct name, got: %s", evt.Message)
		}
		if evt.Type != EventUpdate {
			t.Fatalf("expected update event")
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("event not received")
	}
}

func TestNotifyDelete(t *testing.T) {
	sse := NewSSEService()
	client := newTestClient("client-1")
	sse.RegisterClient(client)
	time.Sleep(50 * time.Millisecond)

	sse.NotifyDelete("user", 123, "user-1")

	select {
	case evt := <-client.Channel:
		if evt.Type != EventDelete {
			t.Fatalf("expected delete event")
		}
		if evt.Data.(map[string]interface{})["id"] != 123 {
			t.Fatalf("expected id 123")
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("event not received")
	}
}

func TestGenerateMessage_NoName(t *testing.T) {
	sse := NewSSEService()

	msg := sse.generateMessage("user", EventCreate, map[string]interface{}{})
	if msg != "User baru berhasil ditambahkan" {
		t.Fatalf("unexpected message: %s", msg)
	}
}

func TestFormatSSEMessage(t *testing.T) {
	event := SSEEvent{
		Type:     EventCreate,
		Resource: "user",
		Message:  "User baru",
		UserID:   "1",
	}

	msg := FormatSSEMessage(event)

	if !strings.HasPrefix(msg, "data: ") {
		t.Fatal("SSE message must start with 'data: '")
	}

	var decoded SSEEvent
	jsonPart := strings.TrimPrefix(strings.TrimSpace(msg), "data: ")
	if err := json.Unmarshal([]byte(jsonPart), &decoded); err != nil {
		t.Fatalf("invalid JSON SSE payload: %v", err)
	}
}
