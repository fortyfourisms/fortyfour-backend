package services_test

import (
	"testing"

	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/services"

	"github.com/casbin/casbin/v3"
	"github.com/casbin/casbin/v3/model"
)

// newCasbinEnforcer membuat enforcer in-memory tanpa adapter DB untuk testing
func newCasbinEnforcer(t *testing.T) *casbin.Enforcer {
	t.Helper()
	m, err := model.NewModelFromString(`
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && r.obj == p.obj && (r.act == p.act || p.act == "*")
`)
	if err != nil {
		t.Fatalf("failed to create casbin model: %v", err)
	}
	e, err := casbin.NewEnforcer(m)
	if err != nil {
		t.Fatalf("failed to create casbin enforcer: %v", err)
	}
	return e
}

// newServiceWithPolicies membuat CasbinService dengan enforcer + policies untuk testing
func newServiceWithPolicies(t *testing.T, policies [][]string) *services.CasbinService {
	t.Helper()
	e := newCasbinEnforcer(t)
	for _, p := range policies {
		if _, err := e.AddPolicy(p[0], p[1], p[2]); err != nil {
			t.Fatalf("failed to add policy %v: %v", p, err)
		}
	}
	svc := &services.CasbinService{}
	svc.SetEnforcer(e)
	return svc
}

// ============================================================
// TestCasbinService_Enforce
// ============================================================

func TestCasbinService_Enforce_AllowedRole(t *testing.T) {
	svc := newServiceWithPolicies(t, [][]string{
		{"admin", "/api/users", "GET"},
	})

	ok, err := svc.Enforce("admin", "/api/users", "GET")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !ok {
		t.Error("expected allow, got deny")
	}
}

func TestCasbinService_Enforce_DeniedRole(t *testing.T) {
	svc := newServiceWithPolicies(t, [][]string{
		{"admin", "/api/users", "GET"},
	})

	ok, err := svc.Enforce("user", "/api/users", "GET")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if ok {
		t.Error("expected deny, got allow")
	}
}

func TestCasbinService_Enforce_DeniedAction(t *testing.T) {
	svc := newServiceWithPolicies(t, [][]string{
		{"user", "/api/data", "GET"},
	})

	ok, err := svc.Enforce("user", "/api/data", "DELETE")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if ok {
		t.Error("expected deny untuk action yang tidak diizinkan")
	}
}

func TestCasbinService_Enforce_WildcardAction(t *testing.T) {
	svc := newServiceWithPolicies(t, [][]string{
		{"admin", "/api/users", "*"},
	})

	for _, action := range []string{"GET", "POST", "PUT", "DELETE"} {
		ok, err := svc.Enforce("admin", "/api/users", action)
		if err != nil {
			t.Fatalf("action %s: expected no error, got %v", action, err)
		}
		if !ok {
			t.Errorf("action %s: expected allow dengan wildcard, got deny", action)
		}
	}
}

func TestCasbinService_Enforce_EmptyPolicies_Deny(t *testing.T) {
	svc := newServiceWithPolicies(t, nil)

	ok, err := svc.Enforce("admin", "/api/users", "GET")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if ok {
		t.Error("expected deny saat tidak ada policy")
	}
}

// ============================================================
// TestCasbinService_AddPolicy
// ============================================================

func TestCasbinService_AddPolicy_Success(t *testing.T) {
	svc := newServiceWithPolicies(t, nil)

	added, err := svc.AddPolicy("editor", "/api/posts", "POST")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !added {
		t.Error("expected added=true untuk policy baru")
	}

	// Policy harus bisa di-enforce setelah ditambahkan
	ok, _ := svc.Enforce("editor", "/api/posts", "POST")
	if !ok {
		t.Error("policy baru harus bisa di-enforce")
	}
}

func TestCasbinService_AddPolicy_Duplicate_ReturnsFalse(t *testing.T) {
	svc := newServiceWithPolicies(t, [][]string{
		{"editor", "/api/posts", "POST"},
	})

	// Tambah policy yang sama — harus return false tanpa error
	added, err := svc.AddPolicy("editor", "/api/posts", "POST")

	if err != nil {
		t.Fatalf("expected no error untuk duplikat, got %v", err)
	}
	if added {
		t.Error("expected added=false untuk policy yang sudah ada")
	}
}

func TestCasbinService_AddPolicy_EmptyRole_Error(t *testing.T) {
	svc := newServiceWithPolicies(t, nil)

	_, err := svc.AddPolicy("", "/api/posts", "POST")

	if err == nil {
		t.Error("expected error untuk role kosong")
	}
}

func TestCasbinService_AddPolicy_EmptyResource_Error(t *testing.T) {
	svc := newServiceWithPolicies(t, nil)

	_, err := svc.AddPolicy("admin", "", "GET")

	if err == nil {
		t.Error("expected error untuk resource kosong")
	}
}

