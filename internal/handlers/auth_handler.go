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
