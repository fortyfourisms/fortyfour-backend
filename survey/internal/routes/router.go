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

func InitRouter(
	respondenH *handlers.RespondenHandler,
	risikoH *handlers.RisikoHandler,
	authM *middleware.AuthMiddleware,
) http.Handler {

	mux := http.NewServeMux()

	// PUBLIC
	mux.HandleFunc("/api/health", healthHandler)

	// WRAPPERS
	protected := func(h http.Handler) http.Handler {
		return authM.Authenticate(h)
	}
	getOnly := func(h http.HandlerFunc) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}
			h(w, r)
		})
	}
	postOnly := func(h http.HandlerFunc) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}
			h(w, r)
		})
	}
	getRisikoByRespondentID := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/survey/risiko/" {
			risikoH.GetAllRisiko(w, r)
			return
		}
		risikoH.GetByRespondentID(w, r)
	}
	stepWithOwnedRead := func(post http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				risikoH.GetMe(w, r)
			case http.MethodPost:
				post(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}
	}

	// RESPONDEN
	mux.Handle("/api/survey/responden", protected(respondenH))
	mux.Handle("/api/survey/responden/", protected(respondenH))

	// RISIKO
	mux.Handle("/api/survey/risiko/eligibility", protected(http.HandlerFunc(stepWithOwnedRead(risikoH.SubmitEligibility))))
	mux.Handle("/api/survey/risiko/dampak", protected(http.HandlerFunc(stepWithOwnedRead(risikoH.SubmitDampak))))
	mux.Handle("/api/survey/risiko/pengendalian", protected(http.HandlerFunc(stepWithOwnedRead(risikoH.SubmitPengendalian))))
	mux.Handle("/api/survey/risiko/reason", protected(http.HandlerFunc(stepWithOwnedRead(risikoH.SubmitAlasan))))
	mux.Handle("/api/survey/risiko", protected(getOnly(risikoH.GetAllRisiko)))
	mux.Handle("/api/survey/risiko/me", protected(getOnly(risikoH.GetMe)))
	mux.Handle("/api/survey/risiko/", protected(getOnly(getRisikoByRespondentID)))

	// PROGRESS & NAVIGATION
	mux.Handle("/api/survey/progress", protected(getOnly(risikoH.GetProgress)))
	mux.Handle("/api/survey/navigate", protected(postOnly(risikoH.Navigate)))
	mux.Handle("/api/survey/save-progress", protected(postOnly(risikoH.SaveProgress)))
	mux.Handle("/api/survey/finish", protected(postOnly(risikoH.FinishSurvey)))
	mux.Handle("/api/survey/request-edit", protected(postOnly(risikoH.RequestEdit)))
	mux.Handle("/api/survey/edit-requests", protected(getOnly(risikoH.GetEditRequests)))
	mux.Handle("/api/survey/edit-requests/me", protected(getOnly(risikoH.GetMyEditRequest)))
	mux.Handle("/api/survey/edit-requests/", protected(postOnly(risikoH.ReviewEditRequest)))

	// Swagger UI
	mux.HandleFunc("/swagger/survey/", httpSwagger.WrapHandler)

	// GLOBAL MIDDLEWARE
	return middleware.Logger(
		middleware.Recovery(
			middleware.CORS(
				mux,
			),
		),
	)
}
