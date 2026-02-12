package handlers

import (
	"encoding/json"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/utils"
	"net/http"
)

type AuthHandler struct {
	authService  *services.AuthService
	tokenService *services.TokenService
}

func NewAuthHandler(authService *services.AuthService, tokenService *services.TokenService) *AuthHandler {
	return &AuthHandler{
		authService:  authService,
		tokenService: tokenService,
	}
}

// @Summary Register user baru
// @Description Mendaftarkan user baru. Token dikirim via HTTP-only cookies, BUKAN di response body.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Register data"
// @Success 200 {object} map[string]interface{} "message dan user info (tanpa token)"
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

	// Validate user credentials
	user, err := h.authService.Register(req.Username, req.Password, req.Email, req.RoleID, req.IDJabatan)
	if err != nil {
		utils.RespondError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Generate token pair
	tokens, err := h.tokenService.GenerateTokenPair(user.ID, user.Username, user.RoleName)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to generate tokens")
		return
	}

	// Set secure HTTP-only cookies
	h.tokenService.SetAuthCookies(w, tokens)

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

	// Validate user credentials
	user, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		utils.RespondError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Generate token pair
	tokens, err := h.tokenService.GenerateTokenPair(user.ID, user.Username, user.RoleName)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to generate tokens")
		return
	}

	// Set secure HTTP-only cookies
	h.tokenService.SetAuthCookies(w, tokens)

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
	userID, ok := r.Context().Value("user_id").(string)
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
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	username, _ := r.Context().Value("username").(string)
	role, _ := r.Context().Value("role").(string)

	response := map[string]interface{}{
		"user": map[string]interface{}{
			"id":       userID,
			"username": username,
			"role":     role,
		},
	}

	utils.RespondJSON(w, http.StatusOK, response)
}
