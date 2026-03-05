package repository

import (
	"database/sql"
	"fortyfour-backend/internal/models"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUserRepository(t *testing.T) {
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)

	assert.NotNil(t, repo)
	assert.Equal(t, db, repo.db)
}

func TestUserRepository_Create(t *testing.T) {
	tests := []struct {
		name    string
		user    *models.User
		mockFn  func(mock sqlmock.Sqlmock, user *models.User)
		wantErr bool
	}{
		{
			name: "success - create user with all fields",
			user: &models.User{
				ID:       uuid.New().String(),
				Username: "testuser",
				Password: "hashedpassword",
				Email:    "test@example.com",
			},
			mockFn: func(mock sqlmock.Sqlmock, user *models.User) {
				roleID := "role-123"
				user.RoleID = &roleID

				mock.ExpectExec("INSERT INTO users").
					WithArgs(
						sqlmock.AnyArg(), // id (akan di-set jika kosong)
						user.Username,
						user.Password,
						user.Email,
						user.RoleID,
						user.IDJabatan,
						user.IDPerusahaan,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "success - create user without ID (auto-generate)",
			user: &models.User{
				Username: "newuser",
				Password: "password123",
				Email:    "new@example.com",
			},
			mockFn: func(mock sqlmock.Sqlmock, user *models.User) {
				// Mock untuk get default role
				mock.ExpectQuery("SELECT id FROM roles WHERE name = 'user'").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("default-role-id"))

				mock.ExpectExec("INSERT INTO users").
					WithArgs(
						sqlmock.AnyArg(),
						user.Username,
						user.Password,
						user.Email,
						sqlmock.AnyArg(), // role_id akan di-set ke default
						user.IDJabatan,
						user.IDPerusahaan,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "error - database error on insert",
			user: &models.User{
				ID:       uuid.New().String(),
				Username: "testuser",
				Password: "password",
				Email:    "test@example.com",
			},
			mockFn: func(mock sqlmock.Sqlmock, user *models.User) {
				roleID := "role-123"
				user.RoleID = &roleID

				mock.ExpectExec("INSERT INTO users").
					WithArgs(
						user.ID,
						user.Username,
						user.Password,
						user.Email,
						user.RoleID,
						user.IDJabatan,
						user.IDPerusahaan,
					).
					WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := NewUserRepository(db)

			if tt.mockFn != nil {
				tt.mockFn(mock, tt.user)
			}

			err = repo.Create(tt.user)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, tt.user.ID) // ID should be set
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_FindByID(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		mockFn  func(mock sqlmock.Sqlmock)
		want    *models.User
		wantErr bool
	}{
		{
			name:   "success - find user with all fields",
			userID: "user-123",
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "username", "password", "email",
					"role_id", "role_name", "id_jabatan", "nama_jabatan",
					"id_perusahaan", "foto_profile", "banner",
					"mfa_enabled", "mfa_secret",
					"status", "password_changed_at", "login_attempts",
					"created_at", "updated_at",
				}).AddRow(
					"user-123", "testuser", "hashedpass", "test@example.com",
					"role-1", "admin", "jabatan-1", "Manager",
					"perusahaan-1", "photo.jpg", "banner.jpg",
					true, "secret123",
					"Aktif", time.Now(), 0,
					time.Now(), time.Now(),
				)

				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs("user-123").
					WillReturnRows(rows)
			},
			want: &models.User{
				ID:       "user-123",
				Username: "testuser",
				Email:    "test@example.com",
			},
			wantErr: false,
		},
		{
			name:   "success - find user with null fields",
			userID: "user-456",
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "username", "password", "email",
					"role_id", "role_name", "id_jabatan", "nama_jabatan",
					"id_perusahaan", "foto_profile", "banner",
					"mfa_enabled", "mfa_secret",
					"status", "password_changed_at", "login_attempts",
					"created_at", "updated_at",
				}).AddRow(
					"user-456", "testuser2", "hashedpass2", "test2@example.com",
					nil, nil, nil, nil,
					nil, nil, nil,
					false, nil,
					"Aktif", time.Now(), 0,
					time.Now(), time.Now(),
				)

				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs("user-456").
					WillReturnRows(rows)
			},
			want: &models.User{
				ID:       "user-456",
				Username: "testuser2",
				Email:    "test2@example.com",
			},
			wantErr: false,
		},
		{
			name:   "error - user not found",
			userID: "non-existent",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs("non-existent").
					WillReturnError(sql.ErrNoRows)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "error - database error",
			userID: "user-123",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs("user-123").
					WillReturnError(sql.ErrConnDone)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := NewUserRepository(db)

			if tt.mockFn != nil {
				tt.mockFn(mock)
			}

			user, err := repo.FindByID(tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				if tt.want != nil {
					assert.Equal(t, tt.want.ID, user.ID)
					assert.Equal(t, tt.want.Username, user.Username)
					assert.Equal(t, tt.want.Email, user.Email)
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_FindByUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
		mockFn   func(mock sqlmock.Sqlmock)
		want     *models.User
		wantErr  bool
	}{
		{
			name:     "success - find user by username",
			username: "testuser",
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "username", "password", "email",
					"role_id", "role_name", "id_jabatan", "nama_jabatan",
					"id_perusahaan", "foto_profile", "banner",
					"mfa_enabled", "mfa_secret",
					"status", "password_changed_at", "login_attempts",
					"created_at", "updated_at",
				}).AddRow(
					"user-123", "testuser", "hashedpass", "test@example.com",
					"role-1", "admin", nil, nil,
					nil, nil, nil,
					false, nil,
					"Aktif", time.Now(), 0,
					time.Now(), time.Now(),
				)

				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs("testuser").
					WillReturnRows(rows)
			},
			want: &models.User{
				ID:       "user-123",
				Username: "testuser",
				Email:    "test@example.com",
			},
			wantErr: false,
		},
		{
			name:     "error - user not found",
			username: "nonexistent",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs("nonexistent").
					WillReturnError(sql.ErrNoRows)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := NewUserRepository(db)

			if tt.mockFn != nil {
				tt.mockFn(mock)
			}

			user, err := repo.FindByUsername(tt.username)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				if tt.want != nil {
					assert.Equal(t, tt.want.Username, user.Username)
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_Update(t *testing.T) {
	tests := []struct {
		name    string
		user    *models.User
		mockFn  func(mock sqlmock.Sqlmock, user *models.User)
		wantErr bool
	}{
		{
			name: "success - update user",
			user: &models.User{
				ID:       "user-123",
				Username: "updateduser",
				Email:    "updated@example.com",
			},
			mockFn: func(mock sqlmock.Sqlmock, user *models.User) {
				mock.ExpectExec("UPDATE users").
					WithArgs(
						user.Username,
						user.Email,
						user.RoleID,
						user.IDJabatan,
						user.ID,
					).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name: "error - database error",
			user: &models.User{
				ID:       "user-123",
				Username: "updateduser",
				Email:    "updated@example.com",
			},
			mockFn: func(mock sqlmock.Sqlmock, user *models.User) {
				mock.ExpectExec("UPDATE users").
					WithArgs(
						user.Username,
						user.Email,
						user.RoleID,
						user.IDJabatan,
						user.ID,
					).
					WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := NewUserRepository(db)

			if tt.mockFn != nil {
				tt.mockFn(mock, tt.user)
			}

			err = repo.Update(tt.user)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_UpdateWithPhoto(t *testing.T) {
	tests := []struct {
		name    string
		user    *models.User
		mockFn  func(mock sqlmock.Sqlmock, user *models.User)
		wantErr bool
	}{
		{
			name: "success - update user with photo",
			user: &models.User{
				ID:       "user-123",
				Username: "testuser",
				Email:    "test@example.com",
			},
			mockFn: func(mock sqlmock.Sqlmock, user *models.User) {
				photo := "photo.jpg"
				banner := "banner.jpg"
				user.FotoProfile = &photo
				user.Banner = &banner

				mock.ExpectExec("UPDATE users").
					WithArgs(
						user.Username,
						user.Email,
						user.RoleID,
						user.IDJabatan,
						user.FotoProfile,
						user.Banner,
						user.ID,
					).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := NewUserRepository(db)

			if tt.mockFn != nil {
				tt.mockFn(mock, tt.user)
			}

			err = repo.UpdateWithPhoto(tt.user)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_UpdatePassword(t *testing.T) {
	tests := []struct {
		name     string
		userID   string
		password string
		mockFn   func(mock sqlmock.Sqlmock)
		wantErr  bool
	}{
		{
			name:     "success - update password",
			userID:   "user-123",
			password: "newhashedpassword",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE users").
					WithArgs("newhashedpassword", "user-123").
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name:     "error - database error",
			userID:   "user-123",
			password: "newhashedpassword",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE users").
					WithArgs("newhashedpassword", "user-123").
					WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := NewUserRepository(db)

			if tt.mockFn != nil {
				tt.mockFn(mock)
			}

			err = repo.UpdatePassword(tt.userID, tt.password)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_SetMFA(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		secret  *string
		enabled bool
		mockFn  func(mock sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name:    "success - enable MFA with secret",
			userID:  "user-123",
			secret:  stringPtr("MFASECRET123"),
			enabled: true,
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE users").
					WithArgs(sqlmock.AnyArg(), true, "user-123").
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name:    "success - disable MFA",
			userID:  "user-123",
			secret:  nil,
			enabled: false,
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE users").
					WithArgs(nil, false, "user-123").
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := NewUserRepository(db)

			if tt.mockFn != nil {
				tt.mockFn(mock)
			}

			err = repo.SetMFA(tt.userID, tt.secret, tt.enabled)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_Delete(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		mockFn  func(mock sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name:   "success - delete user",
			userID: "user-123",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM users").
					WithArgs("user-123").
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name:   "error - database error",
			userID: "user-123",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM users").
					WithArgs("user-123").
					WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := NewUserRepository(db)

			if tt.mockFn != nil {
				tt.mockFn(mock)
			}

			err = repo.Delete(tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_FindAll(t *testing.T) {
	tests := []struct {
		name    string
		mockFn  func(mock sqlmock.Sqlmock)
		want    int // expected number of users
		wantErr bool
	}{
		{
			name: "success - find multiple users",
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "username", "email",
					"role_id", "role_name", "id_jabatan", "nama_jabatan",
					"id_perusahaan", "foto_profile", "banner",
					"mfa_enabled", "created_at", "updated_at",
				}).
					AddRow("user-1", "user1", "user1@example.com",
						"role-1", "admin", nil, nil, nil, nil, nil,
						false, time.Now(), time.Now()).
					AddRow("user-2", "user2", "user2@example.com",
						"role-2", "user", nil, nil, nil, nil, nil,
						false, time.Now(), time.Now())

				mock.ExpectQuery("SELECT (.+) FROM users").
					WillReturnRows(rows)
			},
			want:    2,
			wantErr: false,
		},
		{
			name: "success - no users found",
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "username", "email",
					"role_id", "role_name", "id_jabatan", "nama_jabatan",
					"id_perusahaan", "foto_profile", "banner",
					"mfa_enabled", "created_at", "updated_at",
				})

				mock.ExpectQuery("SELECT (.+) FROM users").
					WillReturnRows(rows)
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "error - database error",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM users").
					WillReturnError(sql.ErrConnDone)
			},
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := NewUserRepository(db)

			if tt.mockFn != nil {
				tt.mockFn(mock)
			}

			users, err := repo.FindAll()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, users, tt.want)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_GetPasswordByID(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		mockFn  func(mock sqlmock.Sqlmock)
		want    string
		wantErr bool
	}{
		{
			name:   "success - get password",
			userID: "user-123",
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"password"}).
					AddRow("hashedpassword123")

				mock.ExpectQuery("SELECT password FROM users").
					WithArgs("user-123").
					WillReturnRows(rows)
			},
			want:    "hashedpassword123",
			wantErr: false,
		},
		{
			name:   "error - user not found",
			userID: "non-existent",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT password FROM users").
					WithArgs("non-existent").
					WillReturnError(sql.ErrNoRows)
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := NewUserRepository(db)

			if tt.mockFn != nil {
				tt.mockFn(mock)
			}

			password, err := repo.GetPasswordByID(tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, password)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_EmailExists(t *testing.T) {
	tests := []struct {
		name      string
		email     string
		excludeID *string
		mockFn    func(mock sqlmock.Sqlmock)
		want      bool
		wantErr   bool
	}{
		{
			name:      "email exists",
			email:     "existing@example.com",
			excludeID: nil,
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(1)
				mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users WHERE email = \\?").
					WithArgs("existing@example.com").
					WillReturnRows(rows)
			},
			want:    true,
			wantErr: false,
		},
		{
			name:      "email does not exist",
			email:     "new@example.com",
			excludeID: nil,
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
				mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users WHERE email = \\?").
					WithArgs("new@example.com").
					WillReturnRows(rows)
			},
			want:    false,
			wantErr: false,
		},
		{
			name:      "email exists but exclude current user",
			email:     "existing@example.com",
			excludeID: stringPtr("user-123"),
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
				mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users WHERE email = \\? AND id != \\?").
					WithArgs("existing@example.com", "user-123").
					WillReturnRows(rows)
			},
			want:    false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := NewUserRepository(db)

			if tt.mockFn != nil {
				tt.mockFn(mock)
			}

			exists, err := repo.EmailExists(tt.email, tt.excludeID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, exists)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_UsernameExists(t *testing.T) {
	tests := []struct {
		name      string
		username  string
		excludeID *string
		mockFn    func(mock sqlmock.Sqlmock)
		want      bool
		wantErr   bool
	}{
		{
			name:      "username exists",
			username:  "existinguser",
			excludeID: nil,
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(1)
				mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users WHERE username = \\?").
					WithArgs("existinguser").
					WillReturnRows(rows)
			},
			want:    true,
			wantErr: false,
		},
		{
			name:      "username does not exist",
			username:  "newuser",
			excludeID: nil,
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
				mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users WHERE username = \\?").
					WithArgs("newuser").
					WillReturnRows(rows)
			},
			want:    false,
			wantErr: false,
		},
		{
			name:      "username exists but exclude current user",
			username:  "existinguser",
			excludeID: stringPtr("user-123"),
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
				mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users WHERE username = \\? AND id != \\?").
					WithArgs("existinguser", "user-123").
					WillReturnRows(rows)
			},
			want:    false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := NewUserRepository(db)

			if tt.mockFn != nil {
				tt.mockFn(mock)
			}

			exists, err := repo.UsernameExists(tt.username, tt.excludeID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, exists)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
