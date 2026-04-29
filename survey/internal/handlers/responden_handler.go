package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"survey/internal/dto"
	"survey/internal/services"
	"survey/internal/utils"
)

type RespondenHandler struct {
	service *services.RespondenService
}

func NewRespondenHandler(service *services.RespondenService) *RespondenHandler {
	return &RespondenHandler{service: service}
}

func (h *RespondenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	path := strings.TrimPrefix(r.URL.Path, "/api/survey/responden")
	id := strings.TrimPrefix(path, "/")

	switch r.Method {

	case http.MethodGet:
		if id == "" {
			h.handleGetAll(w)
		} else {
			h.handleGetByID(w, id)
		}

	case http.MethodPost:
		if id != "" {
			utils.RespondError(w, http.StatusBadRequest, "ID tidak diperlukan untuk create")
			return
		}
		h.handleCreate(w, r)

	case http.MethodPut:
		if id == "" {
			utils.RespondError(w, http.StatusBadRequest, "ID wajib")
			return
		}
		h.handleUpdate(w, r, id)

	default:
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method tidak diizinkan")
	}
}

// GET ALL
// GetAllResponden godoc
// @Summary      Ambil semua responden
// @Description  Mengambil seluruh data responden (join users, jabatan, perusahaan)
// @Tags         Responden Survey
// @Produce      json
// @Success      200 {array} dto.RespondenResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /api/survey/responden [get]
func (h *RespondenHandler) handleGetAll(w http.ResponseWriter) {

	data, err := h.service.GetAll()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, data)
}

// GET BY ID
// GetRespondenByID godoc
// @Summary      Ambil responden berdasarkan ID
// @Description  Mengambil detail responden beserta data user, jabatan, dan perusahaan
// @Tags         Responden Survey
// @Produce      json
// @Param        id path int true "Responden ID"
// @Success      200 {object} dto.RespondenResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      404 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /api/survey/responden/{id} [get]
func (h *RespondenHandler) handleGetByID(w http.ResponseWriter, id string) {

	idInt, err := strconv.Atoi(id)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "ID harus berupa angka")
		return
	}

	data, err := h.service.GetByID(idInt)
	if err != nil {

		if err.Error() == "data tidak ditemukan" {
			utils.RespondError(w, http.StatusNotFound, err.Error())
		} else {
			utils.RespondError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.RespondJSON(w, http.StatusOK, data)
}

// CREATE
// CreateResponden godoc
// @Summary      Tambah responden
// @Description  Membuat data responden baru (data utama diambil dari users)
// @Tags         Responden Survey
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateRespondenRequest true "Create Responden Request"
// @Success      201 {object} dto.RespondenResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      404 {object} dto.ErrorResponse
// @Failure      409 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /api/survey/responden [post]
func (h *RespondenHandler) handleCreate(w http.ResponseWriter, r *http.Request) {

	var req dto.CreateRespondenRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.service.Create(req)
	if err != nil {

		switch err.Error() {

		case "user_id wajib diisi",
			"nomor telepon tidak boleh kosong",
			"sektor tidak boleh kosong",
			"sertifikat training tidak boleh kosong",
			"sektor lainnya wajib diisi jika sektor = lainnya":
			utils.RespondError(w, http.StatusBadRequest, err.Error())

		case "user tidak ditemukan":
			utils.RespondError(w, http.StatusNotFound, err.Error())

		case "responden sudah ada":
			utils.RespondError(w, http.StatusConflict, err.Error())

		default:
			utils.RespondError(w, http.StatusInternalServerError, err.Error())
		}

		return
	}

	utils.RespondJSON(w, http.StatusCreated, resp)
}

// UPDATE
// UpdateResponden godoc
// @Summary      Update responden
// @Description  Memperbarui data responden (hanya field tambahan, bukan data users)
// @Tags         Responden Survey
// @Accept       json
// @Produce      json
// @Param        id path int true "Responden ID"
// @Param        request body dto.UpdateRespondenRequest true "Update Responden Request"
// @Success      200 {object} dto.RespondenResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      404 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /api/survey/responden/{id} [put]
func (h *RespondenHandler) handleUpdate(w http.ResponseWriter, r *http.Request, id string) {

	idInt, err := strconv.Atoi(id)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "ID harus berupa angka")
		return
	}

	var req dto.UpdateRespondenRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.service.Update(idInt, req)
	if err != nil {

		switch err.Error() {

		case "data tidak ditemukan":
			utils.RespondError(w, http.StatusNotFound, err.Error())

		case "nomor telepon tidak boleh kosong",
			"sektor tidak boleh kosong",
			"sertifikat training tidak boleh kosong",
			"sektor lainnya wajib diisi jika sektor = lainnya":
			utils.RespondError(w, http.StatusBadRequest, err.Error())

		default:
			utils.RespondError(w, http.StatusInternalServerError, err.Error())
		}

		return
	}

	utils.RespondJSON(w, http.StatusOK, resp)
}
