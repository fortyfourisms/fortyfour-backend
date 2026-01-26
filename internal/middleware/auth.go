// internal/middleware/auth.go
package middleware

import (
	"fortyfour-backend/internal/utils"
	"net/http"
	"strings"

	"github.com/rollbar/rollbar-go"
)

type AuthMiddleware struct {
	jwtSecret string
}

type contextKey struct {
	name string
}

// Struct pointer as key
// Avoiding collisions in Go context keys
var (
	UserIDKey = &contextKey{"user-id"}
	Username  = &contextKey{"username"}
	Role      = &contextKey{"role"}
)

func NewAuthMiddleware(jwtSecret string) *AuthMiddleware {
	return &AuthMiddleware{jwtSecret: jwtSecret}
}

func (m *AuthMiddleware) Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
			rollbar.Error(err)
			utils.RespondError(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		UserIDString := claims.UserID
		r.Header.Set("X-User-ID", UserIDString)
		r.Header.Set("X-Username", claims.Username)
		r.Header.Set("X-User-Role", claims.Role)

		// Set ke context juga
		ctx := SetUserID(r.Context(), UserIDString)
		ctx = SetUsername(ctx, claims.Username)
		ctx = SetRole(ctx, claims.Role)

		next(w, r.WithContext(ctx))
	}
}
