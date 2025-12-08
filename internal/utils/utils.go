package utils

import "net/http"

func ValueOrEmpty(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func AdaptHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	}
}
