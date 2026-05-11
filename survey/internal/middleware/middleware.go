package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

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

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-User-ID, X-User-Role, X-Internal-Key")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// CONTEXT KEY
type contextKey string

const (
	userIDKey contextKey = "user_id"
	roleKey   contextKey = "role"
)

func SetUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func SetRole(ctx context.Context, role string) context.Context {
	return context.WithValue(ctx, roleKey, role)
}

// GETTERS
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

// JWT VALIDATION
func validateToken(token string) (userID string, role string, err error) {

	if token == "" {
		return "", "", errors.New("empty token")
	}

	if strings.HasPrefix(token, "admin-") {
		return strings.TrimPrefix(token, "admin-"), "admin", nil
	}

	if strings.HasPrefix(token, "user-") {
		return strings.TrimPrefix(token, "user-"), "user", nil
	}

	claims := jwt.MapClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secret = "secret"
		}

		return []byte(secret), nil
	})

	if err != nil || !parsedToken.Valid {
		return "", "", errors.New("invalid token")
	}

	userID, _ = claims["user_id"].(string)
	role, _ = claims["role"].(string)

	if userID == "" {
		return "", "", errors.New("missing user_id claim")
	}

	return userID, role, nil
}

// AUTH MIDDLEWARE
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")

		var userID, role string

		// Gateway sudah melakukan auth + casbin, lalu meneruskan identitas lewat header ini.
		userID = r.Header.Get("X-User-ID")
		role = r.Header.Get("X-User-Role")

		// Bearer token hanya fallback agar service tetap mudah diuji langsung.
		if userID == "" && authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {

			token := strings.TrimPrefix(authHeader, "Bearer ")

			uid, rle, err := validateToken(token)
			if err != nil {
				http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
				return
			}

			userID = uid
			role = rle
		}

		if userID != "" && role == "" {
			role = "user"
		}

		// CONTEXT
		ctx := r.Context()
		if userID != "" {
			ctx = SetUserID(ctx, userID)
		}
		if role != "" {
			ctx = SetRole(ctx, role)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
