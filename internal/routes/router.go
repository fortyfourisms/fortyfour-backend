package routes

import (
	"encoding/json"
	_ "fortyfour-backend/docs"
	"fortyfour-backend/internal/handlers"
	"fortyfour-backend/internal/middleware"
	"fortyfour-backend/internal/utils"
	"net/http"
	"strings"
	"time"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @Summary Health check
// @Description Check if the API is running and healthy
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/health [get]
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func InitRouter(
	authH *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	perusahaanH *handlers.PerusahaanHandler,
	picH *handlers.PICHandler,
	jabatanH *handlers.JabatanHandler,
	roleH *handlers.RoleHandler,
	casbinH *handlers.CasbinHandler,
	sseH *handlers.SSEHandler,
	authM *middleware.AuthMiddleware,
	casbinM *middleware.CasbinMiddleware,
	strictLimiter *middleware.RateLimiter,
	moderateLimiter *middleware.RateLimiter,
	lenientLimiter *middleware.RateLimiter,
	csirtH *handlers.CsirtHandler,
	csirtExportH *handlers.CsirtExportHandler,
	sdmCsirtH *handlers.SdmCsirtHandler,
	chatHandler *handlers.ChatHandler,
	sektorH *handlers.SektorHandler,
	subsectorH *handlers.SubSektorHandler,
	seH *handlers.SEHandler,
	seExportH *handlers.SEExportHandler,
	dashboardH *handlers.DashboardHandler,
	notificationH *handlers.NotificationHandler,
	ikasProxyH *handlers.ProxyHandler,
	lmsH *handlers.LMSHandler,
) http.Handler {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/api/health", healthHandler)

	// Routes Auth
	mux.HandleFunc("/api/register", strictLimiter.LimitByIP(authH.Register))
	mux.HandleFunc("/api/login", strictLimiter.LimitByIP(authH.Login))
	mux.HandleFunc("/api/refresh", strictLimiter.LimitByIP(authH.Refresh))
	mux.HandleFunc("/api/logout", authH.Logout)
	mux.HandleFunc("/api/logout-all", authM.Authenticate(authH.LogoutAll))

	// Route Me — hanya untuk user yang sedang login (GET: lihat profil, PUT: update profil sendiri)
	mux.HandleFunc("/api/me", authM.Authenticate(moderateLimiter.LimitByUser(authH.MeRouter)))
	mux.HandleFunc("/api/me/", authM.Authenticate(moderateLimiter.LimitByUser(authH.MeRouter)))

	// MFA endpoints
	// setup & enable -> protected (require Authorization header)
	mux.HandleFunc("/api/mfa/setup", strictLimiter.LimitByIP(authH.SetupMFA))
	mux.HandleFunc("/api/mfa/enable", strictLimiter.LimitByIP(authH.EnableMFA))
	// verify -> public (called with mfa_token)
	mux.HandleFunc("/api/mfa/verify", strictLimiter.LimitByIP(authH.VerifyMFA))

	// Routes Users
	mux.HandleFunc("/api/users", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(userHandler)))))
	mux.HandleFunc("/api/users/", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(userHandler)))))

	// Routes Casbin Management (only admin)
	mux.HandleFunc("/api/casbin/policies", authM.Authenticate(casbinH.GetAllPolicies))
	mux.HandleFunc("/api/casbin/policies/add", authM.Authenticate(casbinH.AddPolicy))
	mux.HandleFunc("/api/casbin/policies/bulk", authM.Authenticate(casbinH.BulkAddPolicies))
	mux.HandleFunc("/api/casbin/policies/remove", authM.Authenticate(casbinH.RemovePolicy))
	mux.HandleFunc("/api/casbin/permissions", authM.Authenticate(casbinH.GetRolePermissions))

	// SSE Routes
	mux.HandleFunc("/api/events", authM.Authenticate(sseH.HandleSSE))
	mux.HandleFunc("/api/events/stats", authM.Authenticate(sseH.GetStats))

	// Route Perusahaan
	// GET /api/perusahaan/dropdown -> PUBLIC untuk dropdown register (minimal data: id, nama)
	mux.HandleFunc("/api/perusahaan/dropdown", moderateLimiter.LimitByUser(utils.AdaptHandler(perusahaanH)))

	// GET /api/perusahaan (list all) -> AUTHENTICATED (full data)
	// Other methods & GET with ID -> AUTHENTICATED
	mux.HandleFunc("/api/perusahaan", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(perusahaanH)))))
	mux.HandleFunc("/api/perusahaan/", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(perusahaanH)))))

	// Route PIC
	mux.HandleFunc("/api/pic", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(picH)))))
	mux.HandleFunc("/api/pic/", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(picH)))))

	// Route Jabatan
	mux.HandleFunc("/api/jabatan", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(jabatanH)))))
	mux.HandleFunc("/api/jabatan/", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(jabatanH)))))

	// Route IKAS (Proxy to Microservice)
	mux.Handle("/api/maturity/", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(ikasProxyH.ServeHTTP))))
	mux.Handle("/api/maturity", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(ikasProxyH.ServeHTTP))))

	// Route Role
	mux.HandleFunc("/api/role", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(roleH)))))
	mux.HandleFunc("/api/role/", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(roleH)))))

	// Route CSIRT
	// "/api/csirt/" menangkap semua sub-path termasuk {id}/export-pdf.
	// Di dalam handler, path yang mengandung "export-pdf" diarahkan ke csirtExportH,
	// sisanya ke csirtH (CRUD) — ini satu-satunya cara karena Go ServeMux tidak support wildcard.
	mux.HandleFunc("/api/csirt", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(csirtH)))))
	mux.HandleFunc("/api/csirt/", authM.Authenticate(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "export-pdf") {
			utils.AdaptHandler(csirtExportH)(w, r)
		} else {
			casbinM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(csirtH)))(w, r)
		}
	}))

	// Route SDM_CSIRT
	mux.HandleFunc("/api/sdm_csirt", authM.Authenticate(casbinM.Authorize(utils.AdaptHandler(sdmCsirtH))))
	mux.HandleFunc("/api/sdm_csirt/", authM.Authenticate(casbinM.Authorize(utils.AdaptHandler(sdmCsirtH))))

	// Route Sektor
	mux.HandleFunc("/api/sektor", authM.Authenticate(utils.AdaptHandler(sektorH)))
	mux.HandleFunc("/api/sektor/", authM.Authenticate(utils.AdaptHandler(sektorH)))

	// Route SubSektor
	mux.HandleFunc("/api/sub_sektor", authM.Authenticate(utils.AdaptHandler(subsectorH)))
	mux.HandleFunc("/api/sub_sektor/", authM.Authenticate(utils.AdaptHandler(subsectorH)))

	// Route SE
	// "/api/se/" menangkap semua sub-path termasuk {id}/export-pdf.
	// Di dalam handler, path yang mengandung "export-pdf" diarahkan ke seExportH,
	// sisanya ke seH (CRUD).
	mux.HandleFunc("/api/se", authM.Authenticate(utils.AdaptHandler(seH)))
	mux.HandleFunc("/api/se/", authM.Authenticate(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "export-pdf") {
			utils.AdaptHandler(seExportH)(w, r)
		} else {
			utils.AdaptHandler(seH)(w, r)
		}
	}))

	// Route Dashboard
	// Summary: counts per sektor + ikas + se
	mux.HandleFunc("/api/dashboard/summary", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(dashboardH)))))

	// Routes Notifications
	mux.HandleFunc("/api/notifications", authM.Authenticate(utils.AdaptHandler(notificationH)))
	mux.HandleFunc("/api/notifications/", authM.Authenticate(utils.AdaptHandler(notificationH)))

	// Routes Chat
	mux.HandleFunc("/api/chat", authM.Authenticate(chatHandler.Stream))
	mux.HandleFunc("/api/chat/delete-session", authM.Authenticate(chatHandler.DeleteSession))

	// ── LMS Routes ────────────────────────────────────────────────────────────
	//
	// USER & ADMIN (authenticated):
	//   GET  /api/kelas                              → list kelas (user: published only)
	//   GET  /api/kelas/{id}                         → detail kelas + materi + progress
	//   POST /api/materi/{id}/progress               → update progress video/teks
	//   GET  /api/materi/{id}/file-pendukung         → list file pendukung (PDF)
	//   GET  /api/materi/{id}/diskusi                → list diskusi per materi
	//   POST /api/materi/{id}/diskusi                → buat diskusi/reply
	//   GET  /api/materi/{id}/catatan                → get catatan pribadi
	//   PUT  /api/materi/{id}/catatan                → upsert catatan pribadi
	//   GET  /api/kelas/{id}/kuis                    → list kuis dalam kelas
	//   POST /api/kuis/{id_kuis}/start               → mulai kuis
	//   POST /api/kuis/attempt/{id_attempt}/submit   → submit jawaban kuis
	//   GET  /api/kuis/attempt/{id_attempt}/result   → lihat hasil kuis
	//   GET  /api/kelas/{id}/sertifikat              → cek sertifikat user
	//   POST /api/kelas/{id}/sertifikat/generate     → generate sertifikat
	//   GET  /api/sertifikat/me                      → list sertifikat user
	//   GET  /api/sertifikat/{id}                    → detail sertifikat
	//   GET  /api/sertifikat/{id}/download           → download PDF sertifikat
	//   GET  /api/file-pendukung/{id}/download       → download file pendukung
	//
	// ADMIN ONLY (authenticated + casbin):
	//   POST   /api/kelas                            → buat kelas baru
	//   PUT    /api/kelas/{id}                       → update kelas
	//   DELETE /api/kelas/{id}                       → hapus kelas
	//   POST   /api/kelas/{id}/materi                → tambah materi ke kelas
	//   PUT    /api/materi/{id}                      → update materi
	//   DELETE /api/materi/{id}                      → hapus materi
	//   POST   /api/materi/{id}/file-pendukung       → upload file pendukung
	//   DELETE /api/file-pendukung/{id}              → hapus file pendukung
	//   POST   /api/kelas/{id}/kuis                  → buat kuis (per-materi/final)
	//   PUT    /api/kuis/{id}                        → update kuis
	//   DELETE /api/kuis/{id}                        → hapus kuis
	//   GET    /api/kuis/{id_kuis}/soal              → list soal kuis (admin)
	//   POST   /api/kuis/{id_kuis}/soal              → tambah soal ke kuis
	//   PUT    /api/soal/{id}                        → update soal
	//   DELETE /api/soal/{id}                        → hapus soal

	// /api/kelas — GET (user & admin), POST (admin)
	mux.HandleFunc("/api/kelas", authM.Authenticate(moderateLimiter.LimitByUser(
		func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				lmsH.ServeKelas(w, r)
			case http.MethodPost:
				casbinM.Authorize(lmsH.ServeKelas)(w, r)
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		},
	)))

	// /api/kelas/{id}, /api/kelas/{id}/materi, /api/kelas/{id}/kuis, /api/kelas/{id}/sertifikat
	mux.HandleFunc("/api/kelas/", authM.Authenticate(moderateLimiter.LimitByUser(
		func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path

			// /api/kelas/{id}/materi → POST admin
			if strings.HasSuffix(path, "/materi") {
				if r.Method == http.MethodPost {
					casbinM.Authorize(lmsH.ServeMateriByKelas)(w, r)
					return
				}
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}

			// /api/kelas/{id}/kuis → GET (user & admin), POST (admin)
			if strings.Contains(path, "/kuis") {
				switch r.Method {
				case http.MethodGet:
					lmsH.ServeKelas(w, r)
				case http.MethodPost:
					casbinM.Authorize(lmsH.ServeKelas)(w, r)
				default:
					w.WriteHeader(http.StatusMethodNotAllowed)
				}
				return
			}

			// /api/kelas/{id}/sertifikat → GET/POST user
			if strings.Contains(path, "/sertifikat") {
				lmsH.ServeKelas(w, r)
				return
			}

			// /api/kelas/{id} → GET user & admin, PUT/DELETE admin
			switch r.Method {
			case http.MethodGet:
				lmsH.ServeKelas(w, r)
			case http.MethodPut, http.MethodDelete:
				casbinM.Authorize(lmsH.ServeKelas)(w, r)
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		},
	)))

	// /api/materi/{id}, /api/materi/{id}/progress, /api/materi/{id}/file-pendukung,
	// /api/materi/{id}/diskusi, /api/materi/{id}/catatan
	mux.HandleFunc("/api/materi/", authM.Authenticate(moderateLimiter.LimitByUser(
		func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path

			// /api/materi/{id}/progress → POST user & admin
			if strings.Contains(path, "/progress") {
				lmsH.ServeMateri(w, r)
				return
			}

			// /api/materi/{id}/file-pendukung → GET user, POST admin
			if strings.Contains(path, "/file-pendukung") {
				switch r.Method {
				case http.MethodGet:
					lmsH.ServeMateri(w, r)
				case http.MethodPost:
					casbinM.Authorize(lmsH.ServeMateri)(w, r)
				default:
					w.WriteHeader(http.StatusMethodNotAllowed)
				}
				return
			}

			// /api/materi/{id}/diskusi → GET/POST user
			if strings.Contains(path, "/diskusi") {
				lmsH.ServeMateri(w, r)
				return
			}

			// /api/materi/{id}/catatan → GET/PUT user
			if strings.Contains(path, "/catatan") {
				lmsH.ServeMateri(w, r)
				return
			}

			// /api/materi/{id} → PUT/DELETE admin
			switch r.Method {
			case http.MethodPut, http.MethodDelete:
				casbinM.Authorize(lmsH.ServeMateri)(w, r)
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		},
	)))

	// /api/file-pendukung/{id} → DELETE admin, /api/file-pendukung/{id}/download → GET user
	mux.HandleFunc("/api/file-pendukung/", authM.Authenticate(moderateLimiter.LimitByUser(
		func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			if strings.HasSuffix(path, "/download") && r.Method == http.MethodGet {
				lmsH.ServeFilePendukung(w, r)
				return
			}
			casbinM.Authorize(lmsH.ServeFilePendukung)(w, r)
		},
	)))

	// /api/diskusi/{id} → PUT/DELETE user
	mux.HandleFunc("/api/diskusi/", authM.Authenticate(moderateLimiter.LimitByUser(lmsH.ServeDiskusi)))

	// /api/soal/{id} → PUT/DELETE admin
	mux.HandleFunc("/api/soal/", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(lmsH.ServeSoal))))

	// /api/kuis/ — admin CRUD + user start/submit/result + soal management
	mux.HandleFunc("/api/kuis/", authM.Authenticate(moderateLimiter.LimitByUser(
		func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path

			// User routes: start, attempt/submit, attempt/result
			if strings.Contains(path, "/start") || strings.Contains(path, "/attempt/") {
				lmsH.ServeKuis(w, r)
				return
			}

			// /api/kuis/{id}/soal → GET/POST admin
			if strings.Contains(path, "/soal") {
				casbinM.Authorize(lmsH.ServeKuis)(w, r)
				return
			}

			// PUT/DELETE admin
			switch r.Method {
			case http.MethodPut, http.MethodDelete:
				casbinM.Authorize(lmsH.ServeKuis)(w, r)
			default:
				lmsH.ServeKuis(w, r)
			}
		},
	)))

	// /api/sertifikat → user routes
	mux.HandleFunc("/api/sertifikat/", authM.Authenticate(moderateLimiter.LimitByUser(lmsH.ServeSertifikat)))
	mux.HandleFunc("/api/sertifikat", authM.Authenticate(moderateLimiter.LimitByUser(lmsH.ServeSertifikat)))

	// Swagger UI
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	return mux
}