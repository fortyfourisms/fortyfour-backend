package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/middleware"
	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/utils"
	"fortyfour-backend/internal/validator"
)

// AuthHandler handles authentication-related HTTP endpoints.
type AuthHandler struct {
	authService       *services.AuthService
	tokenService      *services.TokenService
	perusahaanService services.PerusahaanServiceInterface
}

func NewAuthHandler(
	authService *services.AuthService,
	tokenService *services.TokenService,
	perusahaanService services.PerusahaanServiceInterface,
) *AuthHandler {
	return &AuthHandler{
		authService:       authService,
		tokenService:      tokenService,
		perusahaanService: perusahaanService,
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

// @Summary Get current user info
// @Description Mengambil informasi user yang sedang login.
// @Tags Auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} dto.ErrorResponse
// @Router /api/me [get]
// Me returns current user info
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	// Get user info from context (set by auth middleware)
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	username, _ := r.Context().Value(middleware.UsernameKey).(string)
	role, _ := r.Context().Value(middleware.RoleKey).(string)

	response := map[string]interface{}{
		"user": map[string]interface{}{
			"id":       userID,
			"username": username,
			"role":     role,
		},
	}

	utils.RespondJSON(w, http.StatusOK, response)
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
