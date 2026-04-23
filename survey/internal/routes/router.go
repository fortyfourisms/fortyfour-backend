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
	mux.Handle("/api/survey/responden", middleware.Logger(utils.AdaptHandler(respondenH)))
	mux.Handle("/api/survey/responden/", middleware.Logger(utils.AdaptHandler(respondenH)))

	// RISIKO (Intellectual Property Theft Survey)
	mux.HandleFunc("/api/survey/risiko/eligibility", risikoH.SubmitEligibility)
	mux.HandleFunc("POST /api/survey/risiko/dampak", risikoH.SubmitDampak)
	mux.HandleFunc("POST /api/survey/risiko/pengendalian", risikoH.SubmitPengendalian)
	mux.HandleFunc("POST /api/survey/risiko/reason", risikoH.SubmitAlasan)
	mux.HandleFunc("/api/survey/risiko/", risikoH.GetByRespondentID)

	// Progress & navigation
	mux.HandleFunc("/api/survey/progress/", risikoH.GetProgress)
	mux.HandleFunc("/api/survey/navigate", risikoH.Navigate)

	// Lanjutkan Nanti 
	mux.HandleFunc("/api/survey/save-progress", risikoH.SaveProgress)

	// Custom Risiko
	mux.HandleFunc("/api/survey/custom-risk", risikoH.CreateCustomRisiko)
	mux.HandleFunc("/api/survey/finish", risikoH.FinishSurvey)

	return mux
}
