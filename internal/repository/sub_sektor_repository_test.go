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

func TestNewSubSektorRepository(t *testing.T) {
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewSubSektorRepository(db)

	assert.NotNil(t, repo)
	assert.Equal(t, db, repo.db)
}

func TestSubSektorRepository_GetAll(t *testing.T) {
	tests := []struct {
		name    string
		mockFn  func(mock sqlmock.Sqlmock)
		want    int
		wantErr bool
	}{
		{
			name: "success - get all sub_sektor with sektor info",
			mockFn: func(mock sqlmock.Sqlmock) {
				now := time.Now()
				rows := sqlmock.NewRows([]string{
					"id", "nama_sub_sektor", "id_sektor", "nama_sektor", "created_at", "updated_at",
				}).
					AddRow("sub-1", "Perbankan", "sektor-1", "Keuangan", now, now).
					AddRow("sub-2", "Asuransi", "sektor-1", "Keuangan", now, now).
					AddRow("sub-3", "Software", "sektor-2", "Teknologi", now, now)

				mock.ExpectQuery("SELECT ss.id, ss.nama_sub_sektor, ss.id_sektor, s.nama_sektor, ss.created_at, ss.updated_at FROM sub_sektor ss JOIN sektor s").
					WillReturnRows(rows)
			},
			want:    3,
			wantErr: false,
		},
		{
			name: "success - empty result",
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "nama_sub_sektor", "id_sektor", "nama_sektor", "created_at", "updated_at",
				})

				mock.ExpectQuery("SELECT ss.id, ss.nama_sub_sektor, ss.id_sektor, s.nama_sektor, ss.created_at, ss.updated_at FROM sub_sektor ss JOIN sektor s").
					WillReturnRows(rows)
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "error - database error",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT ss.id, ss.nama_sub_sektor, ss.id_sektor, s.nama_sektor, ss.created_at, ss.updated_at FROM sub_sektor ss JOIN sektor s").
					WillReturnError(sql.ErrConnDone)
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "success - single sub_sektor",
			mockFn: func(mock sqlmock.Sqlmock) {
				now := time.Now()
				rows := sqlmock.NewRows([]string{
					"id", "nama_sub_sektor", "id_sektor", "nama_sektor", "created_at", "updated_at",
				}).
					AddRow("sub-1", "Retail Banking", "sektor-1", "Keuangan", now, now)

				mock.ExpectQuery("SELECT ss.id, ss.nama_sub_sektor, ss.id_sektor, s.nama_sektor, ss.created_at, ss.updated_at FROM sub_sektor ss JOIN sektor s").
					WillReturnRows(rows)
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "success - multiple sub_sektor from different sektor",
			mockFn: func(mock sqlmock.Sqlmock) {
				now := time.Now()
				rows := sqlmock.NewRows([]string{
					"id", "nama_sub_sektor", "id_sektor", "nama_sektor", "created_at", "updated_at",
				}).
					AddRow("sub-1", "Jalan Tol", "sektor-1", "Infrastruktur", now, now).
					AddRow("sub-2", "Bandara", "sektor-1", "Infrastruktur", now, now).
					AddRow("sub-3", "Cloud Computing", "sektor-2", "Teknologi", now, now).
					AddRow("sub-4", "AI/ML", "sektor-2", "Teknologi", now, now).
					AddRow("sub-5", "Hospital", "sektor-3", "Kesehatan", now, now)

				mock.ExpectQuery("SELECT ss.id, ss.nama_sub_sektor, ss.id_sektor, s.nama_sektor, ss.created_at, ss.updated_at FROM sub_sektor ss JOIN sektor s").
					WillReturnRows(rows)
			},
			want:    5,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := NewSubSektorRepository(db)

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

				// Verify that each result has sektor info
				if tt.want > 0 {
					for _, sub := range result {
						assert.NotEmpty(t, sub.ID)
						assert.NotEmpty(t, sub.NamaSubSektor)
						assert.NotEmpty(t, sub.IDSektor)
						assert.NotEmpty(t, sub.NamaSektor)
						assert.NotEmpty(t, sub.CreatedAt)
						assert.NotEmpty(t, sub.UpdatedAt)
					}
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestSubSektorRepository_GetByID(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		mockFn  func(mock sqlmock.Sqlmock)
		want    *dto.SubSektorResponse
		wantErr bool
	}{
		{
			name: "success - find sub_sektor with sektor info",
			id:   "sub-123",
			mockFn: func(mock sqlmock.Sqlmock) {
				now := time.Now()
				rows := sqlmock.NewRows([]string{
					"id", "nama_sub_sektor", "id_sektor", "nama_sektor", "created_at", "updated_at",
				}).AddRow("sub-123", "Perbankan", "sektor-1", "Keuangan", now, now)

				mock.ExpectQuery("SELECT ss.id, ss.nama_sub_sektor, ss.id_sektor, s.nama_sektor, ss.created_at, ss.updated_at FROM sub_sektor ss JOIN sektor s").
					WithArgs("sub-123").
					WillReturnRows(rows)
			},
			want: &dto.SubSektorResponse{
				ID:            "sub-123",
				NamaSubSektor: "Perbankan",
				IDSektor:      "sektor-1",
				NamaSektor:    "Keuangan",
			},
			wantErr: false,
		},
		{
			name: "success - find different sub_sektor",
			id:   "sub-456",
			mockFn: func(mock sqlmock.Sqlmock) {
				now := time.Now()
				rows := sqlmock.NewRows([]string{
					"id", "nama_sub_sektor", "id_sektor", "nama_sektor", "created_at", "updated_at",
				}).AddRow("sub-456", "Cloud Computing", "sektor-2", "Teknologi", now, now)

				mock.ExpectQuery("SELECT ss.id, ss.nama_sub_sektor, ss.id_sektor, s.nama_sektor, ss.created_at, ss.updated_at FROM sub_sektor ss JOIN sektor s").
					WithArgs("sub-456").
					WillReturnRows(rows)
			},
			want: &dto.SubSektorResponse{
				ID:            "sub-456",
				NamaSubSektor: "Cloud Computing",
				IDSektor:      "sektor-2",
				NamaSektor:    "Teknologi",
			},
			wantErr: false,
		},
		{
			name: "error - sub_sektor not found",
			id:   "non-existent",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT ss.id, ss.nama_sub_sektor, ss.id_sektor, s.nama_sektor, ss.created_at, ss.updated_at FROM sub_sektor ss JOIN sektor s").
					WithArgs("non-existent").
					WillReturnError(sql.ErrNoRows)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error - database error",
			id:   "sub-123",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT ss.id, ss.nama_sub_sektor, ss.id_sektor, s.nama_sektor, ss.created_at, ss.updated_at FROM sub_sektor ss JOIN sektor s").
					WithArgs("sub-123").
					WillReturnError(sql.ErrConnDone)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "success - infrastructure sub_sektor",
			id:   "sub-789",
			mockFn: func(mock sqlmock.Sqlmock) {
				now := time.Now()
				rows := sqlmock.NewRows([]string{
					"id", "nama_sub_sektor", "id_sektor", "nama_sektor", "created_at", "updated_at",
				}).AddRow("sub-789", "Jalan Tol", "sektor-3", "Infrastruktur", now, now)

				mock.ExpectQuery("SELECT ss.id, ss.nama_sub_sektor, ss.id_sektor, s.nama_sektor, ss.created_at, ss.updated_at FROM sub_sektor ss JOIN sektor s").
					WithArgs("sub-789").
					WillReturnRows(rows)
			},
			want: &dto.SubSektorResponse{
				ID:            "sub-789",
				NamaSubSektor: "Jalan Tol",
				IDSektor:      "sektor-3",
				NamaSektor:    "Infrastruktur",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := NewSubSektorRepository(db)

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
					assert.Equal(t, tt.want.NamaSubSektor, result.NamaSubSektor)
					assert.Equal(t, tt.want.IDSektor, result.IDSektor)
					assert.Equal(t, tt.want.NamaSektor, result.NamaSektor)
					assert.NotEmpty(t, result.CreatedAt)
					assert.NotEmpty(t, result.UpdatedAt)
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestSubSektorRepository_GetBySektorID(t *testing.T) {
	tests := []struct {
		name     string
		sektorID string
		mockFn   func(mock sqlmock.Sqlmock)
		want     int
		wantErr  bool
	}{
		{
			name:     "success - find sub_sektor by sektor_id",
			sektorID: "sektor-1",
			mockFn: func(mock sqlmock.Sqlmock) {
				now := time.Now()
				rows := sqlmock.NewRows([]string{
					"id", "nama_sub_sektor", "id_sektor", "nama_sektor", "created_at", "updated_at",
				}).
					AddRow("sub-1", "Perbankan", "sektor-1", "Keuangan", now, now).
					AddRow("sub-2", "Asuransi", "sektor-1", "Keuangan", now, now).
					AddRow("sub-3", "Pasar Modal", "sektor-1", "Keuangan", now, now)

				mock.ExpectQuery("SELECT ss.id, ss.nama_sub_sektor, ss.id_sektor, s.nama_sektor, ss.created_at, ss.updated_at FROM sub_sektor ss JOIN sektor s").
					WithArgs("sektor-1").
					WillReturnRows(rows)
			},
			want:    3,
			wantErr: false,
		},
		{
			name:     "success - empty result (sektor has no sub_sektor)",
			sektorID: "sektor-empty",
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "nama_sub_sektor", "id_sektor", "nama_sektor", "created_at", "updated_at",
				})

				mock.ExpectQuery("SELECT ss.id, ss.nama_sub_sektor, ss.id_sektor, s.nama_sektor, ss.created_at, ss.updated_at FROM sub_sektor ss JOIN sektor s").
					WithArgs("sektor-empty").
					WillReturnRows(rows)
			},
			want:    0,
			wantErr: false,
		},
		{
			name:     "error - database error",
			sektorID: "sektor-1",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT ss.id, ss.nama_sub_sektor, ss.id_sektor, s.nama_sektor, ss.created_at, ss.updated_at FROM sub_sektor ss JOIN sektor s").
					WithArgs("sektor-1").
					WillReturnError(sql.ErrConnDone)
			},
			want:    0,
			wantErr: true,
		},
		{
			name:     "success - single sub_sektor",
			sektorID: "sektor-2",
			mockFn: func(mock sqlmock.Sqlmock) {
				now := time.Now()
				rows := sqlmock.NewRows([]string{
					"id", "nama_sub_sektor", "id_sektor", "nama_sektor", "created_at", "updated_at",
				}).
					AddRow("sub-10", "Software Development", "sektor-2", "Teknologi", now, now)

				mock.ExpectQuery("SELECT ss.id, ss.nama_sub_sektor, ss.id_sektor, s.nama_sektor, ss.created_at, ss.updated_at FROM sub_sektor ss JOIN sektor s").
					WithArgs("sektor-2").
					WillReturnRows(rows)
			},
			want:    1,
			wantErr: false,
		},
		{
			name:     "success - many sub_sektor (infrastructure)",
			sektorID: "sektor-infra",
			mockFn: func(mock sqlmock.Sqlmock) {
				now := time.Now()
				rows := sqlmock.NewRows([]string{
					"id", "nama_sub_sektor", "id_sektor", "nama_sektor", "created_at", "updated_at",
				}).
					AddRow("sub-20", "Jalan Tol", "sektor-infra", "Infrastruktur", now, now).
					AddRow("sub-21", "Bandara", "sektor-infra", "Infrastruktur", now, now).
					AddRow("sub-22", "Pelabuhan", "sektor-infra", "Infrastruktur", now, now).
					AddRow("sub-23", "Kereta Api", "sektor-infra", "Infrastruktur", now, now).
					AddRow("sub-24", "Jembatan", "sektor-infra", "Infrastruktur", now, now).
					AddRow("sub-25", "Pembangkit Listrik", "sektor-infra", "Infrastruktur", now, now)

				mock.ExpectQuery("SELECT ss.id, ss.nama_sub_sektor, ss.id_sektor, s.nama_sektor, ss.created_at, ss.updated_at FROM sub_sektor ss JOIN sektor s").
					WithArgs("sektor-infra").
					WillReturnRows(rows)
			},
			want:    6,
			wantErr: false,
		},
		{
			name:     "success - verify all sub_sektor belong to same sektor",
			sektorID: "sektor-health",
			mockFn: func(mock sqlmock.Sqlmock) {
				now := time.Now()
				rows := sqlmock.NewRows([]string{
					"id", "nama_sub_sektor", "id_sektor", "nama_sektor", "created_at", "updated_at",
				}).
					AddRow("sub-30", "Rumah Sakit", "sektor-health", "Kesehatan", now, now).
					AddRow("sub-31", "Klinik", "sektor-health", "Kesehatan", now, now)

				mock.ExpectQuery("SELECT ss.id, ss.nama_sub_sektor, ss.id_sektor, s.nama_sektor, ss.created_at, ss.updated_at FROM sub_sektor ss JOIN sektor s").
					WithArgs("sektor-health").
					WillReturnRows(rows)
			},
			want:    2,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := NewSubSektorRepository(db)

			if tt.mockFn != nil {
				tt.mockFn(mock)
			}

			result, err := repo.GetBySektorID(tt.sektorID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.want)

				// Verify all sub_sektor belong to the requested sektor
				if tt.want > 0 {
					for _, sub := range result {
						assert.Equal(t, tt.sektorID, sub.IDSektor)
						assert.NotEmpty(t, sub.ID)
						assert.NotEmpty(t, sub.NamaSubSektor)
						assert.NotEmpty(t, sub.NamaSektor)
						assert.NotEmpty(t, sub.CreatedAt)
						assert.NotEmpty(t, sub.UpdatedAt)
					}
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
