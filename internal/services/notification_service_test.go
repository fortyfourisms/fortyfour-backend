package services

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"fortyfour-backend/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================
// Mock Redis untuk NotificationService
// ============================================================

type notifTestRedis struct {
	data    map[string]string
	setErr  error
	getErr  error
	delErr  error
}

func newNotifTestRedis() *notifTestRedis {
	return &notifTestRedis{data: make(map[string]string)}
}

func (r *notifTestRedis) Set(key string, value interface{}, ttl time.Duration) error {
	if r.setErr != nil {
		return r.setErr
	}
	if v, ok := value.(string); ok {
		r.data[key] = v
	}
	return nil
}

func (r *notifTestRedis) Get(key string) (string, error) {
	if r.getErr != nil {
		return "", r.getErr
	}
	v, ok := r.data[key]
	if !ok {
		return "", errors.New("key not found")
	}
	return v, nil
}

func (r *notifTestRedis) Delete(key string) error {
	if r.delErr != nil {
		return r.delErr
	}
	delete(r.data, key)
	return nil
}

func (r *notifTestRedis) Exists(key string) (bool, error) {
	_, ok := r.data[key]
	return ok, nil
}

func (r *notifTestRedis) Scan(pattern string) ([]string, error) {
	return []string{}, nil
}

func (r *notifTestRedis) Close() error { return nil }

// Helper: simpan notifikasi langsung ke mock redis
func seedNotifications(rc *notifTestRedis, userID string, notifs []models.Notification) {
	data, _ := json.Marshal(notifs)
	rc.data["notif:"+userID] = string(data)
}

// ============================================================
// TEST: redisKey
// ============================================================

func TestNotificationService_RedisKey(t *testing.T) {
	svc := NewNotificationService(newNotifTestRedis())
	assert.Equal(t, "notif:user-123", svc.redisKey("user-123"))
	assert.Equal(t, "notif:abc-def", svc.redisKey("abc-def"))
}

// ============================================================
// TEST: GetAll
// ============================================================

func TestNotificationService_GetAll_EmptyKey_ReturnsEmptySlice(t *testing.T) {
	svc := NewNotificationService(newNotifTestRedis())

	result, err := svc.GetAll("user-1")

	require.NoError(t, err)
	assert.Empty(t, result, "jika key tidak ada harus return slice kosong")
}

func TestNotificationService_GetAll_WithData(t *testing.T) {
	rc := newNotifTestRedis()
	notifs := []models.Notification{
		{ID: "n1", UserID: "u1", Type: models.NotifLoginFailed, Message: "Login gagal", Read: false},
		{ID: "n2", UserID: "u1", Type: models.NotifPasswordExpirySoon, Message: "Password segera expired", Read: true},
	}
	seedNotifications(rc, "u1", notifs)

	svc := NewNotificationService(rc)
	result, err := svc.GetAll("u1")

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "n1", result[0].ID)
	assert.Equal(t, "n2", result[1].ID)
}

func TestNotificationService_GetAll_InvalidJSON_ReturnsError(t *testing.T) {
	rc := newNotifTestRedis()
	rc.data["notif:u1"] = `ini bukan json[`

	svc := NewNotificationService(rc)
	_, err := svc.GetAll("u1")

	assert.Error(t, err, "JSON rusak harus return error")
}

func TestNotificationService_GetAll_RedisGetError_ReturnsEmpty(t *testing.T) {
	rc := newNotifTestRedis()
	rc.getErr = errors.New("connection refused")

	svc := NewNotificationService(rc)
	result, err := svc.GetAll("u1")

	// Ketika key tidak ada (error), return empty slice tanpa error
	require.NoError(t, err)
	assert.Empty(t, result)
}

// ============================================================
// TEST: Push
// ============================================================

func TestNotificationService_Push_AddsNotification(t *testing.T) {
	rc := newNotifTestRedis()
	svc := NewNotificationService(rc)

	err := svc.Push("u1", models.NotifLoginFailed, "Login gagal dari IP baru")

	require.NoError(t, err)

	result, _ := svc.GetAll("u1")
	assert.Len(t, result, 1)
	assert.Equal(t, models.NotifLoginFailed, result[0].Type)
	assert.Equal(t, "Login gagal dari IP baru", result[0].Message)
	assert.Equal(t, "u1", result[0].UserID)
	assert.False(t, result[0].Read, "notifikasi baru harus unread")
	assert.NotEmpty(t, result[0].ID)
}

