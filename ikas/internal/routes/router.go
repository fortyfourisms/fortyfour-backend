package routes

import (
	"ikas/internal/handlers"
	"net/http"
)

func InitRouter(
	ikasH *handlers.IkasHandler,
) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/api/ikas", ikasH)
	mux.Handle("/api/ikas/", ikasH)

	return mux
}
