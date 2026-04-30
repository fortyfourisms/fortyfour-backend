package middleware

import (
	"context"
	"log"
	"net/http"
	"time"
)

// BASIC MIDDLEWARE
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

// CONTEXT HANDLER
type contextKey string

const (
	userIDKey contextKey = "user_id"
	roleKey   contextKey = "role"
)

// SAFE: tidak akan panic walau ctx nil
func ensureContext(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	return ctx
}

// SET
func SetUserID(ctx context.Context, userID string) context.Context {
	ctx = ensureContext(ctx)
	return context.WithValue(ctx, userIDKey, userID)
}

func SetRole(ctx context.Context, role string) context.Context {
	ctx = ensureContext(ctx)
	return context.WithValue(ctx, roleKey, role)
}

// GET
func GetUserID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if val, ok := ctx.Value(userIDKey).(string); ok {
		return val
	}
	return ""
}

func GetRole(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if val, ok := ctx.Value(roleKey).(string); ok {
		return val
	}
	return ""
}

// AUTH MIDDLEWARE
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userID := r.Header.Get("X-User-ID")
		role := r.Header.Get("X-Role")

		ctx := ensureContext(r.Context())
		ctx = SetUserID(ctx, userID)
		ctx = SetRole(ctx, role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
