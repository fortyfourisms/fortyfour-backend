package routes

import (
	"encoding/json"
	"net/http"
	"time"

	"survey/internal/handlers"
	"survey/internal/middleware"
	"survey/internal/utils"
)

// Health handler
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
) *http.ServeMux {

	mux := http.NewServeMux()

	// Health
	mux.HandleFunc("/api/health", healthHandler)

	// RESPONDEN ROUTES
	mux.Handle("/api/responden", middleware.Logger(utils.AdaptHandler(respondenH)))
	mux.Handle("/api/responden/", middleware.Logger(utils.AdaptHandler(respondenH)))

	// RISIKO (Intellectual Property Theft Survey)
	mux.HandleFunc("/api/survey/risk/ip-theft/eligibility", risikoH.SubmitEligibility)
	mux.HandleFunc("/api/survey/risk/ip-theft/detail", risikoH.SubmitDetail)
	mux.HandleFunc("/api/survey/risk/ip-theft/control", risikoH.SubmitControl)
	mux.HandleFunc("/api/survey/risk/ip-theft/reason", risikoH.SubmitReason)
	mux.HandleFunc("/api/survey/risk/ip-theft/", risikoH.GetByRespondentID)

	// Progress & navigation
	mux.HandleFunc("/api/survey/progress/", risikoH.GetProgress)
	mux.HandleFunc("/api/survey/navigate", risikoH.Navigate)

	return mux
}
