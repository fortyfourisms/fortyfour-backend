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
	authM *middleware.AuthMiddleware,
	strictLimiter *middleware.RateLimiter,
	moderateLimiter *middleware.RateLimiter,
	lenientLimiter *middleware.RateLimiter,
	csirtH *handlers.CsirtHandler,
) *http.ServeMux {
	mux := http.NewServeMux()

	// Routes Auth
	mux.HandleFunc("/api/register", strictLimiter.LimitByIP(authH.Register))
	mux.HandleFunc("/api/login", strictLimiter.LimitByIP(authH.Login))
	mux.HandleFunc("/api/refresh", strictLimiter.LimitByIP(authH.RefreshToken))
	mux.HandleFunc("/api/logout", authH.Logout)

	// Routes Posts
	mux.HandleFunc("/api/posts", authM.Authenticate(postH.GetPosts))
	mux.HandleFunc("/api/posts/single", authM.Authenticate(postH.GetPost))
	mux.HandleFunc("/api/posts/create", authM.Authenticate(postH.CreatePost))
	mux.HandleFunc("/api/posts/update", authM.Authenticate(postH.UpdatePost))
	mux.HandleFunc("/api/posts/delete", authM.Authenticate(postH.DeletePost))

	// Route Perusahaan
	mux.HandleFunc("/api/perusahaan", authM.Authenticate(moderateLimiter.LimitByUser(utils.AdaptHandler(perusahaanH))))
	mux.HandleFunc("/api/perusahaan/", authM.Authenticate(moderateLimiter.LimitByUser(utils.AdaptHandler(perusahaanH))))

	// Route PIC
	mux.HandleFunc("/api/pic", authM.Authenticate(moderateLimiter.LimitByUser(utils.AdaptHandler(picH))))
	mux.HandleFunc("/api/pic/", authM.Authenticate(moderateLimiter.LimitByUser(utils.AdaptHandler(picH))))

	// Route Perusahaan
	mux.HandleFunc("/api/jabatan", authM.Authenticate(moderateLimiter.LimitByUser(utils.AdaptHandler(jabatanH))))
	mux.HandleFunc("/api/jabatan/", authM.Authenticate(moderateLimiter.LimitByUser(utils.AdaptHandler(jabatanH))))

	// Route Identifikasi
	mux.HandleFunc("/api/identifikasi", authM.Authenticate(moderateLimiter.LimitByUser(utils.AdaptHandler(identifikasiH))))
	mux.HandleFunc("/api/identifikasi/", authM.Authenticate(moderateLimiter.LimitByUser(utils.AdaptHandler(identifikasiH))))

	// Route Gulih
	mux.HandleFunc("/api/gulih", authM.Authenticate(utils.AdaptHandler(gulihH)))
	mux.HandleFunc("/api/gulih/", authM.Authenticate(utils.AdaptHandler(gulihH)))

	// Route Proteksi
	mux.HandleFunc("/api/proteksi", authM.Authenticate(utils.AdaptHandler(proteksiH)))
	mux.HandleFunc("/api/proteksi/", authM.Authenticate(utils.AdaptHandler(proteksiH)))

	// Route Deteksi
	mux.HandleFunc("/api/deteksi", authM.Authenticate(utils.AdaptHandler(deteksiH)))
	mux.HandleFunc("/api/deteksi/", authM.Authenticate(utils.AdaptHandler(deteksiH)))

	// Route Ikas
	mux.HandleFunc("/api/ikas", authM.Authenticate(utils.AdaptHandler(ikasH)))
	mux.HandleFunc("/api/ikas/", authM.Authenticate(utils.AdaptHandler(ikasH)))

	// Route CSIRT
	mux.HandleFunc("/api/csirt", authM.Authenticate(utils.AdaptHandler(csirtH)))
	mux.HandleFunc("/api/csirt/", authM.Authenticate(utils.AdaptHandler(csirtH)))

	return mux
}