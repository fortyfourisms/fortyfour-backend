package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/utils"
	"fortyfour-backend/internal/validator"

	"github.com/google/uuid"
)

type UserHandler struct {
	service    *services.UserService
	uploadPath string
	sseService *services.SSEService
}

func NewUserHandler(service *services.UserService, uploadPath string, sseService *services.SSEService) *UserHandler {
	return &UserHandler{
		service:    service,
		uploadPath: uploadPath,
		sseService: sseService,
	}
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/users")
	path = strings.TrimPrefix(path, "/")

	if strings.HasSuffix(path, "/password") {
		h.handleUpdatePassword(w, r)
		return
	}
	if strings.HasSuffix(path, "/profile-photo") {
		h.handleUpdateProfilePhoto(w, r)
		return
	}
	if strings.HasSuffix(path, "/banner") {
		h.handleUpdateBanner(w, r)
		return
	}

	id := path
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

func (h *UserHandler) handleGetAll(w http.ResponseWriter, _ *http.Request) {
	data, err := h.service.GetAll()
	if err != nil {
		utils.RespondError(w, 500, err.Error())
		return
	}
	utils.RespondJSON(w, 200, data)
}

func (h *UserHandler) handleGetByID(w http.ResponseWriter, _ *http.Request, id string) {
	data, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondError(w, 404, "User tidak ditemukan")
		return
	}
	utils.RespondJSON(w, 200, data)
}

func (h *UserHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, 400, "Invalid request body")
		return
	}

	// Trim spaces untuk mencegah string kosong
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

	// Validasi menggunakan validator
	if err := validator.Validate(req); err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	}

	resp, err := h.service.Create(req)
	if err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	}

	userID := h.getUserID(r)
	h.sseService.NotifyCreate("users", resp, userID)

	utils.RespondJSON(w, 201, resp)
}

func (h *UserHandler) handleUpdate(w http.ResponseWriter, r *http.Request, id string) {
	currentUserID := h.getUserID(r)
	isAdmin := h.isAdmin(r)

	if !isAdmin && currentUserID != id {
		utils.RespondError(w, 403, "Anda hanya bisa update data diri sendiri")
		return
	}

	var req dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, 400, "Invalid request body")
		return
	}

	// Trim spaces jika field diisi
	if req.Username != nil {
		trimmed := strings.TrimSpace(*req.Username)
		req.Username = &trimmed
	}
	if req.Email != nil {
		trimmed := strings.TrimSpace(*req.Email)
		req.Email = &trimmed
	}

	// Validasi menggunakan validator
	if err := validator.Validate(req); err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	}

	if !isAdmin && req.RoleID != nil {
		utils.RespondError(w, 403, "Anda tidak bisa mengubah role")
		return
	}

	resp, err := h.service.Update(id, req)
	if err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	}

	h.sseService.NotifyUpdate("users", resp, currentUserID)

	utils.RespondJSON(w, 200, resp)
}

func (h *UserHandler) handleUpdatePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/users/")
	id := strings.TrimSuffix(path, "/password")

	currentUserID := h.getUserID(r)
	if currentUserID != id {
		utils.RespondError(w, 403, "Anda hanya bisa mengubah password sendiri")
		return
	}

	var req dto.UpdateUserPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, 400, "Invalid request body")
		return
	}

	// Trim spaces
	req.OldPassword = strings.TrimSpace(req.OldPassword)
	req.NewPassword = strings.TrimSpace(req.NewPassword)

	// Validasi
	if err := validator.Validate(req); err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	}

	if err := h.service.UpdatePassword(id, req); err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	}

	utils.RespondJSON(w, 200, map[string]string{"message": "Password berhasil diubah"})
}

func (h *UserHandler) handleUpdateProfilePhoto(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/users/")
	id := strings.TrimSuffix(path, "/profile-photo")

	currentUserID := h.getUserID(r)
	isAdmin := h.isAdmin(r)

	if !isAdmin && currentUserID != id {
		utils.RespondError(w, 403, "Anda hanya bisa update foto profile sendiri")
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.RespondError(w, 400, "File terlalu besar (max 10MB)")
		return
	}

	file, header, err := r.FormFile("profile_photo")
	if err != nil {
		utils.RespondError(w, 400, "File profile_photo wajib diisi")
		return
	}
	defer file.Close()

	if !h.isValidImageType(header.Filename) {
		utils.RespondError(w, 400, "Format file harus jpg, jpeg, atau png")
		return
	}

	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("profile_%s_%d%s", uuid.New().String(), time.Now().Unix(), ext)
	filePath := filepath.Join(h.uploadPath, filename)

	dst, err := os.Create(filePath)
	if err != nil {
		utils.RespondError(w, 500, "Gagal menyimpan file")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		utils.RespondError(w, 500, "Gagal menyimpan file")
		return
	}

	resp, err := h.service.UpdateProfilePhoto(id, filename)
	if err != nil {
		os.Remove(filePath)
		utils.RespondError(w, 400, err.Error())
		return
	}

	h.sseService.NotifyUpdate("users", resp, currentUserID)

	utils.RespondJSON(w, 200, resp)
}

func (h *UserHandler) handleUpdateBanner(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/users/")
	id := strings.TrimSuffix(path, "/banner")

	currentUserID := h.getUserID(r)
	isAdmin := h.isAdmin(r)

	if !isAdmin && currentUserID != id {
		utils.RespondError(w, 403, "Anda hanya bisa update banner sendiri")
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.RespondError(w, 400, "File terlalu besar (max 10MB)")
		return
	}

	file, header, err := r.FormFile("banner")
	if err != nil {
		utils.RespondError(w, 400, "File banner wajib diisi")
		return
	}
	defer file.Close()

	if !h.isValidImageType(header.Filename) {
		utils.RespondError(w, 400, "Format file harus jpg, jpeg, atau png")
		return
	}

	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("banner_%s_%d%s", uuid.New().String(), time.Now().Unix(), ext)
	filePath := filepath.Join(h.uploadPath, filename)

	dst, err := os.Create(filePath)
	if err != nil {
		utils.RespondError(w, 500, "Gagal menyimpan file")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		utils.RespondError(w, 500, "Gagal menyimpan file")
		return
	}

	resp, err := h.service.UpdateBanner(id, filename)
	if err != nil {
		os.Remove(filePath)
		utils.RespondError(w, 400, err.Error())
		return
	}

	h.sseService.NotifyUpdate("users", resp, currentUserID)

	utils.RespondJSON(w, 200, resp)
}

func (h *UserHandler) handleDelete(w http.ResponseWriter, r *http.Request, id string) {
	if !h.isAdmin(r) {
		utils.RespondError(w, 403, "Hanya admin yang bisa menghapus user")
		return
	}

	if err := h.service.Delete(id); err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	}

	userID := h.getUserID(r)
	h.sseService.NotifyDelete("users", id, userID)

	utils.RespondJSON(w, 200, map[string]string{"message": "User berhasil dihapus"})
}

func (h *UserHandler) getUserID(r *http.Request) string {
	if uid := r.Context().Value("user_id"); uid != nil {
		return uid.(string)
	}
	return ""
}

func (h *UserHandler) isAdmin(r *http.Request) bool {
	if role := r.Context().Value("role"); role != nil {
		return role.(string) == "admin"
	}
	return false
}

func (h *UserHandler) isValidImageType(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png"
}