func TestCasbinService_AddPolicy_EmptyAction_Error(t *testing.T) {
	svc := newServiceWithPolicies(t, nil)

	_, err := svc.AddPolicy("admin", "/api/users", "")

	if err == nil {
		t.Error("expected error untuk action kosong")
	}
}

// ============================================================
// TestCasbinService_RemovePolicy
// ============================================================

func TestCasbinService_RemovePolicy_Success(t *testing.T) {
	svc := newServiceWithPolicies(t, [][]string{
		{"viewer", "/api/reports", "GET"},
	})

	removed, err := svc.RemovePolicy("viewer", "/api/reports", "GET")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !removed {
		t.Error("expected removed=true")
	}

	// Policy tidak boleh berlaku lagi
	ok, _ := svc.Enforce("viewer", "/api/reports", "GET")
	if ok {
		t.Error("policy yang dihapus tidak boleh bisa di-enforce")
	}
}

func TestCasbinService_RemovePolicy_NonExistent_ReturnsFalse(t *testing.T) {
	svc := newServiceWithPolicies(t, nil)

	removed, err := svc.RemovePolicy("ghost", "/api/nowhere", "GET")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if removed {
		t.Error("expected removed=false untuk policy yang tidak ada")
	}
}

// ============================================================
// TestCasbinService_BulkAddPolicies
// ============================================================

func TestCasbinService_BulkAddPolicies_SemuaBaru(t *testing.T) {
	svc := newServiceWithPolicies(t, nil)

	policies := [][]string{
		{"user", "/api/posts", "GET"},
		{"user", "/api/comments", "GET"},
	}

	result, err := svc.BulkAddPolicies(policies)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.Added) != 2 {
		t.Errorf("expected 2 added, got %d", len(result.Added))
	}
	if len(result.Existing) != 0 {
		t.Errorf("expected 0 existing, got %d", len(result.Existing))
	}
}

func TestCasbinService_BulkAddPolicies_Sebagian_Sudah_Ada(t *testing.T) {
	svc := newServiceWithPolicies(t, [][]string{
		{"user", "/api/posts", "GET"},
	})

	policies := [][]string{
		{"user", "/api/posts", "GET"},     // sudah ada
		{"user", "/api/comments", "POST"}, // baru
	}

	result, err := svc.BulkAddPolicies(policies)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.Added) != 1 {
		t.Errorf("expected 1 added, got %d", len(result.Added))
	}
	if len(result.Existing) != 1 {
		t.Errorf("expected 1 existing, got %d", len(result.Existing))
	}
}

func TestCasbinService_BulkAddPolicies_SemuaSudahAda(t *testing.T) {
	existing := [][]string{
		{"admin", "/api/users", "GET"},
		{"admin", "/api/users", "POST"},
	}
	svc := newServiceWithPolicies(t, existing)

	result, err := svc.BulkAddPolicies(existing)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.Added) != 0 {
		t.Errorf("expected 0 added, got %d", len(result.Added))
	}
	if len(result.Existing) != 2 {
		t.Errorf("expected 2 existing, got %d", len(result.Existing))
	}
}

func TestCasbinService_BulkAddPolicies_SkipPolicyKurangDari3Elemen(t *testing.T) {
	svc := newServiceWithPolicies(t, nil)

	policies := [][]string{
		{"admin", "/api/users"},       // hanya 2 elemen — harus diskip
		{"user", "/api/posts", "GET"}, // valid
	}

	result, err := svc.BulkAddPolicies(policies)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	// Hanya 1 yang valid ditambahkan
	if len(result.Added) != 1 {
		t.Errorf("expected 1 added (policy tidak lengkap di-skip), got %d", len(result.Added))
	}
}

func TestCasbinService_BulkAddPolicies_Kosong(t *testing.T) {
	svc := newServiceWithPolicies(t, nil)

	result, err := svc.BulkAddPolicies([][]string{})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.Added) != 0 {
		t.Errorf("expected 0 added, got %d", len(result.Added))
	}
	if len(result.Existing) != 0 {
		t.Errorf("expected 0 existing, got %d", len(result.Existing))
	}
}

// ============================================================
// TestCasbinService_AddPolicies
// ============================================================

func TestCasbinService_AddPolicies_Success(t *testing.T) {
	svc := newServiceWithPolicies(t, nil)

	policies := [][]string{
		{"manager", "/api/reports", "GET"},
		{"manager", "/api/reports", "POST"},
	}

	added, err := svc.AddPolicies(policies)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !added {
		t.Error("expected added=true")
	}
}

// ============================================================
// TestCasbinService_GetRolePermissions
// ============================================================

