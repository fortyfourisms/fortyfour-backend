package routes

import (
	"net/http"

	"fortyfour-backend/internal/handlers"
	"fortyfour-backend/internal/middleware"
)

func InitRouter(authH *handlers.AuthHandler, postH *handlers.PostHandler, perusahaanH *handlers.PerusahaanHandler, authM *middleware.AuthMiddleware) *http.ServeMux {
	mux := http.NewServeMux()

	// Public routes (Auth)
	mux.HandleFunc("/api/register", authH.Register)
	mux.HandleFunc("/api/login", authH.Login)

	// Public routes (Posts)
	mux.HandleFunc("/api/posts", postH.GetPosts)
	mux.HandleFunc("/api/posts/single", postH.GetPost)

	mux.HandleFunc("/api/perusahaan", perusahaanH.GetAll)
	mux.HandleFunc("/api/perusahaan/single", perusahaanH.GetByID)

	// Protected routes (Posts CRUD)
	mux.HandleFunc("/api/posts/create", authM.Authenticate(postH.CreatePost))
	mux.HandleFunc("/api/posts/update", authM.Authenticate(postH.UpdatePost))
	mux.HandleFunc("/api/posts/delete", authM.Authenticate(postH.DeletePost))

	mux.HandleFunc("/api/perusahaan/create", authM.Authenticate(perusahaanH.Create))
	mux.HandleFunc("/api/perusahaan/update", authM.Authenticate(perusahaanH.Update))
	mux.HandleFunc("/api/perusahaan/delete", authM.Authenticate(perusahaanH.Delete))

	return mux
}
