package utils

import "net/http"

func ValueOrNull(s *string) interface{} {
	if s != nil {
		return *s
	}
	return nil
}

func ValueOrEmpty(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func IntOrNull(i *int) interface{} {
	if i != nil {
		return *i
	}
	return nil
}

func BoolOrNull(b *bool) interface{} {
	if b != nil {
		return *b
	}
	return nil
}

func AdaptHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	}
}
