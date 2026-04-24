package seeder

import (
	"fortyfour-backend/internal/services"
	"fortyfour-backend/pkg/logger"
)

// Policy mendefinisikan satu baris policy casbin: (role, resource, action)
type Policy struct {
	Role     string
	Resource string
	Action   string
}

// defaultPolicies adalah daftar semua policy yang harus selalu ada di sistem.
// Policy ini akan selalu di-sync setiap kali server restart:
//   - Jika belum ada → ditambahkan
//   - Jika ada di database tapi tidak ada di sini → dihapus (kecuali admin & staff)
var defaultPolicies = []Policy{

	// ── ADMIN ────────────────────────────────────────────────────────────────
	// Admin punya akses penuh ke semua endpoint
	{"admin", "/api/*", "*"},

	// ── USER PIC / PIC ───────────────────────────────────────────────────────
	// PIC Perusahaan — akses ke layanan SE, CSIRT, IKAS, Maturity, dan LMS

	// SE (Sistem Elektronik)
	{"user_pic", "/api/se", "GET"},
	{"user_pic", "/api/se", "POST"},
	{"user_pic", "/api/se/:id", "GET"},
	{"user_pic", "/api/se/:id/request-edit", "POST"},
	{"user_pic", "/api/se/edit-requests", "GET"},

	// CSIRT
	{"user_pic", "/api/csirt", "GET"},
	{"user_pic", "/api/csirt", "POST"},
	{"user_pic", "/api/csirt/:id", "GET"},
	{"user_pic", "/api/csirt/:id", "PUT"},
	{"user_pic", "/api/csirt/:id", "DELETE"},
	{"user_pic", "/api/csirt/:id/pgp-download", "GET"},

	// SDM CSIRT
	{"user_pic", "/api/sdm_csirt", "GET"},
	{"user_pic", "/api/sdm_csirt", "POST"},
	{"user_pic", "/api/sdm_csirt/:id", "GET"},
	{"user_pic", "/api/sdm_csirt/:id", "PUT"},
	{"user_pic", "/api/sdm_csirt/:id", "DELETE"},

	// PIC Perusahaan
	{"user_pic", "/api/pic", "GET"},
	{"user_pic", "/api/pic", "POST"},
	{"user_pic", "/api/pic/:id", "GET"},
	{"user_pic", "/api/pic/:id", "PUT"},
	{"user_pic", "/api/pic/:id", "DELETE"},

	// Perusahaan (user_pic hanya bisa lihat dan update miliknya sendiri)
	{"user_pic", "/api/perusahaan", "GET"},
	{"user_pic", "/api/perusahaan/:id", "GET"},
	{"user_pic", "/api/perusahaan/:id", "PUT"},

	// ── MATURITY ─────────────────────────────────────────────────────────────

	// Ruang Lingkup (master — read only untuk user_pic)
	{"user_pic", "/api/maturity/ruang-lingkup", "GET"},
	{"user_pic", "/api/maturity/ruang-lingkup/:id", "GET"},

	// Domain (master — read only untuk user_pic)
	{"user_pic", "/api/maturity/domain", "GET"},
	{"user_pic", "/api/maturity/domain/:id", "GET"},

	// Kategori (master — read only untuk user_pic)
	{"user_pic", "/api/maturity/kategori", "GET"},
	{"user_pic", "/api/maturity/kategori/:id", "GET"},

	// Sub Kategori (master — read only untuk user_pic)
	{"user_pic", "/api/maturity/sub-kategori", "GET"},
	{"user_pic", "/api/maturity/sub-kategori/:id", "GET"},

	// IKAS
	{"user_pic", "/api/maturity/ikas", "GET"},
	{"user_pic", "/api/maturity/ikas", "POST"},
	{"user_pic", "/api/maturity/ikas/:id", "GET"},
	{"user_pic", "/api/maturity/ikas/:id", "PUT"},
	{"user_pic", "/api/maturity/ikas/:id", "DELETE"},

	// Domain Identifikasi (read only untuk user_pic)
	{"user_pic", "/api/maturity/identifikasi", "GET"},
	{"user_pic", "/api/maturity/identifikasi/:id", "GET"},

	// Domain Proteksi (read only untuk user_pic)
	{"user_pic", "/api/maturity/proteksi", "GET"},
	{"user_pic", "/api/maturity/proteksi/:id", "GET"},

	// Domain Deteksi (read only untuk user_pic)
	{"user_pic", "/api/maturity/deteksi", "GET"},
	{"user_pic", "/api/maturity/deteksi/:id", "GET"},

	// Domain Gulih (read only untuk user_pic)
	{"user_pic", "/api/maturity/gulih", "GET"},
	{"user_pic", "/api/maturity/gulih/:id", "GET"},

	// Maturity (Pertanyaan — read only untuk user_pic)
	{"user_pic", "/api/maturity/pertanyaan-identifikasi", "GET"},
	{"user_pic", "/api/maturity/pertanyaan-identifikasi/:id", "GET"},

	{"user_pic", "/api/maturity/pertanyaan-proteksi", "GET"},
	{"user_pic", "/api/maturity/pertanyaan-proteksi/:id", "GET"},

	{"user_pic", "/api/maturity/pertanyaan-deteksi", "GET"},
	{"user_pic", "/api/maturity/pertanyaan-deteksi/:id", "GET"},

	{"user_pic", "/api/maturity/pertanyaan-gulih", "GET"},
	{"user_pic", "/api/maturity/pertanyaan-gulih/:id", "GET"},

	// Maturity (Jawaban)
	{"user_pic", "/api/maturity/jawaban-identifikasi", "GET"},
	{"user_pic", "/api/maturity/jawaban-identifikasi", "POST"},
	{"user_pic", "/api/maturity/jawaban-identifikasi/:id", "GET"},
	{"user_pic", "/api/maturity/jawaban-identifikasi/:id", "PUT"},

	{"user_pic", "/api/maturity/jawaban-proteksi", "GET"},
	{"user_pic", "/api/maturity/jawaban-proteksi", "POST"},
	{"user_pic", "/api/maturity/jawaban-proteksi/:id", "GET"},
	{"user_pic", "/api/maturity/jawaban-proteksi/:id", "PUT"},

	{"user_pic", "/api/maturity/jawaban-deteksi", "GET"},
	{"user_pic", "/api/maturity/jawaban-deteksi", "POST"},
	{"user_pic", "/api/maturity/jawaban-deteksi/:id", "GET"},
	{"user_pic", "/api/maturity/jawaban-deteksi/:id", "PUT"},

	{"user_pic", "/api/maturity/jawaban-gulih", "GET"},
	{"user_pic", "/api/maturity/jawaban-gulih", "POST"},
	{"user_pic", "/api/maturity/jawaban-gulih/:id", "GET"},
	{"user_pic", "/api/maturity/jawaban-gulih/:id", "PUT"},

	// ── LMS (user_pic juga bisa akses LMS) ──────────────────────────────────

	// Kelas
	{"user_pic", "/api/kelas", "GET"},
	{"user_pic", "/api/kelas/:id", "GET"},
	{"user_pic", "/api/kelas/:id/kuis", "GET"},
	{"user_pic", "/api/kelas/:id/sertifikat", "GET"},
	{"user_pic", "/api/kelas/:id/sertifikat/generate", "POST"},

	// Materi
	{"user_pic", "/api/materi/:id/progress", "POST"},
	{"user_pic", "/api/materi/:id/file-pendukung", "GET"},
	{"user_pic", "/api/materi/:id/diskusi", "GET"},
	{"user_pic", "/api/materi/:id/diskusi", "POST"},
	{"user_pic", "/api/diskusi/:id", "PUT"},
	{"user_pic", "/api/diskusi/:id", "DELETE"},
	{"user_pic", "/api/materi/:id/catatan", "GET"},
	{"user_pic", "/api/materi/:id/catatan", "PUT"},

	// File pendukung
	{"user_pic", "/api/file-pendukung/:id/download", "GET"},

	// Kuis
	{"user_pic", "/api/kuis/:id_kuis/start", "POST"},
	{"user_pic", "/api/kuis/attempt/:id_attempt/submit", "POST"},
	{"user_pic", "/api/kuis/attempt/:id_attempt/result", "GET"},

	// Sertifikat
	{"user_pic", "/api/sertifikat/me", "GET"},
	{"user_pic", "/api/sertifikat/:id", "GET"},
	{"user_pic", "/api/sertifikat/:id/download", "GET"},

	// ── USER ─────────────────────────────────────────────────────────────────
	// User biasa — hanya bisa mengakses LMS

	// Kelas (user bisa lihat list & detail)
	{"user", "/api/kelas", "GET"},
	{"user", "/api/kelas/:id", "GET"},

	// Kelas → kuis list (user bisa lihat)
	{"user", "/api/kelas/:id/kuis", "GET"},

	// Kelas → sertifikat (user)
	{"user", "/api/kelas/:id/sertifikat", "GET"},
	{"user", "/api/kelas/:id/sertifikat/generate", "POST"},

	// Materi — progress update (user)
	{"user", "/api/materi/:id/progress", "POST"},

	// Materi — file pendukung (user bisa lihat)
	{"user", "/api/materi/:id/file-pendukung", "GET"},

	// Materi — diskusi (user bisa CRUD diskusi sendiri)
	{"user", "/api/materi/:id/diskusi", "GET"},
	{"user", "/api/materi/:id/diskusi", "POST"},
	{"user", "/api/diskusi/:id", "PUT"},
	{"user", "/api/diskusi/:id", "DELETE"},

	// Materi — catatan pribadi (user)
	{"user", "/api/materi/:id/catatan", "GET"},
	{"user", "/api/materi/:id/catatan", "PUT"},

	// File pendukung — download (user)
	{"user", "/api/file-pendukung/:id/download", "GET"},

	// Kuis — start, submit, result (user)
	{"user", "/api/kuis/:id_kuis/start", "POST"},
	{"user", "/api/kuis/attempt/:id_attempt/submit", "POST"},
	{"user", "/api/kuis/attempt/:id_attempt/result", "GET"},

	// Sertifikat (user)
	{"user", "/api/sertifikat/me", "GET"},
	{"user", "/api/sertifikat/:id", "GET"},
	{"user", "/api/sertifikat/:id/download", "GET"},

	// Notifications (user)
	{"user", "/api/notifications", "GET"},
	{"user", "/api/notifications/read-all", "PATCH"},
	{"user", "/api/notifications/:id/read", "PATCH"},

	// Perusahaan (user hanya bisa baca perusahaan miliknya sendiri, dibatasi lagi di handler)
	{"user", "/api/perusahaan", "GET"},
	{"user", "/api/perusahaan/:id", "GET"},
}

