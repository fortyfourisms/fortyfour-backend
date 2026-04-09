package services

import (
	"fmt"
	"strings"

	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/repository"
	"fortyfour-backend/pkg/logger"
)

// STRExpiryService mengecek tanggal kadaluarsa dan registrasi ulang STR CSIRT
// dan mengirim notifikasi ke user terkait.
type STRExpiryService struct {
	csirtRepo repository.CsirtRepositoryInterface
	notifSvc  *NotificationService
}

func NewSTRExpiryService(
	csirtRepo repository.CsirtRepositoryInterface,
	notifSvc *NotificationService,
) *STRExpiryService {
	return &STRExpiryService{
		csirtRepo: csirtRepo,
		notifSvc:  notifSvc,
	}
}

// CheckAndNotify mengecek tanggal-tanggal STR untuk CSIRT milik perusahaan user
// dan push notifikasi jika mendekati atau sudah melewati kadaluarsa.
//
// Dipanggil saat login (dalam goroutine) untuk user yang punya perusahaan.
// Karena setiap perusahaan hanya memiliki satu akun user, notifikasi
// langsung dikirim ke userID yang login.
func (s *STRExpiryService) CheckAndNotify(userID, idPerusahaan string) {
	// Recover dari panic agar goroutine tidak crash silent
	defer func() {
		if r := recover(); r != nil {
			logger.Warnf("STR expiry check panic: %v", r)
		}
	}()

	logger.Infof("STR expiry check: userID=%s, perusahaan=%s", userID, idPerusahaan)

	// Ambil data CSIRT milik perusahaan
	csirt, err := s.csirtRepo.GetByPerusahaanModel(idPerusahaan)
	if err != nil {
		logger.Infof("STR expiry check: perusahaan %s tidak punya CSIRT (%v)", idPerusahaan, err)
		return
	}

	logger.Infof("STR expiry check: CSIRT ditemukan — nama=%s, kadaluarsa=%v, reg_ulang=%v",
		csirt.NamaCsirt,
		csirt.TanggalKadaluarsa,
		csirt.TanggalRegistrasiUlang,
	)

	// ── Cek tanggal_kadaluarsa ──────────────────────────────────────────
	s.checkTanggalKadaluarsa(userID, csirt)

	// ── Cek tanggal_registrasi_ulang ────────────────────────────────────
	s.checkTanggalRegistrasiUlang(userID, csirt)
}

// checkTanggalKadaluarsa mengecek apakah STR sudah atau akan kadaluarsa
func (s *STRExpiryService) checkTanggalKadaluarsa(userID string, csirt *models.Csirt) {
	if csirt.TanggalKadaluarsa == nil || *csirt.TanggalKadaluarsa == "" {
		logger.Info("STR expiry check: tanggal_kadaluarsa kosong, skip")
		return
	}

	logger.Infof("STR expiry check: tanggal_kadaluarsa=%s, expired=%v, expiringSoon=%v",
		*csirt.TanggalKadaluarsa, csirt.IsSTRExpired(), csirt.IsSTRExpiringSoon())

	if csirt.IsSTRExpired() {
		// STR sudah expired — push notif expired (cek duplikasi)
		hasNotif, _ := s.hasNotifByType(userID, models.NotifSTRExpired, "kadaluarsa")
		if !hasNotif {
			msg := fmt.Sprintf(
				"STR CSIRT \"%s\" telah melewati tanggal kadaluarsa (%s). Segera lakukan perpanjangan.",
				csirt.NamaCsirt, *csirt.TanggalKadaluarsa,
			)
			if err := s.notifSvc.Push(userID, models.NotifSTRExpired, msg); err != nil {
				logger.Error(err, "failed to push STR expired notification")
			} else {
				logger.Info("STR expiry check: notif STR expired berhasil di-push")
			}
		} else {
			logger.Info("STR expiry check: notif STR expired sudah ada, skip duplikasi")
		}
		return
	}

	if csirt.IsSTRExpiringSoon() {
		// STR akan expired dalam 30 hari — push notif expiry soon (cek duplikasi)
		hasNotif, _ := s.notifSvc.HasSTRExpirySoonNotif(userID)
		if !hasNotif {
			days := csirt.DaysUntilSTRExpiry()
			msg := fmt.Sprintf(
				"STR CSIRT \"%s\" akan kadaluarsa dalam %d hari (tanggal %s). Segera lakukan perpanjangan.",
				csirt.NamaCsirt, days, *csirt.TanggalKadaluarsa,
			)
			if err := s.notifSvc.Push(userID, models.NotifSTRExpirySoon, msg); err != nil {
				logger.Error(err, "failed to push STR expiry soon notification")
			} else {
				logger.Infof("STR expiry check: notif STR expiry soon berhasil di-push (sisa %d hari)", days)
			}
		} else {
			logger.Info("STR expiry check: notif STR expiry soon sudah ada, skip duplikasi")
		}
	}
}

