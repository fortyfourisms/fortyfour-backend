package models

import "time"

// UserStatus mendefinisikan status akun user
type UserStatus string

const (
	UserStatusAktif    UserStatus = "Aktif"
	UserStatusSuspend  UserStatus = "Suspend"
	UserStatusNonaktif UserStatus = "Nonaktif"
)

// PasswordExpiryDays adalah masa berlaku password (3 bulan = 90 hari)
const PasswordExpiryDays = 90

// PasswordExpiryWarnDays adalah batas peringatan sebelum expired (1 minggu = 7 hari)
const PasswordExpiryWarnDays = 7

// MaxLoginAttempts adalah batas maksimal percobaan login sebelum suspend
const MaxLoginAttempts = 5

type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Password     string    `json:"-"`
	Email        string    `json:"email"`
	RoleID       *string   `json:"role_id,omitempty"`
	RoleName     string    `json:"role_name,omitempty"`
	IDJabatan    *string   `json:"id_jabatan,omitempty"`
	JabatanName  *string   `json:"jabatan_name"`
	IDPerusahaan *string   `json:"id_perusahaan,omitempty"`
	FotoProfile  *string   `json:"foto_profile"`
	Banner       *string   `json:"banner"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	MFAEnabled bool    `json:"mfa_enabled"`
	MFASecret  *string `json:"-"`

	// Security fields
	Status            UserStatus `json:"status"`
	PasswordChangedAt time.Time  `json:"password_changed_at"`
	LoginAttempts     int        `json:"-"`
}

// IsPasswordExpired mengecek apakah password sudah melewati masa berlaku
func (u *User) IsPasswordExpired() bool {
	return time.Since(u.PasswordChangedAt) > time.Duration(PasswordExpiryDays)*24*time.Hour
}

// IsPasswordExpiringSoon mengecek apakah password akan expired dalam 7 hari ke depan
func (u *User) IsPasswordExpiringSoon() bool {
	elapsed := time.Since(u.PasswordChangedAt)
	warnThreshold := time.Duration(PasswordExpiryDays-PasswordExpiryWarnDays) * 24 * time.Hour
	expiryThreshold := time.Duration(PasswordExpiryDays) * 24 * time.Hour
	return elapsed >= warnThreshold && elapsed < expiryThreshold
}

// DaysUntilPasswordExpiry mengembalikan sisa hari sebelum password expired
func (u *User) DaysUntilPasswordExpiry() int {
	expiry := u.PasswordChangedAt.Add(time.Duration(PasswordExpiryDays) * 24 * time.Hour)
	remaining := time.Until(expiry)
	if remaining < 0 {
		return 0
	}
	return int(remaining.Hours() / 24)
}
