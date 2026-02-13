package repository

import (
	"context"
	"database/sql"
	"fortyfour-backend/internal/dto"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDashboardRepository(t *testing.T) {
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewDashboardRepository(db)

	assert.NotNil(t, repo)
	assert.Equal(t, db, repo.db)
}

func TestDashboardRepository_CountPerSektor(t *testing.T) {
	tests := []struct {
		name    string
		from    *string
		to      *string
		mockFn  func(mock sqlmock.Sqlmock, from, to *string)
		want    int
		wantErr bool
	}{
		{
			name: "success - count per sektor without date filter (default current month)",
			from: nil,
			to:   nil,
			mockFn: func(mock sqlmock.Sqlmock, from, to *string) {
				rows := sqlmock.NewRows([]string{"id", "nama_sektor", "total", "this_month"}).
					AddRow("sektor-1", "Keuangan", int64(25), int64(5)).
					AddRow("sektor-2", "Teknologi", int64(18), int64(3)).
					AddRow("sektor-3", "Kesehatan", int64(12), int64(2))

				// When no date filter, query uses DATE_FORMAT(CURDATE()...)
				mock.ExpectQuery("SELECT s.id, s.nama_sektor, COUNT\\(p.id\\) AS total, SUM\\(CASE WHEN p.created_at >= DATE_FORMAT\\(CURDATE\\(\\), '%Y-%m-01'\\) THEN 1 ELSE 0 END\\) AS this_month FROM perusahaan p JOIN sub_sektor ss").
					WillReturnRows(rows)
			},
			want:    3,
			wantErr: false,
		},
		{
			name: "success - count per sektor with date filter",
			from: stringPtr("2024-01-01"),
			to:   stringPtr("2024-01-31"),
			mockFn: func(mock sqlmock.Sqlmock, from, to *string) {
				rows := sqlmock.NewRows([]string{"id", "nama_sektor", "total", "this_month"}).
					AddRow("sektor-1", "Keuangan", int64(25), int64(8)).
					AddRow("sektor-2", "Teknologi", int64(18), int64(6))

				// When date filter provided, query uses BETWEEN ? AND ?
				mock.ExpectQuery("SELECT s.id, s.nama_sektor, COUNT\\(p.id\\) AS total, SUM\\(CASE WHEN p.created_at BETWEEN \\? AND \\? THEN 1 ELSE 0 END\\) AS this_month FROM perusahaan p JOIN sub_sektor ss").
					WithArgs("2024-01-01", "2024-01-31").
					WillReturnRows(rows)
			},
			want:    2,
			wantErr: false,
		},
		{
			name: "success - empty result (no perusahaan)",
			from: nil,
			to:   nil,
			mockFn: func(mock sqlmock.Sqlmock, from, to *string) {
				rows := sqlmock.NewRows([]string{"id", "nama_sektor", "total", "this_month"})

				mock.ExpectQuery("SELECT s.id, s.nama_sektor, COUNT\\(p.id\\) AS total, SUM\\(CASE WHEN p.created_at >= DATE_FORMAT\\(CURDATE\\(\\), '%Y-%m-01'\\) THEN 1 ELSE 0 END\\) AS this_month FROM perusahaan p JOIN sub_sektor ss").
					WillReturnRows(rows)
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "success - single sektor",
			from: nil,
			to:   nil,
			mockFn: func(mock sqlmock.Sqlmock, from, to *string) {
				rows := sqlmock.NewRows([]string{"id", "nama_sektor", "total", "this_month"}).
					AddRow("sektor-1", "Keuangan", int64(30), int64(10))

				mock.ExpectQuery("SELECT s.id, s.nama_sektor, COUNT\\(p.id\\) AS total, SUM\\(CASE WHEN p.created_at >= DATE_FORMAT\\(CURDATE\\(\\), '%Y-%m-01'\\) THEN 1 ELSE 0 END\\) AS this_month FROM perusahaan p JOIN sub_sektor ss").
					WillReturnRows(rows)
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "success - multiple sectors with zero this_month",
			from: stringPtr("2024-12-01"),
			to:   stringPtr("2024-12-31"),
			mockFn: func(mock sqlmock.Sqlmock, from, to *string) {
				rows := sqlmock.NewRows([]string{"id", "nama_sektor", "total", "this_month"}).
					AddRow("sektor-1", "Keuangan", int64(25), int64(0)).
					AddRow("sektor-2", "Teknologi", int64(18), int64(0)).
					AddRow("sektor-3", "Kesehatan", int64(12), int64(1))

				mock.ExpectQuery("SELECT s.id, s.nama_sektor, COUNT\\(p.id\\) AS total, SUM\\(CASE WHEN p.created_at BETWEEN \\? AND \\? THEN 1 ELSE 0 END\\) AS this_month FROM perusahaan p JOIN sub_sektor ss").
					WithArgs("2024-12-01", "2024-12-31").
					WillReturnRows(rows)
			},
			want:    3,
			wantErr: false,
		},
		{
			name: "error - database error on query",
			from: nil,
			to:   nil,
			mockFn: func(mock sqlmock.Sqlmock, from, to *string) {
				mock.ExpectQuery("SELECT s.id, s.nama_sektor, COUNT\\(p.id\\) AS total, SUM\\(CASE WHEN p.created_at >= DATE_FORMAT\\(CURDATE\\(\\), '%Y-%m-01'\\) THEN 1 ELSE 0 END\\) AS this_month FROM perusahaan p JOIN sub_sektor ss").
					WillReturnError(sql.ErrConnDone)
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "error - scan error on row",
			from: nil,
			to:   nil,
			mockFn: func(mock sqlmock.Sqlmock, from, to *string) {
				// Return wrong column types to cause scan error
				rows := sqlmock.NewRows([]string{"id", "nama_sektor", "total", "this_month"}).
					AddRow("sektor-1", "Keuangan", "invalid", "invalid") // String instead of int64

				mock.ExpectQuery("SELECT s.id, s.nama_sektor, COUNT\\(p.id\\) AS total, SUM\\(CASE WHEN p.created_at >= DATE_FORMAT\\(CURDATE\\(\\), '%Y-%m-01'\\) THEN 1 ELSE 0 END\\) AS this_month FROM perusahaan p JOIN sub_sektor ss").
					WillReturnRows(rows)
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

			repo := NewDashboardRepository(db)

			if tt.mockFn != nil {
				tt.mockFn(mock, tt.from, tt.to)
			}

			ctx := context.Background()
			result, err := repo.CountPerSektor(ctx, tt.from, tt.to)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.want)

				// Verify structure if we have results
				if tt.want > 0 {
					for _, sector := range result {
						assert.NotEmpty(t, sector.ID)
						assert.NotEmpty(t, sector.Nama)
						assert.GreaterOrEqual(t, sector.Total, int64(0))
						assert.GreaterOrEqual(t, sector.ThisMonth, int64(0))
					}
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDashboardRepository_SeGlobalAgg(t *testing.T) {
	tests := []struct {
		name    string
		mockFn  func(mock sqlmock.Sqlmock)
		want    dto.SeAgg
		wantErr bool
	}{
		{
			name: "success - get total SE count",
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"total_se"}).
					AddRow(int64(42))

				mock.ExpectQuery("SELECT COUNT\\(id\\) as total_se FROM se").
					WillReturnRows(rows)
			},
			want: dto.SeAgg{
				TotalSE: int64(42),
			},
			wantErr: false,
		},
		{
			name: "success - zero SE (empty table)",
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"total_se"}).
					AddRow(int64(0))

				mock.ExpectQuery("SELECT COUNT\\(id\\) as total_se FROM se").
					WillReturnRows(rows)
			},
			want: dto.SeAgg{
				TotalSE: int64(0),
			},
			wantErr: false,
		},
		{
			name: "success - large number of SE",
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"total_se"}).
					AddRow(int64(1500))

				mock.ExpectQuery("SELECT COUNT\\(id\\) as total_se FROM se").
					WillReturnRows(rows)
			},
			want: dto.SeAgg{
				TotalSE: int64(1500),
			},
			wantErr: false,
		},
		{
			name: "success - no rows returns zero (ErrNoRows handled)",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT COUNT\\(id\\) as total_se FROM se").
					WillReturnError(sql.ErrNoRows)
			},
			want: dto.SeAgg{
				TotalSE: int64(0),
			},
			wantErr: false,
		},
		{
			name: "error - database error",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT COUNT\\(id\\) as total_se FROM se").
					WillReturnError(sql.ErrConnDone)
			},
			want: dto.SeAgg{
				TotalSE: int64(0),
			},
			wantErr: true,
		},
		{
			name: "error - scan error (wrong column type)",
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"total_se"}).
					AddRow("invalid") // String instead of int64

				mock.ExpectQuery("SELECT COUNT\\(id\\) as total_se FROM se").
					WillReturnRows(rows)
			},
			want: dto.SeAgg{
				TotalSE: int64(0),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := NewDashboardRepository(db)

			if tt.mockFn != nil {
				tt.mockFn(mock)
			}

			ctx := context.Background()
			result, err := repo.SeGlobalAgg(ctx)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.TotalSE, result.TotalSE)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

