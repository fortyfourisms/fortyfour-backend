package routes

import (
	"encoding/json"
	_ "fortyfour-backend/docs"
	"fortyfour-backend/internal/handlers"
	"fortyfour-backend/internal/middleware"
	"fortyfour-backend/internal/utils"
	"net/http"
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
	identifikasiH *handlers.IdentifikasiHandler,
	deteksiH *handlers.DeteksiHandler,
	gulihH *handlers.GulihHandler,
	ikasH *handlers.IkasHandler,
	proteksiH *handlers.ProteksiHandler,
	roleH *handlers.RoleHandler,
	casbinH *handlers.CasbinHandler,
	sseH *handlers.SSEHandler,
	authM *middleware.AuthMiddleware,
	casbinM *middleware.CasbinMiddleware,
	strictLimiter *middleware.RateLimiter,
	moderateLimiter *middleware.RateLimiter,
	lenientLimiter *middleware.RateLimiter,
	csirtH *handlers.CsirtHandler,
	sdmCsirtH *handlers.SdmCsirtHandler,
	sektorH *handlers.SektorHandler,
	subsectorH *handlers.SubSektorHandler,
	seH *handlers.SEHandler,
	dashboardH *handlers.DashboardHandler,
) http.Handler {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/api/health", healthHandler)

	// Routes Auth
	mux.HandleFunc("/api/register", strictLimiter.LimitByIP(authH.Register))
	mux.HandleFunc("/api/login", strictLimiter.LimitByIP(authH.Login))
	mux.HandleFunc("/api/refresh", strictLimiter.LimitByIP(authH.Refresh))
	mux.HandleFunc("/api/logout", authH.Logout)

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

	// Route Identifikasi
	mux.HandleFunc("/api/identifikasi", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(identifikasiH)))))
	mux.HandleFunc("/api/identifikasi/", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(identifikasiH)))))

	// Route Gulih
	mux.HandleFunc("/api/gulih", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(gulihH)))))
	mux.HandleFunc("/api/gulih/", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(gulihH)))))

	// Route Proteksi
	mux.HandleFunc("/api/proteksi", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(proteksiH)))))
	mux.HandleFunc("/api/proteksi/", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(proteksiH)))))

	// Route Deteksi
	mux.HandleFunc("/api/deteksi", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(deteksiH)))))
	mux.HandleFunc("/api/deteksi/", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(deteksiH)))))

	// Route Ikas
	mux.HandleFunc("/api/ikas", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(ikasH)))))
	mux.HandleFunc("/api/ikas/", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(ikasH)))))

	// Route Role
	mux.HandleFunc("/api/role", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(roleH)))))
	mux.HandleFunc("/api/role/", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(roleH)))))

	// Route CSIRT
	mux.HandleFunc("/api/csirt", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(csirtH)))))
	mux.HandleFunc("/api/csirt/", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(csirtH)))))

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
	mux.HandleFunc("/api/se", authM.Authenticate(utils.AdaptHandler(seH)))
	mux.HandleFunc("/api/se/", authM.Authenticate(utils.AdaptHandler(seH)))

	// Route Dashboard 
	// Summary: counts per sektor + ikas + se
	mux.HandleFunc("/api/dashboard/summary", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(dashboardH)))))

	// Swagger UI
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	return (mux)
}