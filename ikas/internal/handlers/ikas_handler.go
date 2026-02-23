package handlers

import (
	"encoding/json"
	"ikas/internal/dto"
	"ikas/internal/services"
	"ikas/internal/utils"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/rollbar/rollbar-go"
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
	path := strings.TrimPrefix(r.URL.Path, "/api/ikas")

	// Handle import endpoint
	if path == "/import" && r.Method == http.MethodPost {
		h.handleImport(w, r)
		return
	}

	id := strings.TrimPrefix(path, "/")

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
// @Router       /api/ikas [get]
func (h *IkasHandler) handleGetAll(w http.ResponseWriter, _ *http.Request) {
	data, err := h.service.GetAll()
	if err != nil {
		rollbar.Error(err)
		utils.RespondError(w, 500, err.Error())
		return
	}
	utils.RespondJSON(w, 200, data)
}

// GetIkasByID godoc
// @Summary      Ambil ikas berdasarkan ID
// @Description  Mengambil satu data ikas
// @Tags         Ikas
// @Produce      json
// @Param        id   path      string  true  "Ikas ID"
// @Success      200  {object} dto.IkasResponse
// @Failure      404  {object} dto.ErrorResponse
// @Router       /api/ikas/{id} [get]
func (h *IkasHandler) handleGetByID(w http.ResponseWriter, _ *http.Request, id string) {
	data, err := h.service.GetByID(id)
	if err != nil {
		rollbar.Error(err)
		utils.RespondError(w, 404, "Data tidak ditemukan")
		return
	}
	utils.RespondJSON(w, 200, data)
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
// @Router       /api/ikas [post]
func (h *IkasHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateIkasRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rollbar.Error(err)
		utils.RespondError(w, 400, "Invalid request body")
		return
	}

	// Generate UUID untuk ID baru
	newID := uuid.New().String()

	// Create dengan ID
	if err := h.service.Create(req, newID); err != nil {
		rollbar.Error(err)
		utils.RespondError(w, 400, err.Error())
		return
	}

	// Ambil data yang baru dibuat (dengan JOIN)
	resp, err := h.service.GetByID(newID)
	if err != nil {
		rollbar.Error(err)
		utils.RespondError(w, 500, "Data berhasil dibuat tapi gagal diambil")
		return
	}

	utils.RespondJSON(w, 201, resp)
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
// @Router       /api/ikas/{id} [put]
func (h *IkasHandler) handleUpdate(w http.ResponseWriter, r *http.Request, id string) {
	var req dto.UpdateIkasRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rollbar.Error(err)
		utils.RespondError(w, 400, "Invalid request body")
		return
	}

	resp, err := h.service.Update(id, req)
	if err != nil {
		rollbar.Error(err)
		utils.RespondError(w, 400, err.Error())
		return
	}

	utils.RespondJSON(w, 200, resp)
}

// DeleteIkas godoc
// @Summary      Hapus ikas
// @Description  Menghapus data ikas berdasarkan ID
// @Tags         Ikas
// @Produce      json
// @Param        id  path  string  true  "Ikas ID"
// @Success      200  {object} dto.MessageResponse
// @Failure      400  {object} dto.ErrorResponse
// @Router       /api/ikas/{id} [delete]
func (h *IkasHandler) handleDelete(w http.ResponseWriter, _ *http.Request, id string) {
	if err := h.service.Delete(id); err != nil {
		rollbar.Error(err)
		utils.RespondError(w, 400, err.Error())
		return
	}

	utils.RespondJSON(w, 200, map[string]string{"message": "Delete success"})
}

// Tambahkan method baru di struct IkasHandler

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
// @Router       /api/ikas/import [post]
func (h *IkasHandler) handleImport(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form (max 10MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		rollbar.Error(err)
		utils.RespondError(w, 400, "Gagal parse form data")
		return
	}

	// Ambil file dari form
	file, header, err := r.FormFile("file")
	if err != nil {
		rollbar.Error(err)
		utils.RespondError(w, 400, "File 'file' tidak ditemukan")
		return
	}
	defer file.Close()

	// Validasi extension
	if !strings.HasSuffix(strings.ToLower(header.Filename), ".xlsx") {
		utils.RespondError(w, 400, "File harus berformat .xlsx")
		return
	}

	// Baca file ke memory
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		rollbar.Error(err)
		utils.RespondError(w, 400, "Gagal membaca file")
		return
	}

	// Import data - semua data diambil dari Excel
	resp, err := h.service.ImportFromExcel(fileBytes)
	if err != nil {
		rollbar.Error(err)
		response := dto.ImportIkasResponse{
			Success: false,
			Message: "Import gagal",
			Errors:  []string{err.Error()},
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	// Success response
	response := dto.ImportIkasResponse{
		Success: true,
		Message: "Import berhasil",
		Data:    resp,
	}
	utils.RespondJSON(w, 201, response)
}
