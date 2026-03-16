package middleware

import (
	"context"
	"ikas/internal/utils"
	"net/http"
)

type AuthMiddleware struct {
	jwtSecret          string
	internalGatewayKey string
}

type contextKey struct {
	name string
}

// context keys (sama konsep, beda package aman)
var (
	UserIDKey = &contextKey{"user-id"}
	Username  = &contextKey{"username"}
	Role      = &contextKey{"role"}
)

func NewAuthMiddleware(jwtSecret string, internalKey string) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret:          jwtSecret,
		internalGatewayKey: internalKey,
	}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for internal gateway key first
		internalKey := r.Header.Get("X-Internal-Key")
		if internalKey != m.internalGatewayKey {
			utils.RespondError(w, http.StatusUnauthorized, "Unauthorized: Direct access not allowed or invalid internal key")
			return
		}

		// Extract user info from headers injected by Gateway
		userID := r.Header.Get("X-User-ID")
		username := r.Header.Get("X-Username") // Optional, falls back to empty if not set
		role := r.Header.Get("X-User-Role")

		if userID == "" || role == "" {
			utils.RespondError(w, http.StatusUnauthorized, "Unauthorized: User identification missing in gateway headers")
			return
		}

		// Inject into context for downstream handlers/services
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		ctx = context.WithValue(ctx, Username, username)
		ctx = context.WithValue(ctx, Role, role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
