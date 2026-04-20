package services
 
 import (
 	"errors"
 	"testing"
 
 	"fortyfour-backend/internal/models"
 
 	"github.com/stretchr/testify/assert"
 	"github.com/stretchr/testify/require"
 )
 
 // ============================================================
 // Mock Repository untuk NotificationService
 // ============================================================
 
 type mockNotifRepo struct {
 	data         map[string][]models.Notification
 	createErr    error
 	findAllErr   error
 	markReadErr  error
 	markAllErr   error
 	deleteErr    error
 	deleteAllErr error
 }
 
 func newMockNotifRepo() *mockNotifRepo {
 	return &mockNotifRepo{data: make(map[string][]models.Notification)}
 }
 
 func (m *mockNotifRepo) Create(notif *models.Notification) error {
 	if m.createErr != nil {
 		return m.createErr
 	}
 	m.data[notif.UserID] = append([]models.Notification{*notif}, m.data[notif.UserID]...)
 	return nil
 }
 
 func (m *mockNotifRepo) FindAllByUserID(userID string) ([]models.Notification, error) {
 	if m.findAllErr != nil {
 		return nil, m.findAllErr
 	}
 	notifs, ok := m.data[userID]
 	if !ok {
 		return []models.Notification{}, nil
 	}
 	return notifs, nil
 }
 
 func (m *mockNotifRepo) MarkRead(userID, notifID string) error {
 	if m.markReadErr != nil {
 		return m.markReadErr
 	}
 	notifs, ok := m.data[userID]
 	if !ok {
 		return errors.New("notifikasi tidak ditemukan")
 	}
 	found := false
 	for i := range notifs {
 		if notifs[i].ID == notifID {
 			notifs[i].Read = true
 			found = true
 			break
 		}
 	}
 	if !found {
 		return errors.New("notifikasi tidak ditemukan")
 	}
 	m.data[userID] = notifs
 	return nil
 }
 
 func (m *mockNotifRepo) MarkAllRead(userID string) error {
 	if m.markAllErr != nil {
 		return m.markAllErr
 	}
 	notifs := m.data[userID]
 	for i := range notifs {
 		notifs[i].Read = true
 	}
 	m.data[userID] = notifs
 	return nil
 }
 
 func (m *mockNotifRepo) Delete(userID, notifID string) error {
 	if m.deleteErr != nil {
 		return m.deleteErr
 	}
 	notifs, ok := m.data[userID]
 	if !ok {
 		return errors.New("notifikasi tidak ditemukan")
 	}
 	filtered := make([]models.Notification, 0, len(notifs))
 	found := false
 	for _, n := range notifs {
 		if n.ID == notifID {
 			found = true
 			continue
 		}
 		filtered = append(filtered, n)
 	}
 	if !found {
 		return errors.New("notifikasi tidak ditemukan")
 	}
 	m.data[userID] = filtered
 	return nil
 }
 
 func (m *mockNotifRepo) DeleteAllByUserID(userID string) error {
 	if m.deleteAllErr != nil {
 		return m.deleteAllErr
 	}
 	delete(m.data, userID)
 	return nil
 }
 
 // ============================================================
 // TEST: GetAll
 // ============================================================
 
 func TestNotificationService_GetAll_Empty_ReturnsEmptySlice(t *testing.T) {
 	repo := newMockNotifRepo()
 	svc := NewNotificationService(repo)
 
 	result, err := svc.GetAll("user-1")
 
 	require.NoError(t, err)
 	assert.Empty(t, result)
 }
 
 func TestNotificationService_GetAll_WithData(t *testing.T) {
 	repo := newMockNotifRepo()
 	notifs := []models.Notification{
 		{ID: "n1", UserID: "u1", Type: models.NotifLoginFailed, Message: "Login gagal", Read: false},
 		{ID: "n2", UserID: "u1", Type: models.NotifPasswordExpirySoon, Message: "Password segera expired", Read: true},
 	}
 	repo.data["u1"] = notifs
 
 	svc := NewNotificationService(repo)
 	result, err := svc.GetAll("u1")
 
 	require.NoError(t, err)
 	assert.Len(t, result, 2)
 	assert.Equal(t, "n1", result[0].ID)
 }
 
 // ============================================================
 // TEST: Push
 // ============================================================
 
 func TestNotificationService_Push_AddsNotification(t *testing.T) {
 	repo := newMockNotifRepo()
 	svc := NewNotificationService(repo)
 
 	err := svc.Push("u1", models.NotifLoginFailed, "Login gagal")
 
 	require.NoError(t, err)
 	assert.Len(t, repo.data["u1"], 1)
 	assert.Equal(t, models.NotifLoginFailed, repo.data["u1"][0].Type)
 }
 
 // ============================================================
 // TEST: MarkRead
 // ============================================================
 
 func TestNotificationService_MarkRead_Success(t *testing.T) {
 	repo := newMockNotifRepo()
 	repo.data["u1"] = []models.Notification{
 		{ID: "n1", UserID: "u1", Read: false},
 	}
 
 	svc := NewNotificationService(repo)
 	err := svc.MarkRead("u1", "n1")
 
 	require.NoError(t, err)
 	assert.True(t, repo.data["u1"][0].Read)
 }
 
 // ============================================================
 // TEST: Delete
 // ============================================================
 
 func TestNotificationService_Delete_Success(t *testing.T) {
 	repo := newMockNotifRepo()
 	repo.data["u1"] = []models.Notification{
 		{ID: "n1"}, {ID: "n2"},
 	}
 
 	svc := NewNotificationService(repo)
 	err := svc.Delete("u1", "n1")
 
 	require.NoError(t, err)
 	assert.Len(t, repo.data["u1"], 1)
 	assert.Equal(t, "n2", repo.data["u1"][0].ID)
 }
 
 // ============================================================
 // TEST: DeleteAll
 // ============================================================
 
 func TestNotificationService_DeleteAll_Success(t *testing.T) {
 	repo := newMockNotifRepo()
 	repo.data["u1"] = []models.Notification{{ID: "n1"}}
 
 	svc := NewNotificationService(repo)
 	err := svc.DeleteAll("u1")
 
 	require.NoError(t, err)
 	assert.Empty(t, repo.data["u1"])
 }
