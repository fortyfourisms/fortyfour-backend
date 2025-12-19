package dto

type CreateRoleRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=50"`
	Description string `json:"description" validate:"max=255"`
}

type UpdateRoleRequest struct {
	Name        string `json:"name" validate:"omitempty,min=3,max=50"`
	Description string `json:"description" validate:"omitempty,max=255"`
}

type RoleResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
