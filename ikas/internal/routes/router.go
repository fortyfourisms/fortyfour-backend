package routes

import (
	"ikas/internal/handlers"
	"ikas/internal/middleware"
	"ikas/internal/utils"
	"net/http"
)

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
