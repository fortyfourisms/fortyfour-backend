package dto

type RegisterRequest struct {
	Username  string  `json:"username" validate:"required,min=3,max=50"`
	Password  string  `json:"password" validate:"required,min=8"`
	Email     string  `json:"email" validate:"required,email"`
	RoleID    *string `json:"role_id,omitempty"`
	IDJabatan *string `json:"id_jabatan,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type ErrorResponse struct {
	Message string `json:"message" example:"invalid credentials"`
}

type MessageResponse struct {
	Message string `json:"message" example:"Logged out successfully"`
}