func TestNotificationService_Push_PrependsTerbaru(t *testing.T) {
	rc := newNotifTestRedis()
	svc := NewNotificationService(rc)

	svc.Push("u1", models.NotifLoginFailed, "pertama")
	svc.Push("u1", models.NotifPasswordExpirySoon, "kedua")

	result, _ := svc.GetAll("u1")
	assert.Len(t, result, 2)
	// Push kedua harus di index 0 (prepend = terbaru di atas)
	assert.Equal(t, "kedua", result[0].Message)
	assert.Equal(t, "pertama", result[1].Message)
}

func TestNotificationService_Push_MultipleTypes(t *testing.T) {
	rc := newNotifTestRedis()
	svc := NewNotificationService(rc)

	svc.Push("u1", models.NotifLoginFailed, "msg1")
	svc.Push("u1", models.NotifAccountSuspended, "msg2")
	svc.Push("u1", models.NotifPasswordExpired, "msg3")
	svc.Push("u1", models.NotifPasswordExpirySoon, "msg4")

	result, _ := svc.GetAll("u1")
	assert.Len(t, result, 4)
}

func TestNotificationService_Push_SetError_ReturnsError(t *testing.T) {
	rc := newNotifTestRedis()
	rc.setErr = errors.New("redis full")

	svc := NewNotificationService(rc)
	err := svc.Push("u1", models.NotifLoginFailed, "msg")

	assert.Error(t, err)
}

// ============================================================
// TEST: MarkRead
// ============================================================

func TestNotificationService_MarkRead_Success(t *testing.T) {
	rc := newNotifTestRedis()
	notifs := []models.Notification{
		{ID: "n1", UserID: "u1", Type: models.NotifLoginFailed, Read: false},
		{ID: "n2", UserID: "u1", Type: models.NotifPasswordExpirySoon, Read: false},
	}
	seedNotifications(rc, "u1", notifs)

	svc := NewNotificationService(rc)
	err := svc.MarkRead("u1", "n1")

	require.NoError(t, err)

	result, _ := svc.GetAll("u1")
	var n1 *models.Notification
	for i := range result {
		if result[i].ID == "n1" {
			n1 = &result[i]
		}
	}
	require.NotNil(t, n1)
	assert.True(t, n1.Read, "n1 harus sudah ditandai read")

	// n2 tetap unread
	for _, n := range result {
		if n.ID == "n2" {
			assert.False(t, n.Read)
		}
	}
}

func TestNotificationService_MarkRead_NotFound_ReturnsError(t *testing.T) {
	rc := newNotifTestRedis()
	notifs := []models.Notification{
		{ID: "n1", UserID: "u1", Read: false},
	}
	seedNotifications(rc, "u1", notifs)

	svc := NewNotificationService(rc)
	err := svc.MarkRead("u1", "tidak-ada")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "tidak ditemukan")
}

func TestNotificationService_MarkRead_EmptyUser_ReturnsError(t *testing.T) {
	svc := NewNotificationService(newNotifTestRedis())

	err := svc.MarkRead("u1", "n1")

	assert.Error(t, err, "user tanpa notifikasi harus error karena notif tidak ditemukan")
}

// ============================================================
// TEST: MarkAllRead
// ============================================================

func TestNotificationService_MarkAllRead_Success(t *testing.T) {
	rc := newNotifTestRedis()
	notifs := []models.Notification{
		{ID: "n1", Read: false},
		{ID: "n2", Read: false},
		{ID: "n3", Read: true},
	}
	seedNotifications(rc, "u1", notifs)

	svc := NewNotificationService(rc)
	err := svc.MarkAllRead("u1")

	require.NoError(t, err)

	result, _ := svc.GetAll("u1")
	for _, n := range result {
		assert.True(t, n.Read, "semua notifikasi harus read setelah MarkAllRead")
	}
}

func TestNotificationService_MarkAllRead_EmptyList_NoError(t *testing.T) {
	svc := NewNotificationService(newNotifTestRedis())

	err := svc.MarkAllRead("u1")

	assert.NoError(t, err)
}

// ============================================================
// TEST: Delete
// ============================================================

func TestNotificationService_Delete_Success(t *testing.T) {
	rc := newNotifTestRedis()
	notifs := []models.Notification{
		{ID: "n1", Message: "pertama"},
		{ID: "n2", Message: "kedua"},
	}
	seedNotifications(rc, "u1", notifs)

	svc := NewNotificationService(rc)
	err := svc.Delete("u1", "n1")

	require.NoError(t, err)

	result, _ := svc.GetAll("u1")
	assert.Len(t, result, 1)
	assert.Equal(t, "n2", result[0].ID)
}

