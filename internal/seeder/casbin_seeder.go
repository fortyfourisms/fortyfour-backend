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

// defaultPolicies adalah daftar semua policy yang harus selalu ada di sistem
var defaultPolicies = []Policy{

	// ── ADMIN ────────────────────────────────────────────────────────────────
	// Admin punya akses penuh ke semua endpoint
	{"admin", "/api/*", "*"},

	// ── USER ─────────────────────────────────────────────────────────────────
	// SE (Sistem Elektronik)
	{"user", "/api/se", "GET"},
	{"user", "/api/se", "POST"},
	{"user", "/api/se/:id", "GET"},
	{"user", "/api/se/:id", "PUT"},
	{"user", "/api/se/:id", "DELETE"},

	// CSIRT
	{"user", "/api/csirt", "GET"},
	{"user", "/api/csirt", "POST"},
	{"user", "/api/csirt/:id", "GET"},
	{"user", "/api/csirt/:id", "PUT"},
	{"user", "/api/csirt/:id", "DELETE"},
	{"user", "/api/csirt/:id/pgp-download", "GET"},

	// SDM CSIRT
	{"user", "/api/sdm_csirt", "GET"},
	{"user", "/api/sdm_csirt", "POST"},
	{"user", "/api/sdm_csirt/:id", "GET"},
	{"user", "/api/sdm_csirt/:id", "PUT"},
	{"user", "/api/sdm_csirt/:id", "DELETE"},

	// PIC Perusahaan
	{"user", "/api/pic", "GET"},
	{"user", "/api/pic", "POST"},
	{"user", "/api/pic/:id", "GET"},
	{"user", "/api/pic/:id", "PUT"},
	{"user", "/api/pic/:id", "DELETE"},

	// Jabatan (user bisa lihat list dan tambah jabatan baru untuk dropdown profil)
	{"user", "/api/jabatan", "GET"},
	{"user", "/api/jabatan", "POST"},

	// Perusahaan (user hanya bisa lihat dan update miliknya sendiri)
	{"user", "/api/perusahaan", "GET"},
	{"user", "/api/perusahaan/:id", "GET"},
	{"user", "/api/perusahaan/:id", "PUT"},

	// ── MATURITY ─────────────────────────────────────────────────────────────

	// Ruang Lingkup (master — read only untuk user)
	{"user", "/api/maturity/ruang-lingkup", "GET"},
	{"user", "/api/maturity/ruang-lingkup/:id", "GET"},

	// Domain (master — read only untuk user)
	{"user", "/api/maturity/domain", "GET"},
	{"user", "/api/maturity/domain/:id", "GET"},

	// Kategori (master — read only untuk user)
	{"user", "/api/maturity/kategori", "GET"},
	{"user", "/api/maturity/kategori/:id", "GET"},

	// Sub Kategori (master — read only untuk user)
	{"user", "/api/maturity/sub-kategori", "GET"},
	{"user", "/api/maturity/sub-kategori/:id", "GET"},

	// IKAS
	{"user", "/api/maturity/ikas", "GET"},
	{"user", "/api/maturity/ikas", "POST"},
	{"user", "/api/maturity/ikas/:id", "GET"},
	{"user", "/api/maturity/ikas/:id", "PUT"},
	{"user", "/api/maturity/ikas/:id", "DELETE"},

	// Domain Identifikasi (read only untuk user)
	{"user", "/api/maturity/identifikasi", "GET"},
	{"user", "/api/maturity/identifikasi/:id", "GET"},

	// Domain Proteksi (read only untuk user)
	{"user", "/api/maturity/proteksi", "GET"},
	{"user", "/api/maturity/proteksi/:id", "GET"},

	// Domain Deteksi (read only untuk user)
	{"user", "/api/maturity/deteksi", "GET"},
	{"user", "/api/maturity/deteksi/:id", "GET"},

	// Domain Gulih (read only untuk user)
	{"user", "/api/maturity/gulih", "GET"},
	{"user", "/api/maturity/gulih/:id", "GET"},

	// Maturity (Pertanyaan — read only untuk user)
	{"user", "/api/maturity/pertanyaan-identifikasi", "GET"},
	{"user", "/api/maturity/pertanyaan-identifikasi/:id", "GET"},

	{"user", "/api/maturity/pertanyaan-proteksi", "GET"},
	{"user", "/api/maturity/pertanyaan-proteksi/:id", "GET"},

	{"user", "/api/maturity/pertanyaan-deteksi", "GET"},
	{"user", "/api/maturity/pertanyaan-deteksi/:id", "GET"},

	{"user", "/api/maturity/pertanyaan-gulih", "GET"},
	{"user", "/api/maturity/pertanyaan-gulih/:id", "GET"},

	// Maturity (Jawaban)
	{"user", "/api/maturity/jawaban-identifikasi", "GET"},
	{"user", "/api/maturity/jawaban-identifikasi", "POST"},
	{"user", "/api/maturity/jawaban-identifikasi/:id", "GET"},
	{"user", "/api/maturity/jawaban-identifikasi/:id", "PUT"},
	// {"user", "/api/maturity/jawaban-identifikasi/:id", "DELETE"},

	{"user", "/api/maturity/jawaban-proteksi", "GET"},
	{"user", "/api/maturity/jawaban-proteksi", "POST"},
	{"user", "/api/maturity/jawaban-proteksi/:id", "GET"},
	{"user", "/api/maturity/jawaban-proteksi/:id", "PUT"},
	// {"user", "/api/maturity/jawaban-proteksi/:id", "DELETE"},

	{"user", "/api/maturity/jawaban-deteksi", "GET"},
	{"user", "/api/maturity/jawaban-deteksi", "POST"},
	{"user", "/api/maturity/jawaban-deteksi/:id", "GET"},
	{"user", "/api/maturity/jawaban-deteksi/:id", "PUT"},
	// {"user", "/api/maturity/jawaban-deteksi/:id", "DELETE"},

	{"user", "/api/maturity/jawaban-gulih", "GET"},
	{"user", "/api/maturity/jawaban-gulih", "POST"},
	{"user", "/api/maturity/jawaban-gulih/:id", "GET"},
	{"user", "/api/maturity/jawaban-gulih/:id", "PUT"},
	// {"user", "/api/maturity/jawaban-gulih/:id", "DELETE"},

	// ── LMS ──────────────────────────────────────────────────────────────────
	// user bisa melihat kelas, detail kelas, update progress, dan mengikuti kuis.
 
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
}

// SeedCasbinPolicies memastikan semua default policy ada di database.
// Aman dijalankan berulang kali — tidak akan duplikat karena pakai HasPolicy check.
func SeedCasbinPolicies(casbinService *services.CasbinService) {
	added := 0
	skipped := 0

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

	logger.Infof("Casbin seeder selesai: %d policy ditambahkan, %d sudah ada", added, skipped)
}