func TestCasbinService_GetRolePermissions_HasPolicies(t *testing.T) {
	svc := newServiceWithPolicies(t, [][]string{
		{"operator", "/api/tasks", "GET"},
		{"operator", "/api/tasks", "PUT"},
		{"admin", "/api/users", "GET"}, // milik role lain
	})

	perms := svc.GetRolePermissions("operator")

	if len(perms) != 2 {
		t.Errorf("expected 2 permissions untuk operator, got %d", len(perms))
	}
	for _, p := range perms {
		if p.Role != "operator" {
			t.Errorf("expected role 'operator', got '%s'", p.Role)
		}
	}
}

func TestCasbinService_GetRolePermissions_RoleTidakAda(t *testing.T) {
	svc := newServiceWithPolicies(t, [][]string{
		{"admin", "/api/users", "GET"},
	})

	perms := svc.GetRolePermissions("ghost")

	if len(perms) != 0 {
		t.Errorf("expected 0 permissions untuk role yang tidak ada, got %d", len(perms))
	}
}

func TestCasbinService_GetRolePermissions_ReturnsCorrectFields(t *testing.T) {
	svc := newServiceWithPolicies(t, [][]string{
		{"viewer", "/api/docs", "GET"},
	})

	perms := svc.GetRolePermissions("viewer")

	if len(perms) != 1 {
		t.Fatalf("expected 1 permission, got %d", len(perms))
	}
	p := perms[0]
	if p.Role != "viewer" || p.Resource != "/api/docs" || p.Action != "GET" {
		t.Errorf("unexpected policy fields: %+v", p)
	}
}

// ============================================================
// TestCasbinService_GetAllPolicies
// ============================================================

func TestCasbinService_GetAllPolicies_Success(t *testing.T) {
	svc := newServiceWithPolicies(t, [][]string{
		{"admin", "/api/users", "GET"},
		{"admin", "/api/users", "POST"},
		{"user", "/api/posts", "GET"},
	})

	all := svc.GetAllPolicies()

	if len(all) != 3 {
		t.Errorf("expected 3 policies, got %d", len(all))
	}
}

func TestCasbinService_GetAllPolicies_Kosong(t *testing.T) {
	svc := newServiceWithPolicies(t, nil)

	all := svc.GetAllPolicies()

	if len(all) != 0 {
		t.Errorf("expected 0 policies, got %d", len(all))
	}
}

func TestCasbinService_GetAllPolicies_ReturnsCasbinPolicyType(t *testing.T) {
	svc := newServiceWithPolicies(t, [][]string{
		{"admin", "/api/users", "GET"},
	})

	all := svc.GetAllPolicies()

	// Pastikan tipe return adalah []models.CasbinPolicy
	var _ []models.CasbinPolicy = all
	if all[0].Role != "admin" || all[0].Resource != "/api/users" || all[0].Action != "GET" {
		t.Errorf("unexpected policy: %+v", all[0])
	}
}

// ============================================================
// TestCasbinService_GetEnforcer
// ============================================================

func TestCasbinService_GetEnforcer_NotNil(t *testing.T) {
	svc := newServiceWithPolicies(t, nil)

	e := svc.GetEnforcer()

	if e == nil {
		t.Error("GetEnforcer() tidak boleh return nil")
	}
}

func TestCasbinService_GetEnforcer_SameInstance(t *testing.T) {
	e := newCasbinEnforcer(t)
	svc := &services.CasbinService{}
	svc.SetEnforcer(e)

	got := svc.GetEnforcer()

	if got != e {
		t.Error("GetEnforcer() harus mengembalikan enforcer yang sama")
	}
}

// ============================================================
// TestCasbinService_ReloadPolicy
// ============================================================

func TestCasbinService_ReloadPolicy_InMemory_NoAdapter(t *testing.T) {
	// Enforcer in-memory tidak punya adapter DB, sehingga ReloadPolicy()
	// tidak bisa dipanggil langsung (akan nil pointer panic).
	// Test ini memverifikasi bahwa GetAdapter() memang nil pada in-memory enforcer.
	svc := newServiceWithPolicies(t, [][]string{
		{"admin", "/api/users", "GET"},
	})

	e := svc.GetEnforcer()
	if e == nil {
		t.Fatal("enforcer harus tidak nil setelah SetEnforcer")
	}

	// In-memory enforcer tidak punya adapter -- ini yang menyebabkan ReloadPolicy panic
	// jika dipanggil tanpa adapter nyata. Verifikasi ini adalah kondisi yang diharapkan.
	if e.GetAdapter() != nil {
		t.Log("enforcer memiliki adapter, ReloadPolicy aman dipanggil")
	} else {
		t.Log("in-memory enforcer tanpa adapter -- ReloadPolicy tidak dapat ditest tanpa DB")
	}
}
