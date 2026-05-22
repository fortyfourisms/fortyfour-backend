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
	if notif.ID == 0 {
		notif.ID = int64(len(m.data[notif.UserID]) + 1)
	}
	m.data[notif.UserID] = append([]models.Notification{*notif}, m.data[notif.UserID]...)
	return nil
}

func (m *mockNotifRepo) FindAll() ([]models.Notification, error) {
	if m.findAllErr != nil {
		return nil, m.findAllErr
	}
	var all []models.Notification
	for _, notifs := range m.data {
		all = append(all, notifs...)
	}
	return all, nil
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

func (m *mockNotifRepo) MarkRead(userID string, notifID int64) error {
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

func (m *mockNotifRepo) MarkAllReadGlobal() error {
	if m.markAllErr != nil {
		return m.markAllErr
	}
	for userID, notifs := range m.data {
		for i := range notifs {
			notifs[i].Read = true
		}
		m.data[userID] = notifs
	}
	return nil
}

func (m *mockNotifRepo) Delete(userID string, notifID int64) error {
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

func (m *mockNotifRepo) DeleteAll() error {
	if m.deleteAllErr != nil {
		return m.deleteAllErr
	}
	m.data = make(map[string][]models.Notification)
	return nil
}

// ============================================================
// TEST: GetAll
// ============================================================

func TestNotificationService_GetAll_Empty_ReturnsEmptySlice(t *testing.T) {
	repo := newMockNotifRepo()
	svc := NewNotificationService(repo)

	result, err := svc.GetAll("user-1", "user")

	require.NoError(t, err)
	assert.Empty(t, result)
}

func TestNotificationService_GetAll_WithData(t *testing.T) {
	repo := newMockNotifRepo()
	notifs := []models.Notification{
		{ID: 1, UserID: "u1", Type: models.NotifLoginFailed, Message: "Login gagal", Read: false},
		{ID: 2, UserID: "u1", Type: models.NotifPasswordExpirySoon, Message: "Password segera expired", Read: true},
	}
	repo.data["u1"] = notifs

	svc := NewNotificationService(repo)
	result, err := svc.GetAll("u1", "user")

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(1), result[0].ID)
}

func TestNotificationService_GetAll_AdminReturnsAll(t *testing.T) {
	repo := newMockNotifRepo()
	repo.data["u1"] = []models.Notification{{ID: 1, UserID: "u1"}}
	repo.data["u2"] = []models.Notification{{ID: 2, UserID: "u2"}}

	svc := NewNotificationService(repo)
	result, err := svc.GetAll("admin-id", "admin")

	require.NoError(t, err)
	assert.Len(t, result, 2, "admin harus melihat semua notifikasi")
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
	assert.NotEqual(t, int64(0), repo.data["u1"][0].ID)
}

// ============================================================
// TEST: MarkRead
// ============================================================

func TestNotificationService_MarkRead_Success(t *testing.T) {
	repo := newMockNotifRepo()
	repo.data["u1"] = []models.Notification{
		{ID: 10, UserID: "u1", Read: false},
	}

	svc := NewNotificationService(repo)
	err := svc.MarkRead("u1", 10)

	require.NoError(t, err)
	assert.True(t, repo.data["u1"][0].Read)
}

// ============================================================
// TEST: Delete
// ============================================================

func TestNotificationService_Delete_Success(t *testing.T) {
	repo := newMockNotifRepo()
	repo.data["u1"] = []models.Notification{
		{ID: 1}, {ID: 2},
	}

	svc := NewNotificationService(repo)
	err := svc.Delete("u1", 1)

	require.NoError(t, err)
	assert.Len(t, repo.data["u1"], 1)
	assert.Equal(t, int64(2), repo.data["u1"][0].ID)
}

// ============================================================
// TEST: DeleteAll (user)
// ============================================================

func TestNotificationService_DeleteAll_User_Success(t *testing.T) {
	repo := newMockNotifRepo()
	repo.data["u1"] = []models.Notification{{ID: 100}}
	repo.data["u2"] = []models.Notification{{ID: 200}}

	svc := NewNotificationService(repo)
	err := svc.DeleteAll("u1", "user")

	require.NoError(t, err)
	assert.Empty(t, repo.data["u1"], "u1 harus terhapus")
	assert.Len(t, repo.data["u2"], 1, "u2 tidak boleh terpengaruh")
}

// ============================================================
// TEST: DeleteAll (admin)
// ============================================================

func TestNotificationService_DeleteAll_Admin_DeletesAll(t *testing.T) {
	repo := newMockNotifRepo()
	repo.data["u1"] = []models.Notification{{ID: 1}}
	repo.data["u2"] = []models.Notification{{ID: 2}}

	svc := NewNotificationService(repo)
	err := svc.DeleteAll("admin-id", "admin")

	require.NoError(t, err)
	assert.Empty(t, repo.data, "admin harus menghapus semua notifikasi semua user")
}

// ============================================================
// TEST: MarkAllRead (user)
// ============================================================

func TestNotificationService_MarkAllRead_User_OnlyOwn(t *testing.T) {
	repo := newMockNotifRepo()
	repo.data["u1"] = []models.Notification{{ID: 1, Read: false}, {ID: 2, Read: false}}
	repo.data["u2"] = []models.Notification{{ID: 3, Read: false}}

	svc := NewNotificationService(repo)
	err := svc.MarkAllRead("u1", "user")

	require.NoError(t, err)
	assert.True(t, repo.data["u1"][0].Read)
	assert.True(t, repo.data["u1"][1].Read)
	assert.False(t, repo.data["u2"][0].Read, "u2 tidak boleh terpengaruh")
}

// ============================================================
// TEST: MarkAllRead (admin)
// ============================================================

func TestNotificationService_MarkAllRead_Admin_AllUsers(t *testing.T) {
	repo := newMockNotifRepo()
	repo.data["u1"] = []models.Notification{{ID: 1, Read: false}}
	repo.data["u2"] = []models.Notification{{ID: 2, Read: false}}

	svc := NewNotificationService(repo)
	err := svc.MarkAllRead("admin-id", "admin")

	require.NoError(t, err)
	assert.True(t, repo.data["u1"][0].Read, "u1 harus ter-mark read oleh admin")
	assert.True(t, repo.data["u2"][0].Read, "u2 harus ter-mark read oleh admin")
}

// ============================================================
// TEST: Staff permissions (GetAll & MarkAllRead)
// ============================================================

func TestNotificationService_GetAll_StaffReturnsAll(t *testing.T) {
	repo := newMockNotifRepo()
	repo.data["u1"] = []models.Notification{{ID: 1, UserID: "u1"}}
	repo.data["u2"] = []models.Notification{{ID: 2, UserID: "u2"}}

	svc := NewNotificationService(repo)
	result, err := svc.GetAll("staff-id", "staff")

	require.NoError(t, err)
	assert.Len(t, result, 2, "staff harus melihat semua notifikasi")
}

func TestNotificationService_MarkAllRead_Staff_AllUsers(t *testing.T) {
	repo := newMockNotifRepo()
	repo.data["u1"] = []models.Notification{{ID: 1, Read: false}}
	repo.data["u2"] = []models.Notification{{ID: 2, Read: false}}

	svc := NewNotificationService(repo)
	err := svc.MarkAllRead("staff-id", "staff")

	require.NoError(t, err)
	assert.True(t, repo.data["u1"][0].Read, "u1 harus ter-mark read oleh staff")
	assert.True(t, repo.data["u2"][0].Read, "u2 harus ter-mark read oleh staff")
}
