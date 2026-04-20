package services
 
 import (
 	"time"
 
 	"fortyfour-backend/internal/models"
 	"fortyfour-backend/internal/repository"
 )
 
 // NotificationService mengelola notifikasi user di MySQL (history)
 type NotificationService struct {
 	repo repository.NotificationRepositoryInterface
 }
 
 func NewNotificationService(repo repository.NotificationRepositoryInterface) *NotificationService {
 	return &NotificationService{repo: repo}
 }
 
 // GetAll mengambil semua notifikasi milik user dari database
 func (s *NotificationService) GetAll(userID string) ([]models.Notification, error) {
 	return s.repo.FindAllByUserID(userID)
 }
 
 // Push menambahkan notifikasi baru ke database
 func (s *NotificationService) Push(userID string, notifType models.NotificationType, message string) error {
 	newNotif := &models.Notification{
 		UserID:    userID,
 		Type:      notifType,
 		Message:   message,
 		Read:      false,
 		CreatedAt: time.Now(),
 	}
 
 	return s.repo.Create(newNotif)
 }
 
 // MarkRead menandai satu notifikasi sebagai sudah dibaca
 func (s *NotificationService) MarkRead(userID string, notifID int64) error {
 	return s.repo.MarkRead(userID, notifID)
 }
 
 // MarkAllRead menandai semua notifikasi user sebagai sudah dibaca
 func (s *NotificationService) MarkAllRead(userID string) error {
 	return s.repo.MarkAllRead(userID)
 }
 
 // Delete menghapus satu notifikasi berdasarkan ID
 func (s *NotificationService) Delete(userID string, notifID int64) error {
 	return s.repo.Delete(userID, notifID)
 }
 
 // DeleteAll menghapus semua notifikasi milik user
 func (s *NotificationService) DeleteAll(userID string) error {
 	return s.repo.DeleteAllByUserID(userID)
 }
 
 // HasPasswordExpirySoonNotif mengecek apakah sudah ada notifikasi password expiry soon
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
