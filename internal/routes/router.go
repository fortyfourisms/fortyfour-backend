package routes

import (
	"net/http"

	"fortyfour-backend/internal/handlers"
	"fortyfour-backend/internal/middleware"
)

func InitRouter(authH *handlers.AuthHandler, postH *handlers.PostHandler, perusahaanH *handlers.PerusahaanHandler, picH *handlers.PICPerusahaanHandler, authM *middleware.AuthMiddleware) *http.ServeMux {
	mux := http.NewServeMux()

	// Routes Auth
	mux.HandleFunc("/api/register", authH.Register)
	mux.HandleFunc("/api/login", authH.Login)

	// Routes Posts
	mux.HandleFunc("/api/posts", postH.GetPosts)
	mux.HandleFunc("/api/posts/single", postH.GetPost)

	mux.HandleFunc("/api/posts/create", authM.Authenticate(postH.CreatePost))
	mux.HandleFunc("/api/posts/update", authM.Authenticate(postH.UpdatePost))
	mux.HandleFunc("/api/posts/delete", authM.Authenticate(postH.DeletePost))

	// Route Perusahaan
	mux.HandleFunc("/api/perusahaan", perusahaanH.GetAll)
	mux.HandleFunc("/api/perusahaan/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			perusahaanH.GetByID(w, r)
		case http.MethodPut:
			perusahaanH.Update(w, r)
		case http.MethodDelete:
			perusahaanH.Delete(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/api/perusahaan/create", perusahaanH.Create)

	// Routes PIC Perusahaan
	mux.HandleFunc("/api/pic", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			picH.GetAll(w, r)
		case http.MethodPost:
			picH.Create(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/api/pic/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			picH.GetByID(w, r)
		case http.MethodPut:
			picH.Update(w, r)
		case http.MethodDelete:
			picH.Delete(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	return mux
}
