package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/middleware"
	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================
// Mock Redis untuk NotificationService di handler test
// ============================================================

type notifHandlerRedis struct {
	data   map[string]string
	setErr error
	getErr error
	delErr error
}

func newNotifHandlerRedis() *notifHandlerRedis {
	return &notifHandlerRedis{data: make(map[string]string)}
}

func (r *notifHandlerRedis) Set(key string, value interface{}, ttl time.Duration) error {
	if r.setErr != nil {
		return r.setErr
	}
	if v, ok := value.(string); ok {
		r.data[key] = v
	}
	return nil
}

func (r *notifHandlerRedis) Get(key string) (string, error) {
	if r.getErr != nil {
		return "", r.getErr
	}
	v, ok := r.data[key]
	if !ok {
		return "", errors.New("key not found")
	}
	return v, nil
}

func (r *notifHandlerRedis) Delete(key string) error {
	if r.delErr != nil {
		return r.delErr
	}
	delete(r.data, key)
	return nil
}

func (r *notifHandlerRedis) Exists(key string) (bool, error) {
	_, ok := r.data[key]
	return ok, nil
}

func (r *notifHandlerRedis) Scan(pattern string) ([]string, error) { return nil, nil }
func (r *notifHandlerRedis) Close() error                          { return nil }

// ============================================================
// Helper setup
// ============================================================

func setupNotificationHandler() (*NotificationHandler, *services.NotificationService, *notifHandlerRedis) {
	rc := newNotifHandlerRedis()
	svc := services.NewNotificationService(rc)
	handler := NewNotificationHandler(svc)
	return handler, svc, rc
}

func reqWithUserID(method, path, userID string) *http.Request {
	req := httptest.NewRequest(method, path, nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
	return req.WithContext(ctx)
}

func seedNotifInHandler(rc *notifHandlerRedis, userID string, notifs []models.Notification) {
	data, _ := json.Marshal(notifs)
	rc.data["notif:"+userID] = string(data)
}

// ============================================================
// TEST: GetAll
// ============================================================

func TestNotificationHandler_GetAll_Unauthorized(t *testing.T) {
	handler, _, _ := setupNotificationHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/notifications", nil)
	rr := httptest.NewRecorder()
	handler.GetAll(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestNotificationHandler_GetAll_EmptyList(t *testing.T) {
	handler, _, _ := setupNotificationHandler()

	req := reqWithUserID(http.MethodGet, "/api/notifications", "user-1")
	rr := httptest.NewRecorder()
	handler.GetAll(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp dto.NotificationListResponse
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&resp))
	assert.Empty(t, resp.Notifications)
	assert.Equal(t, 0, resp.UnreadCount)
}

func TestNotificationHandler_GetAll_WithUnreadNotifs(t *testing.T) {
	handler, _, rc := setupNotificationHandler()

	notifs := []models.Notification{
		{ID: "n1", Type: models.NotifLoginFailed, Message: "Login gagal", Read: false, CreatedAt: time.Now()},
		{ID: "n2", Type: models.NotifPasswordExpirySoon, Message: "Password mau expired", Read: true, CreatedAt: time.Now()},
		{ID: "n3", Type: models.NotifAccountSuspended, Message: "Akun suspend", Read: false, CreatedAt: time.Now()},
	}
	seedNotifInHandler(rc, "user-1", notifs)

	req := reqWithUserID(http.MethodGet, "/api/notifications", "user-1")
	rr := httptest.NewRecorder()
	handler.GetAll(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp dto.NotificationListResponse
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&resp))
	assert.Len(t, resp.Notifications, 3)
	assert.Equal(t, 2, resp.UnreadCount, "harus ada 2 notif unread")
}

func TestNotificationHandler_GetAll_ContentTypeJSON(t *testing.T) {
	handler, _, _ := setupNotificationHandler()

	req := reqWithUserID(http.MethodGet, "/api/notifications", "user-1")
	rr := httptest.NewRecorder()
	handler.GetAll(rr, req)

	assert.Contains(t, rr.Header().Get("Content-Type"), "application/json")
}

func TestNotificationHandler_GetAll_IsolatesPerUser(t *testing.T) {
	handler, svc, _ := setupNotificationHandler()

	svc.Push("user-A", models.NotifLoginFailed, "untuk A")
	svc.Push("user-B", models.NotifPasswordExpired, "untuk B")

	req := reqWithUserID(http.MethodGet, "/api/notifications", "user-A")
	rr := httptest.NewRecorder()
	handler.GetAll(rr, req)

	var resp dto.NotificationListResponse
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&resp))
	assert.Len(t, resp.Notifications, 1)
	assert.Equal(t, "untuk A", resp.Notifications[0].Message)
}

