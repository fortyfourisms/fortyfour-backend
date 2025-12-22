package models

import "time"

type User struct {
	ID          string    `json:"id"`
	Username    string    `json:"username"`
	Password    string    `json:"-"`
	Email       string    `json:"email"`
	RoleID      *string   `json:"role_id,omitempty"`
	RoleName    string    `json:"role_name,omitempty"`
	IDJabatan   *string   `json:"id_jabatan,omitempty"`
	JabatanName *string   `json:"jabatan_name"`
	FotoProfile *string   `json:"foto_profile"`
	Banner      *string   `json:"banner"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
