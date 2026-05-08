package handlers

import (
	"encoding/json"
	"ikas/internal/dto"
	"ikas/internal/middleware"
	"ikas/internal/services"
	"ikas/internal/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/rollbar/rollbar-go"
)

type JawabanGulihHandler struct {
	service *services.JawabanGulihService
}

func NewJawabanGulihHandler(service *services.JawabanGulihService) *JawabanGulihHandler {
	return &JawabanGulihHandler{service: service}
}

func (h *JawabanGulihHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	method := r.Method

	switch {
	case method == http.MethodPost && path == "/api/maturity/jawaban-gulih":
		h.handleCreate(w, r)
	case method == http.MethodGet && path == "/api/maturity/jawaban-gulih":
		h.handleGetAll(w, r)
	case method == http.MethodGet && strings.HasPrefix(path, "/api/maturity/jawaban-gulih/"):
		h.handleGetByID(w, r)
	case method == http.MethodPut && strings.HasPrefix(path, "/api/maturity/jawaban-gulih/"):
		h.handleUpdate(w, r)
	case method == http.MethodDelete && strings.HasPrefix(path, "/api/maturity/jawaban-gulih/"):
		h.handleDelete(w, r)
	default:
		utils.RespondError(w, http.StatusNotFound, "Endpoint tidak ditemukan")
	}
}

// @Summary		Create Jawaban Gulih
// @Description	Membuat record jawaban gulih baru (dikirim ke buffer RabbitMQ)
// @Tags			Jawaban Gulih
// @Accept			json
// @Produce		json
// @Param			request	body		dto.CreateJawabanGulihRequest	true	"Jawaban Gulih Request"
// @Success		201		{object}	utils.JSONResponse
// @Router			/api/maturity/jawaban-gulih [post]
func (h *JawabanGulihHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateJawabanGulihRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	userRole, _ := r.Context().Value(middleware.Role).(string)
	userPerusahaanID, _ := r.Context().Value(middleware.PerusahaanIDKey).(string)

	msg, err := h.service.Create(req, userRole, userPerusahaanID)
	if err != nil {
		rollbar.Error(err)
		switch err.Error() {
		case "pertanyaan_gulih_id tidak valid",
			"ikas_id tidak boleh kosong",
			"format ikas_id tidak valid",
			"jawaban_gulih harus bernilai antara 0 sampai 5, atau null untuk N/A",
			"validasi hanya boleh diisi jika evidence ada",
			"validasi hanya boleh berisi 'yes' atau 'no'":
			utils.RespondError(w, 400, err.Error())
		case "pertanyaan_gulih_id tidak ditemukan",
			"ikas_id tidak ditemukan":
			utils.RespondError(w, 404, err.Error())
		case "pertanyaan ini sudah pernah diisi untuk asesmen ini":
			utils.RespondError(w, 409, err.Error())
		default:
			utils.RespondError(w, 500, err.Error())
		}
		return
	}

	utils.RespondSuccess(w, 201, msg, nil)
}

// @Summary		Get All Jawaban Gulih
// @Description	Mengambil seluruh data jawaban gulih. Jika ikas_id diberikan, mengembalikan Unified Response (Main + Buffer) dengan metrik penyelesaian.
// @Tags			Jawaban Gulih
// @Produce		json
// @Param			ikas_id				query		string	false	"Filter by IKAS ID (Unified API)"
// @Param			perusahaan_id		query		string	false	"Filter by Perusahaan ID"
// @Param			pertanyaan_gulih_id	query		int		false	"Filter by Pertanyaan Gulih ID"
// @Success		200					{object}	dto.UnifiedJawabanGulihResponse
// @Router			/api/maturity/jawaban-gulih [get]
func (h *JawabanGulihHandler) handleGetAll(w http.ResponseWriter, r *http.Request) {
	userRole, _ := r.Context().Value(middleware.Role).(string)
	userPerusahaanID, _ := r.Context().Value(middleware.PerusahaanIDKey).(string)

	if userRole != "admin" && userRole != "staff" && (userPerusahaanID == "" || userPerusahaanID == "null") {
		utils.RespondSuccess(w, 200, "Berhasil mengambil data jawaban gulih", dto.UnifiedJawabanGulihResponse{Data: []dto.JawabanGulihResponse{}, Count: 0})
		return
	}

	ikasID := r.URL.Query().Get("ikas_id")
	pertanyaanIDStr := r.URL.Query().Get("pertanyaan_gulih_id")

	var data []dto.JawabanGulihResponse
	var err error

	if ikasID != "" {
		data, err := h.service.GetByIkasID(ikasID, userRole, userPerusahaanID)
		if err != nil {
			status := http.StatusInternalServerError
			if err.Error() == "format ikas_id tidak valid" {
				status = http.StatusBadRequest
			}
			utils.RespondError(w, status, err.Error())
			return
		}
		utils.RespondSuccess(w, 200, "Berhasil mengambil data jawaban gulih", data)
		return
	} else if pertanyaanIDStr != "" {
		pID, _ := strconv.Atoi(pertanyaanIDStr)
		data, err = h.service.GetByPertanyaan(pID)
	} else {
		if userRole != "admin" && userRole != "staff" {
			data, err = h.service.GetByPerusahaanID(userPerusahaanID, userRole, userPerusahaanID)
		} else {
			data, err = h.service.GetAll(userRole)
		}
	}

	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondListData(w, 200, "Berhasil mengambil data jawaban gulih", data, len(data))
}

