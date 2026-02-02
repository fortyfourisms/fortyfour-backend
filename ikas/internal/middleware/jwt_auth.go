package middleware

import (
	"context"
	"ikas/internal/utils"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	jwtSecret string
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

func NewAuthMiddleware(jwtSecret string) *AuthMiddleware {
	return &AuthMiddleware{jwtSecret: jwtSecret}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.RespondError(w, http.StatusUnauthorized, "Missing authorization header")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.RespondError(w, http.StatusUnauthorized, "Invalid authorization format")
			return
		}

		claims, err := utils.VerifyToken(parts[1], m.jwtSecret)
		if err != nil {
			utils.RespondError(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		// inject context
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, Username, claims.Username)
		ctx = context.WithValue(ctx, Role, claims.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
