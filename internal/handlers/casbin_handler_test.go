package handlers

import (
	"context"
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

func TestCasbinHandler_AddPolicy(t *testing.T) {
	service := newTestCasbinService(t)
	h := &CasbinHandler{casbinService: service}

	payload := `{"role":"editor","resource":"posts","action":"create"}`
	req := httptest.NewRequest(http.MethodPost, "/casbin/policy", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	// inject context required by handler
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "test-user-id")
	ctx = context.WithValue(ctx, middleware.Role, "admin")
	req = req.WithContext(ctx)

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

	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "test-user-id")
	ctx = context.WithValue(ctx, middleware.Role, "user") // bukan admin
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	h.AddPolicy(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", w.Code)
	}
}

func TestCasbinHandler_AddPolicy_Conflict(t *testing.T) {
	service := newTestCasbinService(t)
	h := &CasbinHandler{casbinService: service}

	// policy ini SUDAH ADA dari newTestCasbinService
	payload := `{"role":"admin","resource":"users","action":"read"}`
	req := httptest.NewRequest(http.MethodPost, "/casbin/policy", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "test-user-id")
	ctx = context.WithValue(ctx, middleware.Role, "admin")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	h.AddPolicy(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected status 409, got %d", w.Code)
	}
}

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

	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "test-user-id")
	ctx = context.WithValue(ctx, middleware.Role, "user") // bukan admin
	req = req.WithContext(ctx)

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

	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "test-user-id")
	ctx = context.WithValue(ctx, middleware.Role, "admin")
	req = req.WithContext(ctx)

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

	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "test-user-id")
	ctx = context.WithValue(ctx, middleware.Role, "admin")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	h.BulkAddPolicies(w, req)

	if w.Code != http.StatusPartialContent {
		t.Fatalf("expected status 206, got %d", w.Code)
	}
}

func TestCasbinHandler_RemovePolicy(t *testing.T) {
	service := newTestCasbinService(t)
	h := &CasbinHandler{casbinService: service}

	// pastikan policy ada sebelum dihapus
	policies := service.GetRolePermissions("admin")
	if len(policies) == 0 {
		t.Fatal("expected existing policy before removal")
	}

	payload := `{"role":"admin","resource":"users","action":"read"}`
	req := httptest.NewRequest(http.MethodDelete, "/casbin/policy", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	// inject context
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "test-user-id")
	ctx = context.WithValue(ctx, middleware.Role, "admin")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	h.RemovePolicy(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	// pastikan policy benar-benar terhapus
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

	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "test-user-id")
	ctx = context.WithValue(ctx, middleware.Role, "admin")
	req = req.WithContext(ctx)

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

	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "test-user-id")
	ctx = context.WithValue(ctx, middleware.Role, "user") // bukan admin
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	h.RemovePolicy(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", w.Code)
	}
}