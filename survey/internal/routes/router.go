package routes

import (
	"encoding/json"
	"net/http"
	"time"

	_ "survey/docs"
	"survey/internal/handlers"
	"survey/internal/middleware"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @Summary		Health check
// @Description	Check if the Survey API is running and healthy
// @Tags			Health
// @Produce		json
// @Success		200	{object}	map[string]string
// @Router			/api/health [get]
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

	// RESPONDEN
	mux.Handle("/api/survey/responden", protected(respondenH))
	mux.Handle("/api/survey/responden/", protected(respondenH))

	// RISIKO
	mux.Handle("/api/survey/risiko/eligibility", protected(http.HandlerFunc(risikoH.SubmitEligibility)))
	mux.Handle("/api/survey/risiko/dampak", protected(http.HandlerFunc(risikoH.SubmitDampak)))
	mux.Handle("/api/survey/risiko/pengendalian", protected(http.HandlerFunc(risikoH.SubmitPengendalian)))
	mux.Handle("/api/survey/risiko/reason", protected(http.HandlerFunc(risikoH.SubmitAlasan)))
	mux.Handle("/api/survey/risiko/", protected(http.HandlerFunc(risikoH.GetByRespondentID)))

	// PROGRESS & NAVIGATION
	mux.Handle("/api/survey/navigate", protected(http.HandlerFunc(risikoH.Navigate)))
	mux.Handle("/api/survey/save-progress", protected(http.HandlerFunc(risikoH.SaveProgress)))
	mux.Handle("/api/survey/finish", protected(http.HandlerFunc(risikoH.FinishSurvey)))

	// Swagger UI
	mux.HandleFunc("/swagger/survey/", httpSwagger.WrapHandler)

	return mux
}
