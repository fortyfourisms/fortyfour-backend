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
	return &RespondenHandler{
		service: service,
	}
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

	case http.MethodDelete:
		if id == "" {
			utils.RespondError(w, http.StatusBadRequest, "ID wajib")
			return
		}
		h.handleDelete(w, id)

	default:
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method tidak diizinkan")
	}
}

// GetAllResponden godoc
// @Summary      Ambil semua responden
// @Description  Mengambil seluruh data responden
// @Tags         Responden Survey
// @Produce      json
// @Success      200  {array}   dto.RespondenResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /api/survey/responden [get]
func (h *RespondenHandler) handleGetAll(w http.ResponseWriter) {

	data, err := h.service.GetAll()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, data)
}

// GetRespondenByID godoc
// @Summary      Ambil responden berdasarkan ID
// @Description  Mengambil data responden berdasarkan ID
// @Tags         Responden Survey
// @Produce      json
// @Param        id   path      int  true  "Responden ID"
// @Success      200  {object}  dto.RespondenResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
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

// CreateResponden godoc
// @Summary      Tambah responden
// @Description  Membuat data responden baru
// @Tags         Responden Survey
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateRespondenRequest true "Data responden"
// @Success      201  {object}  dto.RespondenResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      409  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /api/survey/responden [post]
func (h *RespondenHandler) handleCreate(w http.ResponseWriter, r *http.Request) {

	var req dto.CreateRespondenRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	resp, err := h.service.Create(req)
	if err != nil {

		switch err.Error() {

		case "nama lengkap tidak boleh kosong",
			"jabatan tidak boleh kosong",
			"perusahaan tidak boleh kosong",
			"email tidak boleh kosong",
			"format email tidak valid",
			"nomor telepon tidak boleh kosong",
			"sektor tidak boleh kosong":
			utils.RespondError(w, http.StatusBadRequest, err.Error())

		case "email sudah terdaftar":
			utils.RespondError(w, http.StatusConflict, err.Error())

		default:
			utils.RespondError(w, http.StatusInternalServerError, err.Error())
		}

		return
	}

	utils.RespondJSON(w, http.StatusCreated, resp)
}

// UpdateResponden godoc
// @Summary      Update responden
// @Description  Memperbarui data responden berdasarkan ID
// @Tags         Responden Survey
// @Accept       json
// @Produce      json
// @Param        id      path  int                          true  "Responden ID"
// @Param        request body  dto.UpdateRespondenRequest   true  "Data responden"
// @Success      200  {object}  dto.RespondenResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      409  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /api/survey/responden/{id} [put]
func (h *RespondenHandler) handleUpdate(w http.ResponseWriter, r *http.Request, id string) {

	idInt, err := strconv.Atoi(id)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "ID harus berupa angka")
		return
	}

	var req dto.UpdateRespondenRequest

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	resp, err := h.service.Update(idInt, req)
	if err != nil {

		switch err.Error() {

		case "data tidak ditemukan":
			utils.RespondError(w, http.StatusNotFound, err.Error())

		case "nama lengkap tidak boleh kosong",
			"jabatan tidak boleh kosong",
			"perusahaan tidak boleh kosong",
			"format email tidak valid",
			"sektor tidak boleh kosong":
			utils.RespondError(w, http.StatusBadRequest, err.Error())

		case "email sudah terdaftar":
			utils.RespondError(w, http.StatusConflict, err.Error())

		default:
			utils.RespondError(w, http.StatusInternalServerError, err.Error())
		}

		return
	}

	utils.RespondJSON(w, http.StatusOK, resp)
}

// DeleteResponden godoc
// @Summary      Hapus responden
// @Description  Menghapus data responden berdasarkan ID
// @Tags         Responden Survey
// @Produce      json
// @Param        id   path      int  true  "Responden ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /api/survey/responden/{id} [delete]
func (h *RespondenHandler) handleDelete(w http.ResponseWriter, id string) {

	idInt, err := strconv.Atoi(id)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "ID harus berupa angka")
		return
	}

	err = h.service.Delete(idInt)
	if err != nil {

		if err.Error() == "data tidak ditemukan" {
			utils.RespondError(w, http.StatusNotFound, err.Error())
		} else {
			utils.RespondError(w, http.StatusInternalServerError, err.Error())
		}

		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "Delete success",
	})
}
