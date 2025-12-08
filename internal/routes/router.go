package routes

import (
	"net/http"

	"fortyfour-backend/internal/handlers"
	"fortyfour-backend/internal/middleware"
)

func InitRouter(authH *handlers.AuthHandler, postH *handlers.PostHandler, perusahaanH *handlers.PerusahaanHandler, ikasH *handlers.IkasHandler, proteksiH *handlers.ProteksiHandler, authM *middleware.AuthMiddleware) *http.ServeMux {
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

	// Route IKAS
	mux.HandleFunc("/api/ikas", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			ikasH.GetAllIkas(w, r) // Read all
		case http.MethodPost:
			ikasH.Create(w, r) // Create
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/ikas/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			ikasH.GetIkasByID(w, r) // Read by ID
		case http.MethodPut:
			ikasH.UpdateIkas(w, r) // Update
		case http.MethodDelete:
			ikasH.DeleteIkas(w, r) // Delete
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

		mux.HandleFunc("/api/protelsi/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			proteksiH.GetProteksiByID(w, r) // Read by ID
		case http.MethodPut:
			proteksiH.UpdateProteksi(w, r) // Update
		case http.MethodDelete:
			proteksiH.DeleteProteksi(w, r) // Delete
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	return mux
}
