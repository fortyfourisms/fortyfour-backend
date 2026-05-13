package utils

import (
	"net/http"
	"regexp"
	"strings"
)

func Slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	// Remove other special characters for safety
	reg := regexp.MustCompile("[^a-z0-9-]+")
	s = reg.ReplaceAllString(s, "")
	// Remove multiple consecutive hyphens
	reg = regexp.MustCompile("-+")
	s = reg.ReplaceAllString(s, "-")
	// Trim hyphens from ends
	s = strings.Trim(s, "-")
	return s
}

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
