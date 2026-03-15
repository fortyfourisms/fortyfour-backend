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

	// Identifikasi
	{"user", "/api/identifikasi", "GET"},
	{"user", "/api/identifikasi", "POST"},
	{"user", "/api/identifikasi/:id", "GET"},
	{"user", "/api/identifikasi/:id", "PUT"},
	{"user", "/api/identifikasi/:id", "DELETE"},

	// Proteksi
	{"user", "/api/proteksi", "GET"},
	{"user", "/api/proteksi", "POST"},
	{"user", "/api/proteksi/:id", "GET"},
	{"user", "/api/proteksi/:id", "PUT"},
	{"user", "/api/proteksi/:id", "DELETE"},

	// Deteksi
	{"user", "/api/deteksi", "GET"},
	{"user", "/api/deteksi", "POST"},
	{"user", "/api/deteksi/:id", "GET"},
	{"user", "/api/deteksi/:id", "PUT"},
	{"user", "/api/deteksi/:id", "DELETE"},

	// Gulih
	{"user", "/api/gulih", "GET"},
	{"user", "/api/gulih", "POST"},
	{"user", "/api/gulih/:id", "GET"},
	{"user", "/api/gulih/:id", "PUT"},
	{"user", "/api/gulih/:id", "DELETE"},

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

	// Profile & data diri sendiri
	{"user", "/api/users/profile", "GET"},
	{"user", "/api/users/profile", "PUT"},
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
