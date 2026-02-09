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

func TestNewSektorRepository(t *testing.T) {
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewSektorRepository(db)

	assert.NotNil(t, repo)
	assert.Equal(t, db, repo.db)
}

func TestSektorRepository_GetAll(t *testing.T) {
	tests := []struct {
		name    string
		mockFn  func(mock sqlmock.Sqlmock)
		want    int // expected number of sektors
		wantErr bool
	}{
		{
			name: "success - get all sektor with sub_sektor",
			mockFn: func(mock sqlmock.Sqlmock) {
				now := time.Now()
				rows := sqlmock.NewRows([]string{
					"id", "nama_sektor", "created_at", "updated_at",
					"sub_id", "nama_sub_sektor", "id_sektor", "sub_created_at", "sub_updated_at",
				}).
					// Sektor 1 with 2 sub_sektor
					AddRow("sektor-1", "Keuangan", now, now, "sub-1", "Perbankan", "sektor-1", now, now).
					AddRow("sektor-1", "Keuangan", now, now, "sub-2", "Asuransi", "sektor-1", now, now).
					// Sektor 2 with 1 sub_sektor
					AddRow("sektor-2", "Teknologi", now, now, "sub-3", "Software", "sektor-2", now, now).
					// Sektor 3 without sub_sektor
					AddRow("sektor-3", "Kesehatan", now, now, nil, nil, nil, nil, nil)

				mock.ExpectQuery("SELECT (.+) FROM sektor (.+) LEFT JOIN sub_sektor").
					WillReturnRows(rows)
			},
			want:    3,
			wantErr: false,
		},
		{
			name: "success - empty result",
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "nama_sektor", "created_at", "updated_at",
					"sub_id", "nama_sub_sektor", "id_sektor", "sub_created_at", "sub_updated_at",
				})

				mock.ExpectQuery("SELECT (.+) FROM sektor (.+) LEFT JOIN sub_sektor").
					WillReturnRows(rows)
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "success - sektor without sub_sektor",
			mockFn: func(mock sqlmock.Sqlmock) {
				now := time.Now()
				rows := sqlmock.NewRows([]string{
					"id", "nama_sektor", "created_at", "updated_at",
					"sub_id", "nama_sub_sektor", "id_sektor", "sub_created_at", "sub_updated_at",
				}).
					AddRow("sektor-1", "Pertanian", now, now, nil, nil, nil, nil, nil)

				mock.ExpectQuery("SELECT (.+) FROM sektor (.+) LEFT JOIN sub_sektor").
					WillReturnRows(rows)
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "error - database error",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM sektor (.+) LEFT JOIN sub_sektor").
					WillReturnError(sql.ErrConnDone)
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "success - single sektor with multiple sub_sektor",
			mockFn: func(mock sqlmock.Sqlmock) {
				now := time.Now()
				rows := sqlmock.NewRows([]string{
					"id", "nama_sektor", "created_at", "updated_at",
					"sub_id", "nama_sub_sektor", "id_sektor", "sub_created_at", "sub_updated_at",
				}).
					AddRow("sektor-1", "Infrastruktur", now, now, "sub-1", "Jalan Tol", "sektor-1", now, now).
					AddRow("sektor-1", "Infrastruktur", now, now, "sub-2", "Bandara", "sektor-1", now, now).
					AddRow("sektor-1", "Infrastruktur", now, now, "sub-3", "Pelabuhan", "sektor-1", now, now)

				mock.ExpectQuery("SELECT (.+) FROM sektor (.+) LEFT JOIN sub_sektor").
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

			repo := NewSektorRepository(db)

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

				// Additional assertions for the first test case
				if tt.name == "success - get all sektor with sub_sektor" && len(result) > 0 {
					// Find sektor-1 (Keuangan) and verify it has 2 sub_sektor
					for _, sektor := range result {
						if sektor.ID == "sektor-1" {
							assert.Len(t, sektor.SubSektor, 2)
						}
					}
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestSektorRepository_GetByID(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		mockFn  func(mock sqlmock.Sqlmock)
		want    *dto.SektorResponse
		wantErr bool
	}{
		{
			name: "success - find sektor with sub_sektor",
			id:   "sektor-123",
			mockFn: func(mock sqlmock.Sqlmock) {
				now := time.Now()

				// Mock main sektor query
				sektorRows := sqlmock.NewRows([]string{
					"id", "nama_sektor", "created_at", "updated_at",
				}).AddRow("sektor-123", "Keuangan", now, now)

				mock.ExpectQuery("SELECT id, nama_sektor, created_at, updated_at FROM sektor WHERE id=\\?").
					WithArgs("sektor-123").
					WillReturnRows(sektorRows)

				// Mock sub_sektor query
				subSektorRows := sqlmock.NewRows([]string{
					"id", "nama_sub_sektor", "id_sektor", "created_at", "updated_at",
				}).
					AddRow("sub-1", "Perbankan", "sektor-123", now, now).
					AddRow("sub-2", "Asuransi", "sektor-123", now, now)

				mock.ExpectQuery("SELECT id, nama_sub_sektor, id_sektor, created_at, updated_at FROM sub_sektor WHERE id_sektor=\\?").
					WithArgs("sektor-123").
					WillReturnRows(subSektorRows)
			},
			want: &dto.SektorResponse{
				ID:         "sektor-123",
				NamaSektor: "Keuangan",
			},
			wantErr: false,
		},
		{
			name: "success - find sektor without sub_sektor",
			id:   "sektor-456",
			mockFn: func(mock sqlmock.Sqlmock) {
				now := time.Now()

				// Mock main sektor query
				sektorRows := sqlmock.NewRows([]string{
					"id", "nama_sektor", "created_at", "updated_at",
				}).AddRow("sektor-456", "Pertanian", now, now)

				mock.ExpectQuery("SELECT id, nama_sektor, created_at, updated_at FROM sektor WHERE id=\\?").
					WithArgs("sektor-456").
					WillReturnRows(sektorRows)

				// Mock empty sub_sektor query
				subSektorRows := sqlmock.NewRows([]string{
					"id", "nama_sub_sektor", "id_sektor", "created_at", "updated_at",
				})

				mock.ExpectQuery("SELECT id, nama_sub_sektor, id_sektor, created_at, updated_at FROM sub_sektor WHERE id_sektor=\\?").
					WithArgs("sektor-456").
					WillReturnRows(subSektorRows)
			},
			want: &dto.SektorResponse{
				ID:         "sektor-456",
				NamaSektor: "Pertanian",
			},
			wantErr: false,
		},
		{
			name: "error - sektor not found",
			id:   "non-existent",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, nama_sektor, created_at, updated_at FROM sektor WHERE id=\\?").
					WithArgs("non-existent").
					WillReturnError(sql.ErrNoRows)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error - database error on sektor query",
			id:   "sektor-123",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, nama_sektor, created_at, updated_at FROM sektor WHERE id=\\?").
					WithArgs("sektor-123").
					WillReturnError(sql.ErrConnDone)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error - database error on sub_sektor query",
			id:   "sektor-123",
			mockFn: func(mock sqlmock.Sqlmock) {
				now := time.Now()

				// Mock main sektor query succeeds
				sektorRows := sqlmock.NewRows([]string{
					"id", "nama_sektor", "created_at", "updated_at",
				}).AddRow("sektor-123", "Keuangan", now, now)

				mock.ExpectQuery("SELECT id, nama_sektor, created_at, updated_at FROM sektor WHERE id=\\?").
					WithArgs("sektor-123").
					WillReturnRows(sektorRows)

				// Mock sub_sektor query fails
				mock.ExpectQuery("SELECT id, nama_sub_sektor, id_sektor, created_at, updated_at FROM sub_sektor WHERE id_sektor=\\?").
					WithArgs("sektor-123").
					WillReturnError(sql.ErrConnDone)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "success - sektor with many sub_sektor",
			id:   "sektor-789",
			mockFn: func(mock sqlmock.Sqlmock) {
				now := time.Now()

				// Mock main sektor query
				sektorRows := sqlmock.NewRows([]string{
					"id", "nama_sektor", "created_at", "updated_at",
				}).AddRow("sektor-789", "Infrastruktur", now, now)

				mock.ExpectQuery("SELECT id, nama_sektor, created_at, updated_at FROM sektor WHERE id=\\?").
					WithArgs("sektor-789").
					WillReturnRows(sektorRows)

				// Mock sub_sektor query with 5 items
				subSektorRows := sqlmock.NewRows([]string{
					"id", "nama_sub_sektor", "id_sektor", "created_at", "updated_at",
				}).
					AddRow("sub-1", "Jalan Tol", "sektor-789", now, now).
					AddRow("sub-2", "Bandara", "sektor-789", now, now).
					AddRow("sub-3", "Pelabuhan", "sektor-789", now, now).
					AddRow("sub-4", "Kereta Api", "sektor-789", now, now).
					AddRow("sub-5", "Jembatan", "sektor-789", now, now)

				mock.ExpectQuery("SELECT id, nama_sub_sektor, id_sektor, created_at, updated_at FROM sub_sektor WHERE id_sektor=\\?").
					WithArgs("sektor-789").
					WillReturnRows(subSektorRows)
			},
			want: &dto.SektorResponse{
				ID:         "sektor-789",
				NamaSektor: "Infrastruktur",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := NewSektorRepository(db)

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
				if tt.want != nil {
					assert.Equal(t, tt.want.ID, result.ID)
					assert.Equal(t, tt.want.NamaSektor, result.NamaSektor)
					assert.NotEmpty(t, result.CreatedAt)
					assert.NotEmpty(t, result.UpdatedAt)
					assert.NotNil(t, result.SubSektor) // should be initialized even if empty

					// Additional assertion for test with sub_sektor
					if tt.name == "success - find sektor with sub_sektor" {
						assert.Len(t, result.SubSektor, 2)
					}
					// Test with many sub_sektor
					if tt.name == "success - sektor with many sub_sektor" {
						assert.Len(t, result.SubSektor, 5)
					}
					// Test without sub_sektor
					if tt.name == "success - find sektor without sub_sektor" {
						assert.Len(t, result.SubSektor, 0)
					}
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}