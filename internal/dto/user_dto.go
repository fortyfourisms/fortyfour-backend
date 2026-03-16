package dto

type CreateUserRequest struct {
	Username  string  `json:"username" validate:"required,min=3,max=50"`
	Password  string  `json:"password" validate:"required,min=8"`
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

// UpdateMeRequest adalah DTO untuk user memperbarui data dirinya sendiri.
// Tidak mengizinkan perubahan role_id (hanya admin yang bisa).
type UpdateMeRequest struct {
	Username  *string `json:"username" validate:"omitempty,min=3,max=50"`
	Email     *string `json:"email" validate:"omitempty,email"`
	IDJabatan *string `json:"id_jabatan"`
}

type UpdateUserPasswordRequest struct {
	OldPassword        string `json:"old_password" validate:"required,min=8"`
	NewPassword        string `json:"new_password" validate:"required,min=8"`
	ConfirmNewPassword string `json:"confirm_new_password" validate:"required,min=8,eqfield=NewPassword"`
}

type UserResponse struct {
	ID           string  `json:"id"`
	Username     string  `json:"username"`
	Email        string  `json:"email"`
	RoleID       *string `json:"role_id"`
	RoleName     string  `json:"role_name"`
	IDJabatan    *string `json:"id_jabatan"`
	JabatanName  *string `json:"jabatan_name"`
	IDPerusahaan *string `json:"id_perusahaan"`
	FotoProfile  *string `json:"foto_profile"`
	Banner       *string `json:"banner"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}
