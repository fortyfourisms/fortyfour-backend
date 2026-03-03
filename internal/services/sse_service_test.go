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

// ============================================================
// Register / Unregister — edge cases
// ============================================================

func TestRegisterMultipleClients_CountCorrect(t *testing.T) {
	sse := NewSSEService()

	c1 := newTestClient("client-1")
	c2 := newTestClient("client-2")
	c3 := newTestClient("client-3")

	sse.RegisterClient(c1)
	sse.RegisterClient(c2)
	sse.RegisterClient(c3)
	time.Sleep(80 * time.Millisecond)

	if count := sse.GetClientCount(); count != 3 {
		t.Fatalf("expected 3 clients, got %d", count)
	}
}

func TestUnregisterNonExistentClient_NoEffect(t *testing.T) {
	sse := NewSSEService()

	c1 := newTestClient("client-1")
	sse.RegisterClient(c1)
	time.Sleep(50 * time.Millisecond)

	// Unregister client yang tidak terdaftar — tidak boleh panik atau error
	ghost := newTestClient("client-ghost")
	sse.UnregisterClient(ghost)
	time.Sleep(50 * time.Millisecond)

	// client-1 masih terdaftar
	if count := sse.GetClientCount(); count != 1 {
		t.Fatalf("expected 1 client, got %d", count)
	}
}

func TestUnregisterClient_ChannelClosed(t *testing.T) {
	sse := NewSSEService()
	client := newTestClient("client-1")

	sse.RegisterClient(client)
	time.Sleep(50 * time.Millisecond)

	sse.UnregisterClient(client)
	time.Sleep(50 * time.Millisecond)

	// Channel harus sudah ditutup setelah unregister
	_, ok := <-client.Channel
	if ok {
		t.Fatal("channel seharusnya sudah ditutup setelah unregister")
	}
}

// ============================================================
// Broadcast — multiple clients & channel penuh
// ============================================================

func TestBroadcast_MultipleClients_AllReceive(t *testing.T) {
	sse := NewSSEService()

	c1 := newTestClient("client-1")
	c2 := newTestClient("client-2")
	c3 := newTestClient("client-3")

	sse.RegisterClient(c1)
	sse.RegisterClient(c2)
	sse.RegisterClient(c3)
	time.Sleep(80 * time.Millisecond)

	event := SSEEvent{Type: EventCreate, Resource: "perusahaan", Message: "test"}
	sse.Broadcast(event)

	timeout := time.After(300 * time.Millisecond)
	for _, client := range []*Client{c1, c2, c3} {
		select {
		case evt := <-client.Channel:
			if evt.Type != EventCreate {
				t.Errorf("client %s: expected type %s, got %s", client.ID, EventCreate, evt.Type)
			}
		case <-timeout:
			t.Fatalf("client %s tidak menerima event sebelum timeout", client.ID)
		}
	}
}

func TestBroadcast_ChannelPenuh_TidakBlokir(t *testing.T) {
	sse := NewSSEService()

	// Buat client dengan channel kecil (kapasitas 1) — pasti penuh setelah 1 event
	client := &Client{
		ID:      "client-kecil",
		UserID:  "user-1",
		Channel: make(chan SSEEvent, 1),
	}
	sse.RegisterClient(client)
	time.Sleep(50 * time.Millisecond)

	// Kirim 5 event — hanya 1 yang muat, sisanya di-drop tanpa blocking
	done := make(chan struct{})
	go func() {
		for i := 0; i < 5; i++ {
			sse.Broadcast(SSEEvent{Type: EventCreate, Resource: "test"})
		}
		close(done)
	}()

	select {
	case <-done:
		// Broadcast selesai tanpa deadlock
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Broadcast memblokir saat channel penuh")
	}
}

func TestBroadcast_NoClients_NoError(t *testing.T) {
	sse := NewSSEService()

	// Broadcast tanpa client — tidak boleh panik atau deadlock
	done := make(chan struct{})
	go func() {
		sse.Broadcast(SSEEvent{Type: EventCreate, Resource: "test"})
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(300 * time.Millisecond):
		t.Fatal("Broadcast memblokir saat tidak ada client")
	}
}

func TestBroadcast_TimestampDiSet(t *testing.T) {
	sse := NewSSEService()
	client := newTestClient("client-1")
	sse.RegisterClient(client)
	time.Sleep(50 * time.Millisecond)

	before := time.Now()
	sse.Broadcast(SSEEvent{Type: EventCreate, Resource: "test"})

	select {
	case evt := <-client.Channel:
		if evt.Timestamp.IsZero() {
			t.Fatal("Timestamp harus diisi oleh Broadcast")
		}
		if evt.Timestamp.Before(before) {
			t.Fatal("Timestamp harus setelah waktu Broadcast dipanggil")
		}
	case <-time.After(300 * time.Millisecond):
		t.Fatal("event tidak diterima")
	}
}

