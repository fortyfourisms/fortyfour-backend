package handlers

import (
	"encoding/json"
	"net/http"

	"survey/internal/dto"
	"survey/internal/services"
	"survey/internal/utils"
)

type RisikoHandler struct {
	service *services.RisikoService
}

func NewRisikoHandler(s *services.RisikoService) *RisikoHandler {
	return &RisikoHandler{service: s}
}

func (h *RisikoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case http.MethodPost:
		h.createJawaban(w, r)

	default:
		utils.RespondError(w, 405, "Method tidak diizinkan")
	}
}

func (h *RisikoHandler) createJawaban(w http.ResponseWriter, r *http.Request) {

	var req dto.CreateRisikoJawabanRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, 400, "Invalid body")
		return
	}

	if err := h.service.CreateJawaban(req); err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	}

	utils.RespondJSON(w, 201, map[string]string{
		"message": "Jawaban risiko berhasil disimpan",
	})
}