// staffInitialPolicies adalah policy awal untuk role staff.
// Policy ini HANYA di-seed sekali (saat staff belum punya policy apapun).
// Setelah itu, admin bisa menambah/menghapus policy staff secara dinamis via API
// tanpa khawatir seeder akan mengembalikan policy yang sudah dihapus.
var staffInitialPolicies = []Policy{
	// SE
	{"staff", "/api/se", "GET"},
	{"staff", "/api/se", "POST"},
	{"staff", "/api/se/:id", "GET"},
	{"staff", "/api/se/:id", "PUT"},
	{"staff", "/api/se/export-pdf", "GET"},
	{"staff", "/api/se/:id/export-pdf", "GET"},
	{"staff", "/api/se/:id/request-edit", "POST"},
	{"staff", "/api/se/edit-requests", "GET"},

	// CSIRT
	{"staff", "/api/csirt", "GET"},
	{"staff", "/api/csirt", "POST"},
	{"staff", "/api/csirt/:id", "GET"},
	{"staff", "/api/csirt/:id", "PUT"},
	{"staff", "/api/csirt/export-pdf", "GET"},
	{"staff", "/api/csirt/:id/export-pdf", "GET"},
	{"staff", "/api/csirt/:id/pgp-download", "GET"},

	// SDM CSIRT
	{"staff", "/api/sdm_csirt", "GET"},
	{"staff", "/api/sdm_csirt", "POST"},
	{"staff", "/api/sdm_csirt/:id", "GET"},
	{"staff", "/api/sdm_csirt/:id", "PUT"},

	// PIC Perusahaan
	{"staff", "/api/pic", "GET"},
	{"staff", "/api/pic", "POST"},
	{"staff", "/api/pic/:id", "GET"},
	{"staff", "/api/pic/:id", "PUT"},

	// Perusahaan
	{"staff", "/api/perusahaan", "GET"},
	{"staff", "/api/perusahaan", "POST"},
	{"staff", "/api/perusahaan/:id", "GET"},
	{"staff", "/api/perusahaan/:id", "PUT"},

	// Sektor & Sub Sektor
	{"staff", "/api/sektor", "GET"},
	{"staff", "/api/sektor/:id", "GET"},
	{"staff", "/api/sub_sektor", "GET"},
	{"staff", "/api/sub_sektor/:id", "GET"},

	// Maturity — master data
	{"staff", "/api/maturity/ruang-lingkup", "GET"},
	{"staff", "/api/maturity/ruang-lingkup/:id", "GET"},
	{"staff", "/api/maturity/domain", "GET"},
	{"staff", "/api/maturity/domain/:id", "GET"},
	{"staff", "/api/maturity/kategori", "GET"},
	{"staff", "/api/maturity/kategori/:id", "GET"},
	{"staff", "/api/maturity/sub-kategori", "GET"},
	{"staff", "/api/maturity/sub-kategori/:id", "GET"},

	// IKAS
	{"staff", "/api/maturity/ikas", "GET"},
	{"staff", "/api/maturity/ikas", "POST"},
	{"staff", "/api/maturity/ikas/:id", "GET"},
	{"staff", "/api/maturity/ikas/:id", "PUT"},

	// Domain sub-resources
	{"staff", "/api/maturity/identifikasi", "GET"},
	{"staff", "/api/maturity/identifikasi/:id", "GET"},
	{"staff", "/api/maturity/proteksi", "GET"},
	{"staff", "/api/maturity/proteksi/:id", "GET"},
	{"staff", "/api/maturity/deteksi", "GET"},
	{"staff", "/api/maturity/deteksi/:id", "GET"},
	{"staff", "/api/maturity/gulih", "GET"},
	{"staff", "/api/maturity/gulih/:id", "GET"},

	// Pertanyaan
	{"staff", "/api/maturity/pertanyaan-identifikasi", "GET"},
	{"staff", "/api/maturity/pertanyaan-identifikasi/:id", "GET"},
	{"staff", "/api/maturity/pertanyaan-proteksi", "GET"},
	{"staff", "/api/maturity/pertanyaan-proteksi/:id", "GET"},
	{"staff", "/api/maturity/pertanyaan-deteksi", "GET"},
	{"staff", "/api/maturity/pertanyaan-deteksi/:id", "GET"},
	{"staff", "/api/maturity/pertanyaan-gulih", "GET"},
	{"staff", "/api/maturity/pertanyaan-gulih/:id", "GET"},

	// Jawaban
	{"staff", "/api/maturity/jawaban-identifikasi", "GET"},
	{"staff", "/api/maturity/jawaban-identifikasi", "POST"},
	{"staff", "/api/maturity/jawaban-identifikasi/:id", "GET"},
	{"staff", "/api/maturity/jawaban-identifikasi/:id", "PUT"},
	{"staff", "/api/maturity/jawaban-proteksi", "GET"},
	{"staff", "/api/maturity/jawaban-proteksi", "POST"},
	{"staff", "/api/maturity/jawaban-proteksi/:id", "GET"},
	{"staff", "/api/maturity/jawaban-proteksi/:id", "PUT"},
	{"staff", "/api/maturity/jawaban-deteksi", "GET"},
	{"staff", "/api/maturity/jawaban-deteksi", "POST"},
	{"staff", "/api/maturity/jawaban-deteksi/:id", "GET"},
	{"staff", "/api/maturity/jawaban-deteksi/:id", "PUT"},
	{"staff", "/api/maturity/jawaban-gulih", "GET"},
	{"staff", "/api/maturity/jawaban-gulih", "POST"},
	{"staff", "/api/maturity/jawaban-gulih/:id", "GET"},
	{"staff", "/api/maturity/jawaban-gulih/:id", "PUT"},

	// LMS — staff bisa manage kelas, materi, kuis, soal (GET, POST, PUT)
	{"staff", "/api/kelas", "GET"},
	{"staff", "/api/kelas", "POST"},
	{"staff", "/api/kelas/:id", "GET"},
	{"staff", "/api/kelas/:id", "PUT"},
	{"staff", "/api/kelas/:id/kuis", "GET"},
	{"staff", "/api/kelas/:id/sertifikat", "GET"},
	{"staff", "/api/materi/:id/progress", "POST"},
	{"staff", "/api/materi/:id/file-pendukung", "GET"},
	{"staff", "/api/materi/:id/file-pendukung", "POST"},
	{"staff", "/api/materi/:id/diskusi", "GET"},
	{"staff", "/api/materi/:id/diskusi", "POST"},
	{"staff", "/api/materi/:id/catatan", "GET"},
	{"staff", "/api/materi/:id/catatan", "PUT"},
	{"staff", "/api/file-pendukung/:id/download", "GET"},
	{"staff", "/api/kuis/:id_kuis/start", "POST"},
	{"staff", "/api/kuis/attempt/:id_attempt/submit", "POST"},
	{"staff", "/api/kuis/attempt/:id_attempt/result", "GET"},
	{"staff", "/api/sertifikat/me", "GET"},
	{"staff", "/api/sertifikat/:id", "GET"},
	{"staff", "/api/sertifikat/:id/download", "GET"},

	// Dashboard
	{"staff", "/api/dashboard/summary", "GET"},

	// Notifications
	{"staff", "/api/notifications", "GET"},
	{"staff", "/api/notifications/:id", "PUT"},

	// Casbin — staff bisa lihat policies (tapi tidak bisa manage)
	{"staff", "/api/casbin/policies", "GET"},
	{"staff", "/api/casbin/permissions", "GET"},
}

