package repository

import (
	"context"
	"database/sql"
	"database/sql/driver"
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

/*
=====================================
 TEST buildDateRange
=====================================
*/

func TestBuildDateRange(t *testing.T) {
	tests := []struct {
		name       string
		filter     dto.DashboardFilter
		wantFrom   string
		wantTo     string
		wantNilPtr bool
	}{
		{
			name:       "no filter → nil",
			filter:     dto.DashboardFilter{},
			wantNilPtr: true,
		},
		{
			name:     "explicit from+to",
			filter:   dto.DashboardFilter{From: stringPtr("2024-03-01"), To: stringPtr("2024-03-31")},
			wantFrom: "2024-03-01",
			wantTo:   "2024-03-31",
		},
		{
			name:     "year only → full year",
			filter:   dto.DashboardFilter{Year: stringPtr("2024")},
			wantFrom: "2024-01-01",
			wantTo:   "2024-12-31",
		},
		{
			name:     "year + Q1",
			filter:   dto.DashboardFilter{Year: stringPtr("2024"), Quarter: stringPtr("1")},
			wantFrom: "2024-01-01",
			wantTo:   "2024-03-31",
		},
		{
			name:     "year + Q2",
			filter:   dto.DashboardFilter{Year: stringPtr("2024"), Quarter: stringPtr("2")},
			wantFrom: "2024-04-01",
			wantTo:   "2024-06-30",
		},
		{
			name:     "year + Q3",
			filter:   dto.DashboardFilter{Year: stringPtr("2024"), Quarter: stringPtr("3")},
			wantFrom: "2024-07-01",
			wantTo:   "2024-09-30",
		},
		{
			name:     "year + Q4",
			filter:   dto.DashboardFilter{Year: stringPtr("2024"), Quarter: stringPtr("4")},
			wantFrom: "2024-10-01",
			wantTo:   "2024-12-31",
		},
		{
			name:     "year + invalid quarter → full year",
			filter:   dto.DashboardFilter{Year: stringPtr("2024"), Quarter: stringPtr("9")},
			wantFrom: "2024-01-01",
			wantTo:   "2024-12-31",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			from, to := buildDateRange(tt.filter)
			if tt.wantNilPtr {
				assert.Nil(t, from)
				assert.Nil(t, to)
			} else {
				require.NotNil(t, from)
				require.NotNil(t, to)
				assert.Equal(t, tt.wantFrom, *from)
				assert.Equal(t, tt.wantTo, *to)
			}
		})
	}
}

/*
=====================================
 TEST CountPerSektor
=====================================
*/

