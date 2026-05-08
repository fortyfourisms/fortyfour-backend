package handlers

import (
	"encoding/json"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/utils"
	"net"
	"net/http"
	"strings"
)

type ChallengeHandler struct {
	turnstileValidator *utils.TurnstileValidator
	tokenService       *services.TokenService
}

func NewChallengeHandler(turnstileValidator *utils.TurnstileValidator, tokenService *services.TokenService) *ChallengeHandler {
	return &ChallengeHandler{
		turnstileValidator: turnstileValidator,
		tokenService:       tokenService,
	}
}

// Verify handles the verification of Turnstile token and sets the gate cookie
//
//	@Summary		Verifikasi Challenge Gate
//	@Description	Validasi token Cloudflare Turnstile untuk mendapatkan akses ke aplikasi. Memberikan cookie gate_verified.
//	@Tags			Challenge
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.ChallengeRequest	true	"Turnstile Token"
//	@Success		200		{object}	map[string]string		"message"
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/api/challenge/verify [post]
func (h *ChallengeHandler) Verify(w http.ResponseWriter, r *http.Request) {
	var req dto.ChallengeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if h.turnstileValidator == nil {
		utils.RespondError(w, http.StatusInternalServerError, "Turnstile validator not configured")
		return
	}

	// Get Remote IP
	remoteIP := r.Header.Get("CF-Connecting-IP")
	if remoteIP == "" {
		remoteIP = r.Header.Get("X-Forwarded-For")
	}
	if remoteIP != "" {
		remoteIP = strings.TrimSpace(strings.Split(remoteIP, ",")[0])
	} else {
		remoteIP = r.RemoteAddr
		if host, _, err := net.SplitHostPort(remoteIP); err == nil {
			remoteIP = host
		}
	}

	// Validate Turnstile
	turnstileResp, err := h.turnstileValidator.Validate(req.TurnstileToken, remoteIP)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Gagal memvalidasi Turnstile")
		return
	}

	if !turnstileResp.Success {
		utils.RespondError(w, http.StatusBadRequest, "Verifikasi Turnstile gagal: "+strings.Join(turnstileResp.ErrorCodes, ", "))
		return
	}

	// Set Gate Cookie
	h.tokenService.SetGateCookie(w)

	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "Challenge verified successfully. Access granted.",
	})
}
