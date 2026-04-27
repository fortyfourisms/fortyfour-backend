package middleware

import (
	"context"
	"log"
	"net/http"
	"time"
)

// EXISTING MIDDLEWARE (TIDAK DIUBAH)
func AdaptHandler(h http.Handler) http.Handler {
	return h
}

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-User-ID, X-Role")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// TAMBAHAN: CONTEXT HANDLER
type contextKey string

const (
	userIDKey contextKey = "user_id"
	roleKey   contextKey = "role"
)

// SET ke context
func SetUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func SetRole(ctx context.Context, role string) context.Context {
	return context.WithValue(ctx, roleKey, role)
}

// GET dari context (dipakai handler)
func GetUserID(ctx context.Context) string {
	if val, ok := ctx.Value(userIDKey).(string); ok {
		return val
	}
	return ""
}

func GetRole(ctx context.Context) string {
	if val, ok := ctx.Value(roleKey).(string); ok {
		return val
	}
	return ""
}

// TAMBAHAN: AUTH MIDDLEWARE (SIMPLE)
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Ambil dari header (sementara, bisa diganti JWT nanti)
		userID := r.Header.Get("X-User-ID")
		role := r.Header.Get("X-Role")

		// inject ke context
		ctx := r.Context()
		ctx = SetUserID(ctx, userID)
		ctx = SetRole(ctx, role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}