func TestDashboardRepository_CountPerSektor(t *testing.T) {
	tests := []struct {
		name    string
		filter  dto.DashboardFilter
		mockFn  func(mock sqlmock.Sqlmock)
		want    int
		wantErr bool
	}{
		{
			name:   "success - no filter (default current month)",
			filter: dto.DashboardFilter{},
			mockFn: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "nama_sektor", "total", "this_month"}).
					AddRow("sektor-1", "Keuangan", int64(25), int64(5)).
					AddRow("sektor-2", "Teknologi", int64(18), int64(3))

				m.ExpectQuery("SELECT s.id, s.nama_sektor, COUNT").
					WillReturnRows(rows)
			},
			want: 2,
		},
		{
			name:   "success - with year filter",
			filter: dto.DashboardFilter{Year: stringPtr("2024")},
			mockFn: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "nama_sektor", "total", "this_month"}).
					AddRow("sektor-1", "Keuangan", int64(25), int64(20))

				m.ExpectQuery("SELECT s.id, s.nama_sektor, COUNT").
					WithArgs("2024-01-01", "2024-12-31").
					WillReturnRows(rows)
			},
			want: 1,
		},
		{
			name:   "success - with year+quarter filter",
			filter: dto.DashboardFilter{Year: stringPtr("2024"), Quarter: stringPtr("2")},
			mockFn: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "nama_sektor", "total", "this_month"}).
					AddRow("sektor-1", "Keuangan", int64(10), int64(5))

				m.ExpectQuery("SELECT s.id, s.nama_sektor, COUNT").
					WithArgs("2024-04-01", "2024-06-30").
					WillReturnRows(rows)
			},
			want: 1,
		},
		{
			name:   "success - with sub_sektor_id filter",
			filter: dto.DashboardFilter{SubSektorID: stringPtr("sub-uuid-123")},
			mockFn: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "nama_sektor", "total", "this_month"}).
					AddRow("sektor-1", "Keuangan", int64(5), int64(1))

				m.ExpectQuery("SELECT s.id, s.nama_sektor, COUNT").
					WithArgs("sub-uuid-123").
					WillReturnRows(rows)
			},
			want: 1,
		},
		{
			name:   "success - empty result",
			filter: dto.DashboardFilter{},
			mockFn: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT s.id, s.nama_sektor, COUNT").
					WillReturnRows(sqlmock.NewRows([]string{"id", "nama_sektor", "total", "this_month"}))
			},
			want: 0,
		},
		{
			name:    "error - db error",
			filter:  dto.DashboardFilter{},
			mockFn:  func(m sqlmock.Sqlmock) { m.ExpectQuery("SELECT s.id, s.nama_sektor, COUNT").WillReturnError(sql.ErrConnDone) },
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
			tt.mockFn(mock)

			result, err := repo.CountPerSektor(context.Background(), tt.filter)
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

/*
=====================================
 TEST SeGlobalAgg
=====================================
*/

func TestDashboardRepository_SeGlobalAgg(t *testing.T) {
	tests := []struct {
		name    string
		filter  dto.DashboardFilter
		mockFn  func(mock sqlmock.Sqlmock)
		want    dto.SeAgg
		wantErr bool
	}{
		{
			name:   "success - no filter",
			filter: dto.DashboardFilter{},
			mockFn: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT").
					WillReturnRows(sqlmock.NewRows([]string{"total_se", "this_month", "strategis", "tinggi", "rendah"}).
						AddRow(int64(42), int64(5), int64(20), int64(15), int64(7)))
			},
			want: dto.SeAgg{TotalSE: 42, ThisMonth: 5, Strategis: 20, Tinggi: 15, Rendah: 7},
		},
		{
			name:   "success - with year filter",
			filter: dto.DashboardFilter{Year: stringPtr("2024")},
			mockFn: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT").
					WithArgs("2024-01-01", "2024-12-31").
					WillReturnRows(sqlmock.NewRows([]string{"total_se", "this_month", "strategis", "tinggi", "rendah"}).
						AddRow(int64(100), int64(100), int64(40), int64(35), int64(25)))
			},
			want: dto.SeAgg{TotalSE: 100, ThisMonth: 100, Strategis: 40, Tinggi: 35, Rendah: 25},
		},
		{
			name:   "success - with kategori_se filter",
			filter: dto.DashboardFilter{KategoriSE: stringPtr("Strategis")},
			mockFn: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT").
					WithArgs("Strategis").
					WillReturnRows(sqlmock.NewRows([]string{"total_se", "this_month", "strategis", "tinggi", "rendah"}).
						AddRow(int64(30), int64(3), int64(30), int64(0), int64(0)))
			},
			want: dto.SeAgg{TotalSE: 30, ThisMonth: 3, Strategis: 30, Tinggi: 0, Rendah: 0},
		},
		{
			name:   "success - zero rows returns empty struct",
			filter: dto.DashboardFilter{},
			mockFn: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT").WillReturnError(sql.ErrNoRows)
			},
			want: dto.SeAgg{},
		},
		{
			name:    "error - db error",
			filter:  dto.DashboardFilter{},
			mockFn:  func(m sqlmock.Sqlmock) { m.ExpectQuery("SELECT").WillReturnError(sql.ErrConnDone) },
			want:    dto.SeAgg{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := NewDashboardRepository(db)
			tt.mockFn(mock)

			result, err := repo.SeGlobalAgg(context.Background(), tt.filter)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.TotalSE, result.TotalSE)
				assert.Equal(t, tt.want.ThisMonth, result.ThisMonth)
				assert.Equal(t, tt.want.Strategis, result.Strategis)
				assert.Equal(t, tt.want.Tinggi, result.Tinggi)
				assert.Equal(t, tt.want.Rendah, result.Rendah)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

/*
=====================================
 TEST SeStatusCount
=====================================
*/

func TestDashboardRepository_SeStatusCount(t *testing.T) {
	tests := []struct {
		name    string
		filter  dto.DashboardFilter
		mockFn  func(mock sqlmock.Sqlmock)
		want    dto.SeStatusCount
		wantErr bool
	}{
		{
			name:   "success - no filter",
			filter: dto.DashboardFilter{},
			mockFn: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT").
					WillReturnRows(sqlmock.NewRows([]string{"total_perusahaan", "sudah_mengisi_kse", "belum_mengisi_kse"}).
						AddRow(int64(200), int64(150), int64(50)))
			},
			want: dto.SeStatusCount{TotalPerusahaan: 200, SudahMengisiKSE: 150, BelumMengisiKSE: 50},
		},
		{
			name:   "success - with sub_sektor_id filter",
			filter: dto.DashboardFilter{SubSektorID: stringPtr("sub-uuid-abc")},
			mockFn: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT").
					WithArgs("sub-uuid-abc").
					WillReturnRows(sqlmock.NewRows([]string{"total_perusahaan", "sudah_mengisi_kse", "belum_mengisi_kse"}).
						AddRow(int64(50), int64(40), int64(10)))
			},
			want: dto.SeStatusCount{TotalPerusahaan: 50, SudahMengisiKSE: 40, BelumMengisiKSE: 10},
		},
		{
			name:   "success - no perusahaan mengisi",
			filter: dto.DashboardFilter{},
			mockFn: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT").
					WillReturnRows(sqlmock.NewRows([]string{"total_perusahaan", "sudah_mengisi_kse", "belum_mengisi_kse"}).
						AddRow(int64(100), int64(0), int64(100)))
			},
			want: dto.SeStatusCount{TotalPerusahaan: 100, SudahMengisiKSE: 0, BelumMengisiKSE: 100},
		},
		{
			name:   "success - zero rows",
			filter: dto.DashboardFilter{},
			mockFn: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT").WillReturnError(sql.ErrNoRows)
			},
			want: dto.SeStatusCount{},
		},
		{
			name:    "error - db error",
			filter:  dto.DashboardFilter{},
			mockFn:  func(m sqlmock.Sqlmock) { m.ExpectQuery("SELECT").WillReturnError(sql.ErrConnDone) },
			want:    dto.SeStatusCount{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := NewDashboardRepository(db)
			tt.mockFn(mock)

			result, err := repo.SeStatusCount(context.Background(), tt.filter)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.TotalPerusahaan, result.TotalPerusahaan)
				assert.Equal(t, tt.want.SudahMengisiKSE, result.SudahMengisiKSE)
				assert.Equal(t, tt.want.BelumMengisiKSE, result.BelumMengisiKSE)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}


