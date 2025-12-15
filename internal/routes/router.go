package routes

import (
	"fortyfour-backend/internal/handlers"
	"fortyfour-backend/internal/middleware"
	"fortyfour-backend/internal/utils"
	"net/http"
)

func InitRouter(
	authH *handlers.AuthHandler,
	postH *handlers.PostHandler,
	perusahaanH *handlers.PerusahaanHandler,
	picH *handlers.PICHandler,
	jabatanH *handlers.JabatanHandler,
	identifikasiH *handlers.IdentifikasiHandler,
	deteksiH *handlers.DeteksiHandler,
	gulihH *handlers.GulihHandler,
	ikasH *handlers.IkasHandler,
	proteksiH *handlers.ProteksiHandler,
	roleH *handlers.RoleHandler,
	authM *middleware.AuthMiddleware,
	authzM *middleware.AuthorizationMiddleware,
	strictLimiter *middleware.RateLimiter,
	moderateLimiter *middleware.RateLimiter,
	lenientLimiter *middleware.RateLimiter,
) *http.ServeMux {
	mux := http.NewServeMux()

	// Routes Auth
	mux.HandleFunc("/api/register", strictLimiter.LimitByIP(authH.Register))
	mux.HandleFunc("/api/login", strictLimiter.LimitByIP(authH.Login))
	mux.HandleFunc("/api/refresh", strictLimiter.LimitByIP(authH.RefreshToken))
	mux.HandleFunc("/api/logout", authH.Logout)

	// Routes Posts
	mux.HandleFunc("/api/posts", authM.Authenticate(authzM.Authorize(postH.GetPosts)))
	mux.HandleFunc("/api/posts/single", authM.Authenticate(authzM.Authorize(postH.GetPost)))
	mux.HandleFunc("/api/posts/create", authM.Authenticate(authzM.Authorize(postH.CreatePost)))
	mux.HandleFunc("/api/posts/update", authM.Authenticate(authzM.Authorize(postH.UpdatePost)))
	mux.HandleFunc("/api/posts/delete", authM.Authenticate(authzM.Authorize(postH.DeletePost)))

	// Route Perusahaan
	mux.HandleFunc("/api/perusahaan", authM.Authenticate(authzM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(perusahaanH)))))
	mux.HandleFunc("/api/perusahaan/", authM.Authenticate(authzM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(perusahaanH)))))

	// Route PIC
	mux.HandleFunc("/api/pic", authM.Authenticate(authzM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(picH)))))
	mux.HandleFunc("/api/pic/", authM.Authenticate(authzM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(picH)))))

	// Route Jabatan
	mux.HandleFunc("/api/jabatan", authM.Authenticate(authzM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(jabatanH)))))
	mux.HandleFunc("/api/jabatan/", authM.Authenticate(authzM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(jabatanH)))))

	// Route Identifikasi
	mux.HandleFunc("/api/identifikasi", authM.Authenticate(authzM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(identifikasiH)))))
	mux.HandleFunc("/api/identifikasi/", authM.Authenticate(authzM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(identifikasiH)))))

	// Route Gulih
	mux.HandleFunc("/api/gulih", authM.Authenticate(authzM.Authorize(utils.AdaptHandler(gulihH))))
	mux.HandleFunc("/api/gulih/", authM.Authenticate(authzM.Authorize(utils.AdaptHandler(gulihH))))

	// Route Proteksi
	mux.HandleFunc("/api/proteksi", authM.Authenticate(authzM.Authorize(utils.AdaptHandler(proteksiH))))
	mux.HandleFunc("/api/proteksi/", authM.Authenticate(authzM.Authorize(utils.AdaptHandler(proteksiH))))

	// Route Deteksi
	mux.HandleFunc("/api/deteksi", authM.Authenticate(authzM.Authorize(utils.AdaptHandler(deteksiH))))
	mux.HandleFunc("/api/deteksi/", authM.Authenticate(authzM.Authorize(utils.AdaptHandler(deteksiH))))

	// Route Ikas
	mux.HandleFunc("/api/ikas", authM.Authenticate(authzM.Authorize(utils.AdaptHandler(ikasH))))
	mux.HandleFunc("/api/ikas/", authM.Authenticate(authzM.Authorize(utils.AdaptHandler(ikasH))))

	// Route Role
	mux.HandleFunc("/api/role", authM.Authenticate(authzM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(roleH)))))
	mux.HandleFunc("/api/role/", authM.Authenticate(authzM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(roleH)))))

	return mux
}
