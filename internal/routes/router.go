package routes

import (
	"net/http"

	"fortyfour-backend/internal/handlers"
	"fortyfour-backend/internal/middleware"
)

func InitRouter(authH *handlers.AuthHandler, postH *handlers.PostHandler, perusahaanH *handlers.PerusahaanHandler, picH *handlers.PICPerusahaanHandler,
	identifikasiH *handlers.IdentifikasiHandler, proteksiH *handlers.ProteksiHandler, deteksiH *handlers.DeteksiHandler, gulihH *handlers.GulihHandler, ikasH *handlers.IkasHandler, authM *middleware.AuthMiddleware) *http.ServeMux {
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

	// Route Identifikasi
	mux.HandleFunc("/api/identifikasi", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			identifikasiH.GetAll(w, r) // Read all
		case http.MethodPost:
			identifikasiH.Create(w, r) // Create
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/identifikasi/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			identifikasiH.GetByID(w, r) // Read by ID
		case http.MethodPut:
			identifikasiH.Update(w, r) // Update
		case http.MethodDelete:
			identifikasiH.Delete(w, r) // Delete
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
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// Route Deteksi
	mux.HandleFunc("/api/deteksi", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			deteksiH.GetAll(w, r) // Read all
		case http.MethodPost:
			deteksiH.Create(w, r) // Create
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/deteksi/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			deteksiH.GetByID(w, r) // Read by ID
		case http.MethodPut:
			deteksiH.Update(w, r) // Update
		case http.MethodDelete:
			deteksiH.Delete(w, r) // Delete
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// Route Gulih
	mux.HandleFunc("/api/gulih", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			gulihH.GetAll(w, r) // Read all
		case http.MethodPost:
			gulihH.Create(w, r) // Create
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/gulih/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			gulihH.GetByID(w, r) // Read by ID
		case http.MethodPut:
			gulihH.Update(w, r) // Update
		case http.MethodDelete:
			gulihH.Delete(w, r) // Delete
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// Route IKAS
	mux.HandleFunc("/api/ikas", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			ikasH.GetAllIkas(w, r) 
		case http.MethodPost:
			ikasH.Create(w, r) 
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

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

	return mux
}