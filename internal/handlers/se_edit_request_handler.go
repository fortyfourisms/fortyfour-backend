package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/middleware"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/utils"
)

type SEEditRequestHandler struct {
	service services.SEEditRequestService
}

func NewSEEditRequestHandler(service services.SEEditRequestService) *SEEditRequestHandler {
	return &SEEditRequestHandler{service: service}
}

func (h *SEEditRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/se/edit-requests")
	path = strings.TrimPrefix(path, "/")

	switch {
	// GET /api/se/edit-requests — admin list pending, user list own
	case path == "" && r.Method == http.MethodGet:
		h.handleList(w, r)

	// PUT /api/se/edit-requests/{id}/review — admin approve/reject
	case strings.HasSuffix(path, "/review") && r.Method == http.MethodPut:
		id := strings.TrimSuffix(path, "/review")
		h.handleReview(w, r, id)

	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (h *SEEditRequestHandler) handleList(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())

	if role == "admin" {
		data, err := h.service.GetPending()
		if err != nil {
			utils.RespondError(w, 500, err.Error())
			return
		}
		utils.RespondJSON(w, 200, data)
		return
	}

	// User: list own requests
	userID := ""
	if uid := r.Context().Value(middleware.UserIDKey); uid != nil {
		userID = uid.(string)
	}
	if userID == "" {
		utils.RespondError(w, 401, "Unauthorized")
		return
	}

	data, err := h.service.GetByUser(userID)
	if err != nil {
		utils.RespondError(w, 500, err.Error())
		return
	}
	utils.RespondJSON(w, 200, data)
}

func (h *SEEditRequestHandler) handleReview(w http.ResponseWriter, r *http.Request, id string) {
	role := middleware.GetRole(r.Context())
	if role != "admin" {
		utils.RespondError(w, 403, "Hanya admin yang dapat mereview request")
		return
	}

	var req dto.ReviewSEEditRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, 400, "Invalid request body")
		return
	}

	resp, err := h.service.Review(id, req)
	if err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	}

	utils.RespondJSON(w, 200, resp)
}

// HandleRequestEdit handles POST /api/se/{id}/request-edit
func (h *SEEditRequestHandler) HandleRequestEdit(w http.ResponseWriter, r *http.Request, idSE string) {
	userID := ""
	if uid := r.Context().Value(middleware.UserIDKey); uid != nil {
		userID = uid.(string)
	}
	if userID == "" {
		utils.RespondError(w, 401, "Unauthorized")
		return
	}

	var req dto.CreateSEEditRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, 400, "Invalid request body")
		return
	}

	resp, err := h.service.CreateRequest(userID, idSE, req)
	if err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	}

	utils.RespondJSON(w, 201, resp)
}
