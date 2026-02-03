package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/middleware"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/utils"
	"fortyfour-backend/internal/validator"

	"github.com/rollbar/rollbar-go"
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

// Register godoc
// @Summary      Register new user
// @Description  Create account with optional company creation/selection and return JWT tokens
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
		rollbar.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		rollbar.Error(err)
		return
	}

	// Trim spaces untuk mencegah string kosong
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

	// Validasi menggunakan validator
	if err := validator.Validate(req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		rollbar.Error(err)
		return
	}

	user, tokens, err := h.authService.Register(req, h.perusahaanService)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		rollbar.Error(err)
		return
	}

	// tokens is *dto.TokenPair with ExpiresAt already formatted string
	utils.RespondJSON(w, http.StatusCreated, dto.AuthResponse{
		User:         *user,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    tokens.ExpiresAt,
	})
}

// Login godoc
// @Summary      Login user
// @Description  Authenticate user and return JWT tokens or MFA pending response
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        login body dto.LoginRequest true "Login payload"
// @Success      200  {object} dto.AuthResponse
// @Success      200  {object} map[string]interface{} "mfa_required response"
// @Failure      400  {object} dto.ErrorResponse
// @Failure      401  {object} dto.ErrorResponse
// @Router       /api/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		rollbar.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		rollbar.Error(err)
		return
	}

	// Trim spaces
	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)

	// Validasi menggunakan validator
	if err := validator.Validate(req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		rollbar.Error(err)
		return
	}

	user, tokens, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		utils.RespondError(w, http.StatusUnauthorized, err.Error())
		rollbar.Error(err)
		return
	}

	// Jika tokens == nil & user.MFAEnabled true -> return mfa_required with pending token
	if tokens == nil && user != nil && user.MFAEnabled {
		mfaToken, err := h.authService.CreateMFAPending(user.ID)
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "failed to create mfa token")
			rollbar.Error(err)
			return
		}
		utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
			"mfa_required": true,
			"mfa_token":    mfaToken,
			"message":      "MFA required. Please verify using your MFA code.",
		})
		return
	}

	// tokens is *dto.TokenPair with ExpiresAt already formatted string
	utils.RespondJSON(w, http.StatusOK, dto.AuthResponse{
		User:         *user,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    tokens.ExpiresAt,
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
		rollbar.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req dto.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		rollbar.Error(err)
		return
	}

	// Trim spaces
	req.RefreshToken = strings.TrimSpace(req.RefreshToken)

	// Validasi menggunakan validator
	if err := validator.Validate(req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		rollbar.Error(err)
		return
	}

	// Assume tokenService.RefreshAccessToken returns *models.TokenPair (with time.Time ExpiresAt)
	modelTokens, err := h.tokenService.RefreshAccessToken(req.RefreshToken)
	if err != nil {
		utils.RespondError(w, http.StatusUnauthorized, err.Error())
		rollbar.Error(err)
		return
	}

	// Map to DTO (string ExpiresAt)
	dtoTokens := dto.TokenPair{
		AccessToken:  modelTokens.AccessToken,
		RefreshToken: modelTokens.RefreshToken,
		ExpiresAt:    modelTokens.ExpiresAt.Format(time.RFC3339),
	}

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"access_token":  dtoTokens.AccessToken,
		"refresh_token": dtoTokens.RefreshToken,
		"expires_at":    dtoTokens.ExpiresAt,
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
// UNCHANGED from original
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		rollbar.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req dto.LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		rollbar.Error(err)
		return
	}

	// Trim spaces
	req.RefreshToken = strings.TrimSpace(req.RefreshToken)

	// Validasi menggunakan validator
	if err := validator.Validate(req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		rollbar.Error(err)
		return
	}

	if err := h.authService.Logout(req.RefreshToken); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		rollbar.Error(err)
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "Logged out successfully",
	})
}

/* ===================== MFA HANDLERS ===================== */

// SetupMFA ...
func (h *AuthHandler) SetupMFA(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		rollbar.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		utils.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	uri, secret, err := h.authService.SetupMFA(userID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		rollbar.Error(err)
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"provisioning_uri": uri,
		"secret":           secret,
	})
}

// EnableMFA ...
func (h *AuthHandler) EnableMFA(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		rollbar.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		utils.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req struct {
		Code string `json:"code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid body")
		rollbar.Error(err)
		return
	}

	if err := h.authService.EnableMFA(userID, req.Code); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		rollbar.Error(err)
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "mfa enabled"})
}

// VerifyMFA ...
func (h *AuthHandler) VerifyMFA(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		rollbar.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		MFAToken string `json:"mfa_token"`
		Code     string `json:"code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid body")
		rollbar.Error(err)
		return
	}

	user, tokens, err := h.authService.VerifyMFA(req.MFAToken, req.Code)
	if err != nil {
		utils.RespondError(w, http.StatusUnauthorized, err.Error())
		rollbar.Error(err)
		return
	}

	// tokens is *dto.TokenPair with ExpiresAt string
	utils.RespondJSON(w, http.StatusOK, dto.AuthResponse{
		User:         *user,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    tokens.ExpiresAt,
	})
}