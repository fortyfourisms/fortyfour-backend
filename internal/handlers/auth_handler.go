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
	"fortyfour-backend/internal/middleware"
	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/utils"
	"fortyfour-backend/internal/validator"

	"github.com/google/uuid"
)

// AuthHandler handles authentication-related HTTP endpoints.
type AuthHandler struct {
	authService       *services.AuthService
	tokenService      *services.TokenService
	perusahaanService services.PerusahaanServiceInterface
	userService       *services.UserService
	uploadPath        string
}

func NewAuthHandler(
	authService *services.AuthService,
	tokenService *services.TokenService,
	perusahaanService services.PerusahaanServiceInterface,
	userService *services.UserService,
	uploadPath string,
) *AuthHandler {
	return &AuthHandler{
		authService:       authService,
		tokenService:      tokenService,
		perusahaanService: perusahaanService,
		userService:       userService,
		uploadPath:        uploadPath,
	}
}

// @Summary Register user baru
// @Description Mendaftarkan user baru. Token dikirim via HTTP-only cookies, BUKAN di response body.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Register data"
// @Success 201 {object} map[string]interface{} "message dan user info (tanpa token)"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Trim spaces untuk mencegah string kosong
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

	// Validasi menggunakan validator
	if err := validator.Validate(req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, tokens, err := h.authService.Register(req, h.perusahaanService)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Convert DTO tokens to models.TokenPair for cookie setting
	modelTokens := &models.TokenPair{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}
	// Parse ExpiresAt back to time.Time
	if expiresAt, err := time.Parse(time.RFC3339, tokens.ExpiresAt); err == nil {
		modelTokens.ExpiresAt = expiresAt
	}

	// Set secure HTTP-only cookies
	h.tokenService.SetAuthCookies(w, modelTokens)

	// Return user info (without tokens in response body for security)
	response := map[string]interface{}{
		"message": "Registration successful",
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.RoleName,
		},
	}

	utils.RespondJSON(w, http.StatusCreated, response)
}

// @Summary Login user
// @Description Autentikasi user. Token dikirim via HTTP-only cookies, BUKAN di response body.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} map[string]interface{} "message dan user info (tanpa token)"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/login [post]
// Login authenticates user and sets secure cookies
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Trim spaces
	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)

	// Validasi menggunakan validator
	if err := validator.Validate(req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, tokens, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		utils.RespondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Check if MFA is NOT enabled yet
	if !user.MFAEnabled {
		// Force user to setup MFA - return mfa_setup_required
		setupToken, err := h.authService.CreateMFASetupToken(user.ID)
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "failed to create setup token")
			return
		}
		utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
			"mfa_setup_required": true,
			"setup_token":        setupToken,
			"message":            "MFA setup is required. Please setup MFA to continue.",
		})
		return
	}

	// If MFA is enabled, check if tokens are nil (means need MFA verification)
	if tokens == nil && user != nil && user.MFAEnabled {
		mfaToken, err := h.authService.CreateMFAPending(user.ID)
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "failed to create mfa token")
			return
		}
		utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
			"mfa_required": true,
			"mfa_token":    mfaToken,
			"message":      "MFA required. Please verify using your MFA code.",
		})
		return
	}

	// Convert DTO tokens to models.TokenPair for cookie setting
	modelTokens := &models.TokenPair{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}
	// Parse ExpiresAt back to time.Time
	if expiresAt, err := time.Parse(time.RFC3339, tokens.ExpiresAt); err == nil {
		modelTokens.ExpiresAt = expiresAt
	}

	// Set secure HTTP-only cookies
	h.tokenService.SetAuthCookies(w, modelTokens)

	// Return user info (without tokens in response body)
	response := map[string]interface{}{
		"message": "Login successful",
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.RoleName,
		},
	}

	utils.RespondJSON(w, http.StatusOK, response)
}

