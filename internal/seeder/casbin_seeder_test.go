package seeder

import (
	"testing"

	"fortyfour-backend/internal/services"

	"github.com/casbin/casbin/v3"
	"github.com/casbin/casbin/v3/model"
)

func newSeederTestEnforcer(t *testing.T) *casbin.Enforcer {
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

	enforcer, err := casbin.NewEnforcer(m)
	if err != nil {
		t.Fatalf("failed to create casbin enforcer: %v", err)
	}

	return enforcer
}

func newSeederTestService(t *testing.T, policies [][]string) *services.CasbinService {
	t.Helper()

	enforcer := newSeederTestEnforcer(t)
	for _, p := range policies {
		if _, err := enforcer.AddPolicy(p[0], p[1], p[2]); err != nil {
			t.Fatalf("failed to add initial policy %v: %v", p, err)
		}
	}

	svc := &services.CasbinService{}
	svc.SetEnforcer(enforcer)
	return svc
}

func hasPolicy(t *testing.T, svc *services.CasbinService, role, resource, action string) bool {
	t.Helper()

	ok, err := svc.Enforce(role, resource, action)
	if err != nil {
		t.Fatalf("failed to enforce policy (%s, %s, %s): %v", role, resource, action, err)
	}

	return ok
}

func TestSeedCasbinPolicies_AddsDefaultPoliciesAndCleansStaleNonStaff(t *testing.T) {
	svc := newSeederTestService(t, [][]string{
		{"user", "/api/legacy", "DELETE"},
		{"staff", "/api/custom-staff", "PATCH"},
		{"admin", "/api/custom-admin", "PATCH"},
	})

	SeedCasbinPolicies(svc)

	if !hasPolicy(t, svc, "admin", "/api/*", "GET") {
		t.Fatal("expected admin wildcard policy to be seeded")
	}

	if !hasPolicy(t, svc, "user", "/api/kelas", "GET") {
		t.Fatal("expected one of default user policies to be seeded")
	}

	if hasPolicy(t, svc, "user", "/api/legacy", "DELETE") {
		t.Fatal("expected stale non-default policy to be removed")
	}

	if !hasPolicy(t, svc, "staff", "/api/custom-staff", "PATCH") {
		t.Fatal("expected custom staff policy to be preserved")
	}

	if !hasPolicy(t, svc, "admin", "/api/custom-admin", "PATCH") {
		t.Fatal("expected custom admin policy to be preserved")
	}
}

func TestSeedStaffInitialPolicies_SeedsWhenStaffHasNoPolicies(t *testing.T) {
	svc := newSeederTestService(t, nil)

	added := seedStaffInitialPolicies(svc)

	if added != len(staffInitialPolicies) {
		t.Fatalf("expected %d staff policies added, got %d", len(staffInitialPolicies), added)
	}

	if !hasPolicy(t, svc, "staff", "/api/se", "GET") {
		t.Fatal("expected staff SE policy to be seeded")
	}

	if !hasPolicy(t, svc, "staff", "/api/dashboard/summary", "GET") {
		t.Fatal("expected staff dashboard policy to be seeded")
	}
}

func TestSeedStaffInitialPolicies_SkipsWhenStaffAlreadyHasPolicies(t *testing.T) {
	svc := newSeederTestService(t, [][]string{
		{"staff", "/api/custom", "GET"},
	})

	added := seedStaffInitialPolicies(svc)

	if added != 0 {
		t.Fatalf("expected no staff policy added, got %d", added)
	}

	policies := svc.GetRolePermissions("staff")
	if len(policies) != 1 {
		t.Fatalf("expected existing staff policies to remain untouched, got %d", len(policies))
	}

	if !hasPolicy(t, svc, "staff", "/api/custom", "GET") {
		t.Fatal("expected existing staff policy to remain")
	}

	if hasPolicy(t, svc, "staff", "/api/se", "GET") {
		t.Fatal("expected initial staff policies to be skipped when staff already has policies")
	}
}

func TestCleanupStalePolicies_RemovesOnlyNonDefaultNonAdminNonStaff(t *testing.T) {
	svc := newSeederTestService(t, [][]string{
		{"user", "/api/kelas", "GET"},
		{"user", "/api/legacy", "DELETE"},
		{"staff", "/api/custom-staff", "PATCH"},
		{"admin", "/api/custom-admin", "PATCH"},
	})

	removed := cleanupStalePolicies(svc)

	if removed != 1 {
		t.Fatalf("expected 1 stale policy removed, got %d", removed)
	}

	if hasPolicy(t, svc, "user", "/api/legacy", "DELETE") {
		t.Fatal("expected stale user policy to be removed")
	}

	if !hasPolicy(t, svc, "user", "/api/kelas", "GET") {
		t.Fatal("expected valid default user policy to remain")
	}

	if !hasPolicy(t, svc, "staff", "/api/custom-staff", "PATCH") {
		t.Fatal("expected staff policy to remain")
	}

	if !hasPolicy(t, svc, "admin", "/api/custom-admin", "PATCH") {
		t.Fatal("expected admin policy to remain")
	}
}