// SeedCasbinPolicies memastikan semua default policy ada di database
// dan menghapus policy lama yang sudah tidak ada di defaultPolicies.
// Untuk role staff, policy hanya di-seed sekali (initial setup).
// Aman dijalankan berulang kali.
func SeedCasbinPolicies(casbinService *services.CasbinService) {
	added := 0
	skipped := 0

	// 1. Tambah default policy (admin, user_pic, user) yang belum ada
	for _, p := range defaultPolicies {
		ok, err := casbinService.AddPolicy(p.Role, p.Resource, p.Action)
		if err != nil {
			logger.Errorf(err, "Casbin seeder: gagal tambah policy (%s, %s, %s)",
				p.Role, p.Resource, p.Action)
			continue
		}
		if ok {
			added++
		} else {
			skipped++ // sudah ada, skip
		}
	}

	// 2. Seed initial staff policies (hanya jika staff belum punya policy sama sekali)
	staffAdded := seedStaffInitialPolicies(casbinService)

	// 3. Hapus policy lama yang sudah tidak ada di defaultPolicies
	// (skip admin & staff agar policy dinamis staff tidak terhapus)
	removed := cleanupStalePolicies(casbinService)

	logger.Infof("Casbin seeder selesai: %d ditambahkan, %d sudah ada, %d staff initial, %d dihapus",
		added, skipped, staffAdded, removed)
}