// Integration test for dashboard repository
func TestDashboardRepository_Integration(t *testing.T) {
	t.Run("count per sektor and se global agg", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewDashboardRepository(db)
		ctx := context.Background()

		// Test CountPerSektor
		mock.ExpectQuery("SELECT s.id, s.nama_sektor, COUNT\\(p.id\\) AS total, SUM\\(CASE WHEN p.created_at >= DATE_FORMAT\\(CURDATE\\(\\), '%Y-%m-01'\\) THEN 1 ELSE 0 END\\) AS this_month FROM perusahaan p JOIN sub_sektor ss").
			WillReturnRows(sqlmock.NewRows([]string{"id", "nama_sektor", "total", "this_month"}).
				AddRow("s1", "Keuangan", int64(10), int64(2)))

		sectors, err := repo.CountPerSektor(ctx, nil, nil)
		require.NoError(t, err)
		assert.Len(t, sectors, 1)
		assert.Equal(t, "Keuangan", sectors[0].Nama)

		// Test SeGlobalAgg
		mock.ExpectQuery("SELECT COUNT\\(id\\) as total_se FROM se").
			WillReturnRows(sqlmock.NewRows([]string{"total_se"}).AddRow(int64(25)))

		seAgg, err := repo.SeGlobalAgg(ctx)
		require.NoError(t, err)
		assert.Equal(t, int64(25), seAgg.TotalSE)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// Test with different date range scenarios
func TestDashboardRepository_CountPerSektor_DateRanges(t *testing.T) {
	tests := []struct {
		name        string
		from        *string
		to          *string
		description string
	}{
		{
			name:        "current month (default)",
			from:        nil,
			to:          nil,
			description: "Should use DATE_FORMAT for current month",
		},
		{
			name:        "specific month - January 2024",
			from:        stringPtr("2024-01-01"),
			to:          stringPtr("2024-01-31"),
			description: "Should use BETWEEN with provided dates",
		},
		{
			name:        "specific month - December 2024",
			from:        stringPtr("2024-12-01"),
			to:          stringPtr("2024-12-31"),
			description: "Should handle end of year",
		},
		{
			name:        "custom date range",
			from:        stringPtr("2024-06-15"),
			to:          stringPtr("2024-07-15"),
			description: "Should handle custom ranges spanning months",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := NewDashboardRepository(db)
			ctx := context.Background()

			rows := sqlmock.NewRows([]string{"id", "nama_sektor", "total", "this_month"}).
				AddRow("s1", "Test Sektor", int64(5), int64(1))

			if tt.from != nil && tt.to != nil {
				// Expect query with date parameters
				mock.ExpectQuery("SELECT s.id, s.nama_sektor, COUNT\\(p.id\\) AS total, SUM\\(CASE WHEN p.created_at BETWEEN \\? AND \\? THEN 1 ELSE 0 END\\) AS this_month").
					WithArgs(*tt.from, *tt.to).
					WillReturnRows(rows)
			} else {
				// Expect query without parameters (uses CURDATE())
				mock.ExpectQuery("SELECT s.id, s.nama_sektor, COUNT\\(p.id\\) AS total, SUM\\(CASE WHEN p.created_at >= DATE_FORMAT\\(CURDATE\\(\\), '%Y-%m-01'\\) THEN 1 ELSE 0 END\\) AS this_month").
					WillReturnRows(rows)
			}

			result, err := repo.CountPerSektor(ctx, tt.from, tt.to)
			require.NoError(t, err)
			assert.Len(t, result, 1)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}