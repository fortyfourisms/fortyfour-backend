package repository

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"testing"

	"fortyfour-backend/internal/dto"

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

func TestBuildDateRange(t *testing.T) {
	tests := []struct {
		name       string
		filter     dto.DashboardFilter
		wantFrom   string
		wantTo     string
		wantNilPtr bool
	}{
		{name: "no filter", filter: dto.DashboardFilter{}, wantNilPtr: true},
		{name: "explicit from to", filter: dto.DashboardFilter{From: dashboardStringPtr("2024-03-01"), To: dashboardStringPtr("2024-03-31")}, wantFrom: "2024-03-01", wantTo: "2024-03-31"},
		{name: "year only", filter: dto.DashboardFilter{Year: dashboardStringPtr("2024")}, wantFrom: "2024-01-01", wantTo: "2024-12-31"},
		{name: "year quarter", filter: dto.DashboardFilter{Year: dashboardStringPtr("2024"), Quarter: dashboardStringPtr("2")}, wantFrom: "2024-04-01", wantTo: "2024-06-30"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			from, to := buildDateRange(tt.filter)
			if tt.wantNilPtr {
				assert.Nil(t, from)
				assert.Nil(t, to)
				return
			}
			require.NotNil(t, from)
			require.NotNil(t, to)
			assert.Equal(t, tt.wantFrom, *from)
			assert.Equal(t, tt.wantTo, *to)
		})
	}
}

