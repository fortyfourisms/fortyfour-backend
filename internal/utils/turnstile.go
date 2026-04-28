package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// TurnstileResponse represents the response from Cloudflare Siteverify API
type TurnstileResponse struct {
	Success     bool      `json:"success"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []string  `json:"error-codes"`
	Action      string    `json:"action"`
	CData       string    `json:"cdata"`
}

// TurnstileValidator handles validation of Cloudflare Turnstile tokens
type TurnstileValidator struct {
	SecretKey string
	Client    *http.Client
}

// NewTurnstileValidator creates a new TurnstileValidator
func NewTurnstileValidator(secretKey string) *TurnstileValidator {
	return &TurnstileValidator{
		SecretKey: secretKey,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Validate verifies the turnstile token with Cloudflare
func (v *TurnstileValidator) Validate(token string, remoteIP string) (*TurnstileResponse, error) {
	if token == "" {
		return &TurnstileResponse{Success: false, ErrorCodes: []string{"missing-input-response"}}, nil
	}

	// Prepare form data
	data := url.Values{}
	data.Set("secret", v.SecretKey)
	data.Set("response", token)
	if remoteIP != "" {
		data.Set("remoteip", remoteIP)
	}

	// Call Cloudflare API
	resp, err := v.Client.PostForm("https://challenges.cloudflare.com/turnstile/v0/siteverify", data)
	if err != nil {
		return nil, fmt.Errorf("failed to call turnstile api: %w", err)
	}
	defer resp.Body.Close()

	var result TurnstileResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode turnstile response: %w", err)
	}

	return &result, nil
}
