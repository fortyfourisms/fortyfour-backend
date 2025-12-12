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
	gulihH *handlers.GulihHandler,
	ikasH *handlers.IkasHandler,
	proteksiH *handlers.ProteksiHandler,
	authM *middleware.AuthMiddleware,
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
	mux.HandleFunc("/api/posts", authM.Authenticate(postH.GetPosts))
	mux.HandleFunc("/api/posts/single", authM.Authenticate(postH.GetPost))
	mux.HandleFunc("/api/posts/create", authM.Authenticate(postH.CreatePost))
	mux.HandleFunc("/api/posts/update", authM.Authenticate(postH.UpdatePost))
	mux.HandleFunc("/api/posts/delete", authM.Authenticate(postH.DeletePost))

	// Route Perusahaan
	mux.HandleFunc("/api/perusahaan", authM.Authenticate(utils.AdaptHandler(perusahaanH)))
	mux.HandleFunc("/api/perusahaan/", authM.Authenticate(utils.AdaptHandler(perusahaanH)))

	// Route PIC
	mux.HandleFunc("/api/pic", authM.Authenticate(utils.AdaptHandler(picH)))
	mux.HandleFunc("/api/pic/", authM.Authenticate(utils.AdaptHandler(picH)))

	// Route Perusahaan
	mux.HandleFunc("/api/jabatan", authM.Authenticate(utils.AdaptHandler(jabatanH)))
	mux.HandleFunc("/api/jabatan/", authM.Authenticate(utils.AdaptHandler(jabatanH)))

	// Route Identifikasi
	mux.HandleFunc("/api/identifikasi", authM.Authenticate(utils.AdaptHandler(identifikasiH)))
	mux.HandleFunc("/api/identifikasi/", authM.Authenticate(utils.AdaptHandler(identifikasiH)))
	// Route IKAS
	mux.HandleFunc("/api/ikas", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			ikasH.GetAllIkas(w, r)
		case http.MethodPost:
			ikasH.Create(w, r)
		}
	})
	// Route Gulih
	mux.HandleFunc("/api/gulih", authM.Authenticate(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			gulihH.GetAll(w, r)
		case http.MethodPost:
			gulihH.Create(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc("/api/gulih/", authM.Authenticate(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			gulihH.GetByID(w, r)
		case http.MethodPut:
			gulihH.Update(w, r)
		case http.MethodDelete:
			gulihH.Delete(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc("/api/ikas/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			ikasH.GetIkasByID(w, r)
		case http.MethodPut:
			ikasH.UpdateIkas(w, r)
		case http.MethodDelete:
			ikasH.DeleteIkas(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// Route Proteksi
	mux.HandleFunc("/api/proteksi", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			proteksiH.GetAllProteksi(w, r)
		case http.MethodPost:
			proteksiH.Create(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/proteksi/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			proteksiH.GetProteksiByID(w, r)
		case http.MethodPut:
			proteksiH.UpdateProteksi(w, r)
		case http.MethodDelete:
			proteksiH.DeleteProteksi(w, r)
		}
	})

	return mux
}
