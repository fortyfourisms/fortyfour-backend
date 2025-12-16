// internal/services/casbin_service.go
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

// NewCasbinService creates a new Casbin service with GORM adapter
func NewCasbinService(dsn, modelPath string) (*CasbinService, error) {
	// Initialize GORM
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Initialize GORM Adapter
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create casbin adapter: %w", err)
	}

	// Create enforcer with model file
	enforcer, err := casbin.NewEnforcer(modelPath, adapter)
	if err != nil {
		return nil, fmt.Errorf("failed to create casbin enforcer: %w", err)
	}

	// Load policies from database
	if err := enforcer.LoadPolicy(); err != nil {
		return nil, fmt.Errorf("failed to load policy: %w", err)
	}

	// Check if policies exist, if not add defaults
	// PERBAIKAN: GetPolicy() mengembalikan 2 nilai ([][]string, error)
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

		success, err := enforcer.AddPolicies(defaultPolicies)
		if err != nil {
			log.Printf("Error adding default policies: %v", err)
		}
		if success {
			if err := enforcer.SavePolicy(); err != nil {
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

// Enforce checks if a role has permission
func (s *CasbinService) Enforce(role, resource, action string) (bool, error) {
	ok, err := s.enforcer.Enforce(role, resource, action)
	return ok, err
}

// AddPolicy adds a new permission and saves to database
func (s *CasbinService) AddPolicy(role, resource, action string) (bool, error) {
	added, err := s.enforcer.AddPolicy(role, resource, action)
	if err != nil {
		return false, fmt.Errorf("failed to add policy: %w", err)
	}

	if added {
		if err := s.enforcer.SavePolicy(); err != nil {
			return false, fmt.Errorf("failed to save policy: %w", err)
		}
	}

	return added, nil
}

// AddPolicies adds multiple policies at once
func (s *CasbinService) AddPolicies(policies [][]string) (bool, error) {
	added, err := s.enforcer.AddPolicies(policies)
	if err != nil {
		return false, fmt.Errorf("failed to add policies: %w", err)
	}

	if added {
		if err := s.enforcer.SavePolicy(); err != nil {
			return false, fmt.Errorf("failed to save policies: %w", err)
		}
	}

	return added, nil
}

// RemovePolicy removes a permission
func (s *CasbinService) RemovePolicy(role, resource, action string) (bool, error) {
	removed, err := s.enforcer.RemovePolicy(role, resource, action)
	if err != nil {
		return false, fmt.Errorf("failed to remove policy: %w", err)
	}

	if removed {
		if err := s.enforcer.SavePolicy(); err != nil {
			return false, fmt.Errorf("failed to save policy: %w", err)
		}
	}

	return removed, nil
}

// GetRolePermissions gets all permissions for a specific role
func (s *CasbinService) GetRolePermissions(role string) []models.CasbinPolicy {
	// PERBAIKAN: GetFilteredPolicy() mengembalikan 2 nilai ([][]string, error)
	filteredPolicies, err := s.enforcer.GetFilteredPolicy(0, role)
	if err != nil {
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
	// PERBAIKAN: GetPolicy() mengembalikan 2 nilai ([][]string, error)
	allPolicies, err := s.enforcer.GetPolicy()
	if err != nil {
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
