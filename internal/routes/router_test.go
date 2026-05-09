package routes

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"fortyfour-backend/internal/handlers"
	"fortyfour-backend/internal/middleware"
)

func TestHealthHandler_ReturnsHealthyStatus(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	w := httptest.NewRecorder()

	healthHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if got := w.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected application/json content type, got %q", got)
	}

	var body map[string]string
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if body["status"] != "healthy" {
		t.Fatalf("expected healthy status, got %q", body["status"])
	}
	if body["timestamp"] == "" {
		t.Fatal("expected timestamp to be present")
	}
}

func TestInitRouter_HealthRoute_IsAccessibleWithoutDependencies(t *testing.T) {
	var (
		authH           *handlers.AuthHandler
		userH           *handlers.UserHandler
		perusahaanH     *handlers.PerusahaanHandler
		picH            *handlers.PICHandler
		roleH           *handlers.RoleHandler
		casbinH         *handlers.CasbinHandler
		sseH            *handlers.SSEHandler
		authM           *middleware.AuthMiddleware
		casbinM         *middleware.CasbinMiddleware
		strictLimiter   *middleware.RateLimiter
		moderateLimiter *middleware.RateLimiter
		lenientLimiter  *middleware.RateLimiter
		csirtH          *handlers.CsirtHandler
		csirtExportH    *handlers.CsirtExportHandler
		sdmCsirtH       *handlers.SdmCsirtHandler
		chatH           *handlers.ChatHandler
		sektorH         *handlers.SektorHandler
		subSektorH      *handlers.SubSektorHandler
		seH             *handlers.SEHandler
		seExportH       *handlers.SEExportHandler
		seEditReqH      *handlers.SEEditRequestHandler
		dashboardH      *handlers.DashboardHandler
		notificationH   *handlers.NotificationHandler
		ikasProxyH      *handlers.ProxyHandler
		lmsH            *handlers.LMSHandler
		beritaH         *handlers.BeritaHandler
		eventH          *handlers.EventHandler
		aktivitasH      *handlers.AktivitasHandler
		surveyProxyH    *handlers.ProxyHandler
		konversiH       *handlers.KonversiHandler
	)

	router := InitRouter(
		authH,
		userH,
		perusahaanH,
		picH,
		roleH,
		casbinH,
		sseH,
		authM,
		casbinM,
		strictLimiter,
		moderateLimiter,
		lenientLimiter,
		csirtH,
		csirtExportH,
		sdmCsirtH,
		chatH,
		sektorH,
		subSektorH,
		seH,
		seExportH,
		seEditReqH,
		dashboardH,
		notificationH,
		ikasProxyH,
		surveyProxyH,
		lmsH,
		beritaH,
		eventH,
		aktivitasH,
		konversiH,
	)

	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var body map[string]string
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	if body["status"] != "healthy" {
		t.Fatalf("expected healthy status, got %q", body["status"])
	}
}

func TestInitRouter_UnknownRoute_ReturnsNotFound(t *testing.T) {
	router := InitRouter(
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	)

	req := httptest.NewRequest(http.MethodGet, "/api/does-not-exist", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}
