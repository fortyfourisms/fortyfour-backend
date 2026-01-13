package services

import (
	"fmt"
	"fortyfour-backend/internal/models"
	"log"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type CasbinService struct {
	enforcer *casbin.Enforcer
}

// BulkAddResult contains the result of bulk add operation
type BulkAddResult struct {
	Added    []models.CasbinPolicy `json:"added"`
	Existing []models.CasbinPolicy `json:"existing"`
}

// NewCasbinService creates a new Casbin service with GORM adapter
func NewCasbinService(dsn, modelPath string) (*CasbinService, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Cleaning up invalid policies...")
	result := db.Exec(`
		DELETE FROM casbin_rule
		WHERE v0 = '' OR v1 = '' OR v2 = ''
		OR v0 IS NULL OR v1 IS NULL OR v2 IS NULL
	`)
	if result.Error != nil {
		log.Printf("Warning: failed to clean invalid policies: %v", result.Error)
	} else if result.RowsAffected > 0 {
		log.Printf("Cleaned up %d invalid policy records", result.RowsAffected)
	}

	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create casbin adapter: %w", err)
	}

	enforcer, err := casbin.NewEnforcer(modelPath, adapter)
	if err != nil {
		return nil, fmt.Errorf("failed to create casbin enforcer: %w", err)
	}

	if err := enforcer.LoadPolicy(); err != nil {
		return nil, fmt.Errorf("failed to load policy: %w", err)
	}

	allPolicies, err := enforcer.GetPolicy()
	if err != nil {
		return nil, fmt.Errorf("failed to get policies: %w", err)
	}

	if len(allPolicies) == 0 {
		log.Println("No policies found, adding default policies...")
		defaultPolicies := [][]string{
			{"admin", "*", "*"},
			{"user", "/api/posts", "GET"},
			{"user", "/api/posts/single", "GET"},
		}

		added, err := enforcer.AddPolicies(defaultPolicies)
		if err != nil {
			log.Printf("Error adding default policies: %v", err)
		}

		if added && enforcer.GetAdapter() != nil {
			if err := enforcer.SavePolicy(); err != nil {
				log.Printf("Error saving default policies: %v", err)
			}
		}
	}

	return &CasbinService{enforcer: enforcer}, nil
}

//
// ======== TEST SUPPORT ========
//

// SetEnforcer allows injecting a custom enforcer (unit tests)
func (s *CasbinService) SetEnforcer(e *casbin.Enforcer) {
	s.enforcer = e
}

//
// ======== CORE METHODS ========
//

// Enforce checks permission
func (s *CasbinService) Enforce(role, resource, action string) (bool, error) {
	return s.enforcer.Enforce(role, resource, action)
}

// AddPolicy adds a single policy
func (s *CasbinService) AddPolicy(role, resource, action string) (bool, error) {
	if role == "" || resource == "" || action == "" {
		return false, fmt.Errorf("role, resource, and action cannot be empty")
	}

	hasPolicy, err := s.enforcer.HasPolicy([]string{role, resource, action})
	if err != nil {
		return false, fmt.Errorf("failed to check policy existence: %w", err)
	}

	if hasPolicy {
		return false, nil
	}

	added, err := s.enforcer.AddPolicy(role, resource, action)
	if err != nil {
		return false, fmt.Errorf("failed to add policy: %w", err)
	}

	if added && s.enforcer.GetAdapter() != nil {
		if err := s.enforcer.SavePolicy(); err != nil {
			return false, fmt.Errorf("failed to save policy: %w", err)
		}
	}

	return added, nil
}

// AddPolicies adds multiple policies (legacy)
func (s *CasbinService) AddPolicies(policies [][]string) (bool, error) {
	added, err := s.enforcer.AddPolicies(policies)
	if err != nil {
		return false, fmt.Errorf("failed to add policies: %w", err)
	}

	if added && s.enforcer.GetAdapter() != nil {
		if err := s.enforcer.SavePolicy(); err != nil {
			return false, fmt.Errorf("failed to save policies: %w", err)
		}
	}

	return added, nil
}

// BulkAddPolicies adds multiple policies with result detail
func (s *CasbinService) BulkAddPolicies(policies [][]string) (*BulkAddResult, error) {
	result := &BulkAddResult{
		Added:    []models.CasbinPolicy{},
		Existing: []models.CasbinPolicy{},
	}

	for _, policy := range policies {
		if len(policy) < 3 {
			continue
		}

		role, resource, action := policy[0], policy[1], policy[2]

		hasPolicy, err := s.enforcer.HasPolicy(policy)
		if err != nil {
			return nil, fmt.Errorf("failed to check policy (%s,%s,%s): %w", role, resource, action, err)
		}

		if hasPolicy {
			result.Existing = append(result.Existing, models.CasbinPolicy{
				Role: role, Resource: resource, Action: action,
			})
			continue
		}

		added, err := s.enforcer.AddPolicy(role, resource, action)
		if err != nil {
			return nil, fmt.Errorf("failed to add policy (%s,%s,%s): %w", role, resource, action, err)
		}

		if added {
			result.Added = append(result.Added, models.CasbinPolicy{
				Role: role, Resource: resource, Action: action,
			})
		}
	}

	if len(result.Added) > 0 && s.enforcer.GetAdapter() != nil {
		if err := s.enforcer.SavePolicy(); err != nil {
			return nil, fmt.Errorf("failed to save policies: %w", err)
		}
	}

	return result, nil
}

// RemovePolicy removes a policy
func (s *CasbinService) RemovePolicy(role, resource, action string) (bool, error) {
	removed, err := s.enforcer.RemovePolicy(role, resource, action)
	if err != nil {
		return false, fmt.Errorf("failed to remove policy: %w", err)
	}

	if removed && s.enforcer.GetAdapter() != nil {
		if err := s.enforcer.SavePolicy(); err != nil {
			return false, fmt.Errorf("failed to save policy: %w", err)
		}
	}

	return removed, nil
}

// GetRolePermissions returns permissions for a role
func (s *CasbinService) GetRolePermissions(role string) []models.CasbinPolicy {
	policies, err := s.enforcer.GetFilteredPolicy(0, role)
	if err != nil {
		log.Printf("Error getting filtered policy: %v", err)
		return []models.CasbinPolicy{}
	}

	result := make([]models.CasbinPolicy, 0, len(policies))
	for _, p := range policies {
		if len(p) >= 3 {
			result = append(result, models.CasbinPolicy{
				Role: p[0], Resource: p[1], Action: p[2],
			})
		}
	}

	return result
}

// GetAllPolicies returns all policies
func (s *CasbinService) GetAllPolicies() []models.CasbinPolicy {
	policies, err := s.enforcer.GetPolicy()
	if err != nil {
		log.Printf("Error getting all policies: %v", err)
		return []models.CasbinPolicy{}
	}

	result := make([]models.CasbinPolicy, 0, len(policies))
	for _, p := range policies {
		if len(p) >= 3 {
			result = append(result, models.CasbinPolicy{
				Role: p[0], Resource: p[1], Action: p[2],
			})
		}
	}

	return result
}

// ReloadPolicy reloads policies
func (s *CasbinService) ReloadPolicy() error {
	return s.enforcer.LoadPolicy()
}

// GetEnforcer exposes enforcer (middleware)
func (s *CasbinService) GetEnforcer() *casbin.Enforcer {
	return s.enforcer
}
