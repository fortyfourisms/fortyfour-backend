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

// @Summary		List SE edit requests
// @Description	Admin melihat semua request pending. User melihat request miliknya.
// @Tags			SE - Edit Request
// @Produce		json
// @Security		BearerAuth
// @Success		200	{array}		dto.SEEditRequestResponse
// @Failure		500	{object}	dto.ErrorResponse
// @Router			/api/se/edit-requests [get]
func (h *SEEditRequestHandler) handleList(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())

	if role == "admin" || role == "staff" {
		data, err := h.service.GetPending()
		if err != nil {
			utils.RespondError(w, 500, err.Error())
			return
		}
		utils.RespondJSON(w, 200, data)
		return
	}

	// User: list own requests
	userID := middleware.GetUserID(r.Context())
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

// @Summary		Review SE edit request
// @Description	Admin approve atau reject request edit SE. Approve akan otomatis update data SE.
// @Tags			SE - Edit Request
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			id		path		string						true	"ID Edit Request"
// @Param			request	body		dto.ReviewSEEditRequestDTO	true	"Status dan catatan admin"
// @Success		200		{object}	dto.SEEditRequestResponse
// @Failure		400		{object}	dto.ErrorResponse
// @Failure		403		{object}	dto.ErrorResponse
// @Router			/api/se/edit-requests/{id}/review [put]
func (h *SEEditRequestHandler) handleReview(w http.ResponseWriter, r *http.Request, id string) {
	role := middleware.GetRole(r.Context())
	if role != "admin" && role != "staff" {
		utils.RespondError(w, 403, "Hanya admin dan staff yang dapat mereview request")
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
//
//	@Summary		Submit request edit SE
//	@Description	User mengajukan permintaan perubahan data SE beserta alasan.
//	@Tags			SE - Edit Request
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string						true	"ID SE"
//	@Param			request	body		dto.CreateSEEditRequestDTO	true	"Data perubahan dan alasan"
//	@Success		201		{object}	dto.SEEditRequestResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Router			/api/se/{id}/request-edit [post]
func (h *SEEditRequestHandler) HandleRequestEdit(w http.ResponseWriter, r *http.Request, idSE string) {
	userID := middleware.GetUserID(r.Context())
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
