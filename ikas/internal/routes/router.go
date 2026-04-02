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

const baseURL = "/api/maturity"

func handle(mux *http.ServeMux, path string, handler http.Handler) {
	mux.Handle(baseURL+path, handler)
	mux.Handle(baseURL+path+"/", handler)
}

func InitRouter(
	ikasH *handlers.IkasHandler,
	ruangLingkupH *handlers.RuangLingkupHandler,
	domainH *handlers.DomainHandler,
	kategoriH *handlers.KategoriHandler,
	subKategoriH *handlers.SubKategoriHandler,
	identifikasiH *handlers.IdentifikasiHandler,
	proteksiH *handlers.ProteksiHandler,
	deteksiH *handlers.DeteksiHandler,
	gulihH *handlers.GulihHandler,
	pertanyaanIdentifikasiH *handlers.PertanyaanIdentifikasiHandler,
	pertanyaanProteksiH *handlers.PertanyaanProteksiHandler,
	pertanyaanDeteksiH *handlers.PertanyaanDeteksiHandler,
	pertanyaanGulihH *handlers.PertanyaanGulihHandler,
	jawabanIdentifikasiH *handlers.JawabanIdentifikasiHandler,
	jawabanProteksiH *handlers.JawabanProteksiHandler,
	jawabanDeteksiH *handlers.JawabanDeteksiHandler,
	jawabanGulihH *handlers.JawabanGulihHandler,
	authM *middleware.AuthMiddleware,
	casbinM *middleware.CasbinMiddleware,
) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc(baseURL+"/health", healthHandler)
	mux.HandleFunc("/swagger/maturity/", httpSwagger.WrapHandler)

	withAuth := func(h http.HandlerFunc) http.Handler {
		return authM.Authenticate(casbinM.Authorize(h))
	}

	handle(mux, "/ikas", withAuth(utils.AdaptHandler(ikasH)))
	handle(mux, "/ruang-lingkup", withAuth(utils.AdaptHandler(ruangLingkupH)))
	handle(mux, "/domain", withAuth(utils.AdaptHandler(domainH)))
	handle(mux, "/kategori", withAuth(utils.AdaptHandler(kategoriH)))
	handle(mux, "/sub-kategori", withAuth(utils.AdaptHandler(subKategoriH)))
	handle(mux, "/identifikasi", withAuth(utils.AdaptHandler(identifikasiH)))
	handle(mux, "/proteksi", withAuth(utils.AdaptHandler(proteksiH)))
	handle(mux, "/deteksi", withAuth(utils.AdaptHandler(deteksiH)))
	handle(mux, "/gulih", withAuth(utils.AdaptHandler(gulihH)))
	handle(mux, "/pertanyaan-identifikasi", withAuth(utils.AdaptHandler(pertanyaanIdentifikasiH)))
	handle(mux, "/pertanyaan-proteksi", withAuth(utils.AdaptHandler(pertanyaanProteksiH)))
	handle(mux, "/pertanyaan-deteksi", withAuth(utils.AdaptHandler(pertanyaanDeteksiH)))
	handle(mux, "/pertanyaan-gulih", withAuth(utils.AdaptHandler(pertanyaanGulihH)))
	handle(mux, "/jawaban-identifikasi", withAuth(utils.AdaptHandler(jawabanIdentifikasiH)))
	handle(mux, "/jawaban-proteksi", withAuth(utils.AdaptHandler(jawabanProteksiH)))
	handle(mux, "/jawaban-deteksi", withAuth(utils.AdaptHandler(jawabanDeteksiH)))
	handle(mux, "/jawaban-gulih", withAuth(utils.AdaptHandler(jawabanGulihH)))

	return mux
}
