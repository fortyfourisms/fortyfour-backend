package middleware

import (
	"context"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/utils"
	"net/http"
)

type contextKey string

const (
	UserIDKey       contextKey = "user_id"
	UsernameKey     contextKey = "username"
	RoleKey         contextKey = "role"
	IDPerusahaanKey contextKey = "id_perusahaan"
)

type AuthMiddleware struct {
	tokenService *services.TokenService
}

func NewAuthMiddleware(tokenService *services.TokenService) *AuthMiddleware {
	return &AuthMiddleware{
		tokenService: tokenService,
	}
}

// Authenticate validates the access token from cookie and sets user context
func (m *AuthMiddleware) Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get access token from cookie
		accessToken, err := m.tokenService.GetAccessTokenFromCookie(r)
		if err != nil {
			utils.RespondError(w, http.StatusUnauthorized, "Authentication required")
			return
		}

		// Validate access token
		claims, err := utils.ValidateAccessToken(accessToken, m.tokenService.JWTSecret)
		if err != nil {
			utils.RespondError(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		// Set user info in context
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, UsernameKey, claims.Username)
		ctx = context.WithValue(ctx, RoleKey, claims.Role)
		ctx = context.WithValue(ctx, IDPerusahaanKey, claims.IDPerusahaan)

		// Set role in header for Casbin middleware compatibility
		r.Header.Set("X-User-Role", claims.Role)

		next(w, r.WithContext(ctx))
	}
}

// OptionalAuth is similar to Authenticate but doesn't fail if token is missing
func (m *AuthMiddleware) OptionalAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accessToken, err := m.tokenService.GetAccessTokenFromCookie(r)
		if err == nil {
			claims, err := utils.ValidateAccessToken(accessToken, m.tokenService.JWTSecret)
			if err == nil {
				ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
				ctx = context.WithValue(ctx, UsernameKey, claims.Username)
				ctx = context.WithValue(ctx, RoleKey, claims.Role)
				r.Header.Set("X-User-Role", claims.Role)
				r = r.WithContext(ctx)
			}
		}

		next(w, r)
	}
}

// AutoRefresh automatically refreshes the access token if expired
func (m *AuthMiddleware) AutoRefresh(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := m.tokenService.ValidateAndRefreshIfNeeded(w, r)
		if err != nil {
			utils.RespondError(w, http.StatusUnauthorized, "Authentication required")
			return
		}

		// Set user info in context
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, UsernameKey, claims.Username)
		ctx = context.WithValue(ctx, RoleKey, claims.Role)

		// Set role in header for Casbin middleware compatibility
		r.Header.Set("X-User-Role", claims.Role)

		next(w, r.WithContext(ctx))
	}
}
