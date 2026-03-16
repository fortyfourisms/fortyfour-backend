package handlers

import (
	"encoding/json"
	"ikas/internal/dto"
	"ikas/internal/services"
	"ikas/internal/utils"
	"io"
	"net/http"
	"strings"

	"ikas/internal/middleware"
	"fortyfour-backend/pkg/logger"

	"github.com/google/uuid"
)

type IkasHandler struct {
	service *services.IkasService
}

func NewIkasHandler(service *services.IkasService) *IkasHandler {
	return &IkasHandler{
		service: service,
	}
}

func (h *IkasHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	suffix := utils.ExtractID(r.URL.Path, "ikas")

	if suffix == "import" && r.Method == http.MethodPost {
		h.handleImport(w, r)
		return
	}

	id := suffix

	switch r.Method {
	case http.MethodGet:
		if id == "" {
			h.handleGetAll(w, r)
		} else {
			h.handleGetByID(w, r, id)
		}
	case http.MethodPost:
		if id != "" {
			utils.RespondError(w, 400, "ID tidak diperlukan untuk create")
			return
		}
		h.handleCreate(w, r)
	case http.MethodPut:
		if id == "" {
			utils.RespondError(w, 400, "ID wajib")
			return
		}
		h.handleUpdate(w, r, id)
	case http.MethodDelete:
		if id == "" {
			utils.RespondError(w, 400, "ID wajib")
			return
		}
		h.handleDelete(w, r, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// GetAllIkas godoc
// @Summary      List semua ikas
// @Description  Mengambil seluruh data ikas
// @Tags         Ikas
// @Produce      json
// @Success      200  {array}  dto.IkasResponse
// @Failure      500  {object} dto.ErrorResponse
// @Router       /api/maturity/ikas [get]
func (h *IkasHandler) handleGetAll(w http.ResponseWriter, _ *http.Request) {
	data, err := h.service.GetAll()
	if err != nil {
		logger.Error(err, "operation failed")
		utils.RespondError(w, 500, err.Error())
		return
	}
	utils.RespondJSON(w, 200, map[string]interface{}{
		"message": "Berhasil mengambil data",
		"data":    data,
		"total":   len(data),
	})
}

// GetIkasByID godoc
// @Summary      Ambil ikas berdasarkan ID
// @Description  Mengambil satu data ikas
// @Tags         Ikas
// @Produce      json
// @Param        id   path      string  true  "Ikas ID"
// @Success      200  {object} dto.IkasResponse
// @Failure      404  {object} dto.ErrorResponse
// @Router       /api/maturity/ikas/{id} [get]
func (h *IkasHandler) handleGetByID(w http.ResponseWriter, _ *http.Request, id string) {
	data, err := h.service.GetByID(id)
	if err != nil {
		logger.Error(err, "operation failed")
		utils.RespondError(w, 404, "Data tidak ditemukan")
		return
	}
	utils.RespondJSON(w, 200, map[string]interface{}{
		"message": "Berhasil mengambil data",
		"data":    data,
	})
}

// CreateIkas godoc
// @Summary      Tambah ikas baru
// @Description  Membuat record ikas
// @Tags         Ikas
// @Accept       json
// @Produce      json
// @Param        ikas body dto.CreateIkasRequest true "Data ikas"
// @Success      201  {object} dto.IkasResponse
// @Failure      400  {object} dto.ErrorResponse
// @Router       /api/maturity/ikas [post]
func (h *IkasHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateIkasRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error(err, "operation failed")
		utils.RespondError(w, 400, "Invalid request body")
		return
	}

	newID := uuid.New().String()
	userID, _ := r.Context().Value(middleware.UserIDKey).(string)

	if err := h.service.Create(r.Context(), req, newID, userID); err != nil {
		logger.Error(err, "operation failed")
		utils.RespondError(w, 400, err.Error())
		return
	}

	utils.RespondJSON(w, 201, map[string]interface{}{
		"message": "Berhasil menyimpan data",
		"id":      newID,
	})
}

// UpdateIkas godoc
// @Summary      Update ikas
// @Description  Mengubah data ikas berdasarkan ID
// @Tags         Ikas
// @Accept       json
// @Produce      json
// @Param        id      path      string  true  "Ikas ID"
// @Param        ikas body      dto.UpdateIkasRequest true "Data update"
// @Success      200  {object} dto.IkasResponse
// @Failure      400  {object} dto.ErrorResponse
// @Router       /api/maturity/ikas/{id} [put]
func (h *IkasHandler) handleUpdate(w http.ResponseWriter, r *http.Request, id string) {
	var req dto.UpdateIkasRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error(err, "operation failed")
		utils.RespondError(w, 400, "Invalid request body")
		return
	}

	userID, _ := r.Context().Value(middleware.UserIDKey).(string)
	err := h.service.Update(r.Context(), id, req, userID)
	if err != nil {
		logger.Error(err, "operation failed")
		if strings.Contains(err.Error(), "no rows") {
			utils.RespondError(w, 404, "Data tidak ditemukan")
		} else {
			utils.RespondError(w, 400, err.Error())
		}
		return
	}

	utils.RespondJSON(w, 200, map[string]interface{}{
		"message": "Berhasil menyimpan data",
		"id":      id,
	})
}

// DeleteIkas godoc
// @Summary      Hapus ikas
// @Description  Menghapus data ikas berdasarkan ID
// @Tags         Ikas
// @Produce      json
// @Param        id  path  string  true  "Ikas ID"
// @Success      200  {object} dto.MessageResponse
// @Failure      400  {object} dto.ErrorResponse
// @Router       /api/maturity/ikas/{id} [delete]
func (h *IkasHandler) handleDelete(w http.ResponseWriter, r *http.Request, id string) {
	userID, _ := r.Context().Value(middleware.UserIDKey).(string)
	if err := h.service.Delete(r.Context(), id, userID); err != nil {
		logger.Error(err, "operation failed")
		if strings.Contains(err.Error(), "no rows") {
			utils.RespondError(w, 404, "Data tidak ditemukan")
		} else {
			utils.RespondError(w, 400, err.Error())
		}
		return
	}

	utils.RespondJSON(w, 200, map[string]interface{}{
		"message": "Berhasil menghapus data",
		"id":      id,
	})
}

// ImportIkas godoc
// @Summary      Import IKAS dari Excel
// @Description  Import data IKAS dari file Excel (sheet ke-7)
// @Tags         Ikas
// @Accept       multipart/form-data
// @Produce      json
// @Param        file formData file true "File Excel (.xlsx)"
// @Param        id_perusahaan formData string true "ID Perusahaan"
// @Param        tanggal formData string true "Tanggal (YYYY-MM-DD)"
// @Param        responden formData string true "Nama Responden"
// @Param        telepon formData string true "Nomor Telepon"
// @Param        jabatan formData string true "Jabatan"
// @Success      201  {object} dto.ImportIkasResponse
// @Failure      400  {object} dto.ErrorResponse
// @Router       /api/maturity/ikas/import [post]
func (h *IkasHandler) handleImport(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		logger.Error(err, "operation failed")
		utils.RespondError(w, 400, "Gagal parse form data")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		logger.Error(err, "operation failed")
		utils.RespondError(w, 400, "File 'file' tidak ditemukan")
		return
	}
	defer file.Close()

	if !strings.HasSuffix(strings.ToLower(header.Filename), ".xlsx") {
		utils.RespondError(w, 400, "File harus berformat .xlsx")
		return
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		logger.Error(err, "operation failed")
		utils.RespondError(w, 400, "Gagal membaca file")
		return
	}

	userID, _ := r.Context().Value(middleware.UserIDKey).(string)
	newID, err := h.service.ImportFromExcel(r.Context(), fileBytes, userID)
	if err != nil {
		logger.Error(err, "operation failed")
		response := struct {
			Success bool     `json:"success"`
			Message string   `json:"message"`
			Errors  []string `json:"errors,omitempty"`
		}{
			Success: false,
			Message: "Import gagal",
			Errors:  []string{err.Error()},
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	utils.RespondJSON(w, 201, map[string]interface{}{
		"success": true,
		"message": "Berhasil menyimpan data",
		"id":      newID,
	})
}
