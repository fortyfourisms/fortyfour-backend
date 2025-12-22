package dto

type CreateUserRequest struct {
	Username  string  `json:"username" validate:"required,min=3,max=50"`
	Password  string  `json:"password" validate:"required,min=6"`
	Email     string  `json:"email" validate:"required,email"`
	RoleID    *string `json:"role_id"`
	IDJabatan *string `json:"id_jabatan"`
}

type UpdateUserRequest struct {
	Username  *string `json:"username" validate:"omitempty,min=3,max=50"`
	Email     *string `json:"email" validate:"omitempty,email"`
	RoleID    *string `json:"role_id"`
	IDJabatan *string `json:"id_jabatan"`
}

type UpdateUserPasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required,min=6"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

type UserResponse struct {
	ID          string  `json:"id"`
	Username    string  `json:"username"`
	Email       string  `json:"email"`
	RoleID      *string `json:"role_id"`
	RoleName    string  `json:"role_name"`
	IDJabatan   *string `json:"id_jabatan"`
	FotoProfile *string `json:"foto_profile"`
	Banner      *string `json:"banner"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}
