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
			id, username, password, email, role_id, id_jabatan, id_perusahaan,
			mfa_enabled, mfa_secret, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, false, NULL, NOW(), NOW())
	`

	_, err := r.db.Exec(
		query,
		user.ID,
		user.Username,
		user.Password,
		user.Email,
		user.RoleID,
		user.IDJabatan,
		user.IDPerusahaan,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) FindByID(id string) (*models.User, error) {
	query := `
		SELECT
			u.id, u.username, u.password, u.email,
			u.role_id, r.name AS role_name,
			u.id_jabatan, j.nama_jabatan,
			u.id_perusahaan,
			u.foto_profile, u.banner,
			u.mfa_enabled, u.mfa_secret,
			u.created_at, u.updated_at
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		LEFT JOIN jabatan j ON u.id_jabatan = j.id
		WHERE u.id = ?
	`

	user := &models.User{}
	var (
		roleID        sql.NullString
		roleName      sql.NullString
		idJabatan     sql.NullString
		jabatanName   sql.NullString
		idPerusahaan  sql.NullString
		fotoProfile   sql.NullString
		banner        sql.NullString
		mfaSecret     sql.NullString
	)

	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
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
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
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

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	query := `
		SELECT
			u.id, u.username, u.password, u.email,
			u.role_id, r.name AS role_name,
			u.id_jabatan, j.nama_jabatan,
			u.id_perusahaan,
			u.foto_profile, u.banner,
			u.mfa_enabled, u.mfa_secret,
			u.created_at, u.updated_at
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		LEFT JOIN jabatan j ON u.id_jabatan = j.id
		WHERE u.username = ?
	`

	user := &models.User{}
	var (
		roleID        sql.NullString
		roleName      sql.NullString
		idJabatan     sql.NullString
		jabatanName   sql.NullString
		idPerusahaan  sql.NullString
		fotoProfile   sql.NullString
		banner        sql.NullString
		mfaSecret     sql.NullString
	)

	err := r.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
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
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
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

func (r *UserRepository) Update(user *models.User) error {
	query := `
		UPDATE users
		SET username = ?, email = ?, role_id = ?, id_jabatan = ?, updated_at = NOW()
		WHERE id = ?
	`
	_, err := r.db.Exec(
		query,
		user.Username,
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
		SET username = ?, email = ?, role_id = ?, id_jabatan = ?,
		    foto_profile = ?, banner = ?, updated_at = NOW()
		WHERE id = ?
	`
	_, err := r.db.Exec(
		query,
		user.Username,
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

func (r *UserRepository) FindAll() ([]models.User, error) {
	query := `
		SELECT
			u.id, u.username, u.email,
			u.role_id, r.name AS role_name,
			u.id_jabatan, j.nama_jabatan,
			u.id_perusahaan,
			u.foto_profile, u.banner,
			u.mfa_enabled,
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
		var roleID, roleName, idJabatan, jabatanName, idPerusahaan, fotoProfile, banner sql.NullString

		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&roleID,
			&roleName,
			&idJabatan,
			&jabatanName,
			&idPerusahaan,
			&fotoProfile,
			&banner,
			&user.MFAEnabled,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
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