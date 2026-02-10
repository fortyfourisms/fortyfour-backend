package routes

import (
	"encoding/json"
	_ "fortyfour-backend/docs"
	"fortyfour-backend/internal/handlers"
	"fortyfour-backend/internal/middleware"
	"fortyfour-backend/internal/utils"
	"net/http"
	"time"

	"github.com/rs/cors"
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
	seCsirtH *handlers.SeCsirtHandler,
) http.Handler {
	mux := http.NewServeMux()

	// In main():
	mux.HandleFunc("/api/health", healthHandler)

	// Routes Auth
	mux.HandleFunc("/api/register", strictLimiter.LimitByIP(authH.Register))
	mux.HandleFunc("/api/login", strictLimiter.LimitByIP(authH.Login))
	mux.HandleFunc("/api/refresh", strictLimiter.LimitByIP(authH.Refresh))
	mux.HandleFunc("/api/logout", authH.Logout)

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

	// Route SE_CSIRT
	mux.HandleFunc("/api/se_csirt", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(seCsirtH)))))
	mux.HandleFunc("/api/se_csirt/", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(seCsirtH)))))

	// Swagger UI
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	// CORS Configuration
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "https://admin.kssindustri.site", "https://fortyfouris.netlify.app"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "Origin", "Accept"},
		AllowCredentials: true,
	})

	return c.Handler(mux)
}
