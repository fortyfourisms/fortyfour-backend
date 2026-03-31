package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/casbin/casbin/v3"
	"github.com/casbin/casbin/v3/model"

	"fortyfour-backend/internal/middleware"
	"fortyfour-backend/internal/services"
)

func newTestCasbinService(t *testing.T) *services.CasbinService {
	m, err := model.NewModelFromString(`
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && r.obj == p.obj && r.act == p.act
`)
	if err != nil {
		t.Fatalf("failed to create model: %v", err)
	}

	e, err := casbin.NewEnforcer(m)
	if err != nil {
		t.Fatalf("failed to create enforcer: %v", err)
	}

	_, err = e.AddPolicy("admin", "users", "read")
	if err != nil {
		t.Fatalf("failed to add policy: %v", err)
	}

	svc := &services.CasbinService{}
	svc.SetEnforcer(e)
	return svc
}

// setupCasbinHandler membuat handler dengan service dan SSEService
func setupCasbinHandler(t *testing.T) *CasbinHandler {
	t.Helper()
	svc := newTestCasbinService(t)
	sseService := services.NewSSEService()
	return NewCasbinHandler(svc, sseService)
}

// adminCtx menyuntikkan admin context ke request
func adminCtx(req *http.Request) *http.Request {
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "test-user-id")
	ctx = context.WithValue(ctx, middleware.RoleKey, "admin")
	return req.WithContext(ctx)
}

// userCtx menyuntikkan non-admin context ke request
func userCtx(req *http.Request) *http.Request {
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "test-user-id")
	ctx = context.WithValue(ctx, middleware.RoleKey, "user")
	return req.WithContext(ctx)
}

// =========================
// GetRolePermissions
// =========================