/*
=====================================
 TEST CountPerSektor - FILTER KOMBINASI
=====================================
*/

func TestDashboardRepository_CountPerSektor_CombinedFilter(t *testing.T) {
	tests := []struct {
		name   string
		filter dto.DashboardFilter
		args   []driver.Value
	}{
		{
			name:   "year + sub_sektor_id",
			filter: dto.DashboardFilter{Year: stringPtr("2024"), SubSektorID: stringPtr("sub-uuid")},
			// args: from, to (dari year), lalu sub_sektor_id
			args: []driver.Value{"2024-01-01", "2024-12-31", "sub-uuid"},
		},
		{
			name:   "quarter + sub_sektor_id",
			filter: dto.DashboardFilter{Year: stringPtr("2024"), Quarter: stringPtr("3"), SubSektorID: stringPtr("sub-uuid")},
			args:   []driver.Value{"2024-07-01", "2024-09-30", "sub-uuid"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := NewDashboardRepository(db)

			rows := sqlmock.NewRows([]string{"id", "nama_sektor", "total", "this_month"}).
				AddRow("s1", "ILMATE", int64(10), int64(5))

			mock.ExpectQuery("SELECT s.id, s.nama_sektor, COUNT").
				WithArgs(tt.args...).
				WillReturnRows(rows)

			result, err := repo.CountPerSektor(context.Background(), tt.filter)
			assert.NoError(t, err)
			assert.Len(t, result, 1)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDashboardRepository_CountPerSektor_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewDashboardRepository(db)

	// Return kolom bertipe salah untuk trigger scan error
	rows := sqlmock.NewRows([]string{"id", "nama_sektor", "total", "this_month"}).
		AddRow("s1", "ILMATE", "bukan_angka", "bukan_angka")

	mock.ExpectQuery("SELECT s.id, s.nama_sektor, COUNT").WillReturnRows(rows)

	result, err := repo.CountPerSektor(context.Background(), dto.DashboardFilter{})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

/*
=====================================
 TEST SeGlobalAgg - FILTER KOMBINASI
=====================================
*/

func TestDashboardRepository_SeGlobalAgg_CombinedFilter(t *testing.T) {
	tests := []struct {
		name   string
		filter dto.DashboardFilter
		args   []driver.Value
	}{
		{
			name:   "sub_sektor_id + tahun",
			filter: dto.DashboardFilter{SubSektorID: stringPtr("sub-uuid"), Year: stringPtr("2024")},
			// WHERE: sub_sektor_id=?, lalu this_month: from?, to?
			args: []driver.Value{"sub-uuid", "2024-01-01", "2024-12-31"},
		},
		{
			name:   "sub_sektor_id + kategori_se + quarter",
			filter: dto.DashboardFilter{SubSektorID: stringPtr("sub-uuid"), KategoriSE: stringPtr("Strategis"), Year: stringPtr("2024"), Quarter: stringPtr("2")},
			args:   []driver.Value{"sub-uuid", "Strategis", "2024-04-01", "2024-06-30"},
		},
		{
			name:   "kategori_se saja",
			filter: dto.DashboardFilter{KategoriSE: stringPtr("Rendah")},
			// WHERE: kategori_se=?, this_month: CURDATE (no extra args)
			args: []driver.Value{"Rendah"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := NewDashboardRepository(db)

			rows := sqlmock.NewRows([]string{"total_se", "this_month", "strategis", "tinggi", "rendah"}).
				AddRow(int64(10), int64(3), int64(5), int64(3), int64(2))

			mock.ExpectQuery("SELECT").
				WithArgs(tt.args...).
				WillReturnRows(rows)

			result, err := repo.SeGlobalAgg(context.Background(), tt.filter)
			assert.NoError(t, err)
			assert.Equal(t, int64(10), result.TotalSE)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDashboardRepository_SeGlobalAgg_KategoriBreakdownSumEqualsTotal(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewDashboardRepository(db)

	mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{"total_se", "this_month", "strategis", "tinggi", "rendah"}).
			AddRow(int64(100), int64(20), int64(40), int64(35), int64(25)))

	result, err := repo.SeGlobalAgg(context.Background(), dto.DashboardFilter{})
	assert.NoError(t, err)

	// Invariant: Strategis + Tinggi + Rendah harus == TotalSE
	assert.Equal(t, result.TotalSE, result.Strategis+result.Tinggi+result.Rendah)
	assert.NoError(t, mock.ExpectationsWereMet())
}

/*
=====================================
 TEST SeStatusCount - INVARIANT
=====================================
*/

func TestDashboardRepository_SeStatusCount_SudahPlusBelumEqualsTotal(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewDashboardRepository(db)

	mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{"total_perusahaan", "sudah_mengisi_kse", "belum_mengisi_kse"}).
			AddRow(int64(200), int64(120), int64(80)))

	result, err := repo.SeStatusCount(context.Background(), dto.DashboardFilter{})
	assert.NoError(t, err)

	// Invariant: Sudah + Belum == Total
	assert.Equal(t, result.TotalPerusahaan, result.SudahMengisiKSE+result.BelumMengisiKSE)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDashboardRepository_SeStatusCount_AllFilled(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewDashboardRepository(db)

	mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{"total_perusahaan", "sudah_mengisi_kse", "belum_mengisi_kse"}).
			AddRow(int64(50), int64(50), int64(0)))

	result, err := repo.SeStatusCount(context.Background(), dto.DashboardFilter{})
	assert.NoError(t, err)
	assert.Equal(t, int64(50), result.TotalPerusahaan)
	assert.Equal(t, int64(50), result.SudahMengisiKSE)
	assert.Equal(t, int64(0), result.BelumMengisiKSE)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDashboardRepository_SeStatusCount_NoneFilled(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewDashboardRepository(db)

	mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{"total_perusahaan", "sudah_mengisi_kse", "belum_mengisi_kse"}).
			AddRow(int64(100), int64(0), int64(100)))

	result, err := repo.SeStatusCount(context.Background(), dto.DashboardFilter{})
	assert.NoError(t, err)
	assert.Equal(t, int64(0), result.SudahMengisiKSE)
	assert.Equal(t, result.TotalPerusahaan, result.BelumMengisiKSE)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDashboardRepository_SeStatusCount_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewDashboardRepository(db)

	rows := sqlmock.NewRows([]string{"total_perusahaan", "sudah_mengisi_kse", "belum_mengisi_kse"}).
		AddRow("bukan_angka", "bukan_angka", "bukan_angka")

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	_, err = repo.SeStatusCount(context.Background(), dto.DashboardFilter{})
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

/*
=====================================
 TEST INTEGRASI — semua repo bersama
=====================================
*/

func TestDashboardRepository_Integration_AllMethods(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewDashboardRepository(db)
	ctx := context.Background()
	f := dto.DashboardFilter{}

	// CountPerSektor
	mock.ExpectQuery("SELECT s.id, s.nama_sektor, COUNT").
		WillReturnRows(sqlmock.NewRows([]string{"id", "nama_sektor", "total", "this_month"}).
			AddRow("s1", "ILMATE", int64(97), int64(8)).
			AddRow("s2", "Industri Agro", int64(91), int64(10)).
			AddRow("s3", "IKFT", int64(63), int64(6)))

	// SeGlobalAgg
	mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{"total_se", "this_month", "strategis", "tinggi", "rendah"}).
			AddRow(int64(77), int64(8), int64(30), int64(28), int64(19)))

	// SeStatusCount
	mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{"total_perusahaan", "sudah_mengisi_kse", "belum_mengisi_kse"}).
			AddRow(int64(251), int64(77), int64(174)))

	sectors, err := repo.CountPerSektor(ctx, f)
	require.NoError(t, err)
	assert.Len(t, sectors, 3)
	assert.Equal(t, int64(97), sectors[0].Total)

	seAgg, err := repo.SeGlobalAgg(ctx, f)
	require.NoError(t, err)
	assert.Equal(t, int64(77), seAgg.TotalSE)
	assert.Equal(t, seAgg.TotalSE, seAgg.Strategis+seAgg.Tinggi+seAgg.Rendah)

	seStatus, err := repo.SeStatusCount(ctx, f)
	require.NoError(t, err)
	assert.Equal(t, int64(251), seStatus.TotalPerusahaan)
	assert.Equal(t, seStatus.TotalPerusahaan, seStatus.SudahMengisiKSE+seStatus.BelumMengisiKSE)

	assert.NoError(t, mock.ExpectationsWereMet())
}