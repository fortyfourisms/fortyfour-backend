package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"
	"survey/internal/utils"
	"time"
)

type AuthMiddleware struct {
	internalGatewayKey string
}

type contextKey struct {
	name string
}

var (
	UserIDKey       = &contextKey{"user-id"}
	UsernameKey     = &contextKey{"username"}
	RoleKey         = &contextKey{"role"}
	PerusahaanIDKey = &contextKey{"perusahaan-id"}
)

func NewAuthMiddleware(internalKey string) *AuthMiddleware {
	return &AuthMiddleware{
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
		userID := strings.TrimSpace(r.Header.Get("X-User-ID"))
		username := strings.TrimSpace(r.Header.Get("X-Username"))
		role := strings.ToLower(strings.TrimSpace(r.Header.Get("X-User-Role")))
		perusahaanID := strings.TrimSpace(r.Header.Get("X-Perusahaan-ID"))

		if userID == "" || role == "" {
			utils.RespondError(w, http.StatusUnauthorized, "Unauthorized: User identification missing in gateway headers")
			return
		}

		// Inject into context for downstream handlers/services
		ctx := r.Context()
		ctx = context.WithValue(ctx, UserIDKey, userID)
		ctx = context.WithValue(ctx, UsernameKey, username)
		ctx = context.WithValue(ctx, RoleKey, role)
		ctx = context.WithValue(ctx, PerusahaanIDKey, perusahaanID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GETTERS
func GetUserID(ctx context.Context) string {
	if val, ok := ctx.Value(UserIDKey).(string); ok {
		return val
	}
	return ""
}

func GetUsername(ctx context.Context) string {
	if val, ok := ctx.Value(UsernameKey).(string); ok {
		return val
	}
	return ""
}

func GetRole(ctx context.Context) string {
	if val, ok := ctx.Value(RoleKey).(string); ok {
		return val
	}
	return ""
}

func GetPerusahaanID(ctx context.Context) string {
	if val, ok := ctx.Value(PerusahaanIDKey).(string); ok {
		return val
	}
	return ""
}

// HasRole checks if the user has at least one of the required roles.
// It handles comma-separated roles from the gateway.
func HasRole(ctx context.Context, requiredRoles ...string) bool {
	userRole := GetRole(ctx)
	if userRole == "" {
		return false
	}

	// Split by comma in case multiple roles are passed
	roles := strings.Split(userRole, ",")
	for _, r := range roles {
		r = strings.TrimSpace(strings.ToLower(r))
		for _, req := range requiredRoles {
			if r == strings.ToLower(req) {
				return true
			}
		}
	}
	return false
}

// BASIC MIDDLEWARE
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("PANIC:", err)
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
		log.Printf("[REQUEST] %s %s %v",
			r.Method,
			r.URL.Path,
			time.Since(start),
		)
	})
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-User-ID, X-User-Role, X-Perusahaan-ID, X-Internal-Key, X-Requested-With")
		w.Header().Set("Access-Control-Max-Age", "3600")

		// Handle preflight
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