// seedStaffInitialPolicies menyemai policy awal untuk role staff.
// Policy hanya ditambahkan jika staff belum memiliki policy APAPUN di database.
// Ini memastikan:
//   - Pertama kali: staff mendapat set policy default
//   - Admin menghapus policy → tetap terhapus (tidak di-seed ulang)
//   - Admin menambah policy → tetap ada (tidak dihapus oleh cleanup)
func seedStaffInitialPolicies(casbinService *services.CasbinService) int {
	// Cek apakah staff sudah punya policy di database
	existingPolicies := casbinService.GetRolePermissions("staff")
	if len(existingPolicies) > 0 {
		logger.Infof("Casbin seeder: staff sudah punya %d policy, skip initial seed", len(existingPolicies))
		return 0
	}

	// Staff belum punya policy → seed initial
	added := 0
	for _, p := range staffInitialPolicies {
		ok, err := casbinService.AddPolicy(p.Role, p.Resource, p.Action)
		if err != nil {
			logger.Errorf(err, "Casbin seeder: gagal tambah staff initial policy (%s, %s, %s)",
				p.Role, p.Resource, p.Action)
			continue
		}
		if ok {
			added++
		}
	}

	logger.Infof("Casbin seeder: staff initial seed — %d policy ditambahkan", added)
	return added
}

