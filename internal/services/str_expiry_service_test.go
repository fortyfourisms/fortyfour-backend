package services
 
 import (
 	
 	"errors"
 	"fortyfour-backend/internal/dto"
 	"fortyfour-backend/internal/models"
 	"fortyfour-backend/internal/testhelpers"
 	"testing"
 	"time"
 
 	"github.com/stretchr/testify/assert"
 	"github.com/stretchr/testify/require"
 )
 
 // ============================================================
 // Mock CSIRT Repo khusus STRExpiry (butuh GetByPerusahaanModel)
 // ============================================================
 
 type mockCsirtRepoSTR struct {
 	GetByPerusahaanModelFn func(idPerusahaan string) (*models.Csirt, error)
 }
 
 func (m *mockCsirtRepoSTR) Create(req dto.CreateCsirtRequest, id string) error {
 	return nil
 }
 func (m *mockCsirtRepoSTR) ExistsByPerusahaan(idPerusahaan string) (bool, error) {
 	return false, nil
 }
 func (m *mockCsirtRepoSTR) GetByID(id string) (*models.Csirt, error) {
 	return nil, errors.New("not found")
 }
 func (m *mockCsirtRepoSTR) GetAllWithPerusahaan() ([]dto.CsirtResponse, error) {
 	return nil, nil
 }
 func (m *mockCsirtRepoSTR) GetByIDWithPerusahaan(id string) (*dto.CsirtResponse, error) {
 	return nil, errors.New("not found")
 }
 func (m *mockCsirtRepoSTR) Update(id string, csirt models.Csirt) error { return nil }
 func (m *mockCsirtRepoSTR) Delete(id string) error                     { return nil }
 func (m *mockCsirtRepoSTR) GetByPerusahaan(idPerusahaan string) ([]dto.CsirtResponse, error) {
 	return nil, nil
 }
 func (m *mockCsirtRepoSTR) GetByPerusahaanModel(idPerusahaan string) (*models.Csirt, error) {
 	return m.GetByPerusahaanModelFn(idPerusahaan)
 }
 
 // ============================================================
 // Helper: membuat tanggal string
 // ============================================================
 
 func dateStr(t time.Time) string {
 	return t.Format("2006-01-02")
 }
 
 func ptrStr(s string) *string {
 	return &s
 }
 
 // ============================================================
 // TEST: CheckAndNotify
 // ============================================================
 
 func TestSTRExpiryService_CheckAndNotify_NoCsirt_NoNotification(t *testing.T) {
 	csirtRepo := &mockCsirtRepoSTR{
 		GetByPerusahaanModelFn: func(idPerusahaan string) (*models.Csirt, error) {
 			return nil, errors.New("not found")
 		},
 	}
 	repo := testhelpers.NewMockNotificationRepository()
 	notifSvc := NewNotificationService(repo)
 	svc := NewSTRExpiryService(csirtRepo, notifSvc)
 
 	// Tidak boleh panic
 	svc.CheckAndNotify("user-1", "perusahaan-1")
 
 	// Tidak ada notifikasi yang di-push
 	notifs, err := notifSvc.GetAll("user-1")
 	require.NoError(t, err)
 	assert.Empty(t, notifs)
 }
 
 func TestSTRExpiryService_CheckAndNotify_TanggalKadaluarsaNil_Skip(t *testing.T) {
 	csirtRepo := &mockCsirtRepoSTR{
 		GetByPerusahaanModelFn: func(idPerusahaan string) (*models.Csirt, error) {
 			return &models.Csirt{
 				ID:                "csirt-1",
 				NamaCsirt:         "CSIRT Test",
 				TanggalKadaluarsa: nil, // nil
 			}, nil
 		},
 	}
 	repo := testhelpers.NewMockNotificationRepository()
 	notifSvc := NewNotificationService(repo)
 	svc := NewSTRExpiryService(csirtRepo, notifSvc)
 
 	svc.CheckAndNotify("user-1", "perusahaan-1")
 
 	notifs, _ := notifSvc.GetAll("user-1")
 	assert.Empty(t, notifs, "tanggal nil → tidak push notif")
 }
 
 func TestSTRExpiryService_CheckAndNotify_TanggalKadaluarsaKosong_Skip(t *testing.T) {
 	csirtRepo := &mockCsirtRepoSTR{
 		GetByPerusahaanModelFn: func(idPerusahaan string) (*models.Csirt, error) {
 			return &models.Csirt{
 				ID:                "csirt-1",
 				NamaCsirt:         "CSIRT Test",
 				TanggalKadaluarsa: ptrStr(""), // string kosong
 			}, nil
 		},
 	}
 	repo := testhelpers.NewMockNotificationRepository()
 	notifSvc := NewNotificationService(repo)
 	svc := NewSTRExpiryService(csirtRepo, notifSvc)
 
 	svc.CheckAndNotify("user-1", "perusahaan-1")
 
 	notifs, _ := notifSvc.GetAll("user-1")
 	assert.Empty(t, notifs, "tanggal kosong → tidak push notif")
 }
 
 func TestSTRExpiryService_CheckAndNotify_STRExpired_PushNotif(t *testing.T) {
 	// Tanggal kadaluarsa kemarin → sudah expired
 	yesterday := dateStr(time.Now().AddDate(0, 0, -1))
 	csirtRepo := &mockCsirtRepoSTR{
 		GetByPerusahaanModelFn: func(idPerusahaan string) (*models.Csirt, error) {
 			return &models.Csirt{
 				ID:                "csirt-1",
 				NamaCsirt:         "CSIRT ABC",
 				TanggalKadaluarsa: &yesterday,
 			}, nil
 		},
 	}
 	repo := testhelpers.NewMockNotificationRepository()
 	notifSvc := NewNotificationService(repo)
 	svc := NewSTRExpiryService(csirtRepo, notifSvc)
 
 	svc.CheckAndNotify("user-1", "perusahaan-1")
 
 	notifs, _ := notifSvc.GetAll("user-1")
 	require.Len(t, notifs, 1)
 	assert.Equal(t, models.NotifSTRExpired, notifs[0].Type)
 	assert.Contains(t, notifs[0].Message, "CSIRT ABC")
 	assert.Contains(t, notifs[0].Message, "melewati tanggal kadaluarsa")
 }
 
 func TestSTRExpiryService_CheckAndNotify_STRExpiringSoon_PushNotif(t *testing.T) {
 	// Tanggal kadaluarsa 10 hari lagi → expiring soon
 	future := dateStr(time.Now().AddDate(0, 0, 10))
 	csirtRepo := &mockCsirtRepoSTR{
 		GetByPerusahaanModelFn: func(idPerusahaan string) (*models.Csirt, error) {
 			return &models.Csirt{
 				ID:                "csirt-1",
 				NamaCsirt:         "CSIRT DEF",
 				TanggalKadaluarsa: &future,
 			}, nil
 		},
 	}
 	repo := testhelpers.NewMockNotificationRepository()
 	notifSvc := NewNotificationService(repo)
 	svc := NewSTRExpiryService(csirtRepo, notifSvc)
 
 	svc.CheckAndNotify("user-1", "perusahaan-1")
 
 	notifs, _ := notifSvc.GetAll("user-1")
 	require.Len(t, notifs, 1)
 	assert.Equal(t, models.NotifSTRExpirySoon, notifs[0].Type)
 	assert.Contains(t, notifs[0].Message, "CSIRT DEF")
 	assert.Contains(t, notifs[0].Message, "akan kadaluarsa")
 }
 
 func TestSTRExpiryService_CheckAndNotify_STRMasihJauh_NoNotif(t *testing.T) {
 	// Tanggal kadaluarsa 1 tahun lagi → masih sangat jauh, bukan expiring soon
 	farFuture := dateStr(time.Now().AddDate(1, 0, 0))
 	csirtRepo := &mockCsirtRepoSTR{
 		GetByPerusahaanModelFn: func(idPerusahaan string) (*models.Csirt, error) {
 			return &models.Csirt{
 				ID:                "csirt-1",
 				NamaCsirt:         "CSIRT GHI",
 				TanggalKadaluarsa: &farFuture,
 			}, nil
 		},
 	}
 	repo := testhelpers.NewMockNotificationRepository()
 	notifSvc := NewNotificationService(repo)
 	svc := NewSTRExpiryService(csirtRepo, notifSvc)
 
 	svc.CheckAndNotify("user-1", "perusahaan-1")
 
 	notifs, _ := notifSvc.GetAll("user-1")
 	assert.Empty(t, notifs, "kadaluarsa masih jauh → tidak push notif")
 }
 
 func TestSTRExpiryService_CheckAndNotify_DuplikasiExpired_TidakPushUlang(t *testing.T) {
 	yesterday := dateStr(time.Now().AddDate(0, 0, -1))
 	csirtRepo := &mockCsirtRepoSTR{
 		GetByPerusahaanModelFn: func(idPerusahaan string) (*models.Csirt, error) {
 			return &models.Csirt{
 				ID:                "csirt-1",
 				NamaCsirt:         "CSIRT Dup",
 				TanggalKadaluarsa: &yesterday,
 			}, nil
 		},
 	}
 	repo := testhelpers.NewMockNotificationRepository()
 	notifSvc := NewNotificationService(repo)
 
 	// Seed notifikasi expired yang sudah ada (unread)
 	existing := []models.Notification{
 		{
 			ID:      "existing-1",
 			UserID:  "user-1",
 			Type:    models.NotifSTRExpired,
 			Message: "STR CSIRT \"CSIRT Dup\" telah melewati tanggal kadaluarsa (" + yesterday + "). Segera lakukan perpanjangan.",
 			Read:    false,
 		},
 	}
 	repo.Notifications["user-1"] = existing
 
 	svc := NewSTRExpiryService(csirtRepo, notifSvc)
 	svc.CheckAndNotify("user-1", "perusahaan-1")
 
 	notifs, _ := notifSvc.GetAll("user-1")
 	assert.Len(t, notifs, 1, "notif duplikat tidak boleh di-push ulang")
 	assert.Equal(t, "existing-1", notifs[0].ID, "notif lama tetap ada")
 }
 
 func TestSTRExpiryService_CheckAndNotify_DuplikasiExpiringSoon_TidakPushUlang(t *testing.T) {
 	future := dateStr(time.Now().AddDate(0, 0, 10))
 	csirtRepo := &mockCsirtRepoSTR{
 		GetByPerusahaanModelFn: func(idPerusahaan string) (*models.Csirt, error) {
 			return &models.Csirt{
 				ID:                "csirt-1",
 				NamaCsirt:         "CSIRT Dup Soon",
 				TanggalKadaluarsa: &future,
 			}, nil
 		},
 	}
 	repo := testhelpers.NewMockNotificationRepository()
 	notifSvc := NewNotificationService(repo)
 
 	// Seed notifikasi expirySoon yang sudah ada (unread)
 	existing := []models.Notification{
 		{
 			ID:      "existing-soon-1",
 			UserID:  "user-1",
 			Type:    models.NotifSTRExpirySoon,
 			Message: "STR hampir kadaluarsa",
 			Read:    false,
 		},
 	}
 	repo.Notifications["user-1"] = existing
 
 	svc := NewSTRExpiryService(csirtRepo, notifSvc)
 	svc.CheckAndNotify("user-1", "perusahaan-1")
 
 	notifs, _ := notifSvc.GetAll("user-1")
 	assert.Len(t, notifs, 1, "notif expiring soon duplikat tidak boleh di-push ulang")
 }
 
 func TestSTRExpiryService_CheckAndNotify_ExpiredNotifSudahDibaca_PushUlang(t *testing.T) {
 	yesterday := dateStr(time.Now().AddDate(0, 0, -1))
 	csirtRepo := &mockCsirtRepoSTR{
 		GetByPerusahaanModelFn: func(idPerusahaan string) (*models.Csirt, error) {
 			return &models.Csirt{
 				ID:                "csirt-1",
 				NamaCsirt:         "CSIRT Read",
 				TanggalKadaluarsa: &yesterday,
 			}, nil
 		},
 	}
 	repo := testhelpers.NewMockNotificationRepository()
 	notifSvc := NewNotificationService(repo)
 
 	// Seed notifikasi expired yang sudah dibaca (Read: true)
 	existing := []models.Notification{
 		{
 			ID:      "read-1",
 			UserID:  "user-1",
 			Type:    models.NotifSTRExpired,
 			Message: "pesan lama kadaluarsa",
 			Read:    true, // sudah dibaca
 		},
 	}
 	repo.Notifications["user-1"] = existing
 
 	svc := NewSTRExpiryService(csirtRepo, notifSvc)
 	svc.CheckAndNotify("user-1", "perusahaan-1")
 
 	notifs, _ := notifSvc.GetAll("user-1")
 	// Notif lama (read) masih ada + notif baru di-push
 	assert.Len(t, notifs, 2, "jika notif lama sudah dibaca, push notif baru")
 }
 
 func TestSTRExpiryService_CheckAndNotify_PanicRecovery(t *testing.T) {
 	// Test bahwa panic dalam goroutine tidak crash
 	csirtRepo := &mockCsirtRepoSTR{
 		GetByPerusahaanModelFn: func(idPerusahaan string) (*models.Csirt, error) {
 			panic("unexpected panic in repo")
 		},
 	}
 	repo := testhelpers.NewMockNotificationRepository()
 	notifSvc := NewNotificationService(repo)
 	svc := NewSTRExpiryService(csirtRepo, notifSvc)
 
 	// Harus tidak panic
 	assert.NotPanics(t, func() {
 		svc.CheckAndNotify("user-1", "perusahaan-1")
 	})
 }
 
 // ============================================================
 // TEST: hasNotifByType
 // ============================================================
 
 func TestSTRExpiryService_HasNotifByType_FoundUnread(t *testing.T) {
 	repo := testhelpers.NewMockNotificationRepository()
 	notifSvc := NewNotificationService(repo)
 	svc := NewSTRExpiryService(nil, notifSvc)
 
 	notifs := []models.Notification{
 		{ID: "n1", Type: models.NotifSTRExpired, Message: "STR kadaluarsa sudah lewat", Read: false},
 	}
 	repo.Notifications["user-1"] = notifs
 
 	has, err := svc.hasNotifByType("user-1", models.NotifSTRExpired, "kadaluarsa")
 	require.NoError(t, err)
 	assert.True(t, has)
 }
 
 func TestSTRExpiryService_HasNotifByType_AlreadyRead_ReturnsFalse(t *testing.T) {
 	repo := testhelpers.NewMockNotificationRepository()
 	notifSvc := NewNotificationService(repo)
 	svc := NewSTRExpiryService(nil, notifSvc)
 
 	notifs := []models.Notification{
 		{ID: "n1", Type: models.NotifSTRExpired, Message: "STR kadaluarsa sudah lewat", Read: true},
 	}
 	repo.Notifications["user-1"] = notifs
 
 	has, err := svc.hasNotifByType("user-1", models.NotifSTRExpired, "kadaluarsa")
 	require.NoError(t, err)
 	assert.False(t, has, "notif sudah dibaca → return false")
 }
 
 func TestSTRExpiryService_HasNotifByType_DifferentKeyword_ReturnsFalse(t *testing.T) {
 	repo := testhelpers.NewMockNotificationRepository()
 	notifSvc := NewNotificationService(repo)
 	svc := NewSTRExpiryService(nil, notifSvc)
 
 	notifs := []models.Notification{
 		{ID: "n1", Type: models.NotifSTRExpired, Message: "STR registrasi ulang sudah lewat", Read: false},
 	}
 	repo.Notifications["user-1"] = notifs
 
 	// Cari keyword "kadaluarsa" tapi message-nya "registrasi ulang"
 	has, err := svc.hasNotifByType("user-1", models.NotifSTRExpired, "kadaluarsa")
 	require.NoError(t, err)
 	assert.False(t, has, "keyword berbeda → return false")
 }
 
 func TestSTRExpiryService_HasNotifByType_NoNotif_ReturnsFalse(t *testing.T) {
 	repo := testhelpers.NewMockNotificationRepository()
 	notifSvc := NewNotificationService(repo)
 	svc := NewSTRExpiryService(nil, notifSvc)
 
 	has, err := svc.hasNotifByType("user-1", models.NotifSTRExpired, "kadaluarsa")
 	require.NoError(t, err)
 	assert.False(t, has)
 }
 
 func TestSTRExpiryService_HasNotifByType_DifferentType_ReturnsFalse(t *testing.T) {
 	repo := testhelpers.NewMockNotificationRepository()
 	notifSvc := NewNotificationService(repo)
 	svc := NewSTRExpiryService(nil, notifSvc)
 
 	notifs := []models.Notification{
 		{ID: "n1", Type: models.NotifLoginFailed, Message: "Login gagal kadaluarsa", Read: false},
 	}
 	repo.Notifications["user-1"] = notifs
 
 	has, err := svc.hasNotifByType("user-1", models.NotifSTRExpired, "kadaluarsa")
 	require.NoError(t, err)
 	assert.False(t, has, "type berbeda → return false meskipun keyword match")
 }
 
 func TestSTRExpiryService_HasNotifByType_GetAllError(t *testing.T) {
 	repo := testhelpers.NewMockNotificationRepository()
 	notifSvc := NewNotificationService(repo)
 	svc := NewSTRExpiryService(nil, notifSvc)
 
 	// FindAllByUserID dengan JSON rusak (simulasi manual)
 	repo.FindAllByUserIDFn = func(userID string) ([]models.Notification, error) {
 		return nil, errors.New("database error")
 	}
 
 	_, err := svc.hasNotifByType("user-1", models.NotifSTRExpired, "kadaluarsa")
 	assert.Error(t, err)
 }
 
 // ============================================================
 // TEST: CheckAndNotify — integrasi cek duplikasi via hasNotifByType
 // ============================================================
 
 func TestSTRExpiryService_CheckAndNotify_ExpiredNotifSudahAda_DenganKeyword_TidakPush(t *testing.T) {
 	yesterday := dateStr(time.Now().AddDate(0, 0, -1))
 	csirtRepo := &mockCsirtRepoSTR{
 		GetByPerusahaanModelFn: func(idPerusahaan string) (*models.Csirt, error) {
 			return &models.Csirt{
 				ID:                "csirt-1",
 				NamaCsirt:         "Test CSIRT",
 				TanggalKadaluarsa: &yesterday,
 			}, nil
 		},
 	}
 	repo := testhelpers.NewMockNotificationRepository()
 	notifSvc := NewNotificationService(repo)
 
 	// Seed notifikasi expired yang sudah ada (dengan keyword "kadaluarsa"), unread
 	existingNotifs := []models.Notification{
 		{
 			ID:      "existing-1",
 			UserID:  "user-1",
 			Type:    models.NotifSTRExpired,
 			Message: "STR kadaluarsa sudah lewat",
 			Read:    false,
 		},
 	}
 	repo.Notifications["user-1"] = existingNotifs
 
 	svc := NewSTRExpiryService(csirtRepo, notifSvc)
 	svc.CheckAndNotify("user-1", "perusahaan-1")
 
 	notifs, _ := notifSvc.GetAll("user-1")
 	assert.Len(t, notifs, 1, "tidak boleh push duplikat expired")
 }
