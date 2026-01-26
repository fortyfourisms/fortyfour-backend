package services

import (
	"fmt"
	"fortyfour-backend/internal/models"
	"log"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/rollbar/rollbar-go"
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
	// Initialize GORM
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		rollbar.Error(err)
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Clean up invalid policies before initializing casbin
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

	// Initialize GORM Adapter
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		rollbar.Error(err)
		return nil, fmt.Errorf("failed to create casbin adapter: %w", err)
	}

	// Create enforcer with model file
	enforcer, err := casbin.NewEnforcer(modelPath, adapter)
	if err != nil {
		rollbar.Error(err)
		return nil, fmt.Errorf("failed to create casbin enforcer: %w", err)
	}

	// Load policies from database
	if err := enforcer.LoadPolicy(); err != nil {
		rollbar.Error(err)
		return nil, fmt.Errorf("failed to load policy: %w", err)
	}

	// Check if policies exist, if not add defaults
	allPolicies, err := enforcer.GetPolicy()
	if err != nil {
		rollbar.Error(err)
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
			rollbar.Error(err)
			log.Printf("Error adding default policies: %v", err)
		}

		if added && enforcer.GetAdapter() != nil {
			if err := enforcer.SavePolicy(); err != nil {
				rollbar.Error(err)
				log.Printf("Error saving default policies: %v", err)
			} else {
				log.Println("Default policies added successfully")
			}
		}
	} else {
		log.Printf("Loaded %d existing policies from database", len(allPolicies))
	}

	return &CasbinService{
		enforcer: enforcer,
	}, nil
}

//
// ======== TEST SUPPORT ========
//

// SetEnforcer allows injecting a custom enforcer (for unit tests)
func (s *CasbinService) SetEnforcer(e *casbin.Enforcer) {
	s.enforcer = e
}

//
// ======== CORE METHODS ========
//

// Enforce checks if a role has permission
func (s *CasbinService) Enforce(role, resource, action string) (bool, error) {
	ok, err := s.enforcer.Enforce(role, resource, action)
	return ok, err
}

// AddPolicy adds a new permission and saves to database
func (s *CasbinService) AddPolicy(role, resource, action string) (bool, error) {
	// Validate input
	if role == "" || resource == "" || action == "" {
		return false, fmt.Errorf("role, resource, and action cannot be empty")
	}

	// Check if policy already exists first
	hasPolicy, err := s.enforcer.HasPolicy([]string{role, resource, action})
	if err != nil {
		rollbar.Error(err)
		return false, fmt.Errorf("failed to check policy existence: %w", err)
	}

	if hasPolicy {
		return false, nil // Policy already exists, return false without error
	}

	// Add the policy
	added, err := s.enforcer.AddPolicy(role, resource, action)
	if err != nil {
		rollbar.Error(err)
		return false, fmt.Errorf("failed to add policy: %w", err)
	}

	if added && s.enforcer.GetAdapter() != nil {
		if err := s.enforcer.SavePolicy(); err != nil {
			rollbar.Error(err)
			return false, fmt.Errorf("failed to save policy: %w", err)
		}
	}

	return added, nil
}

// AddPolicies adds multiple policies at once (legacy method)
func (s *CasbinService) AddPolicies(policies [][]string) (bool, error) {
	added, err := s.enforcer.AddPolicies(policies)
	if err != nil {
		rollbar.Error(err)
		return false, fmt.Errorf("failed to add policies: %w", err)
	}

	if added && s.enforcer.GetAdapter() != nil {
		if err := s.enforcer.SavePolicy(); err != nil {
			rollbar.Error(err)
			return false, fmt.Errorf("failed to save policies: %w", err)
		}
	}

	return added, nil
}

