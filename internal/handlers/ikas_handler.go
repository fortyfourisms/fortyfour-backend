package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/services"

	"github.com/gorilla/mux"
)

type IkasHandler struct {
	service services.IkasService
}

func NewIkasHandler(service services.IkasService) *IkasHandler {
	return &IkasHandler{service: service}
}

func (h *IkasHandler) CreateIkas(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var ikas models.Ikas
	if err := json.NewDecoder(r.Body).Decode(&ikas); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.CreateIkas(&ikas); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Data berhasil ditambahkan",
		"data":    ikas,
	})
}

func (h *IkasHandler) GetAllIkas(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ikasList, err := h.service.GetAllIkas()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(ikasList)
}

func (h *IkasHandler) GetIkasByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	ikas, err := h.service.GetIkasByID(id)
	if err != nil {
		http.Error(w, "Data tidak ditemukan", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(ikas)
}

func (h *IkasHandler) UpdateIkas(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var ikas models.Ikas
	if err := json.NewDecoder(r.Body).Decode(&ikas); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateIkas(id, &ikas); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ikas.ID = id
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Data berhasil diupdate",
		"data":    ikas,
	})
}

func (h *IkasHandler) DeleteIkas(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	if err := h.service.DeleteIkas(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Data berhasil dihapus",
	})
}