func TestNotificationService_Delete_NotFound_ReturnsError(t *testing.T) {
	rc := newNotifTestRedis()
	notifs := []models.Notification{
		{ID: "n1", Message: "ada"},
	}
	seedNotifications(rc, "u1", notifs)

	svc := NewNotificationService(rc)
	err := svc.Delete("u1", "tidak-ada")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "tidak ditemukan")
}

func TestNotificationService_Delete_LastNotif_LeavesEmptySlice(t *testing.T) {
	rc := newNotifTestRedis()
	notifs := []models.Notification{
		{ID: "satu-satunya"},
	}
	seedNotifications(rc, "u1", notifs)

	svc := NewNotificationService(rc)
	err := svc.Delete("u1", "satu-satunya")

	require.NoError(t, err)

	result, _ := svc.GetAll("u1")
	assert.Empty(t, result)
}

// ============================================================
// TEST: DeleteAll
// ============================================================

func TestNotificationService_DeleteAll_Success(t *testing.T) {
	rc := newNotifTestRedis()
	notifs := []models.Notification{
		{ID: "n1"}, {ID: "n2"}, {ID: "n3"},
	}
	seedNotifications(rc, "u1", notifs)

	svc := NewNotificationService(rc)
	err := svc.DeleteAll("u1")

	require.NoError(t, err)

	// Key harus sudah terhapus dari redis
	_, keyExists := rc.data["notif:u1"]
	assert.False(t, keyExists, "key notif:u1 harus sudah dihapus")
}

func TestNotificationService_DeleteAll_EmptyUser_NoError(t *testing.T) {
	svc := NewNotificationService(newNotifTestRedis())

	err := svc.DeleteAll("u1")

	// Delete key yang tidak ada tidak error (redis delete behavior)
	assert.NoError(t, err)
}

func TestNotificationService_DeleteAll_DeleteError_ReturnsError(t *testing.T) {
	rc := newNotifTestRedis()
	rc.delErr = errors.New("redis down")

	svc := NewNotificationService(rc)
	err := svc.DeleteAll("u1")

	assert.Error(t, err)
}

// ============================================================
// TEST: HasPasswordExpirySoonNotif
// ============================================================

func TestNotificationService_HasPasswordExpirySoonNotif_ReturnsTrueIfUnread(t *testing.T) {
	rc := newNotifTestRedis()
	notifs := []models.Notification{
		{ID: "n1", Type: models.NotifPasswordExpirySoon, Read: false},
	}
	seedNotifications(rc, "u1", notifs)

	svc := NewNotificationService(rc)
	has, err := svc.HasPasswordExpirySoonNotif("u1")

	require.NoError(t, err)
	assert.True(t, has)
}

func TestNotificationService_HasPasswordExpirySoonNotif_ReturnsFalseIfRead(t *testing.T) {
	rc := newNotifTestRedis()
	notifs := []models.Notification{
		{ID: "n1", Type: models.NotifPasswordExpirySoon, Read: true}, // sudah dibaca
	}
	seedNotifications(rc, "u1", notifs)

	svc := NewNotificationService(rc)
	has, err := svc.HasPasswordExpirySoonNotif("u1")

	require.NoError(t, err)
	assert.False(t, has, "notif yang sudah dibaca tidak dihitung")
}

func TestNotificationService_HasPasswordExpirySoonNotif_ReturnsFalseIfNoNotif(t *testing.T) {
	svc := NewNotificationService(newNotifTestRedis())

	has, err := svc.HasPasswordExpirySoonNotif("u1")

	require.NoError(t, err)
	assert.False(t, has)
}

func TestNotificationService_HasPasswordExpirySoonNotif_ReturnsFalseForOtherType(t *testing.T) {
	rc := newNotifTestRedis()
	notifs := []models.Notification{
		{ID: "n1", Type: models.NotifLoginFailed, Read: false}, // bukan PasswordExpirySoon
	}
	seedNotifications(rc, "u1", notifs)

	svc := NewNotificationService(rc)
	has, err := svc.HasPasswordExpirySoonNotif("u1")

	require.NoError(t, err)
	assert.False(t, has)
}

// ============================================================
// TEST: IsolatePerUser
// ============================================================

func TestNotificationService_IsolatePerUser(t *testing.T) {
	rc := newNotifTestRedis()
	svc := NewNotificationService(rc)

	svc.Push("user-A", models.NotifLoginFailed, "untuk A")
	svc.Push("user-B", models.NotifPasswordExpired, "untuk B")

	resultA, _ := svc.GetAll("user-A")
	resultB, _ := svc.GetAll("user-B")

	assert.Len(t, resultA, 1)
	assert.Len(t, resultB, 1)
	assert.Equal(t, "untuk A", resultA[0].Message)
	assert.Equal(t, "untuk B", resultB[0].Message)
}