// cleanupStalePolicies menghapus policy di database yang sudah tidak ada di defaultPolicies.
// Ini penting agar policy lama (misal: user PUT/DELETE SE) tidak tersisa.
// Role admin dan staff di-skip:
//   - admin: wildcard policy jangan pernah dihapus
//   - staff: dikelola dinamis oleh admin via API, jangan di-cleanup
func cleanupStalePolicies(casbinService *services.CasbinService) int {
	// Bangun lookup set dari defaultPolicies
	policySet := make(map[string]bool)
	for _, p := range defaultPolicies {
		key := p.Role + "|" + p.Resource + "|" + p.Action
		policySet[key] = true
	}

	// Ambil semua policy yang ada di database
	allPolicies := casbinService.GetAllPolicies()

	removed := 0
	for _, p := range allPolicies {
		// Skip admin & staff — admin wildcard jangan dihapus,
		// staff dikelola dinamis oleh admin via API
		if p.Role == "admin" || p.Role == "staff" {
			continue
		}
		key := p.Role + "|" + p.Resource + "|" + p.Action
		if !policySet[key] {
			ok, err := casbinService.RemovePolicy(p.Role, p.Resource, p.Action)
			if err != nil {
				logger.Errorf(err, "Casbin seeder: gagal hapus stale policy (%s, %s, %s)",
					p.Role, p.Resource, p.Action)
				continue
			}
			if ok {
				logger.Infof("Casbin seeder: hapus stale policy (%s, %s, %s)", p.Role, p.Resource, p.Action)
				removed++
			}
		}
	}

	return removed
}
