package routes

import (
	"encoding/json"
	"net/http"
	"time"

	"survey/internal/handlers"
	"survey/internal/middleware"
)

// HEALTH CHECK
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// INIT ROUTER
func InitRouter(
	respondenH *handlers.RespondenHandler,
	risikoH *handlers.RisikoHandler,
	authMiddleware func(http.Handler) http.Handler,
) *http.ServeMux {

	mux := http.NewServeMux()

	// PUBLIC
	mux.HandleFunc("/api/health", healthHandler)

	// MIDDLEWARE WRAPPER
	protected := func(h http.Handler) http.Handler {
		return middleware.Logger(
			middleware.Recovery(
				authMiddleware(h),
			),
		)
	}

	protectedFunc := func(h http.HandlerFunc) http.Handler {
		return protected(http.HandlerFunc(h))
	}

	// RESPONDEN 
	mux.Handle("/api/survey/responden", protected(respondenH))
	mux.Handle("/api/survey/responden/", protected(respondenH))

	// RISIKO 
	mux.Handle("/api/survey/risiko/eligibility", protectedFunc(risikoH.SubmitEligibility))
	mux.Handle("/api/survey/risiko/dampak", protectedFunc(risikoH.SubmitDampak))
	mux.Handle("/api/survey/risiko/pengendalian", protectedFunc(risikoH.SubmitPengendalian))
	mux.Handle("/api/survey/risiko/reason", protectedFunc(risikoH.SubmitAlasan))
	mux.Handle("/api/survey/risiko/", protectedFunc(risikoH.GetByRespondentID))

	// PROGRESS & NAVIGATION
	mux.Handle("/api/survey/navigate", protectedFunc(risikoH.Navigate))
	mux.Handle("/api/survey/save-progress", protectedFunc(risikoH.SaveProgress))
	mux.Handle("/api/survey/finish", protectedFunc(risikoH.FinishSurvey))

	return mux
}