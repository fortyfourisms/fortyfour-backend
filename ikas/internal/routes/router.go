package routes

import (
	"ikas/internal/handlers"
	"net/http"
)

func InitRouter(
	ikasH *handlers.IkasHandler,
	ruangLingkupH *handlers.RuangLingkupHandler,
	domainH *handlers.DomainHandler,
) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/api/ikas", ikasH)
	mux.Handle("/api/ikas/", ikasH)

	mux.Handle("/api/ruang-lingkup", ruangLingkupH)
	mux.Handle("/api/ruang-lingkup/", ruangLingkupH)

	mux.Handle("/api/domain", domainH)
	mux.Handle("/api/domain/", domainH)

	return mux
}
