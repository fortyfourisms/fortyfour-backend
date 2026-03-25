package middleware

import (
	"ikas/internal/utils"
	"log"
	"net/http"

	"github.com/casbin/casbin/v3"
)

type CasbinMiddleware struct {
	enforcer *casbin.Enforcer
}

func NewCasbinMiddleware(enforcer *casbin.Enforcer) *CasbinMiddleware {
	return &CasbinMiddleware{
		enforcer: enforcer,
	}
}

// Authorize checks if user role has permission to access resource
func (m *CasbinMiddleware) Authorize(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get role from header (set by AuthMiddleware)
		role := r.Header.Get("X-User-Role")
		if role == "" {
			utils.RespondError(w, http.StatusUnauthorized, "Role not found in token")
			return
		}

		// Get resource and action
		resource := r.URL.Path
		action := r.Method

		// Check permission with Casbin
		allowed, err := m.enforcer.Enforce(role, resource, action)
		if err != nil {
			log.Printf("Casbin enforce error: %v", err)
			utils.RespondError(w, http.StatusInternalServerError, "Authorization check failed")
			return
		}

		if !allowed {
			log.Printf("Access denied: role=%s, resource=%s, action=%s", role, resource, action)
			utils.RespondError(w, http.StatusForbidden, "You don't have permission to access this resource")
			return
		}

		log.Printf("Access granted: role=%s, resource=%s, action=%s", role, resource, action)
		next(w, r)
	}
}