// ============================================================
// TEST: MarkRead
// ============================================================

func TestNotificationHandler_MarkRead_Unauthorized(t *testing.T) {
	handler, _, _ := setupNotificationHandler()

	req := httptest.NewRequest(http.MethodPatch, "/api/notifications/n1/read", nil)
	rr := httptest.NewRecorder()
	handler.MarkRead(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestNotificationHandler_MarkRead_MissingNotifID(t *testing.T) {
	handler, _, _ := setupNotificationHandler()

	req := reqWithUserID(http.MethodPatch, "/api/notifications//read", "user-1")
	rr := httptest.NewRecorder()
	handler.MarkRead(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestNotificationHandler_MarkRead_Success(t *testing.T) {
	handler, _, rc := setupNotificationHandler()

	notifs := []models.Notification{
		{ID: "notif-abc", UserID: "user-1", Type: models.NotifLoginFailed, Read: false, CreatedAt: time.Now()},
	}
	seedNotifInHandler(rc, "user-1", notifs)

	req := reqWithUserID(http.MethodPatch, "/api/notifications/notif-abc/read", "user-1")
	rr := httptest.NewRecorder()
	handler.MarkRead(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp map[string]string
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&resp))
	assert.Contains(t, resp["message"], "dibaca")
}

func TestNotificationHandler_MarkRead_NotFound_Returns500(t *testing.T) {
	handler, _, _ := setupNotificationHandler()

	req := reqWithUserID(http.MethodPatch, "/api/notifications/tidak-ada/read", "user-1")
	rr := httptest.NewRecorder()
	handler.MarkRead(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

// ============================================================
// TEST: MarkAllRead
// ============================================================

func TestNotificationHandler_MarkAllRead_Unauthorized(t *testing.T) {
	handler, _, _ := setupNotificationHandler()

	req := httptest.NewRequest(http.MethodPatch, "/api/notifications/read-all", nil)
	rr := httptest.NewRecorder()
	handler.MarkAllRead(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestNotificationHandler_MarkAllRead_Success(t *testing.T) {
	handler, _, rc := setupNotificationHandler()

	notifs := []models.Notification{
		{ID: "n1", Read: false, CreatedAt: time.Now()},
		{ID: "n2", Read: false, CreatedAt: time.Now()},
	}
	seedNotifInHandler(rc, "user-1", notifs)

	req := reqWithUserID(http.MethodPatch, "/api/notifications/read-all", "user-1")
	rr := httptest.NewRecorder()
	handler.MarkAllRead(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestNotificationHandler_MarkAllRead_EmptyList_NoError(t *testing.T) {
	handler, _, _ := setupNotificationHandler()

	req := reqWithUserID(http.MethodPatch, "/api/notifications/read-all", "user-1")
	rr := httptest.NewRecorder()
	handler.MarkAllRead(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

// ============================================================
// TEST: Delete
// ============================================================

func TestNotificationHandler_Delete_Unauthorized(t *testing.T) {
	handler, _, _ := setupNotificationHandler()

	req := httptest.NewRequest(http.MethodDelete, "/api/notifications/n1", nil)
	rr := httptest.NewRecorder()
	handler.Delete(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestNotificationHandler_Delete_MissingID(t *testing.T) {
	handler, _, _ := setupNotificationHandler()

	req := reqWithUserID(http.MethodDelete, "/api/notifications/", "user-1")
	rr := httptest.NewRecorder()
	handler.Delete(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestNotificationHandler_Delete_Success(t *testing.T) {
	handler, _, rc := setupNotificationHandler()

	notifs := []models.Notification{
		{ID: "notif-del", CreatedAt: time.Now()},
	}
	seedNotifInHandler(rc, "user-1", notifs)

	req := reqWithUserID(http.MethodDelete, "/api/notifications/notif-del", "user-1")
	rr := httptest.NewRecorder()
	handler.Delete(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp map[string]string
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&resp))
	assert.Contains(t, resp["message"], "dihapus")
}

func TestNotificationHandler_Delete_NotFound_Returns500(t *testing.T) {
	handler, _, _ := setupNotificationHandler()

	req := reqWithUserID(http.MethodDelete, "/api/notifications/tidak-ada", "user-1")
	rr := httptest.NewRecorder()
	handler.Delete(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

// ============================================================
// TEST: DeleteAll
// ============================================================

func TestNotificationHandler_DeleteAll_Unauthorized(t *testing.T) {
	handler, _, _ := setupNotificationHandler()

	req := httptest.NewRequest(http.MethodDelete, "/api/notifications", nil)
	rr := httptest.NewRecorder()
	handler.DeleteAll(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestNotificationHandler_DeleteAll_Success(t *testing.T) {
	handler, _, rc := setupNotificationHandler()

	notifs := []models.Notification{
		{ID: "n1"}, {ID: "n2"},
	}
	seedNotifInHandler(rc, "user-1", notifs)

	req := reqWithUserID(http.MethodDelete, "/api/notifications", "user-1")
	rr := httptest.NewRecorder()
	handler.DeleteAll(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp map[string]string
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&resp))
	assert.Contains(t, resp["message"], "dihapus")
}

// ============================================================
// TEST: ServeHTTP routing
// ============================================================

func TestNotificationHandler_ServeHTTP_GetAll(t *testing.T) {
	handler, _, _ := setupNotificationHandler()

	req := reqWithUserID(http.MethodGet, "/api/notifications", "user-1")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestNotificationHandler_ServeHTTP_MarkReadAll(t *testing.T) {
	handler, _, _ := setupNotificationHandler()

	req := reqWithUserID(http.MethodPatch, "/api/notifications/read-all", "user-1")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestNotificationHandler_ServeHTTP_MarkRead(t *testing.T) {
	handler, _, rc := setupNotificationHandler()

	notifs := []models.Notification{
		{ID: "notif-xyz", Read: false, CreatedAt: time.Now()},
	}
	seedNotifInHandler(rc, "user-1", notifs)

	req := reqWithUserID(http.MethodPatch, "/api/notifications/notif-xyz/read", "user-1")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestNotificationHandler_ServeHTTP_DeleteAll(t *testing.T) {
	handler, _, _ := setupNotificationHandler()

	req := reqWithUserID(http.MethodDelete, "/api/notifications", "user-1")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestNotificationHandler_ServeHTTP_Delete(t *testing.T) {
	handler, _, rc := setupNotificationHandler()

	notifs := []models.Notification{
		{ID: "n-route", CreatedAt: time.Now()},
	}
	seedNotifInHandler(rc, "user-1", notifs)

	req := reqWithUserID(http.MethodDelete, "/api/notifications/n-route", "user-1")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestNotificationHandler_ServeHTTP_MethodNotAllowed(t *testing.T) {
	handler, _, _ := setupNotificationHandler()

	req := reqWithUserID(http.MethodPost, "/api/notifications", "user-1")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}

func TestNotificationHandler_ServeHTTP_PutNotAllowed(t *testing.T) {
	handler, _, _ := setupNotificationHandler()

	req := reqWithUserID(http.MethodPut, "/api/notifications/n1", "user-1")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}
