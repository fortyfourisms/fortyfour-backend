package handlers

import (
	"encoding/json"
	"strings"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/utils"
	"fortyfour-backend/internal/validator"

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

// Register godoc
// @Summary      Register new user
// @Description  Create account and return JWT tokens
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        register body dto.RegisterRequest true "Registration payload"
// @Success      201  {object} dto.AuthResponse
// @Failure      400  {object} dto.ErrorResponse
// @Router       /api/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

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

	user, tokens, err := h.authService.Register(req.Username, req.Password, req.Email, req.RoleID, req.IDJabatan)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusCreated, dto.AuthResponse{
		User:         *user,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    tokens.ExpiresAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

// Login godoc
// @Summary      Login user
// @Description  Authenticate user and return JWT tokens
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        login body dto.LoginRequest true "Login payload"
// @Success      200  {object} dto.AuthResponse
// @Failure      400  {object} dto.ErrorResponse
// @Failure      401  {object} dto.ErrorResponse
// @Router       /api/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Trim spaces untuk mencegah string kosong
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

	utils.RespondJSON(w, http.StatusOK, dto.AuthResponse{
		User:         *user,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    tokens.ExpiresAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

// RefreshToken godoc
// @Summary      Refresh access token
// @Description  Generate new access & refresh token pair
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        refresh body dto.RefreshTokenRequest true "Refresh token payload"
// @Success      200  {object} dto.TokenPair
// @Failure      400  {object} dto.ErrorResponse
// @Failure      401  {object} dto.ErrorResponse
// @Router       /api/refresh [post]
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req dto.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Trim spaces
	req.RefreshToken = strings.TrimSpace(req.RefreshToken)

	// Validasi menggunakan validator
	if err := validator.Validate(req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	tokens, err := h.tokenService.RefreshAccessToken(req.RefreshToken)
	if err != nil {
		utils.RespondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
		"expires_at":    tokens.ExpiresAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

// Logout godoc
// @Summary      Logout user
// @Description  Revoke refresh token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        logout body dto.LogoutRequest true "Logout payload"
// @Success      200  {object} dto.MessageResponse
// @Failure      400  {object} dto.ErrorResponse
// @Router       /api/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req dto.LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Trim spaces
	req.RefreshToken = strings.TrimSpace(req.RefreshToken)

	// Validasi menggunakan validator
	if err := validator.Validate(req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.authService.Logout(req.RefreshToken); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "Logged out successfully",
	})
}
