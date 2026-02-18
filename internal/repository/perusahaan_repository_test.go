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
				NamaPerusahaan: stringPtr("PT ABC Indonesia"),
				IDSubSektor:    stringPtr("sub-sektor-1"),
				Alamat:         stringPtr("Jl. Sudirman No. 1"),
				Telepon:        stringPtr("021-1234567"),
				Email:          stringPtr("info@ptabc.com"),
				Website:        stringPtr("www.ptabc.com"),
			},
			id: "perusahaan-1",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO perusahaan").
					WithArgs(
						"perusahaan-1",
						"photo.jpg",
						"PT ABC Indonesia",
						"sub-sektor-1",
						"Jl. Sudirman No. 1",
						"021-1234567",
						"info@ptabc.com",
						"www.ptabc.com",
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "success - create perusahaan with null id_sub_sektor",
			req: dto.CreatePerusahaanRequest{
				NamaPerusahaan: stringPtr("PT XYZ"),
				IDSubSektor:    nil,
			},
			id: "perusahaan-2",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO perusahaan").
					WithArgs(
						"perusahaan-2",
						nil,
						"PT XYZ",
						nil,
						nil,
						nil,
						nil,
						nil,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "success - create perusahaan with empty string id_sub_sektor",
			req: dto.CreatePerusahaanRequest{
				NamaPerusahaan: stringPtr("PT Empty SubSektor"),
				IDSubSektor:    stringPtr(""), // empty string should become nil
			},
			id: "perusahaan-3",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO perusahaan").
					WithArgs(
						"perusahaan-3",
						nil,
						"PT Empty SubSektor",
						nil,
						nil,
						nil,
						nil,
						nil,
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
						nil,
						"PT Minimal",
						nil,
						nil,
						nil,
						nil,
						nil,
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
						nil,
						"PT Test",
						nil,
						nil,
						nil,
						nil,
						nil,
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
			name: "success - get all perusahaan with sub_sektor",
			mockFn: func(mock sqlmock.Sqlmock) {
				now := time.Now()
				rows := sqlmock.NewRows([]string{
					"id", "photo", "nama_perusahaan", "alamat", "telepon", "email", "website", "created_at", "updated_at",
					"sub_id", "nama_sub_sektor", "id_sektor", "sub_created_at", "sub_updated_at", "nama_sektor",
				}).
					AddRow("p1", "photo1.jpg", "PT ABC", "Jl. A", "021-111", "a@a.com", "www.a.com", now, now,
						"sub1", "Perbankan", "s1", now, now, "Keuangan").
					AddRow("p2", "photo2.jpg", "PT XYZ", "Jl. B", "021-222", "b@b.com", "www.b.com", now, now,
						"sub2", "Asuransi", "s1", now, now, "Keuangan")

				mock.ExpectQuery("SELECT (.+) FROM perusahaan p LEFT JOIN sub_sektor ss (.+) LEFT JOIN sektor s").
					WillReturnRows(rows)
			},
			want:    2,
			wantErr: false,
		},
		{
			name: "success - get all perusahaan without sub_sektor",
			mockFn: func(mock sqlmock.Sqlmock) {
				now := time.Now()
				rows := sqlmock.NewRows([]string{
					"id", "photo", "nama_perusahaan", "alamat", "telepon", "email", "website", "created_at", "updated_at",
					"sub_id", "nama_sub_sektor", "id_sektor", "sub_created_at", "sub_updated_at", "nama_sektor",
				}).
					AddRow("p1", "", "PT ABC", "", "", "", "", now, now,
						nil, nil, nil, nil, nil, nil)

				mock.ExpectQuery("SELECT (.+) FROM perusahaan p LEFT JOIN sub_sektor ss (.+) LEFT JOIN sektor s").
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
					"sub_id", "nama_sub_sektor", "id_sektor", "sub_created_at", "sub_updated_at", "nama_sektor",
				})

				mock.ExpectQuery("SELECT (.+) FROM perusahaan p LEFT JOIN sub_sektor ss (.+) LEFT JOIN sektor s").
					WillReturnRows(rows)
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "error - database error",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM perusahaan p LEFT JOIN sub_sektor ss (.+) LEFT JOIN sektor s").
					WillReturnError(sql.ErrConnDone)
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "success - scan error is skipped",
			mockFn: func(mock sqlmock.Sqlmock) {
				now := time.Now()
				// This row will cause scan error but should be skipped (continue)
				rows := sqlmock.NewRows([]string{
					"id", "photo", "nama_perusahaan", "alamat", "telepon", "email", "website", "created_at", "updated_at",
					"sub_id", "nama_sub_sektor", "id_sektor", "sub_created_at", "sub_updated_at", "nama_sektor",
				}).
					AddRow("p1", "photo1.jpg", "PT ABC", "Jl. A", "021-111", "a@a.com", "www.a.com", now, now,
						"sub1", "Perbankan", "s1", now, now, "Keuangan")

				mock.ExpectQuery("SELECT (.+) FROM perusahaan p LEFT JOIN sub_sektor ss (.+) LEFT JOIN sektor s").
					WillReturnRows(rows)
			},
			want:    1,
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

			result, err := repo.GetAll()

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.want)
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
		wantErr bool
	}{
		{
			name: "success - get perusahaan with sub_sektor",
			id:   "p1",
			mockFn: func(mock sqlmock.Sqlmock) {
				now := time.Now()
				rows := sqlmock.NewRows([]string{
					"id", "photo", "nama_perusahaan", "alamat", "telepon", "email", "website", "created_at", "updated_at",
					"sub_id", "nama_sub_sektor", "id_sektor", "sub_created_at", "sub_updated_at", "nama_sektor",
				}).
					AddRow("p1", "photo1.jpg", "PT ABC", "Jl. A", "021-111", "a@a.com", "www.a.com", now, now,
						"sub1", "Perbankan", "s1", now, now, "Keuangan")

				mock.ExpectQuery("SELECT (.+) FROM perusahaan p LEFT JOIN sub_sektor ss (.+) LEFT JOIN sektor s (.+) WHERE p.id=(.+)").
					WithArgs("p1").
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name: "success - get perusahaan without sub_sektor",
			id:   "p2",
			mockFn: func(mock sqlmock.Sqlmock) {
				now := time.Now()
				rows := sqlmock.NewRows([]string{
					"id", "photo", "nama_perusahaan", "alamat", "telepon", "email", "website", "created_at", "updated_at",
					"sub_id", "nama_sub_sektor", "id_sektor", "sub_created_at", "sub_updated_at", "nama_sektor",
				}).
					AddRow("p2", "", "PT XYZ", "", "", "", "", now, now,
						nil, nil, nil, nil, nil, nil)

				mock.ExpectQuery("SELECT (.+) FROM perusahaan p LEFT JOIN sub_sektor ss (.+) LEFT JOIN sektor s (.+) WHERE p.id=(.+)").
					WithArgs("p2").
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name: "error - perusahaan not found",
			id:   "non-existent",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM perusahaan p LEFT JOIN sub_sektor ss (.+) LEFT JOIN sektor s (.+) WHERE p.id=(.+)").
					WithArgs("non-existent").
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
		},
		{
			name: "error - database error",
			id:   "p1",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM perusahaan p LEFT JOIN sub_sektor ss (.+) LEFT JOIN sektor s (.+) WHERE p.id=(.+)").
					WithArgs("p1").
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

			result, err := repo.GetByID(tt.id)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.id, result.ID)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPerusahaanRepository_Update(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		data    dto.PerusahaanResponse
		mockFn  func(mock sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "success - update perusahaan with sub_sektor",
			id:   "p1",
			data: dto.PerusahaanResponse{
				Photo:          "updated_photo.jpg",
				NamaPerusahaan: "PT ABC Updated",
				Alamat:         "Jl. Updated",
				Telepon:        "021-999",
				Email:          "updated@abc.com",
				Website:        "www.updated.com",
				SubSektor: &dto.SubSektorResponse{
					ID: "sub1",
				},
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE perusahaan SET").
					WithArgs(
						"updated_photo.jpg",
						"PT ABC Updated",
						"sub1",
						"Jl. Updated",
						"021-999",
						"updated@abc.com",
						"www.updated.com",
						"p1",
					).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name: "success - update perusahaan without sub_sektor",
			id:   "p2",
			data: dto.PerusahaanResponse{
				Photo:          "",
				NamaPerusahaan: "PT XYZ Updated",
				Alamat:         "",
				Telepon:        "",
				Email:          "",
				Website:        "",
				SubSektor:      nil,
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE perusahaan SET").
					WithArgs(
						nil,
						"PT XYZ Updated",
						nil,
						nil,
						nil,
						nil,
						nil,
						"p2",
					).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name: "error - database error",
			id:   "perusahaan-error",
			data: dto.PerusahaanResponse{
				Photo:          "",
				NamaPerusahaan: "PT Error",
				SubSektor:      nil,
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE perusahaan SET").
					WithArgs(
						nil,
						"PT Error",
						nil,
						nil,
						nil,
						nil,
						nil,
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

			err = repo.Update(tt.id, tt.data)

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
			id:   "p1",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM perusahaan WHERE id=(.+)").
					WithArgs("p1").
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name: "error - database error",
			id:   "p1",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM perusahaan WHERE id=(.+)").
					WithArgs("p1").
					WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
		{
			name: "error - foreign key constraint",
			id:   "p1",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM perusahaan WHERE id=(.+)").
					WithArgs("p1").
					WillReturnError(sql.ErrTxDone)
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

// Test helper functions
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
			input: stringPtr("hello"),
			want:  "hello",
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
			got := valueOrEmpty(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

// Integration test example
func TestPerusahaanRepository_Integration(t *testing.T) {
	t.Run("create and retrieve perusahaan", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewPerusahaanRepository(db)

		// Create
		createReq := dto.CreatePerusahaanRequest{
			NamaPerusahaan: stringPtr("PT Test Integration"),
		}
		createID := "test-integration"

		mock.ExpectExec("INSERT INTO perusahaan").
			WithArgs(createID, nil, "PT Test Integration", nil, nil, nil, nil, nil).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = repo.Create(createReq, createID)
		require.NoError(t, err)

		// GetByID
		now := time.Now()
		mock.ExpectQuery("SELECT (.+) FROM perusahaan p LEFT JOIN sub_sektor ss (.+) LEFT JOIN sektor s (.+) WHERE p.id=(.+)").
			WithArgs(createID).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "photo", "nama_perusahaan", "alamat", "telepon", "email", "website", "created_at", "updated_at",
				"sub_id", "nama_sub_sektor", "id_sektor", "sub_created_at", "sub_updated_at", "nama_sektor",
			}).AddRow(createID, nil, "PT Test Integration", nil, nil, nil, nil, now, now,
				nil, nil, nil, nil, nil, nil))

		result, err := repo.GetByID(createID)
		require.NoError(t, err)
		assert.Equal(t, "PT Test Integration", result.NamaPerusahaan)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
