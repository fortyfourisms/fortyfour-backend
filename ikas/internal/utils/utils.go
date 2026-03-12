package utils

import (
	"math"
	"net/http"
	"strconv"
)

func RoundToTwo(f float64) float64 {
	return math.Round(f*100) / 100
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

func StringToInt(s string) (int, error) {
	return strconv.Atoi(s)
}

func ExtractIntID(path, resource string) (int, error) {
	idStr := ExtractID(path, resource)
	if idStr == "" {
		return 0, nil
	}
	return strconv.Atoi(idStr)
}

func AdaptHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	}
}
