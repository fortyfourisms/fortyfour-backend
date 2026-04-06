package repository

import (
	"database/sql"
	"errors"
	"fortyfour-backend/internal/models"

	"github.com/google/uuid"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	// Default role = user
	if user.RoleID == nil || *user.RoleID == "" {
		var defaultRoleID string
		err := r.db.QueryRow("SELECT id FROM roles WHERE name = 'user'").Scan(&defaultRoleID)
		if err == nil {
			user.RoleID = &defaultRoleID
		}
	}

	query := `
		INSERT INTO users (
			id, username, display_name, password, email, role_id, id_jabatan, id_perusahaan,
			mfa_enabled, mfa_secret, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, false, NULL, NOW(), NOW())
	`

	_, err := r.db.Exec(
		query,
		user.ID,
		user.Username,
		user.DisplayName,
		user.Password,
		user.Email,
		user.RoleID,
		user.IDJabatan,
		user.IDPerusahaan,
	)
	return err
}

// selectUserColumns adalah SELECT query yang dipakai bersama oleh FindByID, FindByUsername, FindByEmail
const selectUserColumns = `
	SELECT
		u.id, u.username, u.display_name, u.password, u.email,
		u.role_id, r.name AS role_name,
		u.id_jabatan, j.nama_jabatan,
		u.id_perusahaan,
		u.foto_profile, u.banner,
		u.mfa_enabled, u.mfa_secret,
		u.status, u.password_changed_at, u.login_attempts,
		u.created_at, u.updated_at
	FROM users u
	LEFT JOIN roles r ON u.role_id = r.id
	LEFT JOIN jabatan j ON u.id_jabatan = j.id
`

// scanFullUser melakukan Scan dari row ke *models.User
func scanFullUser(row *sql.Row) (*models.User, error) {
	user := &models.User{}
	var (
		displayName  sql.NullString
		roleID       sql.NullString
		roleName     sql.NullString
		idJabatan    sql.NullString
		jabatanName  sql.NullString
		idPerusahaan sql.NullString
		fotoProfile  sql.NullString
		banner       sql.NullString
		mfaSecret    sql.NullString
	)

	err := row.Scan(
		&user.ID,
		&user.Username,
		&displayName,
		&user.Password,
		&user.Email,
		&roleID,
		&roleName,
		&idJabatan,
		&jabatanName,
		&idPerusahaan,
		&fotoProfile,
		&banner,
		&user.MFAEnabled,
		&mfaSecret,
		&user.Status,
		&user.PasswordChangedAt,
		&user.LoginAttempts,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	if displayName.Valid {
		tmp := displayName.String
		user.DisplayName = &tmp
	}
	if roleID.Valid {
		user.RoleID = &roleID.String
	}
	if roleName.Valid {
		user.RoleName = roleName.String
	}
	if idJabatan.Valid {
		user.IDJabatan = &idJabatan.String
	}
	if jabatanName.Valid {
		tmp := jabatanName.String
		user.JabatanName = &tmp
	}
	if idPerusahaan.Valid {
		tmp := idPerusahaan.String
		user.IDPerusahaan = &tmp
	}
	if fotoProfile.Valid {
		tmp := fotoProfile.String
		user.FotoProfile = &tmp
	}
	if banner.Valid {
		tmp := banner.String
		user.Banner = &tmp
	}
	if mfaSecret.Valid {
		tmp := mfaSecret.String
		user.MFASecret = &tmp
	}

	return user, nil
}

func (r *UserRepository) FindByID(id string) (*models.User, error) {
	row := r.db.QueryRow(selectUserColumns+` WHERE u.id = ?`, id)
	return scanFullUser(row)
}

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	row := r.db.QueryRow(selectUserColumns+` WHERE u.username = ?`, username)
	return scanFullUser(row)
}

// FindByEmail mencari user berdasarkan email — digunakan untuk login by email
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	row := r.db.QueryRow(selectUserColumns+` WHERE u.email = ?`, email)
	return scanFullUser(row)
}

func (r *UserRepository) Update(user *models.User) error {
	query := `
		UPDATE users
		SET username = ?, display_name = ?, email = ?, role_id = ?, id_jabatan = ?, updated_at = NOW()
		WHERE id = ?
	`
	_, err := r.db.Exec(
		query,
		user.Username,
		user.DisplayName,
		user.Email,
		user.RoleID,
		user.IDJabatan,
		user.ID,
	)
	return err
}

func (r *UserRepository) UpdateWithPhoto(user *models.User) error {
	query := `
		UPDATE users
		SET username = ?, display_name = ?, email = ?, role_id = ?, id_jabatan = ?,
		    foto_profile = ?, banner = ?, updated_at = NOW()
		WHERE id = ?
	`
	_, err := r.db.Exec(
		query,
		user.Username,
		user.DisplayName,
		user.Email,
		user.RoleID,
		user.IDJabatan,
		user.FotoProfile,
		user.Banner,
		user.ID,
	)
	return err
}

func (r *UserRepository) UpdatePassword(id, hashedPassword string) error {
	query := `
		UPDATE users
		SET password = ?, updated_at = NOW()
		WHERE id = ?
	`
	_, err := r.db.Exec(query, hashedPassword, id)
	return err
}

