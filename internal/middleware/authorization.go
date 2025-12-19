package middleware

import (
	"net/http"
	"strings"

	"github.com/casbin/casbin/v2"
)

type AuthorizationMiddleware struct {
	enforcer *casbin.Enforcer
}

func NewAuthorizationMiddleware(enforcer *casbin.Enforcer) *AuthorizationMiddleware {
	return &AuthorizationMiddleware{
		enforcer: enforcer,
	}
}

func (m *AuthorizationMiddleware) Authorize(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get role from context (set by AuthMiddleware)
		role, ok := r.Context().Value("role").(string)
		if !ok || role == "" {
			http.Error(w, "Forbidden: No role found", http.StatusForbidden)
			return
		}

		// Get resource path and method
		obj := m.extractResourcePath(r.URL.Path)
		act := r.Method

		// Check permission using Casbin
		allowed, err := m.enforcer.Enforce(role, obj, act)
		if err != nil {
			http.Error(w, "Authorization error", http.StatusInternalServerError)
			return
		}

		if !allowed {
			http.Error(w, "Forbidden: Insufficient permissions", http.StatusForbidden)
			return
		}

		next(w, r)
	}
}

// extractResourcePath removes ID from path for matching
// Example: /api/posts/123 -> /api/posts
func (m *AuthorizationMiddleware) extractResourcePath(path string) string {
	parts := strings.Split(strings.TrimPrefix(path, "/"), "/")
	if len(parts) >= 2 {
		return "/" + parts[0] + "/" + parts[1]
	}
	return path
}
