package services

import (
	"encoding/json"
	"fmt"
	"time"

	"fortyfour-backend/internal/models"
	"fortyfour-backend/pkg/cache"

	"github.com/google/uuid"
)

const (
	// Format key Redis untuk notifikasi per user
	// notif:{user_id} → JSON array of Notification
	notifKeyPrefix = "notif:"
)

// NotificationService mengelola notifikasi user di Redis
type NotificationService struct {
	rc cache.RedisInterface
}

func NewNotificationService(rc cache.RedisInterface) *NotificationService {
	return &NotificationService{rc: rc}
}

// redisKey mengembalikan key Redis untuk notifikasi user
func (s *NotificationService) redisKey(userID string) string {
	return fmt.Sprintf("%s%s", notifKeyPrefix, userID)
}

// GetAll mengambil semua notifikasi milik user
func (s *NotificationService) GetAll(userID string) ([]models.Notification, error) {
	key := s.redisKey(userID)
	raw, err := s.rc.Get(key)
	if err != nil {
		// Jika key tidak ada, kembalikan slice kosong (bukan error)
		return []models.Notification{}, nil
	}

	var notifs []models.Notification
	if err := json.Unmarshal([]byte(raw), &notifs); err != nil {
		return nil, fmt.Errorf("gagal membaca notifikasi: %w", err)
	}
	return notifs, nil
}

// Push menambahkan notifikasi baru ke daftar notifikasi user.
// Notifikasi disimpan tanpa TTL (persistent sampai di-dismiss user).
func (s *NotificationService) Push(userID string, notifType models.NotificationType, message string) error {
	notifs, err := s.GetAll(userID)
	if err != nil {
		return err
	}

	newNotif := models.Notification{
		ID:        uuid.New().String(),
		UserID:    userID,
		Type:      notifType,
		Message:   message,
		Read:      false,
		CreatedAt: time.Now(),
	}

	notifs = append([]models.Notification{newNotif}, notifs...) // prepend agar terbaru di atas

	return s.save(userID, notifs)
}

// MarkRead menandai satu notifikasi sebagai sudah dibaca
func (s *NotificationService) MarkRead(userID, notifID string) error {
	notifs, err := s.GetAll(userID)
	if err != nil {
		return err
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
		return fmt.Errorf("notifikasi tidak ditemukan")
	}

	return s.save(userID, notifs)
}

// MarkAllRead menandai semua notifikasi user sebagai sudah dibaca
func (s *NotificationService) MarkAllRead(userID string) error {
	notifs, err := s.GetAll(userID)
	if err != nil {
		return err
	}

	for i := range notifs {
		notifs[i].Read = true
	}

	return s.save(userID, notifs)
}

// Delete menghapus satu notifikasi berdasarkan ID
func (s *NotificationService) Delete(userID, notifID string) error {
	notifs, err := s.GetAll(userID)
	if err != nil {
		return err
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
		return fmt.Errorf("notifikasi tidak ditemukan")
	}

	return s.save(userID, filtered)
}

// DeleteAll menghapus semua notifikasi milik user
func (s *NotificationService) DeleteAll(userID string) error {
	return s.rc.Delete(s.redisKey(userID))
}

// HasPasswordExpirySoonNotif mengecek apakah sudah ada notifikasi password expiry soon
// untuk menghindari duplikasi notifikasi di setiap login
func (s *NotificationService) HasPasswordExpirySoonNotif(userID string) (bool, error) {
	notifs, err := s.GetAll(userID)
	if err != nil {
		return false, err
	}

	for _, n := range notifs {
		if n.Type == models.NotifPasswordExpirySoon && !n.Read {
			return true, nil
		}
	}
	return false, nil
}

// HasSTRExpirySoonNotif mengecek apakah sudah ada notifikasi STR expiry soon
// yang belum dibaca, untuk menghindari duplikasi notifikasi di setiap login
func (s *NotificationService) HasSTRExpirySoonNotif(userID string) (bool, error) {
	notifs, err := s.GetAll(userID)
	if err != nil {
		return false, err
	}

	for _, n := range notifs {
		if n.Type == models.NotifSTRExpirySoon && !n.Read {
			return true, nil
		}
	}
	return false, nil
}

// save menyimpan kembali slice notifikasi ke Redis tanpa TTL
func (s *NotificationService) save(userID string, notifs []models.Notification) error {
	data, err := json.Marshal(notifs)
	if err != nil {
		return fmt.Errorf("gagal menyimpan notifikasi: %w", err)
	}
	// TTL 0 → tidak expired (persistent)
	return s.rc.Set(s.redisKey(userID), string(data), 0)
}
