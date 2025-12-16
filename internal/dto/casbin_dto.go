package dto

type AddPolicyRequest struct {
	Role     string `json:"role" validate:"required"`
	Resource string `json:"resource" validate:"required"`
	Action   string `json:"action" validate:"required"`
}

type RemovePolicyRequest struct {
	Role     string `json:"role" validate:"required"`
	Resource string `json:"resource" validate:"required"`
	Action   string `json:"action" validate:"required"`
}

type BulkAddPolicyRequest struct {
	Role     string             `json:"role" validate:"required"`
	Policies []PolicyDefinition `json:"policies" validate:"required,dive"`
}

type PolicyDefinition struct {
	Resource string `json:"resource" validate:"required"`
	Action   string `json:"action" validate:"required"`
}

type GetRolePermissionsRequest struct {
	Role string `json:"role" validate:"required"`
}