// SetMFA menyimpan secret and enabled flag
func (r *UserRepository) SetMFA(userID string, secret *string, enabled bool) error {
	query := `
		UPDATE users
		SET mfa_secret = ?, mfa_enabled = ?, updated_at = NOW()
		WHERE id = ?
	`
	_, err := r.db.Exec(query, secret, enabled, userID)
	return err
}

func (r *UserRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM users WHERE id = ?`, id)
	return err
}

// ExistsByPerusahaan cek apakah sudah ada user yang terdaftar di perusahaan ini
func (r *UserRepository) ExistsByPerusahaan(idPerusahaan string) (bool, error) {
	var count int
	err := r.db.QueryRow(
		`SELECT COUNT(*) FROM users WHERE id_perusahaan = ?`,
		idPerusahaan,
	).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *UserRepository) FindAll() ([]models.User, error) {
	query := `
		SELECT
			u.id, u.username, u.display_name, u.email,
			u.role_id, r.name AS role_name,
			u.id_jabatan, j.nama_jabatan,
			u.id_perusahaan,
			u.foto_profile, u.banner,
			u.mfa_enabled,
			u.status,
			u.created_at, u.updated_at
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		LEFT JOIN jabatan j ON u.id_jabatan = j.id
		ORDER BY u.created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		var displayName, roleID, roleName, idJabatan, jabatanName, idPerusahaan, fotoProfile, banner sql.NullString

		err := rows.Scan(
			&user.ID,
			&user.Username,
			&displayName,
			&user.Email,
			&roleID,
			&roleName,
			&idJabatan,
			&jabatanName,
			&idPerusahaan,
			&fotoProfile,
			&banner,
			&user.MFAEnabled,
			&user.Status,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if displayName.Valid {
			tmp := displayName.String
			user.DisplayName = &tmp
		}
		if roleID.Valid {
			user.RoleID = &roleID.String
		}
		if roleName.Valid {
			user.RoleName = roleName.String
		}
		if idJabatan.Valid {
			user.IDJabatan = &idJabatan.String
		}
		if jabatanName.Valid {
			tmp := jabatanName.String
			user.JabatanName = &tmp
		}
		if idPerusahaan.Valid {
			tmp := idPerusahaan.String
			user.IDPerusahaan = &tmp
		}
		if fotoProfile.Valid {
			tmp := fotoProfile.String
			user.FotoProfile = &tmp
		}
		if banner.Valid {
			tmp := banner.String
			user.Banner = &tmp
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) GetPasswordByID(id string) (string, error) {
	var password string
	err := r.db.QueryRow(`SELECT password FROM users WHERE id = ?`, id).Scan(&password)
	if err == sql.ErrNoRows {
		return "", errors.New("user not found")
	}
	return password, err
}

func (r *UserRepository) EmailExists(email string, excludeID *string) (bool, error) {
	var count int
	var query string
	var args []interface{}

	if excludeID != nil {
		query = `SELECT COUNT(*) FROM users WHERE email = ? AND id != ?`
		args = []interface{}{email, *excludeID}
	} else {
		query = `SELECT COUNT(*) FROM users WHERE email = ?`
		args = []interface{}{email}
	}

	err := r.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *UserRepository) UsernameExists(username string, excludeID *string) (bool, error) {
	var count int
	var query string
	var args []interface{}

	if excludeID != nil {
		query = `SELECT COUNT(*) FROM users WHERE username = ? AND id != ?`
		args = []interface{}{username, *excludeID}
	} else {
		query = `SELECT COUNT(*) FROM users WHERE username = ?`
		args = []interface{}{username}
	}

	err := r.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// UpdateStatus mengubah status akun user (Aktif, Suspend, Nonaktif)
func (r *UserRepository) UpdateStatus(userID string, status models.UserStatus) error {
	query := `UPDATE users SET status = ?, updated_at = NOW() WHERE id = ?`
	_, err := r.db.Exec(query, status, userID)
	return err
}

// IncrementLoginAttempts menambah counter gagal login dan mengembalikan nilai terbaru
func (r *UserRepository) IncrementLoginAttempts(userID string) (int, error) {
	query := `UPDATE users SET login_attempts = login_attempts + 1, updated_at = NOW() WHERE id = ?`
	if _, err := r.db.Exec(query, userID); err != nil {
		return 0, err
	}

	var attempts int
	if err := r.db.QueryRow(`SELECT login_attempts FROM users WHERE id = ?`, userID).Scan(&attempts); err != nil {
		return 0, err
	}
	return attempts, nil
}

// ResetLoginAttempts mereset counter gagal login ke 0
func (r *UserRepository) ResetLoginAttempts(userID string) error {
	query := `UPDATE users SET login_attempts = 0, updated_at = NOW() WHERE id = ?`
	_, err := r.db.Exec(query, userID)
	return err
}

// UpdatePasswordChangedAt memperbarui timestamp pergantian password
func (r *UserRepository) UpdatePasswordChangedAt(userID string) error {
	query := `UPDATE users SET password_changed_at = NOW(), updated_at = NOW() WHERE id = ?`
	_, err := r.db.Exec(query, userID)
	return err
}