func TestCasbinHandler_GetRolePermissions_MissingRole(t *testing.T) {
	h := &CasbinHandler{
		casbinService: &services.CasbinService{},
	}

	req := httptest.NewRequest(http.MethodGet, "/casbin/permissions", nil)
	w := httptest.NewRecorder()

	h.GetRolePermissions(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestCasbinHandler_GetRolePermissions(t *testing.T) {
	service := newTestCasbinService(t)
	h := &CasbinHandler{casbinService: service}

	req := httptest.NewRequest(http.MethodGet, "/casbin/permissions?role=admin", nil)
	w := httptest.NewRecorder()

	h.GetRolePermissions(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	if w.Body.String() == "" {
		t.Fatal("expected response body, got empty")
	}
}

func TestCasbinHandler_GetRolePermissions_RoleWithNoPermissions(t *testing.T) {
	service := newTestCasbinService(t)
	h := &CasbinHandler{casbinService: service}

	req := httptest.NewRequest(http.MethodGet, "/casbin/permissions?role=nonexistent", nil)
	w := httptest.NewRecorder()
	h.GetRolePermissions(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for role with no permissions, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)
	count, _ := resp["count"].(float64)
	if count != 0 {
		t.Errorf("expected count=0 for role with no permissions, got %v", count)
	}
}

func TestCasbinHandler_GetRolePermissions_ResponseBodyFields(t *testing.T) {
	service := newTestCasbinService(t)
	h := &CasbinHandler{casbinService: service}

	req := httptest.NewRequest(http.MethodGet, "/casbin/permissions?role=admin", nil)
	w := httptest.NewRecorder()
	h.GetRolePermissions(w, req)

	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)

	if resp["role"] != "admin" {
		t.Errorf("expected role='admin' in response, got %v", resp["role"])
	}
	if _, ok := resp["permissions"]; !ok {
		t.Error("expected 'permissions' field in response")
	}
	if _, ok := resp["count"]; !ok {
		t.Error("expected 'count' field in response")
	}
}

func TestCasbinHandler_GetRolePermissions_MethodNotAllowed(t *testing.T) {
	h := setupCasbinHandler(t)

	req := httptest.NewRequest(http.MethodPost, "/casbin/permissions?role=admin", nil)
	w := httptest.NewRecorder()
	h.GetRolePermissions(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405 for non-GET, got %d", w.Code)
	}
}

// =========================
// GetAllPolicies
// =========================

func TestCasbinHandler_GetAllPolicies_Success(t *testing.T) {
	h := setupCasbinHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/casbin/policies", nil)
	w := httptest.NewRecorder()
	h.GetAllPolicies(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)

	if _, ok := resp["policies"]; !ok {
		t.Error("expected 'policies' field in response")
	}
	if _, ok := resp["count"]; !ok {
		t.Error("expected 'count' field in response")
	}
}

func TestCasbinHandler_GetAllPolicies_ContainsSeededPolicy(t *testing.T) {
	h := setupCasbinHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/casbin/policies", nil)
	w := httptest.NewRecorder()
	h.GetAllPolicies(w, req)

	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)

	count, _ := resp["count"].(float64)
	if count < 1 {
		t.Errorf("expected at least 1 policy (seeded), got %v", count)
	}
}

func TestCasbinHandler_GetAllPolicies_MethodNotAllowed(t *testing.T) {
	h := setupCasbinHandler(t)

	req := httptest.NewRequest(http.MethodPost, "/casbin/policies", nil)
	w := httptest.NewRecorder()
	h.GetAllPolicies(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405 for non-GET, got %d", w.Code)
	}
}

// =========================
// AddPolicy
// =========================

func TestCasbinHandler_AddPolicy(t *testing.T) {
	service := newTestCasbinService(t)
	h := &CasbinHandler{casbinService: service}

	payload := `{"role":"editor","resource":"posts","action":"create"}`
	req := httptest.NewRequest(http.MethodPost, "/casbin/policy", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req = adminCtx(req)
	w := httptest.NewRecorder()

	h.AddPolicy(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}

	policies := service.GetRolePermissions("editor")
	if len(policies) != 1 {
		t.Fatalf("expected 1 policy, got %d", len(policies))
	}
}

func TestCasbinHandler_AddPolicy_NotAdmin(t *testing.T) {
	service := newTestCasbinService(t)
	h := &CasbinHandler{casbinService: service}

	payload := `{"role":"editor","resource":"posts","action":"create"}`
	req := httptest.NewRequest(http.MethodPost, "/casbin/policy", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req = userCtx(req)
	w := httptest.NewRecorder()

	h.AddPolicy(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", w.Code)
	}
}

func TestCasbinHandler_AddPolicy_Conflict(t *testing.T) {
	service := newTestCasbinService(t)
	h := &CasbinHandler{casbinService: service}

	payload := `{"role":"admin","resource":"users","action":"read"}`
	req := httptest.NewRequest(http.MethodPost, "/casbin/policy", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req = adminCtx(req)
	w := httptest.NewRecorder()

	h.AddPolicy(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected status 409, got %d", w.Code)
	}
}

func TestCasbinHandler_AddPolicy_InvalidBody(t *testing.T) {
	h := setupCasbinHandler(t)

	req := httptest.NewRequest(http.MethodPost, "/casbin/policy", strings.NewReader("{invalid json"))
	req.Header.Set("Content-Type", "application/json")
	req = adminCtx(req)
	w := httptest.NewRecorder()

	h.AddPolicy(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid body, got %d", w.Code)
	}
}

func TestCasbinHandler_AddPolicy_MissingFields(t *testing.T) {
	h := setupCasbinHandler(t)

	// Role ada tapi resource dan action kosong
	payload := `{"role":"editor","resource":"","action":""}`
	req := httptest.NewRequest(http.MethodPost, "/casbin/policy", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req = adminCtx(req)
	w := httptest.NewRecorder()

	h.AddPolicy(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing fields, got %d", w.Code)
	}
}

func TestCasbinHandler_AddPolicy_MethodNotAllowed(t *testing.T) {
	h := setupCasbinHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/casbin/policy", nil)
	req = adminCtx(req)
	w := httptest.NewRecorder()

	h.AddPolicy(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405 for non-POST, got %d", w.Code)
	}
}

func TestCasbinHandler_AddPolicy_ResponseBodyFields(t *testing.T) {
	h := setupCasbinHandler(t)

	payload := `{"role":"viewer","resource":"reports","action":"read"}`
	req := httptest.NewRequest(http.MethodPost, "/casbin/policy", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req = adminCtx(req)
	w := httptest.NewRecorder()
	h.AddPolicy(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)

	if resp["message"] == "" {
		t.Error("expected 'message' field in response")
	}
	if _, ok := resp["policy"]; !ok {
		t.Error("expected 'policy' field in response")
	}
}

func TestCasbinHandler_AddPolicy_WithoutContext(t *testing.T) {
	h := setupCasbinHandler(t)

	// Tanpa inject context — getUserFromContext akan return ok=false
	// sehingga check admin dilewati dan request diproses
	payload := `{"role":"guest","resource":"docs","action":"read"}`
	req := httptest.NewRequest(http.MethodPost, "/casbin/policy", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	// Sengaja tidak set context
	w := httptest.NewRecorder()

	h.AddPolicy(w, req)

	// Tanpa context, ok=false → tidak diblok karena kondisi `if ok && userRole != "admin"`
	if w.Code != http.StatusCreated {
		t.Errorf("expected 201 when no context (guard is skipped), got %d", w.Code)
	}
}

// =========================
// BulkAddPolicies
// =========================

func TestCasbinHandler_BulkAddPolicies_NotAdmin(t *testing.T) {
	service := newTestCasbinService(t)
	h := &CasbinHandler{casbinService: service}

	payload := `{
		"role": "editor",
		"policies": [
			{"resource":"posts","action":"create"},
			{"resource":"posts","action":"update"}
		]
	}`

	req := httptest.NewRequest(http.MethodPost, "/casbin/policies/bulk", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req = userCtx(req)
	w := httptest.NewRecorder()

	h.BulkAddPolicies(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", w.Code)
	}
}

func TestCasbinHandler_BulkAddPolicies_AllExisting(t *testing.T) {
	service := newTestCasbinService(t)
	h := &CasbinHandler{casbinService: service}

	payload := `{
		"role": "admin",
		"policies": [
			{"resource":"users","action":"read"}
		]
	}`

	req := httptest.NewRequest(http.MethodPost, "/casbin/policies/bulk", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req = adminCtx(req)
	w := httptest.NewRecorder()

	h.BulkAddPolicies(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected status 409, got %d", w.Code)
	}
}

func TestCasbinHandler_BulkAddPolicies_Partial(t *testing.T) {
	service := newTestCasbinService(t)
	h := &CasbinHandler{casbinService: service}

	payload := `{
		"role": "admin",
		"policies": [
			{"resource":"users","action":"read"},
			{"resource":"users","action":"write"}
		]
	}`

	req := httptest.NewRequest(http.MethodPost, "/casbin/policies/bulk", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req = adminCtx(req)
	w := httptest.NewRecorder()

	h.BulkAddPolicies(w, req)

	if w.Code != http.StatusPartialContent {
		t.Fatalf("expected status 206, got %d", w.Code)
	}
}

func TestCasbinHandler_BulkAddPolicies_AllNew(t *testing.T) {
	h := setupCasbinHandler(t)

	payload := `{
		"role": "newrole",
		"policies": [
			{"resource":"reports","action":"read"},
			{"resource":"reports","action":"write"}
		]
	}`

	req := httptest.NewRequest(http.MethodPost, "/casbin/policies/bulk", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req = adminCtx(req)
	w := httptest.NewRecorder()

	h.BulkAddPolicies(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201 for all-new bulk policies, got %d", w.Code)
	}
}

func TestCasbinHandler_BulkAddPolicies_InvalidBody(t *testing.T) {
	h := setupCasbinHandler(t)

	req := httptest.NewRequest(http.MethodPost, "/casbin/policies/bulk", strings.NewReader("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	req = adminCtx(req)
	w := httptest.NewRecorder()

	h.BulkAddPolicies(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid body, got %d", w.Code)
	}
}

func TestCasbinHandler_BulkAddPolicies_MethodNotAllowed(t *testing.T) {
	h := setupCasbinHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/casbin/policies/bulk", nil)
	req = adminCtx(req)
	w := httptest.NewRecorder()

	h.BulkAddPolicies(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405 for non-POST, got %d", w.Code)
	}
}

func TestCasbinHandler_BulkAddPolicies_ResponseBodyFields(t *testing.T) {
	h := setupCasbinHandler(t)

	payload := `{
		"role": "analyst",
		"policies": [
			{"resource":"dashboard","action":"view"}
		]
	}`

	req := httptest.NewRequest(http.MethodPost, "/casbin/policies/bulk", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req = adminCtx(req)
	w := httptest.NewRecorder()
	h.BulkAddPolicies(w, req)

	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)

	if _, ok := resp["added"]; !ok {
		t.Error("expected 'added' field in response")
	}
	if _, ok := resp["summary"]; !ok {
		t.Error("expected 'summary' field in response")
	}
}

// =========================
// RemovePolicy
// =========================

func TestCasbinHandler_RemovePolicy(t *testing.T) {
	service := newTestCasbinService(t)
	h := &CasbinHandler{casbinService: service}

	policies := service.GetRolePermissions("admin")
	if len(policies) == 0 {
		t.Fatal("expected existing policy before removal")
	}

	payload := `{"role":"admin","resource":"users","action":"read"}`
	req := httptest.NewRequest(http.MethodDelete, "/casbin/policy", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req = adminCtx(req)
	w := httptest.NewRecorder()

	h.RemovePolicy(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	remaining := service.GetRolePermissions("admin")
	if len(remaining) != 0 {
		t.Fatalf("expected policy to be removed, remaining: %d", len(remaining))
	}
}

func TestCasbinHandler_RemovePolicy_NotFound(t *testing.T) {
	service := newTestCasbinService(t)
	h := &CasbinHandler{casbinService: service}

	payload := `{"role":"admin","resource":"posts","action":"delete"}`
	req := httptest.NewRequest(http.MethodDelete, "/casbin/policy", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req = adminCtx(req)
	w := httptest.NewRecorder()

	h.RemovePolicy(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", w.Code)
	}
}

func TestCasbinHandler_RemovePolicy_NotAdmin(t *testing.T) {
	service := newTestCasbinService(t)
	h := &CasbinHandler{casbinService: service}

	payload := `{"role":"admin","resource":"users","action":"read"}`
	req := httptest.NewRequest(http.MethodDelete, "/casbin/policy", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req = userCtx(req)
	w := httptest.NewRecorder()

	h.RemovePolicy(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", w.Code)
	}
}

func TestCasbinHandler_RemovePolicy_InvalidBody(t *testing.T) {
	h := setupCasbinHandler(t)

	req := httptest.NewRequest(http.MethodDelete, "/casbin/policy", strings.NewReader("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	req = adminCtx(req)
	w := httptest.NewRecorder()

	h.RemovePolicy(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid body, got %d", w.Code)
	}
}

func TestCasbinHandler_RemovePolicy_MethodNotAllowed(t *testing.T) {
	h := setupCasbinHandler(t)

	req := httptest.NewRequest(http.MethodPost, "/casbin/policy", nil)
	req = adminCtx(req)
	w := httptest.NewRecorder()

	h.RemovePolicy(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405 for non-DELETE, got %d", w.Code)
	}
}

func TestCasbinHandler_RemovePolicy_ResponseBodyFields(t *testing.T) {
	h := setupCasbinHandler(t)

	payload := `{"role":"admin","resource":"users","action":"read"}`
	req := httptest.NewRequest(http.MethodDelete, "/casbin/policy", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req = adminCtx(req)
	w := httptest.NewRecorder()
	h.RemovePolicy(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)

	if resp["message"] == "" {
		t.Error("expected 'message' field in response body")
	}
	removed, _ := resp["removed"].(bool)
	if !removed {
		t.Error("expected 'removed' to be true in response body")
	}
}

func TestCasbinHandler_RemovePolicy_PolicyRemovedFromService(t *testing.T) {
	service := newTestCasbinService(t)
	h := &CasbinHandler{casbinService: service}

	// Tambah policy baru dulu
	service.AddPolicy("testRole", "testResource", "testAction")

	payload := `{"role":"testRole","resource":"testResource","action":"testAction"}`
	req := httptest.NewRequest(http.MethodDelete, "/casbin/policy", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req = adminCtx(req)
	w := httptest.NewRecorder()
	h.RemovePolicy(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// Verifikasi benar-benar hilang dari service
	remaining := service.GetRolePermissions("testRole")
	if len(remaining) != 0 {
		t.Errorf("expected policy to be removed from service, got %d remaining", len(remaining))
	}
}
