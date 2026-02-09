package repository

import (
	"database/sql"
	"fortyfour-backend/internal/dto"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPerusahaanRepository(t *testing.T) {
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPerusahaanRepository(db)

	assert.NotNil(t, repo)
	assert.Equal(t, db, repo.db)
}

func TestPerusahaanRepository_Create(t *testing.T) {
	tests := []struct {
		name    string
		req     dto.CreatePerusahaanRequest
		id      string
		mockFn  func(mock sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "success - create perusahaan with all fields",
			req: dto.CreatePerusahaanRequest{
				Photo:          stringPtr("photo.jpg"),
				NamaPerusahaan: stringPtr("PT Test"),
				IDSubSektor:    stringPtr("subsektor-123"),
				Alamat:         stringPtr("Jl. Test No. 1"),
				Telepon:        stringPtr("021-12345678"),
				Email:          stringPtr("test@example.com"),
				Website:        stringPtr("https://test.com"),
			},
			id: "perusahaan-123",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO perusahaan").
					WithArgs(
						"perusahaan-123",
						"photo.jpg",
						"PT Test",
						"subsektor-123",
						"Jl. Test No. 1",
						"021-12345678",
						"test@example.com",
						"https://test.com",
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "success - create perusahaan with null id_sub_sektor",
			req: dto.CreatePerusahaanRequest{
				Photo:          stringPtr("photo.jpg"),
				NamaPerusahaan: stringPtr("PT Test"),
				IDSubSektor:    nil, // null sub sektor
				Alamat:         stringPtr("Jl. Test No. 1"),
				Telepon:        stringPtr("021-12345678"),
				Email:          stringPtr("test@example.com"),
				Website:        stringPtr("https://test.com"),
			},
			id: "perusahaan-456",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO perusahaan").
					WithArgs(
						"perusahaan-456",
						"photo.jpg",
						"PT Test",
						nil, // id_sub_sektor is nil
						"Jl. Test No. 1",
						"021-12345678",
						"test@example.com",
						"https://test.com",
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "success - create perusahaan with empty string id_sub_sektor",
			req: dto.CreatePerusahaanRequest{
				Photo:          stringPtr("photo.jpg"),
				NamaPerusahaan: stringPtr("PT Test"),
				IDSubSektor:    stringPtr(""), // empty string
				Alamat:         stringPtr("Jl. Test No. 1"),
				Telepon:        stringPtr("021-12345678"),
				Email:          stringPtr("test@example.com"),
				Website:        stringPtr("https://test.com"),
			},
			id: "perusahaan-789",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO perusahaan").
					WithArgs(
						"perusahaan-789",
						"photo.jpg",
						"PT Test",
						nil, // empty string becomes nil
						"Jl. Test No. 1",
						"021-12345678",
						"test@example.com",
						"https://test.com",
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "success - create perusahaan with minimal fields",
			req: dto.CreatePerusahaanRequest{
				NamaPerusahaan: stringPtr("PT Minimal"),
			},
			id: "perusahaan-min",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO perusahaan").
					WithArgs(
						"perusahaan-min",
						"",  // photo empty
						"PT Minimal",
						nil, // id_sub_sektor nil
						"",  // alamat empty
						"",  // telepon empty
						"",  // email empty
						"",  // website empty
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "error - database error",
			req: dto.CreatePerusahaanRequest{
				NamaPerusahaan: stringPtr("PT Test"),
			},
			id: "perusahaan-error",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO perusahaan").
					WithArgs(
						"perusahaan-error",
						"",
						"PT Test",
						nil,
						"",
						"",
						"",
						"",
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

			repo := NewPerusahaanRepository(db)

			if tt.mockFn != nil {
				tt.mockFn(mock)
			}

			err = repo.Create(tt.req, tt.id)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPerusahaanRepository_GetAll(t *testing.T) {
	tests := []struct {
		name    string
		mockFn  func(mock sqlmock.Sqlmock)
		want    int // expected number of perusahaan
		wantErr bool
	}{
		{
			name: "success - get all perusahaan with sub sektor",
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "photo", "nama_perusahaan", "alamat", "telepon", "email", "website", "created_at", "updated_at",
					"id", "nama_sub_sektor", "id_sektor", "created_at", "updated_at",
					"nama_sektor",
				}).
					AddRow(
						"perusahaan-1", "photo1.jpg", "PT Test 1", "Alamat 1", "021-111", "test1@test.com", "https://test1.com", time.Now(), time.Now(),
						"sub-1", "Sub Sektor 1", "sektor-1", time.Now(), time.Now(),
						"Sektor 1",
					).
					AddRow(
						"perusahaan-2", "photo2.jpg", "PT Test 2", "Alamat 2", "021-222", "test2@test.com", "https://test2.com", time.Now(), time.Now(),
						"sub-2", "Sub Sektor 2", "sektor-2", time.Now(), time.Now(),
						"Sektor 2",
					)

				mock.ExpectQuery("SELECT (.+) FROM perusahaan p").
					WillReturnRows(rows)
			},
			want:    2,
			wantErr: false,
		},
		{
			name: "success - get all perusahaan without sub sektor",
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "photo", "nama_perusahaan", "alamat", "telepon", "email", "website", "created_at", "updated_at",
					"id", "nama_sub_sektor", "id_sektor", "created_at", "updated_at",
					"nama_sektor",
				}).
					AddRow(
						"perusahaan-3", "photo3.jpg", "PT Test 3", "Alamat 3", "021-333", "test3@test.com", "https://test3.com", time.Now(), time.Now(),
						nil, nil, nil, nil, nil, // no sub sektor
						nil,
					)

				mock.ExpectQuery("SELECT (.+) FROM perusahaan p").
					WillReturnRows(rows)
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "success - no perusahaan found",
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "photo", "nama_perusahaan", "alamat", "telepon", "email", "website", "created_at", "updated_at",
					"id", "nama_sub_sektor", "id_sektor", "created_at", "updated_at",
					"nama_sektor",
				})

				mock.ExpectQuery("SELECT (.+) FROM perusahaan p").
					WillReturnRows(rows)
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "error - database error",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM perusahaan p").
					WillReturnError(sql.ErrConnDone)
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "success - scan error is skipped",
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "photo", "nama_perusahaan", "alamat", "telepon", "email", "website", "created_at", "updated_at",
					"id", "nama_sub_sektor", "id_sektor", "created_at", "updated_at",
					"nama_sektor",
				}).
					AddRow(
						"perusahaan-1", "photo1.jpg", "PT Test 1", "Alamat 1", "021-111", "test1@test.com", "https://test1.com", time.Now(), time.Now(),
						"sub-1", "Sub Sektor 1", "sektor-1", time.Now(), time.Now(),
						"Sektor 1",
					).
					AddRow(
						nil, nil, nil, nil, nil, nil, nil, nil, nil, // invalid row - will be skipped
						nil, nil, nil, nil, nil,
						nil,
					).
					AddRow(
						"perusahaan-2", "photo2.jpg", "PT Test 2", "Alamat 2", "021-222", "test2@test.com", "https://test2.com", time.Now(), time.Now(),
						"sub-2", "Sub Sektor 2", "sektor-2", time.Now(), time.Now(),
						"Sektor 2",
					)

				mock.ExpectQuery("SELECT (.+) FROM perusahaan p").
					WillReturnRows(rows)
			},
			want:    2, // only valid rows
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := NewPerusahaanRepository(db)

			if tt.mockFn != nil {
				tt.mockFn(mock)
			}

			perusahaans, err := repo.GetAll()

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, perusahaans)
			} else {
				assert.NoError(t, err)
				assert.Len(t, perusahaans, tt.want)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPerusahaanRepository_GetByID(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		mockFn  func(mock sqlmock.Sqlmock)
		want    *dto.PerusahaanResponse
		wantErr bool
	}{
		{
			name: "success - get perusahaan with sub sektor",
			id:   "perusahaan-123",
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "photo", "nama_perusahaan", "alamat", "telepon", "email", "website", "created_at", "updated_at",
					"id", "nama_sub_sektor", "id_sektor", "created_at", "updated_at",
					"nama_sektor",
				}).AddRow(
					"perusahaan-123", "photo.jpg", "PT Test", "Alamat Test", "021-123", "test@test.com", "https://test.com", time.Now(), time.Now(),
					"sub-123", "Sub Sektor Test", "sektor-123", time.Now(), time.Now(),
					"Sektor Test",
				)

				mock.ExpectQuery("SELECT (.+) FROM perusahaan p (.+) WHERE p.id=\\?").
					WithArgs("perusahaan-123").
					WillReturnRows(rows)
			},
			want: &dto.PerusahaanResponse{
				ID:             "perusahaan-123",
				NamaPerusahaan: "PT Test",
				SubSektor: &dto.SubSektorResponse{
					ID:            "sub-123",
					NamaSubSektor: "Sub Sektor Test",
					IDSektor:      "sektor-123",
				},
			},
			wantErr: false,
		},
		{
			name: "success - get perusahaan without sub sektor",
			id:   "perusahaan-456",
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "photo", "nama_perusahaan", "alamat", "telepon", "email", "website", "created_at", "updated_at",
					"id", "nama_sub_sektor", "id_sektor", "created_at", "updated_at",
					"nama_sektor",
				}).AddRow(
					"perusahaan-456", "photo.jpg", "PT Test 2", "Alamat Test 2", "021-456", "test2@test.com", "https://test2.com", time.Now(), time.Now(),
					nil, nil, nil, nil, nil, // no sub sektor
					nil,
				)

				mock.ExpectQuery("SELECT (.+) FROM perusahaan p (.+) WHERE p.id=\\?").
					WithArgs("perusahaan-456").
					WillReturnRows(rows)
			},
			want: &dto.PerusahaanResponse{
				ID:             "perusahaan-456",
				NamaPerusahaan: "PT Test 2",
				SubSektor:      nil,
			},
			wantErr: false,
		},
		{
			name: "error - perusahaan not found",
			id:   "non-existent",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM perusahaan p (.+) WHERE p.id=\\?").
					WithArgs("non-existent").
					WillReturnError(sql.ErrNoRows)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error - database error",
			id:   "perusahaan-123",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM perusahaan p (.+) WHERE p.id=\\?").
					WithArgs("perusahaan-123").
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

			repo := NewPerusahaanRepository(db)

			if tt.mockFn != nil {
				tt.mockFn(mock)
			}

			perusahaan, err := repo.GetByID(tt.id)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, perusahaan)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, perusahaan)
				if tt.want != nil {
					assert.Equal(t, tt.want.ID, perusahaan.ID)
					assert.Equal(t, tt.want.NamaPerusahaan, perusahaan.NamaPerusahaan)
					if tt.want.SubSektor == nil {
						assert.Nil(t, perusahaan.SubSektor)
					} else {
						assert.NotNil(t, perusahaan.SubSektor)
					}
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPerusahaanRepository_Update(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		req     dto.PerusahaanResponse
		mockFn  func(mock sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "success - update perusahaan with sub sektor",
			id:   "perusahaan-123",
			req: dto.PerusahaanResponse{
				Photo:          "updated-photo.jpg",
				NamaPerusahaan: "PT Updated",
				SubSektor: &dto.SubSektorResponse{
					ID:            "sub-123",
					NamaSubSektor: "Sub Sektor",
					IDSektor:      "sektor-123",
				},
				Alamat:  "Updated Alamat",
				Telepon: "021-999",
				Email:   "updated@test.com",
				Website: "https://updated.com",
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE perusahaan SET").
					WithArgs(
						"updated-photo.jpg",
						"PT Updated",
						"sub-123", // id_sub_sektor from SubSektor
						"Updated Alamat",
						"021-999",
						"updated@test.com",
						"https://updated.com",
						"perusahaan-123",
					).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name: "success - update perusahaan without sub sektor",
			id:   "perusahaan-456",
			req: dto.PerusahaanResponse{
				Photo:          "photo.jpg",
				NamaPerusahaan: "PT No Sub",
				SubSektor:      nil, // no sub sektor
				Alamat:         "Alamat",
				Telepon:        "021-111",
				Email:          "test@test.com",
				Website:        "https://test.com",
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE perusahaan SET").
					WithArgs(
						"photo.jpg",
						"PT No Sub",
						nil, // id_sub_sektor is nil
						"Alamat",
						"021-111",
						"test@test.com",
						"https://test.com",
						"perusahaan-456",
					).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name: "error - database error",
			id:   "perusahaan-error",
			req: dto.PerusahaanResponse{
				NamaPerusahaan: "PT Error",
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE perusahaan SET").
					WithArgs(
						"",
						"PT Error",
						nil,
						"",
						"",
						"",
						"",
						"perusahaan-error",
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

			repo := NewPerusahaanRepository(db)

			if tt.mockFn != nil {
				tt.mockFn(mock)
			}

			err = repo.Update(tt.id, tt.req)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPerusahaanRepository_Delete(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		mockFn  func(mock sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "success - delete perusahaan",
			id:   "perusahaan-123",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM perusahaan WHERE id=\\?").
					WithArgs("perusahaan-123").
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name: "error - database error",
			id:   "perusahaan-error",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM perusahaan WHERE id=\\?").
					WithArgs("perusahaan-error").
					WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
		{
			name: "error - foreign key constraint",
			id:   "perusahaan-fk",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM perusahaan WHERE id=\\?").
					WithArgs("perusahaan-fk").
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := NewPerusahaanRepository(db)

			if tt.mockFn != nil {
				tt.mockFn(mock)
			}

			err = repo.Delete(tt.id)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestValueOrEmpty(t *testing.T) {
	tests := []struct {
		name  string
		input *string
		want  string
	}{
		{
			name:  "nil pointer returns empty string",
			input: nil,
			want:  "",
		},
		{
			name:  "non-nil pointer returns value",
			input: stringPtr("test"),
			want:  "test",
		},
		{
			name:  "empty string pointer returns empty string",
			input: stringPtr(""),
			want:  "",
		},
		{
			name:  "whitespace string returns whitespace",
			input: stringPtr("  "),
			want:  "  ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := valueOrEmpty(tt.input)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestPerusahaanRepository_Integration(t *testing.T) {
	t.Run("create and retrieve perusahaan", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewPerusahaanRepository(db)
		id := "perusahaan-integration"

		// Create
		req := dto.CreatePerusahaanRequest{
			Photo:          stringPtr("photo.jpg"),
			NamaPerusahaan: stringPtr("PT Integration Test"),
			IDSubSektor:    stringPtr("sub-123"),
			Alamat:         stringPtr("Jl. Integration"),
			Telepon:        stringPtr("021-999"),
			Email:          stringPtr("integration@test.com"),
			Website:        stringPtr("https://integration.test"),
		}

		mock.ExpectExec("INSERT INTO perusahaan").
			WithArgs(
				id,
				"photo.jpg",
				"PT Integration Test",
				"sub-123",
				"Jl. Integration",
				"021-999",
				"integration@test.com",
				"https://integration.test",
			).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = repo.Create(req, id)
		require.NoError(t, err)

		// Retrieve
		rows := sqlmock.NewRows([]string{
			"id", "photo", "nama_perusahaan", "alamat", "telepon", "email", "website", "created_at", "updated_at",
			"id", "nama_sub_sektor", "id_sektor", "created_at", "updated_at",
			"nama_sektor",
		}).AddRow(
			id, "photo.jpg", "PT Integration Test", "Jl. Integration", "021-999", "integration@test.com", "https://integration.test", time.Now(), time.Now(),
			"sub-123", "Sub Test", "sektor-123", time.Now(), time.Now(),
			"Sektor Test",
		)

		mock.ExpectQuery("SELECT (.+) FROM perusahaan p (.+) WHERE p.id=\\?").
			WithArgs(id).
			WillReturnRows(rows)

		retrieved, err := repo.GetByID(id)
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.Equal(t, id, retrieved.ID)
		assert.Equal(t, "PT Integration Test", retrieved.NamaPerusahaan)
		assert.NotNil(t, retrieved.SubSektor)
		assert.Equal(t, "sub-123", retrieved.SubSektor.ID)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}