func TestDashboardRepository_CountPerSektor(t *testing.T) {
	tests := []struct {
		name    string
		filter  dto.DashboardFilter
		args    []driver.Value
		want    int
		wantErr bool
	}{
		{name: "no filter", filter: dto.DashboardFilter{}, args: nil, want: 2},
		{name: "with year", filter: dto.DashboardFilter{Year: dashboardStringPtr("2024")}, args: []driver.Value{"2024-01-01", "2024-12-31"}, want: 1},
		{name: "with sub sektor", filter: dto.DashboardFilter{SubSektorID: dashboardStringPtr("sub-uuid-123")}, args: []driver.Value{"sub-uuid-123"}, want: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := NewDashboardRepository(db)
			rows := sqlmock.NewRows([]string{"id", "nama_sektor", "total", "this_month"}).
				AddRow("sektor-1", "Keuangan", int64(25), int64(5))
			if tt.want == 2 {
				rows.AddRow("sektor-2", "Teknologi", int64(18), int64(3))
			}

			q := mock.ExpectQuery("SELECT s.id, s.nama_sektor, COUNT")
			if tt.args != nil {
				q.WithArgs(tt.args...)
			}
			q.WillReturnRows(rows)

			result, err := repo.CountPerSektor(context.Background(), tt.filter)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.want)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDashboardRepository_IkasGlobalAgg(t *testing.T) {
	tests := []struct {
		name    string
		filter  dto.DashboardFilter
		args    []driver.Value
		want    dto.IkasAgg
		wantErr bool
	}{
		{
			name:   "success no filter",
			filter: dto.DashboardFilter{},
			want:   dto.IkasAgg{Total: 25, AvgNilaiKematangan: 2.75, AvgTargetNilai: 4.00},
		},
		{
			name:   "success with sub sektor",
			filter: dto.DashboardFilter{SubSektorID: dashboardStringPtr("sub-uuid-ikas")},
			args:   []driver.Value{"sub-uuid-ikas"},
			want:   dto.IkasAgg{Total: 10, AvgNilaiKematangan: 2.40, AvgTargetNilai: 3.50},
		},
		{
			name:    "db error",
			filter:  dto.DashboardFilter{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := NewDashboardRepository(db)
			q := mock.ExpectQuery("SELECT")
			if tt.args != nil {
				q.WithArgs(tt.args...)
			}
			if tt.wantErr {
				q.WillReturnError(sql.ErrConnDone)
			} else {
				q.WillReturnRows(sqlmock.NewRows([]string{"total_ikas", "avg_nilai_kematangan", "avg_target_nilai"}).
					AddRow(tt.want.Total, tt.want.AvgNilaiKematangan, tt.want.AvgTargetNilai))
			}

			result, err := repo.IkasGlobalAgg(context.Background(), tt.filter)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, result)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDashboardRepository_SeGlobalAgg(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewDashboardRepository(db)
	mock.ExpectQuery("SELECT").
		WithArgs("Strategis").
		WillReturnRows(sqlmock.NewRows([]string{"total_se", "this_month", "strategis", "tinggi", "rendah"}).
			AddRow(int64(30), int64(3), int64(30), int64(0), int64(0)))

	result, err := repo.SeGlobalAgg(context.Background(), dto.DashboardFilter{KategoriSE: dashboardStringPtr("Strategis")})
	assert.NoError(t, err)
	assert.Equal(t, dto.SeAgg{TotalSE: 30, ThisMonth: 3, Strategis: 30, Tinggi: 0, Rendah: 0}, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDashboardRepository_SeStatusCount(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewDashboardRepository(db)
	mock.ExpectQuery("SELECT").
		WithArgs("sub-uuid-abc").
		WillReturnRows(sqlmock.NewRows([]string{"total_perusahaan", "sudah_mengisi_kse", "belum_mengisi_kse"}).
			AddRow(int64(50), int64(40), int64(10)))

	result, err := repo.SeStatusCount(context.Background(), dto.DashboardFilter{SubSektorID: dashboardStringPtr("sub-uuid-abc")})
	assert.NoError(t, err)
	assert.Equal(t, dto.SeStatusCount{TotalPerusahaan: 50, SudahMengisiKSE: 40, BelumMengisiKSE: 10}, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDashboardRepository_IkasStatusCount(t *testing.T) {
	tests := []struct {
		name    string
		filter  dto.DashboardFilter
		args    []driver.Value
		want    dto.IkasStatusCount
		wantErr bool
	}{
		{
			name:   "success no filter",
			filter: dto.DashboardFilter{},
			want:   dto.IkasStatusCount{TotalPerusahaan: 200, SudahMengisiIKAS: 80, BelumMengisiIKAS: 120},
		},
		{
			name:   "success with sub sektor",
			filter: dto.DashboardFilter{SubSektorID: dashboardStringPtr("sub-uuid-ikas")},
			args:   []driver.Value{"sub-uuid-ikas"},
			want:   dto.IkasStatusCount{TotalPerusahaan: 50, SudahMengisiIKAS: 20, BelumMengisiIKAS: 30},
		},
		{
			name:    "db error",
			filter:  dto.DashboardFilter{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := NewDashboardRepository(db)
			q := mock.ExpectQuery("SELECT")
			if tt.args != nil {
				q.WithArgs(tt.args...)
			}
			if tt.wantErr {
				q.WillReturnError(sql.ErrConnDone)
			} else {
				q.WillReturnRows(sqlmock.NewRows([]string{"total_perusahaan", "sudah_mengisi_ikas", "belum_mengisi_ikas"}).
					AddRow(tt.want.TotalPerusahaan, tt.want.SudahMengisiIKAS, tt.want.BelumMengisiIKAS))
			}

			result, err := repo.IkasStatusCount(context.Background(), tt.filter)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, result)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDashboardRepository_Integration_AllMethods(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewDashboardRepository(db)
	ctx := context.Background()
	f := dto.DashboardFilter{}

	mock.ExpectQuery("SELECT s.id, s.nama_sektor, COUNT").
		WillReturnRows(sqlmock.NewRows([]string{"id", "nama_sektor", "total", "this_month"}).
			AddRow("s1", "ILMATE", int64(97), int64(8)))
	mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{"total_ikas", "avg_nilai_kematangan", "avg_target_nilai"}).
			AddRow(int64(80), 2.65, 4.00))
	mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{"total_se", "this_month", "strategis", "tinggi", "rendah"}).
			AddRow(int64(77), int64(8), int64(30), int64(28), int64(19)))
	mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{"total_perusahaan", "sudah_mengisi_kse", "belum_mengisi_kse"}).
			AddRow(int64(251), int64(77), int64(174)))
	mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{"total_perusahaan", "sudah_mengisi_ikas", "belum_mengisi_ikas"}).
			AddRow(int64(251), int64(80), int64(171)))

	sectors, err := repo.CountPerSektor(ctx, f)
	require.NoError(t, err)
	ikasAgg, err := repo.IkasGlobalAgg(ctx, f)
	require.NoError(t, err)
	seAgg, err := repo.SeGlobalAgg(ctx, f)
	require.NoError(t, err)
	seStatus, err := repo.SeStatusCount(ctx, f)
	require.NoError(t, err)
	ikasStatus, err := repo.IkasStatusCount(ctx, f)
	require.NoError(t, err)

	assert.Len(t, sectors, 1)
	assert.Equal(t, int64(80), ikasAgg.Total)
	assert.Equal(t, seAgg.TotalSE, seAgg.Strategis+seAgg.Tinggi+seAgg.Rendah)
	assert.Equal(t, seStatus.TotalPerusahaan, seStatus.SudahMengisiKSE+seStatus.BelumMengisiKSE)
	assert.Equal(t, ikasStatus.TotalPerusahaan, ikasStatus.SudahMengisiIKAS+ikasStatus.BelumMengisiIKAS)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func dashboardStringPtr(s string) *string { return &s }
