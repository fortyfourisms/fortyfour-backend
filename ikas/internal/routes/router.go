package routes

import (
	"encoding/json"
	"ikas/internal/handlers"
	"ikas/internal/middleware"
	"ikas/internal/utils"
	"net/http"
	"time"

	_ "ikas/docs"

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
	ikasH *handlers.IkasHandler,
	ruangLingkupH *handlers.RuangLingkupHandler,
	domainH *handlers.DomainHandler,
	kategoriH *handlers.KategoriHandler,
	subKategoriH *handlers.SubKategoriHandler,
	pertanyaanIdentifikasiH *handlers.PertanyaanIdentifikasiHandler,
	pertanyaanProteksiH *handlers.PertanyaanProteksiHandler,
	jawabanIdentifikasiH *handlers.JawabanIdentifikasiHandler,
	authM *middleware.AuthMiddleware,
	strictLimiter *middleware.RateLimiter,
	moderateLimiter *middleware.RateLimiter,
	lenientLimiter *middleware.RateLimiter,

) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/maturity/health", healthHandler)

	// Swagger UI
	mux.HandleFunc("/swagger/maturity/", httpSwagger.WrapHandler)

	mux.Handle("/api/maturity/ikas", authM.Authenticate(moderateLimiter.LimitByUser(utils.AdaptHandler(ikasH))))
	mux.Handle("/api/maturity/ikas/", authM.Authenticate(moderateLimiter.LimitByUser(utils.AdaptHandler(ikasH))))

	mux.Handle("/api/maturity/ruang-lingkup", authM.Authenticate(moderateLimiter.LimitByUser(utils.AdaptHandler(ruangLingkupH))))
	mux.Handle("/api/maturity/ruang-lingkup/", authM.Authenticate(moderateLimiter.LimitByUser(utils.AdaptHandler(ruangLingkupH))))

	mux.Handle("/api/maturity/domain", authM.Authenticate(moderateLimiter.LimitByUser(utils.AdaptHandler(domainH))))
	mux.Handle("/api/maturity/domain/", authM.Authenticate(moderateLimiter.LimitByUser(utils.AdaptHandler(domainH))))

	mux.Handle("/api/maturity/kategori", authM.Authenticate(moderateLimiter.LimitByUser(utils.AdaptHandler(kategoriH))))
	mux.Handle("/api/maturity/kategori/", authM.Authenticate(moderateLimiter.LimitByUser(utils.AdaptHandler(kategoriH))))

	mux.Handle("/api/maturity/sub-kategori", authM.Authenticate(moderateLimiter.LimitByUser(utils.AdaptHandler(subKategoriH))))
	mux.Handle("/api/maturity/sub-kategori/", authM.Authenticate(moderateLimiter.LimitByUser(utils.AdaptHandler(subKategoriH))))

	mux.Handle("/api/maturity/pertanyaan-identifikasi", authM.Authenticate(moderateLimiter.LimitByUser(utils.AdaptHandler(pertanyaanIdentifikasiH))))
	mux.Handle("/api/maturity/pertanyaan-identifikasi/", authM.Authenticate(moderateLimiter.LimitByUser(utils.AdaptHandler(pertanyaanIdentifikasiH))))

	mux.Handle("/api/maturity/pertanyaan-proteksi", authM.Authenticate(moderateLimiter.LimitByUser(utils.AdaptHandler(pertanyaanProteksiH))))
	mux.Handle("/api/maturity/pertanyaan-proteksi/", authM.Authenticate(moderateLimiter.LimitByUser(utils.AdaptHandler(pertanyaanProteksiH))))

	mux.Handle("/api/jawaban-identifikasi", authM.Authenticate(moderateLimiter.LimitByUser(utils.AdaptHandler(jawabanIdentifikasiH))))
	mux.Handle("/api/jawaban-identifikasi/", authM.Authenticate(moderateLimiter.LimitByUser(utils.AdaptHandler(jawabanIdentifikasiH))))

	return mux
}
