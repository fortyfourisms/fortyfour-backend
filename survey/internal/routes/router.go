package routes

import (
	"encoding/json"
	"net/http"
	"time"

	"survey/internal/handlers"
	"survey/internal/middleware"
	"survey/internal/utils"
)

// HEALTH
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func InitRouter(
	respondenH *handlers.RespondenHandler,
	risikoH *handlers.RisikoHandler,
	authH *handlers.AuthHandler,
) *http.ServeMux {

	mux := http.NewServeMux()

	// PUBLIC ROUTES
	mux.HandleFunc("/api/health", healthHandler)

	// LOGIN
	mux.HandleFunc("/api/auth/login", authH.Login)

	// Logger → Auth
	protected := func(h http.Handler) http.Handler {
		return middleware.Logger(
			middleware.AuthMiddleware(h),
		)
	}

	protectedFunc := func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			protected(http.HandlerFunc(h)).ServeHTTP(w, r)
		}
	}

	// RESPONDEN
	mux.Handle("/api/survey/responden", protected(utils.AdaptHandler(respondenH)))
	mux.Handle("/api/survey/responden/", protected(utils.AdaptHandler(respondenH)))

	// RISIKO
	mux.HandleFunc("/api/survey/risiko/eligibility", protectedFunc(risikoH.SubmitEligibility))
	mux.HandleFunc("/api/survey/risiko/dampak", protectedFunc(risikoH.SubmitDampak))
	mux.HandleFunc("/api/survey/risiko/pengendalian", protectedFunc(risikoH.SubmitPengendalian))
	mux.HandleFunc("/api/survey/risiko/reason", protectedFunc(risikoH.SubmitAlasan))
	mux.HandleFunc("/api/survey/risiko/", protectedFunc(risikoH.GetByRespondentID))

	// PROGRESS
	mux.HandleFunc("/api/survey/progress/", protectedFunc(risikoH.GetProgress))
	mux.HandleFunc("/api/survey/navigate", protectedFunc(risikoH.Navigate))

	// SAVE PROGRESS
	mux.HandleFunc("/api/survey/save-progress", protectedFunc(risikoH.SaveProgress))

	// CUSTOM RISIKO
	mux.HandleFunc("/api/survey/custom-risk", protectedFunc(risikoH.CreateCustomRisiko))
	mux.HandleFunc("/api/survey/finish", protectedFunc(risikoH.FinishSurvey))

	return mux
}