// @Summary Refresh token
// @Description Refresh access token menggunakan refresh token dari cookie. Token baru dikirim via HTTP-only cookies.
// @Tags Auth
// @Produce json
// @Success 200 {object} dto.MessageResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /api/refresh [post]
// Refresh generates new access token using refresh token from cookie
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	// Get refresh token from cookie
	refreshToken, err := h.tokenService.GetRefreshTokenFromCookie(r)
	if err != nil {
		utils.RespondError(w, http.StatusUnauthorized, "Refresh token not found")
		return
	}

	// Generate new token pair
	tokens, err := h.tokenService.RefreshAccessToken(refreshToken)
	if err != nil {
		utils.RespondError(w, http.StatusUnauthorized, "Invalid or expired refresh token")
		return
	}

	// Set new cookies
	h.tokenService.SetAuthCookies(w, tokens)

	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "Token refreshed successfully",
	})
}

// @Summary Logout
// @Description Revoke refresh token dan hapus cookies autentikasi.
// @Tags Auth
// @Produce json
// @Success 200 {object} dto.MessageResponse
// @Router /api/logout [post]
// Logout revokes tokens and clears cookies
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Get refresh token from cookie to revoke it
	refreshToken, err := h.tokenService.GetRefreshTokenFromCookie(r)
	if err == nil {
		// Revoke the refresh token from Redis
		h.tokenService.RevokeRefreshToken(refreshToken)
	}

	// Clear cookies
	h.tokenService.ClearAuthCookies(w)

	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "Logged out successfully",
	})
}

// @Summary Logout dari semua perangkat
// @Description Revoke semua refresh token user dan hapus cookies.
// @Tags Auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.MessageResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/logout-all [post]
// LogoutAll revokes all user's refresh tokens across all devices
func (h *AuthHandler) LogoutAll(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Revoke all user's tokens
	if err := h.tokenService.RevokeAllUserTokens(userID); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to logout from all devices")
		return
	}

	// Clear current cookies
	h.tokenService.ClearAuthCookies(w)

	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "Logged out from all devices successfully",
	})
}

// @Summary Get current user profile
// @Description Mengambil data lengkap user yang sedang login (profil diri sendiri).
// @Tags Me
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.UserResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/me [get]
func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		utils.RespondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	user, err := h.userService.GetByID(userID)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "User tidak ditemukan")
		return
	}

	utils.RespondJSON(w, http.StatusOK, user)
}

// @Summary Update current user profile
// @Description Memperbarui data diri user yang sedang login (username dan email saja).
// Role dan jabatan tidak dapat diubah melalui endpoint ini — hanya admin yang bisa mengubahnya.
// @Tags Me
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdateMeRequest true "Data yang ingin diubah"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /api/me [put]
func (h *AuthHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		utils.RespondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req dto.UpdateMeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Trim spaces
	if req.Username != nil {
		trimmed := strings.TrimSpace(*req.Username)
		req.Username = &trimmed
	}
	if req.Email != nil {
		trimmed := strings.TrimSpace(*req.Email)
		req.Email = &trimmed
	}

	// Validasi
	if err := validator.Validate(req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Mapping ke UpdateUserRequest — role_id dan id_jabatan sengaja tidak diisi
	updateReq := dto.UpdateUserRequest{
		Username: req.Username,
		Email:    req.Email,
	}

	resp, err := h.userService.Update(userID, updateReq)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, resp)
}

