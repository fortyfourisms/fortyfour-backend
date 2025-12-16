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
	casbinH *handlers.CasbinHandler,
	authM *middleware.AuthMiddleware,
	casbinM *middleware.CasbinMiddleware,
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

	// Routes Casbin Management (only admin)
	mux.HandleFunc("/api/casbin/policies", authM.Authenticate(casbinH.GetAllPolicies))
	mux.HandleFunc("/api/casbin/policies/add", authM.Authenticate(casbinH.AddPolicy))
	mux.HandleFunc("/api/casbin/policies/bulk", authM.Authenticate(casbinH.BulkAddPolicies))
	mux.HandleFunc("/api/casbin/policies/remove", authM.Authenticate(casbinH.RemovePolicy))
	mux.HandleFunc("/api/casbin/permissions", authM.Authenticate(casbinH.GetRolePermissions))

	// Routes Posts
	mux.HandleFunc("/api/posts", authM.Authenticate(casbinM.Authorize(postH.GetPosts)))
	mux.HandleFunc("/api/posts/single", authM.Authenticate(casbinM.Authorize(postH.GetPost)))
	mux.HandleFunc("/api/posts/create", authM.Authenticate(casbinM.Authorize(postH.CreatePost)))
	mux.HandleFunc("/api/posts/update", authM.Authenticate(casbinM.Authorize(postH.UpdatePost)))
	mux.HandleFunc("/api/posts/delete", authM.Authenticate(casbinM.Authorize(postH.DeletePost)))

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
	mux.HandleFunc("/api/gulih", authM.Authenticate(casbinM.Authorize(utils.AdaptHandler(gulihH))))
	mux.HandleFunc("/api/gulih/", authM.Authenticate(casbinM.Authorize(utils.AdaptHandler(gulihH))))

	// Route Proteksi
	mux.HandleFunc("/api/proteksi", authM.Authenticate(casbinM.Authorize(utils.AdaptHandler(proteksiH))))
	mux.HandleFunc("/api/proteksi/", authM.Authenticate(casbinM.Authorize(utils.AdaptHandler(proteksiH))))

	// Route Deteksi
	mux.HandleFunc("/api/deteksi", authM.Authenticate(casbinM.Authorize(utils.AdaptHandler(deteksiH))))
	mux.HandleFunc("/api/deteksi/", authM.Authenticate(casbinM.Authorize(utils.AdaptHandler(deteksiH))))

	// Route Ikas
	mux.HandleFunc("/api/ikas", authM.Authenticate(casbinM.Authorize(utils.AdaptHandler(ikasH))))
	mux.HandleFunc("/api/ikas/", authM.Authenticate(casbinM.Authorize(utils.AdaptHandler(ikasH))))

	// Route Role
	mux.HandleFunc("/api/role", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(roleH)))))
	mux.HandleFunc("/api/role/", authM.Authenticate(casbinM.Authorize(moderateLimiter.LimitByUser(utils.AdaptHandler(roleH)))))

	return mux
}