// checkTanggalRegistrasiUlang mengecek apakah tanggal registrasi ulang sudah atau akan jatuh tempo
func (s *STRExpiryService) checkTanggalRegistrasiUlang(userID string, csirt *models.Csirt) {
	if csirt.TanggalRegistrasiUlang == nil || *csirt.TanggalRegistrasiUlang == "" {
		return
	}

	logger.Infof("STR expiry check: tanggal_registrasi_ulang=%s, passed=%v, soon=%v",
		*csirt.TanggalRegistrasiUlang, csirt.IsRegistrasiUlangPassed(), csirt.IsRegistrasiUlangSoon())

	if csirt.IsRegistrasiUlangPassed() {
		// Tanggal registrasi ulang sudah lewat — push notif expired
		hasNotif, _ := s.hasNotifByType(userID, models.NotifSTRExpired, "registrasi ulang")
		if !hasNotif {
			msg := fmt.Sprintf(
				"Tanggal registrasi ulang STR CSIRT \"%s\" telah melewati batas waktu (%s). Segera lakukan registrasi ulang.",
				csirt.NamaCsirt, *csirt.TanggalRegistrasiUlang,
			)
			if err := s.notifSvc.Push(userID, models.NotifSTRExpired, msg); err != nil {
				logger.Error(err, "failed to push registration renewal expired notification")
			} else {
				logger.Info("STR expiry check: notif registrasi ulang expired berhasil di-push")
			}
		}
		return
	}

	if csirt.IsRegistrasiUlangSoon() {
		// Tanggal registrasi ulang mendekati — push notif soon
		hasNotif, _ := s.hasNotifByType(userID, models.NotifSTRExpirySoon, "registrasi ulang")
		if !hasNotif {
			days := csirt.DaysUntilRegistrasiUlang()
			msg := fmt.Sprintf(
				"Tanggal registrasi ulang STR CSIRT \"%s\" tinggal %d hari lagi (tanggal %s). Segera lakukan registrasi ulang.",
				csirt.NamaCsirt, days, *csirt.TanggalRegistrasiUlang,
			)
			if err := s.notifSvc.Push(userID, models.NotifSTRExpirySoon, msg); err != nil {
				logger.Error(err, "failed to push registration renewal soon notification")
			} else {
				logger.Infof("STR expiry check: notif registrasi ulang soon berhasil di-push (sisa %d hari)", days)
			}
		}
	}
}

// hasNotifByType mengecek apakah sudah ada notifikasi unread dengan type dan keyword tertentu
// untuk menghindari duplikasi antara notif kadaluarsa vs registrasi ulang
func (s *STRExpiryService) hasNotifByType(userID string, notifType models.NotificationType, keyword string) (bool, error) {
	notifs, err := s.notifSvc.GetAll(userID)
	if err != nil {
		return false, err
	}

	for _, n := range notifs {
		if n.Type == notifType && !n.Read {
			// Cek keyword dalam message untuk membedakan kadaluarsa vs registrasi ulang
			if strings.Contains(n.Message, keyword) {
				return true, nil
			}
		}
	}
	return false, nil
}
