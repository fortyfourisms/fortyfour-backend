package handlers

import (
	"bytes"
	"encoding/json"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/testhelpers"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthHandler_Register_Success(t *testing.T) {
	repo := testhelpers.NewMockUserRepository()
	service := services.NewAuthService(repo, "test-secret")
	handler := NewAuthHandler(service)

	reqBody := map[string]string{
		"username": "testuser",
		"password": "password123",
		"email":    "test@example.com",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Register(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}

	var response dto.AuthResponse
	json.NewDecoder(w.Body).Decode(&response)

	if response.User.Username != "testuser" {
		t.Errorf("expected username 'testuser', got '%s'", response.User.Username)
	}

	if response.Token == "" {
		t.Error("expected token in response")
	}
}

func TestAuthHandler_Register_InvalidJSON(t *testing.T) {
	repo := testhelpers.NewMockUserRepository()
	service := services.NewAuthService(repo, "test-secret")
	handler := NewAuthHandler(service)

	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Register(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestAuthHandler_Login_Success(t *testing.T) {
	repo := testhelpers.NewMockUserRepository()
	service := services.NewAuthService(repo, "test-secret")
	handler := NewAuthHandler(service)

	service.Register("testuser", "password123", "test@example.com")

	reqBody := map[string]string{
		"username": "testuser",
		"password": "password123",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Login(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response dto.AuthResponse
	json.NewDecoder(w.Body).Decode(&response)

	if response.User.Username != "testuser" {
		t.Errorf("expected username 'testuser', got '%s'", response.User.Username)
	}

	if response.Token == "" {
		t.Error("expected token in response")
	}
}
