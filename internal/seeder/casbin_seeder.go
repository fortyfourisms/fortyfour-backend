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
