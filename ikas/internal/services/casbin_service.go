package services

import (
	"fmt"
	"ikas/internal/models"

	"github.com/casbin/casbin/v3"
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

	return &CasbinService{
		enforcer: enforcer,
	}, nil
}

// Enforce checks if a role has permission
func (s *CasbinService) Enforce(role, resource, action string) (bool, error) {
	ok, err := s.enforcer.Enforce(role, resource, action)
	return ok, err
}

// GetEnforcer returns the casbin enforcer (for middleware)
func (s *CasbinService) GetEnforcer() *casbin.Enforcer {
	return s.enforcer
}
