package rbac

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
)

func NewEnforcer(modelPath, policyPath string) (*casbin.Enforcer, error) {
	// Load model from file
	m, err := model.NewModelFromFile(modelPath)
	if err != nil {
		return nil, err
	}

	// Create file adapter for policy
	a := fileadapter.NewAdapter(policyPath)

	// Create enforcer
	enforcer, err := casbin.NewEnforcer(m, a)
	if err != nil {
		return nil, err
	}

	// Load policy from file
	if err := enforcer.LoadPolicy(); err != nil {
		return nil, err
	}

	return enforcer, nil
}
