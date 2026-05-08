package middleware

import (
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/utils"
	"net/http"
	"strings"
)

type ChallengeMiddleware struct {
	tokenService *services.TokenService
}

func NewChallengeMiddleware(tokenService *services.TokenService) *ChallengeMiddleware {
	return &ChallengeMiddleware{
		tokenService: tokenService,
	}
}

// VerifyChallenge checks if the user has passed the challenge gate
func (m *ChallengeMiddleware) VerifyChallenge(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Public paths that don't require the challenge
		path := r.URL.Path
		if path == "/api/challenge/verify" ||
			path == "/api/health" ||
			strings.HasPrefix(path, "/swagger/") ||
			strings.HasPrefix(path, "/api/public/") {
			next.ServeHTTP(w, r)
			return
		}

		// Check for gate_verified cookie
		if !m.tokenService.GetGateCookie(r) {
			// If missing, return 403 Forbidden with challenge_required flag
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			utils.RespondJSON(w, http.StatusForbidden, map[string]interface{}{
				"challenge_required": true,
				"message":            "Human verification required. Please complete the challenge.",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}