// MeRouter menangani semua route di bawah /api/me
func (h *AuthHandler) MeRouter(w http.ResponseWriter, r *http.Request) {
	sub := strings.TrimPrefix(r.URL.Path, "/api/me")
	sub = strings.TrimPrefix(sub, "/")

	switch {
	case sub == "" || sub == "/":
		switch r.Method {
		case http.MethodGet:
			h.GetMe(w, r)
		case http.MethodPut:
			h.UpdateMe(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	case sub == "password" && r.Method == http.MethodPut:
		h.UpdateMePassword(w, r)
	case sub == "media" && r.Method == http.MethodPost:
		h.UpdateMeMedia(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

// @Summary Update password diri sendiri
// @Description Mengubah password user yang sedang login. Wajib mengisi password lama sebagai verifikasi.
// @Tags Me
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdateUserPasswordRequest true "Password lama dan baru"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /api/me/password [put]
func (h *AuthHandler) UpdateMePassword(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		utils.RespondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req dto.UpdateUserPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	req.OldPassword = strings.TrimSpace(req.OldPassword)
	req.NewPassword = strings.TrimSpace(req.NewPassword)
	req.ConfirmNewPassword = strings.TrimSpace(req.ConfirmNewPassword)

	if err := validator.Validate(req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.userService.UpdatePassword(userID, req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Password berhasil diubah"})
}

// @Summary Update foto profile dan/atau banner diri sendiri
// @Description Upload foto profile dan/atau banner sekaligus. Boleh kirim salah satu atau keduanya.
// @Tags Me
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param profile_photo formData file false "Foto profile (jpg/jpeg/png, max 10MB)"
// @Param banner formData file false "Banner (jpg/jpeg/png, max 10MB)"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /api/me/media [post]
func (h *AuthHandler) UpdateMeMedia(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		utils.RespondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "File terlalu besar (max 10MB)")
		return
	}

	_, hasPhoto := r.MultipartForm.File["profile_photo"]
	_, hasBanner := r.MultipartForm.File["banner"]
	if !hasPhoto && !hasBanner {
		utils.RespondError(w, http.StatusBadRequest, "Kirim minimal satu file: profile_photo atau banner")
		return
	}

	// Ensure upload directory exists
	if err := os.MkdirAll(h.uploadPath, os.ModePerm); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Gagal menyiapkan direktori upload")
		return
	}

	var profileFilename, bannerFilename *string

	if hasPhoto {
		file, header, err := r.FormFile("profile_photo")
		if err != nil {
			utils.RespondError(w, http.StatusBadRequest, "Gagal membaca file profile_photo")
			return
		}
		defer file.Close()

		if !isValidImageType(header.Filename) {
			utils.RespondError(w, http.StatusBadRequest, "profile_photo: format harus jpg, jpeg, atau png")
			return
		}

		ext := filepath.Ext(header.Filename)
		filename := fmt.Sprintf("profile_%s_%d%s", uuid.New().String(), time.Now().Unix(), ext)
		filePath := filepath.Join(h.uploadPath, filename)

		dst, err := os.Create(filePath)
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "Gagal menyimpan profile_photo")
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			os.Remove(filePath)
			utils.RespondError(w, http.StatusInternalServerError, "Gagal menyimpan profile_photo")
			return
		}
		profileFilename = &filename
	}

	if hasBanner {
		file, header, err := r.FormFile("banner")
		if err != nil {
			utils.RespondError(w, http.StatusBadRequest, "Gagal membaca file banner")
			return
		}
		defer file.Close()

		if !isValidImageType(header.Filename) {
			utils.RespondError(w, http.StatusBadRequest, "banner: format harus jpg, jpeg, atau png")
			return
		}

		ext := filepath.Ext(header.Filename)
		filename := fmt.Sprintf("banner_%s_%d%s", uuid.New().String(), time.Now().Unix(), ext)
		filePath := filepath.Join(h.uploadPath, filename)

		dst, err := os.Create(filePath)
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "Gagal menyimpan banner")
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			os.Remove(filePath)
			utils.RespondError(w, http.StatusInternalServerError, "Gagal menyimpan banner")
			return
		}
		bannerFilename = &filename
	}

	var resp *dto.UserResponse
	var err error

	switch {
	case profileFilename != nil && bannerFilename != nil:
		_, err = h.userService.UpdateProfilePhoto(userID, *profileFilename)
		if err != nil {
			utils.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}
		resp, err = h.userService.UpdateBanner(userID, *bannerFilename)
	case profileFilename != nil:
		resp, err = h.userService.UpdateProfilePhoto(userID, *profileFilename)
	case bannerFilename != nil:
		resp, err = h.userService.UpdateBanner(userID, *bannerFilename)
	}

	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, resp)
}

// isValidImageType memvalidasi ekstensi file gambar
func isValidImageType(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png"
}

/* ===================== MFA HANDLERS (MICROSOFT-STYLE) ===================== */

