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
	chatHandler *handlers.ChatHandler,
) http.Handler {
	mux := http.NewServeMux()
	apiV1 := http.NewServeMux()

	// In main():
	apiV1.HandleFunc("/health", healthHandler)

	// Routes Auth
	apiV1.HandleFunc("/register", strictLimiter.LimitByIP(authH.Register))
	apiV1.HandleFunc("/login", strictLimiter.LimitByIP(authH.Login))
	apiV1.HandleFunc("/refresh", strictLimiter.LimitByIP(authH.RefreshToken))
	apiV1.HandleFunc("/logout", authH.Logout)

	// Routes Users
	apiV1.HandleFunc("/users", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(userHandler)))))
	apiV1.HandleFunc("/users/", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(userHandler)))))

	// Routes Casbin Management (only admin)
	apiV1.HandleFunc("/casbin/policies", authM.Authenticate(casbinH.GetAllPolicies))
	apiV1.HandleFunc("/casbin/policies/add", authM.Authenticate(casbinH.AddPolicy))
	apiV1.HandleFunc("/casbin/policies/bulk", authM.Authenticate(casbinH.BulkAddPolicies))
	apiV1.HandleFunc("/casbin/policies/remove", authM.Authenticate(casbinH.RemovePolicy))
	apiV1.HandleFunc("/casbin/permissions", authM.Authenticate(casbinH.GetRolePermissions))

	// SSE Routes
	apiV1.HandleFunc("/events", authM.Authenticate(sseH.HandleSSE))
	apiV1.HandleFunc("/events/stats", authM.Authenticate(sseH.GetStats))

	// Route Perusahaan
	apiV1.HandleFunc("/perusahaan", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(perusahaanH)))))
	apiV1.HandleFunc("/perusahaan/", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(perusahaanH)))))

	// Route PIC
	apiV1.HandleFunc("/pic", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(picH)))))
	apiV1.HandleFunc("/pic/", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(picH)))))

	// Route Jabatan
	apiV1.HandleFunc("/jabatan", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(jabatanH)))))
	apiV1.HandleFunc("/jabatan/", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(jabatanH)))))

	// Route Identifikasi
	apiV1.HandleFunc("/identifikasi", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(identifikasiH)))))
	apiV1.HandleFunc("/identifikasi/", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(identifikasiH)))))

	// Route Gulih
	apiV1.HandleFunc("/gulih", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(gulihH)))))
	apiV1.HandleFunc("/gulih/", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(gulihH)))))

	// Route Proteksi
	apiV1.HandleFunc("/proteksi", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(proteksiH)))))
	apiV1.HandleFunc("/proteksi/", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(proteksiH)))))

	// Route Deteksi
	apiV1.HandleFunc("/deteksi", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(deteksiH)))))
	apiV1.HandleFunc("/deteksi/", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(deteksiH)))))

	// Route Role
	apiV1.HandleFunc("/role", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(roleH)))))
	apiV1.HandleFunc("/role/", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(roleH)))))

	// Route CSIRT
	apiV1.HandleFunc("/csirt", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(csirtH)))))
	apiV1.HandleFunc("/csirt/", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(csirtH)))))

	// Route SDM_CSIRT
	apiV1.HandleFunc("/sdm_csirt", authM.Authenticate(casbinM.Authorize(utils.AdaptHandler(sdmCsirtH))))
	apiV1.HandleFunc("/sdm_csirt/", authM.Authenticate(casbinM.Authorize(utils.AdaptHandler(sdmCsirtH))))

	// Route SE_CSIRT
	apiV1.HandleFunc("/se_csirt", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(seCsirtH)))))
	apiV1.HandleFunc("/se_csirt/", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(seCsirtH)))))

	// Routes Chat
	apiV1.HandleFunc("/chat", authM.Authenticate(chatHandler.Stream))
	apiV1.HandleFunc("/chat/delete-session", authM.Authenticate(chatHandler.DeleteSession))

	// Swagger UI
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	mux.Handle("/api/", http.StripPrefix("/api", apiV1))

	// CORS Configuration
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "https://admin.kssindustri.site", "https://fortyfouris.netlify.app"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "Origin", "Accept"},
		AllowCredentials: true,
	})

	return c.Handler(mux)
}
