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
// @Description  Authenticate user and return JWT tokens or MFA setup/verification required response
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        login body dto.LoginRequest true "Login payload"
// @Success      200  {object} dto.AuthResponse
// @Success      200  {object} map[string]interface{} "mfa_setup_required or mfa_required response"
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

	// Check if MFA is NOT enabled yet
	if !user.MFAEnabled {
		// Force user to setup MFA - return mfa_setup_required
		setupToken, err := h.authService.CreateMFASetupToken(user.ID)
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "failed to create setup token")
			rollbar.Error(err)
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

/* ===================== MFA HANDLERS (UPDATED FOR MICROSOFT-STYLE) ===================== */

// SetupMFA godoc
// @Summary      Setup MFA
// @Description  Generate MFA provisioning URI and secret (Microsoft-style: accepts setup_token)
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        setup body map[string]string false "Setup token (for unauthenticated setup)"
// @Success      200  {object} map[string]string "provisioning_uri and secret"
// @Failure      400  {object} dto.ErrorResponse
// @Failure      401  {object} dto.ErrorResponse
// @Router       /api/mfa/setup [post]
func (h *AuthHandler) SetupMFA(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		rollbar.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Try to get userID from context (authenticated user) or from setup_token
	userID := middleware.GetUserID(r.Context())
	
	// If no userID from context, check for setup_token in body
	if userID == "" {
		var req struct {
			SetupToken string `json:"setup_token"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.RespondError(w, http.StatusBadRequest, "invalid body")
			rollbar.Error(err)
			return
		}

		// Validate setup token and get userID
		var err error
		userID, err = h.authService.ValidateMFASetupToken(req.SetupToken)
		if err != nil {
			utils.RespondError(w, http.StatusUnauthorized, "invalid or expired setup token")
			rollbar.Error(err)
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
		rollbar.Error(err)
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"provisioning_uri": uri,
		"secret":           secret,
	})
}

// EnableMFA godoc
// @Summary      Enable MFA
// @Description  Verify MFA code and enable MFA (Microsoft-style: returns tokens immediately)
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        enable body map[string]string true "MFA code and optional setup_token"
// @Success      200  {object} dto.AuthResponse
// @Failure      400  {object} dto.ErrorResponse
// @Failure      401  {object} dto.ErrorResponse
// @Router       /api/mfa/enable [post]
func (h *AuthHandler) EnableMFA(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		rollbar.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Try to get userID from context (authenticated user) or from setup_token
	userID := middleware.GetUserID(r.Context())
	
	var req struct {
		Code       string  `json:"code"`
		SetupToken *string `json:"setup_token,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid body")
		rollbar.Error(err)
		return
	}

	// If no userID from context but have setup_token, validate it
	if userID == "" && req.SetupToken != nil {
		var err error
		userID, err = h.authService.ValidateMFASetupToken(*req.SetupToken)
		if err != nil {
			utils.RespondError(w, http.StatusUnauthorized, "invalid or expired setup token")
			rollbar.Error(err)
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
		rollbar.Error(err)
		return
	}

	// Return user data and tokens
	utils.RespondJSON(w, http.StatusOK, dto.AuthResponse{
		User:         *user,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    tokens.ExpiresAt,
		Message:      "MFA enabled successfully",
	})
}

// VerifyMFA godoc
// @Summary      Verify MFA
// @Description  Verify MFA code and return access tokens
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        verify body map[string]string true "MFA token and code"
// @Success      200  {object} dto.AuthResponse
// @Failure      400  {object} dto.ErrorResponse
// @Failure      401  {object} dto.ErrorResponse
// @Router       /api/mfa/verify [post]
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