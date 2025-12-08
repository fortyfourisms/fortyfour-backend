package routes

import (
	"fortyfour-backend/internal/handlers"
	"fortyfour-backend/internal/middleware"
	"fortyfour-backend/internal/utils"
	"net/http"
)

func InitRouter(authH *handlers.AuthHandler, postH *handlers.PostHandler, perusahaanH *handlers.PerusahaanHandler, picH *handlers.PICHandler,
	identifikasiH *handlers.IdentifikasiHandler, jabatanH *handlers.JabatanHandler, authM *middleware.AuthMiddleware) *http.ServeMux {
	mux := http.NewServeMux()

	// Routes Auth
	mux.HandleFunc("/api/register", authH.Register)
	mux.HandleFunc("/api/login", authH.Login)

	// Routes Posts
	mux.HandleFunc("/api/posts", authM.Authenticate(postH.GetPosts))
	mux.HandleFunc("/api/posts/single", authM.Authenticate(postH.GetPost))
	mux.HandleFunc("/api/posts/create", authM.Authenticate(postH.CreatePost))
	mux.HandleFunc("/api/posts/update", authM.Authenticate(postH.UpdatePost))
	mux.HandleFunc("/api/posts/delete", authM.Authenticate(postH.DeletePost))

	// Route Perusahaan
	mux.HandleFunc("/api/perusahaan", authM.Authenticate(utils.AdaptHandler(perusahaanH)))
	mux.HandleFunc("/api/pic", authM.Authenticate(utils.AdaptHandler(picH)))

	// Route Perusahaan
	mux.HandleFunc("/api/jabatan", authM.Authenticate(utils.AdaptHandler(jabatanH)))
	mux.HandleFunc("/api/jabatan/", authM.Authenticate(utils.AdaptHandler(jabatanH)))

	// Route Identifikasi
	mux.HandleFunc("/api/identifikasi", authM.Authenticate(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			identifikasiH.GetAll(w, r)
		case http.MethodPost:
			identifikasiH.Create(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc("/api/identifikasi/", authM.Authenticate(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			identifikasiH.GetByID(w, r)
		case http.MethodPut:
			identifikasiH.Update(w, r)
		case http.MethodDelete:
			identifikasiH.Delete(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))

	return mux
}
