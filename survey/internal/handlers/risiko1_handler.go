package handlers

import (
	"encoding/json"
	"net/http"

	"survey/internal/dto"
	"survey/internal/services"
)

type RisikoHandler struct {
	service *services.RisikoService
}

func NewRisikoHandler(s *services.RisikoService) *RisikoHandler {
	return &RisikoHandler{service: s}
}

func (h *RisikoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		h.create(w, r)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *RisikoHandler) create(w http.ResponseWriter, r *http.Request) {

	var req dto.CreateRisikoRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	err = h.service.Create(req)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusCreated)
}