// BulkAddPolicies adds multiple policies and returns detailed result
func (s *CasbinService) BulkAddPolicies(policies [][]string) (*BulkAddResult, error) {
	result := &BulkAddResult{
		Added:    []models.CasbinPolicy{},
		Existing: []models.CasbinPolicy{},
	}

	// Check each policy individually
	for _, policy := range policies {
		if len(policy) < 3 {
			continue
		}

		role := policy[0]
		resource := policy[1]
		action := policy[2]

		// Check if policy already exists
		hasPolicy, err := s.enforcer.HasPolicy(policy)
		if err != nil {
			rollbar.Error(err)
			return nil, fmt.Errorf("failed to check policy existence (%s, %s, %s): %w", role, resource, action, err)
		}

		if hasPolicy {
			// Policy already exists
			result.Existing = append(result.Existing, models.CasbinPolicy{
				Role:     role,
				Resource: resource,
				Action:   action,
			})
		} else {
			// Try to add the policy
			added, err := s.enforcer.AddPolicy(role, resource, action)
			if err != nil {
				rollbar.Error(err)
				return nil, fmt.Errorf("failed to add policy (%s, %s, %s): %w", role, resource, action, err)
			}

			if added {
				result.Added = append(result.Added, models.CasbinPolicy{
					Role:     role,
					Resource: resource,
					Action:   action,
				})
			} else {
				// In case AddPolicy returns false without error
				result.Existing = append(result.Existing, models.CasbinPolicy{
					Role:     role,
					Resource: resource,
					Action:   action,
				})
			}
		}
	}

	// Save to database if any policies were added
	if len(result.Added) > 0 && s.enforcer.GetAdapter() != nil {
		if err := s.enforcer.SavePolicy(); err != nil {
			rollbar.Error(err)
			return nil, fmt.Errorf("failed to save policies: %w", err)
		}
	}

	return result, nil
}

// RemovePolicy removes a permission
func (s *CasbinService) RemovePolicy(role, resource, action string) (bool, error) {
	removed, err := s.enforcer.RemovePolicy(role, resource, action)
	if err != nil {
		rollbar.Error(err)
		return false, fmt.Errorf("failed to remove policy: %w", err)
	}

	if removed && s.enforcer.GetAdapter() != nil {
		if err := s.enforcer.SavePolicy(); err != nil {
			rollbar.Error(err)
			return false, fmt.Errorf("failed to save policy: %w", err)
		}
	}

	return removed, nil
}

// GetRolePermissions gets all permissions for a specific role
func (s *CasbinService) GetRolePermissions(role string) []models.CasbinPolicy {
	filteredPolicies, err := s.enforcer.GetFilteredPolicy(0, role)
	if err != nil {
		rollbar.Error(err)
		log.Printf("Error getting filtered policy: %v", err)
		return []models.CasbinPolicy{}
	}

	policies := make([]models.CasbinPolicy, 0, len(filteredPolicies))
	for _, policy := range filteredPolicies {
		if len(policy) >= 3 {
			policies = append(policies, models.CasbinPolicy{
				Role:     policy[0],
				Resource: policy[1],
				Action:   policy[2],
			})
		}
	}

	return policies
}

// GetAllPolicies gets all policies
func (s *CasbinService) GetAllPolicies() []models.CasbinPolicy {
	allPolicies, err := s.enforcer.GetPolicy()
	if err != nil {
		rollbar.Error(err)
		log.Printf("Error getting all policies: %v", err)
		return []models.CasbinPolicy{}
	}

	policies := make([]models.CasbinPolicy, 0, len(allPolicies))
	for _, policy := range allPolicies {
		if len(policy) >= 3 {
			policies = append(policies, models.CasbinPolicy{
				Role:     policy[0],
				Resource: policy[1],
				Action:   policy[2],
			})
		}
	}

	return policies
}

// ReloadPolicy reloads policies from database
func (s *CasbinService) ReloadPolicy() error {
	return s.enforcer.LoadPolicy()
}

// GetEnforcer returns the casbin enforcer (for middleware)
func (s *CasbinService) GetEnforcer() *casbin.Enforcer {
	return s.enforcer
}