// @Summary		Get Jawaban Gulih by ID
// @Description	Get a specific gulih answer by its ID
// @Tags			Jawaban Gulih
// @Produce		json
// @Param			id	path		int	true	"Jawaban Gulih ID"
// @Success		200	{object}	map[string]interface{}
// @Router			/api/maturity/jawaban-gulih/{id} [get]
func (h *JawabanGulihHandler) handleGetByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ExtractIntID(r.URL.Path, "jawaban-gulih")
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "ID tidak valid")
		return
	}

	userRole, _ := r.Context().Value(middleware.Role).(string)
	userPerusahaanID, _ := r.Context().Value(middleware.PerusahaanIDKey).(string)

	resp, err := h.service.GetByID(id, userRole, userPerusahaanID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "data tidak ditemukan" {
			status = http.StatusNotFound
		}
		utils.RespondError(w, status, err.Error())
		return
	}

	utils.RespondSuccess(w, 200, "Berhasil mengambil data jawaban gulih", resp)
}

// @Summary		Update Jawaban Gulih
// @Description	Mengubah data jawaban gulih berdasarkan ID
// @Tags			Jawaban Gulih
// @Accept			json
// @Produce		json
// @Param			id		path		int								true	"Jawaban Gulih ID"
// @Param			request	body		dto.UpdateJawabanGulihRequest	true	"Update Request"
// @Success		200		{object}	utils.JSONResponse
// @Router			/api/maturity/jawaban-gulih/{id} [put]
func (h *JawabanGulihHandler) handleUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ExtractIntID(r.URL.Path, "jawaban-gulih")
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "ID tidak valid")
		return
	}

	var req dto.UpdateJawabanGulihRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	userID := ""
	if val := r.Context().Value(middleware.UserIDKey); val != nil {
		userID = val.(string)
	}

	userRole := ""
	if val := r.Context().Value(middleware.Role); val != nil {
		userRole = val.(string)
	}

	userPerusahaanID := ""
	if val := r.Context().Value(middleware.PerusahaanIDKey); val != nil {
		userPerusahaanID = val.(string)
	}

	updatedID, msg, err := h.service.Update(id, req, userID, userRole, userPerusahaanID)
	if err != nil {
		rollbar.Error(err)
		switch err.Error() {
		case "data tidak ditemukan":
			utils.RespondError(w, 404, err.Error())
		case "format ID tidak valid",
			"jawaban_gulih tidak boleh kosong",
			"validasi hanya boleh diisi jika evidence ada",
			"validasi hanya boleh berisi 'yes' atau 'no'":
			utils.RespondError(w, 400, err.Error())
		default:
			utils.RespondError(w, 500, err.Error())
		}
		return
	}

	utils.RespondSuccess(w, 200, msg, map[string]interface{}{
		"id": updatedID,
	})
}

// @Summary		Delete Jawaban Gulih
// @Description	Menghapus data jawaban gulih berdasarkan ID
// @Tags			Jawaban Gulih
// @Produce		json
// @Param			id	path		int	true	"Jawaban Gulih ID"
// @Success		200	{object}	utils.JSONResponse
// @Router			/api/maturity/jawaban-gulih/{id} [delete]
func (h *JawabanGulihHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ExtractIntID(r.URL.Path, "jawaban-gulih")
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "ID tidak valid")
		return
	}

	userID := ""
	if val := r.Context().Value(middleware.UserIDKey); val != nil {
		userID = val.(string)
	}

	userRole, _ := r.Context().Value(middleware.Role).(string)
	userPerusahaanID, _ := r.Context().Value(middleware.PerusahaanIDKey).(string)

	if err := h.service.Delete(id, userID, userRole, userPerusahaanID); err != nil {
		rollbar.Error(err)
		if err.Error() == "data tidak ditemukan" {
			utils.RespondError(w, 404, err.Error())
		} else {
			utils.RespondError(w, 500, err.Error())
		}
		return
	}

	utils.RespondSuccess(w, 200, "Berhasil menghapus data", map[string]interface{}{
		"id": id,
	})
}