// @Summary Setup MFA
// @Description Generate MFA provisioning URI and secret (Microsoft-style: accepts setup_token)
// @Tags Auth
// @Accept json
// @Produce json
// @Param setup body map[string]string false "Setup token (for unauthenticated setup)"
// @Success 200 {object} map[string]string "provisioning_uri and secret"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /api/mfa/setup [post]
func (h *AuthHandler) SetupMFA(w http.ResponseWriter, r *http.Request) {
	// Try to get userID from context (authenticated user) or from setup_token
	userID := middleware.GetUserID(r.Context())

	// If no userID from context, check for setup_token in body
	if userID == "" {
		var req struct {
			SetupToken string `json:"setup_token"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.RespondError(w, http.StatusBadRequest, "invalid body")
			return
		}

		// Validate setup token and get userID
		var err error
		userID, err = h.authService.ValidateMFASetupToken(req.SetupToken)
		if err != nil {
			utils.RespondError(w, http.StatusUnauthorized, "invalid or expired setup token")
			return
		}
	}

	if userID == "" {
		utils.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	uri, secret, err := h.authService.SetupMFA(userID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"provisioning_uri": uri,
		"secret":           secret,
	})
}

// @Summary Enable MFA
// @Description Verify MFA code and enable MFA (Microsoft-style: returns tokens immediately)
// @Tags Auth
// @Accept json
// @Produce json
// @Param enable body map[string]string true "MFA code and optional setup_token"
// @Success 200 {object} dto.AuthResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /api/mfa/enable [post]
func (h *AuthHandler) EnableMFA(w http.ResponseWriter, r *http.Request) {
	// Try to get userID from context (authenticated user) or from setup_token
	userID := middleware.GetUserID(r.Context())

	var req struct {
		Code       string  `json:"code"`
		SetupToken *string `json:"setup_token,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid body")
		return
	}

	// If no userID from context but have setup_token, validate it
	if userID == "" && req.SetupToken != nil {
		var err error
		userID, err = h.authService.ValidateMFASetupToken(*req.SetupToken)
		if err != nil {
			utils.RespondError(w, http.StatusUnauthorized, "invalid or expired setup token")
			return
		}
	}

	if userID == "" {
		utils.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Enable MFA and get tokens immediately (Microsoft-style: enable = login)
	user, tokens, err := h.authService.EnableMFAAndLogin(userID, req.Code)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Convert DTO tokens to models.TokenPair for cookie setting
	modelTokens := &models.TokenPair{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}
	// Parse ExpiresAt back to time.Time
	if expiresAt, err := time.Parse(time.RFC3339, tokens.ExpiresAt); err == nil {
		modelTokens.ExpiresAt = expiresAt
	}

	// Set secure HTTP-only cookies
	h.tokenService.SetAuthCookies(w, modelTokens)

	// Return user data (without tokens in response body)
	response := map[string]interface{}{
		"message": "MFA enabled successfully",
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.RoleName,
		},
	}

	utils.RespondJSON(w, http.StatusOK, response)
}

// @Summary Verify MFA
// @Description Verify MFA code and return access tokens
// @Tags Auth
// @Accept json
// @Produce json
// @Param verify body map[string]string true "MFA token and code"
// @Success 200 {object} dto.AuthResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /api/mfa/verify [post]
func (h *AuthHandler) VerifyMFA(w http.ResponseWriter, r *http.Request) {
	var req struct {
		MFAToken string `json:"mfa_token"`
		Code     string `json:"code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid body")
		return
	}

	user, tokens, err := h.authService.VerifyMFA(req.MFAToken, req.Code)
	if err != nil {
		utils.RespondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Convert DTO tokens to models.TokenPair for cookie setting
	modelTokens := &models.TokenPair{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}
	// Parse ExpiresAt back to time.Time
	if expiresAt, err := time.Parse(time.RFC3339, tokens.ExpiresAt); err == nil {
		modelTokens.ExpiresAt = expiresAt
	}

	// Set secure HTTP-only cookies
	h.tokenService.SetAuthCookies(w, modelTokens)

	// Return user info (without tokens in response body)
	response := map[string]interface{}{
		"message": "MFA verification successful",
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.RoleName,
		},
	}

	utils.RespondJSON(w, http.StatusOK, response)
}
