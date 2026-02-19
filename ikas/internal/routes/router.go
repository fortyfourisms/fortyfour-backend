package routes

import (
	"encoding/json"
	"ikas/internal/handlers"
	"ikas/internal/middleware"
	"ikas/internal/utils"
	"net/http"
	"time"
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
	ikasH *handlers.IkasHandler,
	ruangLingkupH *handlers.RuangLingkupHandler,
	domainH *handlers.DomainHandler,
	kategoriH *handlers.KategoriHandler,
	subKategoriH *handlers.SubKategoriHandler,
	authM *middleware.AuthMiddleware,
	strictLimiter *middleware.RateLimiter,
	moderateLimiter *middleware.RateLimiter,
	lenientLimiter *middleware.RateLimiter,

) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/health", healthHandler)
	mux.Handle("/api/ikas", authM.Authenticate(moderateLimiter.LimitByUser(utils.AdaptHandler(ikasH))))
	mux.Handle("/api/ikas/", authM.Authenticate(moderateLimiter.LimitByUser(utils.AdaptHandler(ikasH))))

	mux.Handle("/api/ruang-lingkup", authM.Authenticate(moderateLimiter.LimitByUser(utils.AdaptHandler(ruangLingkupH))))
	mux.Handle("/api/ruang-lingkup/", authM.Authenticate(moderateLimiter.LimitByUser(utils.AdaptHandler(ruangLingkupH))))

	mux.Handle("/api/domain", authM.Authenticate(moderateLimiter.LimitByUser(utils.AdaptHandler(domainH))))
	mux.Handle("/api/domain/", authM.Authenticate(moderateLimiter.LimitByUser(utils.AdaptHandler(domainH))))

	mux.Handle("/api/kategori", authM.Authenticate(moderateLimiter.LimitByUser(utils.AdaptHandler(kategoriH))))
	mux.Handle("/api/kategori/", authM.Authenticate(moderateLimiter.LimitByUser(utils.AdaptHandler(kategoriH))))

	mux.Handle("/api/sub-kategori", authM.Authenticate(moderateLimiter.LimitByUser(utils.AdaptHandler(subKategoriH))))
	mux.Handle("/api/sub-kategori/", authM.Authenticate(moderateLimiter.LimitByUser(utils.AdaptHandler(subKategoriH))))

	return mux
}