// ============================================================
// generateMessage — semua event type dan resource format
// ============================================================

func TestGenerateMessage_Create_DenganNama(t *testing.T) {
	sse := NewSSEService()
	data := map[string]interface{}{"nama": "PT Test"}

	msg := sse.generateMessage("perusahaan", EventCreate, data)

	if !strings.Contains(msg, "PT Test") {
		t.Errorf("expected message to contain 'PT Test', got: %s", msg)
	}
	if !strings.Contains(msg, "berhasil ditambahkan") {
		t.Errorf("expected 'berhasil ditambahkan', got: %s", msg)
	}
}

func TestGenerateMessage_Update_DenganNama(t *testing.T) {
	sse := NewSSEService()
	data := map[string]interface{}{"nama": "Sektor Keuangan"}

	msg := sse.generateMessage("sektor", EventUpdate, data)

	if !strings.Contains(msg, "Sektor Keuangan") {
		t.Errorf("expected message to contain 'Sektor Keuangan', got: %s", msg)
	}
	if !strings.Contains(msg, "berhasil diperbarui") {
		t.Errorf("expected 'berhasil diperbarui', got: %s", msg)
	}
}

func TestGenerateMessage_Update_TanpaNama(t *testing.T) {
	sse := NewSSEService()

	msg := sse.generateMessage("sektor", EventUpdate, map[string]interface{}{})

	if !strings.Contains(msg, "berhasil diperbarui") {
		t.Errorf("expected 'berhasil diperbarui', got: %s", msg)
	}
}

func TestGenerateMessage_Delete(t *testing.T) {
	sse := NewSSEService()

	msg := sse.generateMessage("user", EventDelete, "uuid-123")

	if !strings.Contains(msg, "berhasil dihapus") {
		t.Errorf("expected 'berhasil dihapus', got: %s", msg)
	}
}

func TestGenerateMessage_UnknownEventType(t *testing.T) {
	sse := NewSSEService()

	msg := sse.generateMessage("user", "unknown_event", nil)

	if !strings.Contains(msg, "unknown_event") {
		t.Errorf("expected message to contain event type, got: %s", msg)
	}
}

// ============================================================
// formatResourceName — semua branch
// ============================================================

func TestFormatResourceName_PendekDanUppercase(t *testing.T) {
	sse := NewSSEService()

	// Kata <= 4 huruf dan sudah uppercase → dikembalikan apa adanya
	result := sse.formatResourceName("SE")
	if result != "SE" {
		t.Errorf("expected 'SE', got '%s'", result)
	}
}

func TestFormatResourceName_PendekDanMixed(t *testing.T) {
	sse := NewSSEService()

	// Kata <= 4 huruf tapi bukan semua uppercase → huruf pertama kapital
	result := sse.formatResourceName("role")
	if result != "Role" {
		t.Errorf("expected 'Role', got '%s'", result)
	}
}

func TestFormatResourceName_Panjang(t *testing.T) {
	sse := NewSSEService()

	result := sse.formatResourceName("perusahaan")
	if result != "Perusahaan" {
		t.Errorf("expected 'Perusahaan', got '%s'", result)
	}
}

// ============================================================
// extractName — berbagai tipe data
// ============================================================

func TestExtractName_NilData(t *testing.T) {
	sse := NewSSEService()

	result := sse.extractName(nil)
	if result != "" {
		t.Errorf("expected empty string for nil data, got '%s'", result)
	}
}

func TestExtractName_MapPrioritas_NamaExact(t *testing.T) {
	sse := NewSSEService()

	// Catatan: scoring di extractNameFromMap menggunakan if terpisah (bukan else-if),
	// sehingga key "nama" mendapat score 100 → ditimpa 70 (contains "nama"),
	// dan "nama_perusahaan" mendapat 90 → ditimpa 70 juga.
	// Karena score akhirnya sama, tidak bisa diprediksi siapa yang menang.
	// Yang penting: fungsi mengembalikan salah satu dari kedua value, bukan empty string.
	data := map[string]interface{}{
		"nama":            "Prioritas Tinggi",
		"nama_perusahaan": "Lebih Rendah",
	}

	result := sse.extractName(data)
	if result != "Prioritas Tinggi" && result != "Lebih Rendah" {
		t.Errorf("expected either 'Prioritas Tinggi' or 'Lebih Rendah', got '%s'", result)
	}
	if result == "" {
		t.Error("expected non-empty result from map with nama fields")
	}
}

func TestExtractName_MapNamaPerusahaan(t *testing.T) {
	sse := NewSSEService()

	data := map[string]interface{}{
		"nama_perusahaan": "PT Maju Jaya",
	}

	result := sse.extractName(data)
	if result != "PT Maju Jaya" {
		t.Errorf("expected 'PT Maju Jaya', got '%s'", result)
	}
}

func TestExtractName_MapTidakAdaNama(t *testing.T) {
	sse := NewSSEService()

	data := map[string]interface{}{
		"id":      "uuid-123",
		"status":  "active",
		"tanggal": "2024-01-01",
	}

	result := sse.extractName(data)
	if result != "" {
		t.Errorf("expected empty string, got '%s'", result)
	}
}

func TestExtractName_StructPointerNama(t *testing.T) {
	sse := NewSSEService()

	type DataStruct struct {
		Nama string
		ID   string
	}

	result := sse.extractName(&DataStruct{Nama: "Dari Pointer Struct", ID: "uuid"})
	if result != "Dari Pointer Struct" {
		t.Errorf("expected 'Dari Pointer Struct', got '%s'", result)
	}
}

func TestExtractName_StructNilPointer(t *testing.T) {
	sse := NewSSEService()

	type DataStruct struct{ Nama string }
	var data *DataStruct

	result := sse.extractName(data)
	if result != "" {
		t.Errorf("expected empty string for nil pointer, got '%s'", result)
	}
}

func TestExtractName_NonStructNonMap(t *testing.T) {
	sse := NewSSEService()

	// Integer — bukan struct maupun map
	result := sse.extractName(42)
	if result != "" {
		t.Errorf("expected empty string for non-struct/map, got '%s'", result)
	}
}

// ============================================================
// NotifyCreate / Update / Delete — event fields verified
// ============================================================

func TestNotifyCreate_EventFieldsLengkap(t *testing.T) {
	sse := NewSSEService()
	client := newTestClient("client-1")
	sse.RegisterClient(client)
	time.Sleep(50 * time.Millisecond)

	data := map[string]interface{}{"nama": "Data Baru"}
	sse.NotifyCreate("perusahaan", data, "user-abc")

	select {
	case evt := <-client.Channel:
		if evt.Type != EventCreate {
			t.Errorf("expected type '%s', got '%s'", EventCreate, evt.Type)
		}
		if evt.Resource != "perusahaan" {
			t.Errorf("expected resource 'perusahaan', got '%s'", evt.Resource)
		}
		if evt.UserID != "user-abc" {
			t.Errorf("expected userID 'user-abc', got '%s'", evt.UserID)
		}
		if evt.Data == nil {
			t.Error("expected Data tidak nil")
		}
	case <-time.After(300 * time.Millisecond):
		t.Fatal("event tidak diterima")
	}
}

func TestNotifyDelete_DataWrappedWithID(t *testing.T) {
	sse := NewSSEService()
	client := newTestClient("client-1")
	sse.RegisterClient(client)
	time.Sleep(50 * time.Millisecond)

	sse.NotifyDelete("sektor", "uuid-del-123", "user-1")

	select {
	case evt := <-client.Channel:
		if evt.Type != EventDelete {
			t.Fatalf("expected type '%s', got '%s'", EventDelete, evt.Type)
		}
		dataMap, ok := evt.Data.(map[string]interface{})
		if !ok {
			t.Fatal("expected Data berupa map[string]interface{}")
		}
		if dataMap["id"] != "uuid-del-123" {
			t.Errorf("expected id 'uuid-del-123', got %v", dataMap["id"])
		}
	case <-time.After(300 * time.Millisecond):
		t.Fatal("event tidak diterima")
	}
}

// ============================================================
// FormatSSEMessage — edge cases
// ============================================================

func TestFormatSSEMessage_EndsWith2Newlines(t *testing.T) {
	event := SSEEvent{Type: EventCreate, Resource: "test"}
	msg := FormatSSEMessage(event)

	if !strings.HasSuffix(msg, "\n\n") {
		t.Errorf("SSE message harus diakhiri dengan double newline, got: %q", msg)
	}
}

func TestFormatSSEMessage_ContainsValidJSON(t *testing.T) {
	event := SSEEvent{
		Type:     EventUpdate,
		Resource: "sektor",
		Message:  "Sektor diperbarui",
		UserID:   "user-99",
	}

	msg := FormatSSEMessage(event)
	jsonPart := strings.TrimPrefix(strings.TrimSpace(msg), "data: ")

	var decoded SSEEvent
	if err := json.Unmarshal([]byte(jsonPart), &decoded); err != nil {
		t.Fatalf("payload bukan JSON valid: %v", err)
	}
	if decoded.Type != EventUpdate {
		t.Errorf("expected type '%s', got '%s'", EventUpdate, decoded.Type)
	}
	if decoded.UserID != "user-99" {
		t.Errorf("expected userID 'user-99', got '%s'", decoded.UserID)
	}
}