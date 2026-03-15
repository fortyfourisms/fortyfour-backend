package repository

import (
	"database/sql"
	"fortyfour-backend/internal/dto"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSERepository(t *testing.T) {
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewSERepository(db)

	assert.NotNil(t, repo)
}

func TestSERepository_Create(t *testing.T) {
	tests := []struct {
		name       string
		req        dto.CreateSERequest
		id         string
		totalBobot int
		kategori   string
		mockFn     func(mock sqlmock.Sqlmock, req dto.CreateSERequest, id string, totalBobot int, kategori string)
		wantErr    bool
	}{
		{
			name: "success - create SE with all fields (Strategis)",
			req: dto.CreateSERequest{
				IDPerusahaan:                    "perusahaan-123",
				IDSubSektor:                     "sub-sektor-1",
				IDCsirt:                         "csirt-1",
				NilaiInvestasi:                  "A",
				AnggaranOperasional:             "A",
				KepatuhanPeraturan:              "A",
				TeknikKriptografi:               "A",
				JumlahPengguna:                  "A",
				DataPribadi:                     "A",
				KlasifikasiData:                 "A",
				KekritisanProses:                "A",
				DampakKegagalan:                 "A",
				PotensiKerugiandanDampakNegatif: "A",
				NamaSE:                          "Core Banking System",
				IpSE:                            "192.168.1.100",
				AsNumberSE:                      "AS65000",
				PengelolaSE:                     "IT Department",
				FiturSE:                         "Transaction Processing, Account Management",
			},
			id:         uuid.New().String(),
			totalBobot: 100,
			kategori:   "Strategis",
			mockFn: func(mock sqlmock.Sqlmock, req dto.CreateSERequest, id string, totalBobot int, kategori string) {
				mock.ExpectExec("INSERT INTO se").
					WithArgs(id, req.IDPerusahaan, req.IDSubSektor, req.IDCsirt,
						req.NilaiInvestasi, req.AnggaranOperasional, req.KepatuhanPeraturan,
						req.TeknikKriptografi, req.JumlahPengguna, req.DataPribadi,
						req.KlasifikasiData, req.KekritisanProses, req.DampakKegagalan,
						req.PotensiKerugiandanDampakNegatif, req.NamaSE, req.IpSE,
						req.AsNumberSE, req.PengelolaSE, req.FiturSE, totalBobot, kategori).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "success - create SE with minimal fields (Tinggi)",
			req: dto.CreateSERequest{
				IDPerusahaan:                    "perusahaan-456",
				IDSubSektor:                     "",
				IDCsirt:                         "",
				NilaiInvestasi:                  "B",
				AnggaranOperasional:             "B",
				KepatuhanPeraturan:              "B",
				TeknikKriptografi:               "B",
				JumlahPengguna:                  "B",
				DataPribadi:                     "B",
				KlasifikasiData:                 "B",
				KekritisanProses:                "B",
				DampakKegagalan:                 "B",
				PotensiKerugiandanDampakNegatif: "B",
				NamaSE:                          "Internal Portal",
				IpSE:                            "10.0.0.1",
				AsNumberSE:                      "AS64512",
				PengelolaSE:                     "Internal IT",
				FiturSE:                         "",
			},
			id:         uuid.New().String(),
			totalBobot: 70,
			kategori:   "Tinggi",
			mockFn: func(mock sqlmock.Sqlmock, req dto.CreateSERequest, id string, totalBobot int, kategori string) {
				mock.ExpectExec("INSERT INTO se").
					WithArgs(id, req.IDPerusahaan, req.IDSubSektor, req.IDCsirt,
						req.NilaiInvestasi, req.AnggaranOperasional, req.KepatuhanPeraturan,
						req.TeknikKriptografi, req.JumlahPengguna, req.DataPribadi,
						req.KlasifikasiData, req.KekritisanProses, req.DampakKegagalan,
						req.PotensiKerugiandanDampakNegatif, req.NamaSE, req.IpSE,
						req.AsNumberSE, req.PengelolaSE, req.FiturSE, totalBobot, kategori).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "success - create SE with C values (Rendah)",
			req: dto.CreateSERequest{
				IDPerusahaan:                    "perusahaan-789",
				IDSubSektor:                     "sub-3",
				IDCsirt:                         "csirt-3",
				NilaiInvestasi:                  "C",
				AnggaranOperasional:             "C",
				KepatuhanPeraturan:              "C",
				TeknikKriptografi:               "C",
				JumlahPengguna:                  "C",
				DataPribadi:                     "C",
				KlasifikasiData:                 "C",
				KekritisanProses:                "C",
				DampakKegagalan:                 "C",
				PotensiKerugiandanDampakNegatif: "C",
				NamaSE:                          "Test System",
				IpSE:                            "172.16.0.1",
				AsNumberSE:                      "AS64513",
				PengelolaSE:                     "Test Team",
				FiturSE:                         "Basic Features",
			},
			id:         uuid.New().String(),
			totalBobot: 30,
			kategori:   "Rendah",
			mockFn: func(mock sqlmock.Sqlmock, req dto.CreateSERequest, id string, totalBobot int, kategori string) {
				mock.ExpectExec("INSERT INTO se").
					WithArgs(id, req.IDPerusahaan, req.IDSubSektor, req.IDCsirt,
						req.NilaiInvestasi, req.AnggaranOperasional, req.KepatuhanPeraturan,
						req.TeknikKriptografi, req.JumlahPengguna, req.DataPribadi,
						req.KlasifikasiData, req.KekritisanProses, req.DampakKegagalan,
						req.PotensiKerugiandanDampakNegatif, req.NamaSE, req.IpSE,
						req.AsNumberSE, req.PengelolaSE, req.FiturSE, totalBobot, kategori).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "error - database error on insert",
			req: dto.CreateSERequest{
				IDPerusahaan:                    "perusahaan-123",
				NilaiInvestasi:                  "A",
				AnggaranOperasional:             "A",
				KepatuhanPeraturan:              "A",
				TeknikKriptografi:               "A",
				JumlahPengguna:                  "A",
				DataPribadi:                     "A",
				KlasifikasiData:                 "A",
				KekritisanProses:                "A",
				DampakKegagalan:                 "A",
				PotensiKerugiandanDampakNegatif: "A",
				NamaSE:                          "Test System",
				IpSE:                            "192.168.1.1",
				AsNumberSE:                      "AS65000",
				PengelolaSE:                     "IT",
			},
			id:         uuid.New().String(),
			totalBobot: 100,
			kategori:   "Strategis",
			mockFn: func(mock sqlmock.Sqlmock, req dto.CreateSERequest, id string, totalBobot int, kategori string) {
				mock.ExpectExec("INSERT INTO se").
					WithArgs(id, req.IDPerusahaan, req.IDSubSektor, req.IDCsirt,
						req.NilaiInvestasi, req.AnggaranOperasional, req.KepatuhanPeraturan,
						req.TeknikKriptografi, req.JumlahPengguna, req.DataPribadi,
						req.KlasifikasiData, req.KekritisanProses, req.DampakKegagalan,
						req.PotensiKerugiandanDampakNegatif, req.NamaSE, req.IpSE,
						req.AsNumberSE, req.PengelolaSE, req.FiturSE, totalBobot, kategori).
					WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
		{
			name: "error - foreign key constraint (invalid perusahaan_id)",
			req: dto.CreateSERequest{
				IDPerusahaan:                    "invalid-perusahaan",
				NilaiInvestasi:                  "B",
				AnggaranOperasional:             "B",
				KepatuhanPeraturan:              "B",
				TeknikKriptografi:               "B",
				JumlahPengguna:                  "B",
				DataPribadi:                     "B",
				KlasifikasiData:                 "B",
				KekritisanProses:                "B",
				DampakKegagalan:                 "B",
				PotensiKerugiandanDampakNegatif: "B",
				NamaSE:                          "Test System",
				IpSE:                            "10.0.0.1",
				AsNumberSE:                      "AS64512",
				PengelolaSE:                     "IT",
			},
			id:         uuid.New().String(),
			totalBobot: 70,
			kategori:   "Tinggi",
			mockFn: func(mock sqlmock.Sqlmock, req dto.CreateSERequest, id string, totalBobot int, kategori string) {
				mock.ExpectExec("INSERT INTO se").
					WithArgs(id, req.IDPerusahaan, req.IDSubSektor, req.IDCsirt,
						req.NilaiInvestasi, req.AnggaranOperasional, req.KepatuhanPeraturan,
						req.TeknikKriptografi, req.JumlahPengguna, req.DataPribadi,
						req.KlasifikasiData, req.KekritisanProses, req.DampakKegagalan,
						req.PotensiKerugiandanDampakNegatif, req.NamaSE, req.IpSE,
						req.AsNumberSE, req.PengelolaSE, req.FiturSE, totalBobot, kategori).
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

			repo := NewSERepository(db)

			if tt.mockFn != nil {
				tt.mockFn(mock, tt.req, tt.id, tt.totalBobot, tt.kategori)
			}

			err = repo.Create(tt.req, tt.id, tt.totalBobot, tt.kategori)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestSERepository_GetAll(t *testing.T) {
	tests := []struct {
		name    string
		mockFn  func(mock sqlmock.Sqlmock)
		want    int
		wantErr bool
	}{
		{
			name: "success - get all SE with relations",
			mockFn: func(mock sqlmock.Sqlmock) {
				now := time.Now()
				rows := sqlmock.NewRows([]string{
					"id", "id_perusahaan", "id_sub_sektor", "id_csirt",
					"nilai_investasi", "anggaran_operasional", "kepatuhan_peraturan", "teknik_kriptografi",
					"jumlah_pengguna", "data_pribadi", "klasifikasi_data", "kekritisan_proses",
					"dampak_kegagalan", "potensi_kerugian_dan_dampak_negatif",
					"nama_se", "ip_se", "as_number_se", "pengelola_se", "fitur_se",
					"total_bobot", "kategori_se", "created_at", "updated_at",
					"p_id", "nama_perusahaan",
					"ss_id", "nama_sub_sektor", "s_id", "nama_sektor",
					"c_id", "nama_csirt",
				}).
					AddRow("se-1", "perusahaan-1", "sub-1", "csirt-1",
						"A", "A", "A", "A", "A", "A", "A", "A", "A", "A",
						"Core Banking System", "192.168.1.100", "AS65000", "IT Dept", "Transaction Processing",
						100, "Strategis", now, now,
						"perusahaan-1", "PT Bank ABC",
						"sub-1", "Perbankan", "sektor-1", "Keuangan",
						"csirt-1", "CSIRT Bank ABC")

				mock.ExpectQuery("SELECT (.+) FROM se JOIN perusahaan p (.+) ORDER BY se.created_at DESC").
					WillReturnRows(rows)
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "success - get SE without optional relations",
			mockFn: func(mock sqlmock.Sqlmock) {
				now := time.Now()
				rows := sqlmock.NewRows([]string{
					"id", "id_perusahaan", "id_sub_sektor", "id_csirt",
					"nilai_investasi", "anggaran_operasional", "kepatuhan_peraturan", "teknik_kriptografi",
					"jumlah_pengguna", "data_pribadi", "klasifikasi_data", "kekritisan_proses",
					"dampak_kegagalan", "potensi_kerugian_dan_dampak_negatif",
					"nama_se", "ip_se", "as_number_se", "pengelola_se", "fitur_se",
					"total_bobot", "kategori_se", "created_at", "updated_at",
					"p_id", "nama_perusahaan",
					"ss_id", "nama_sub_sektor", "s_id", "nama_sektor",
					"c_id", "nama_csirt",
				}).
					AddRow("se-1", "perusahaan-1", "", "",
						"C", "C", "C", "C", "C", "C", "C", "C", "C", "C",
						"Internal Portal", "10.0.0.1", "AS64512", "Internal IT", "",
						30, "Rendah", now, now,
						"perusahaan-1", "PT ABC",
						"", "", "", "",
						"", "")

				mock.ExpectQuery("SELECT (.+) FROM se JOIN perusahaan p (.+) ORDER BY se.created_at DESC").
					WillReturnRows(rows)
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "success - multiple SE with different categories",
			mockFn: func(mock sqlmock.Sqlmock) {
				now := time.Now()
				rows := sqlmock.NewRows([]string{
					"id", "id_perusahaan", "id_sub_sektor", "id_csirt",
					"nilai_investasi", "anggaran_operasional", "kepatuhan_peraturan", "teknik_kriptografi",
					"jumlah_pengguna", "data_pribadi", "klasifikasi_data", "kekritisan_proses",
					"dampak_kegagalan", "potensi_kerugian_dan_dampak_negatif",
					"nama_se", "ip_se", "as_number_se", "pengelola_se", "fitur_se",
					"total_bobot", "kategori_se", "created_at", "updated_at",
					"p_id", "nama_perusahaan",
					"ss_id", "nama_sub_sektor", "s_id", "nama_sektor",
					"c_id", "nama_csirt",
				}).
					AddRow("se-1", "p1", "sub-1", "csirt-1",
						"A", "A", "A", "A", "A", "A", "A", "A", "A", "A",
						"System 1", "192.168.1.1", "AS1", "IT", "Features 1",
						100, "Strategis", now, now,
						"p1", "PT 1", "sub-1", "Sub 1", "s1", "Sektor 1", "csirt-1", "CSIRT 1").
					AddRow("se-2", "p2", "sub-2", "csirt-2",
						"B", "B", "B", "B", "B", "B", "B", "B", "B", "B",
						"System 2", "192.168.1.2", "AS2", "IT", "Features 2",
						70, "Tinggi", now, now,
						"p2", "PT 2", "sub-2", "Sub 2", "s2", "Sektor 2", "csirt-2", "CSIRT 2").
					AddRow("se-3", "p3", "", "",
						"C", "C", "C", "C", "C", "C", "C", "C", "C", "C",
						"System 3", "192.168.1.3", "AS3", "IT", "Features 3",
						30, "Rendah", now, now,
						"p3", "PT 3", "", "", "", "", "", "")

				mock.ExpectQuery("SELECT (.+) FROM se JOIN perusahaan p (.+) ORDER BY se.created_at DESC").
					WillReturnRows(rows)
			},
			want:    3,
			wantErr: false,
		},
		{
			name: "success - empty result",
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "id_perusahaan", "id_sub_sektor", "id_csirt",
					"nilai_investasi", "anggaran_operasional", "kepatuhan_peraturan", "teknik_kriptografi",
					"jumlah_pengguna", "data_pribadi", "klasifikasi_data", "kekritisan_proses",
					"dampak_kegagalan", "potensi_kerugian_dan_dampak_negatif",
					"nama_se", "ip_se", "as_number_se", "pengelola_se", "fitur_se",
					"total_bobot", "kategori_se", "created_at", "updated_at",
					"p_id", "nama_perusahaan",
					"ss_id", "nama_sub_sektor", "s_id", "nama_sektor",
					"c_id", "nama_csirt",
				})

				mock.ExpectQuery("SELECT (.+) FROM se JOIN perusahaan p (.+) ORDER BY se.created_at DESC").
					WillReturnRows(rows)
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "error - database error",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM se JOIN perusahaan p (.+) ORDER BY se.created_at DESC").
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

			repo := NewSERepository(db)

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

				if tt.want > 0 {
					for _, se := range result {
						assert.NotEmpty(t, se.ID)
						assert.NotEmpty(t, se.NamaSE)
						assert.NotEmpty(t, se.KategoriSE)
						assert.GreaterOrEqual(t, se.TotalBobot, 0)
						assert.NotNil(t, se.Perusahaan)
						assert.NotEmpty(t, se.Perusahaan.NamaPerusahaan)
						// Verify A/B/C values
						assert.Contains(t, []string{"A", "B", "C"}, se.NilaiInvestasi)
					}
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestSERepository_GetByID(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		mockFn  func(mock sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "success - find SE with all relations",
			id:   "se-123",
			mockFn: func(mock sqlmock.Sqlmock) {
				now := time.Now()
				rows := sqlmock.NewRows([]string{
					"id", "id_perusahaan", "id_sub_sektor", "id_csirt",
					"nilai_investasi", "anggaran_operasional", "kepatuhan_peraturan", "teknik_kriptografi",
					"jumlah_pengguna", "data_pribadi", "klasifikasi_data", "kekritisan_proses",
					"dampak_kegagalan", "potensi_kerugian_dan_dampak_negatif",
					"nama_se", "ip_se", "as_number_se", "pengelola_se", "fitur_se",
					"total_bobot", "kategori_se", "created_at", "updated_at",
					"p_id", "nama_perusahaan",
					"ss_id", "nama_sub_sektor", "s_id", "nama_sektor",
					"c_id", "nama_csirt",
				}).
					AddRow("se-123", "perusahaan-1", "sub-1", "csirt-1",
						"A", "A", "A", "A", "A", "A", "A", "A", "A", "A",
						"Core Banking System", "192.168.1.100", "AS65000", "IT Dept", "Transaction Processing",
						100, "Strategis", now, now,
						"perusahaan-1", "PT Bank ABC",
						"sub-1", "Perbankan", "sektor-1", "Keuangan",
						"csirt-1", "CSIRT Bank ABC")

				mock.ExpectQuery("SELECT (.+) FROM se JOIN perusahaan p (.+) WHERE se.id = \\?").
					WithArgs("se-123").
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name: "success - find SE without optional relations",
			id:   "se-456",
			mockFn: func(mock sqlmock.Sqlmock) {
				now := time.Now()
				rows := sqlmock.NewRows([]string{
					"id", "id_perusahaan", "id_sub_sektor", "id_csirt",
					"nilai_investasi", "anggaran_operasional", "kepatuhan_peraturan", "teknik_kriptografi",
					"jumlah_pengguna", "data_pribadi", "klasifikasi_data", "kekritisan_proses",
					"dampak_kegagalan", "potensi_kerugian_dan_dampak_negatif",
					"nama_se", "ip_se", "as_number_se", "pengelola_se", "fitur_se",
					"total_bobot", "kategori_se", "created_at", "updated_at",
					"p_id", "nama_perusahaan",
					"ss_id", "nama_sub_sektor", "s_id", "nama_sektor",
					"c_id", "nama_csirt",
				}).
					AddRow("se-456", "perusahaan-1", "", "",
						"B", "B", "B", "B", "B", "B", "B", "B", "B", "B",
						"Internal Portal", "10.0.0.1", "AS64512", "Internal IT", "",
						70, "Tinggi", now, now,
						"perusahaan-1", "PT ABC",
						"", "", "", "",
						"", "")

				mock.ExpectQuery("SELECT (.+) FROM se JOIN perusahaan p (.+) WHERE se.id = \\?").
					WithArgs("se-456").
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name: "error - SE not found",
			id:   "non-existent",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM se JOIN perusahaan p (.+) WHERE se.id = \\?").
					WithArgs("non-existent").
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
		},
		{
			name: "error - database error",
			id:   "se-123",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM se JOIN perusahaan p (.+) WHERE se.id = \\?").
					WithArgs("se-123").
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

			repo := NewSERepository(db)

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
				assert.NotEmpty(t, result.NamaSE)
				// Verify A/B/C values
				assert.Contains(t, []string{"A", "B", "C"}, result.NilaiInvestasi)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestSERepository_Update(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		req        dto.UpdateSERequest
		totalBobot int
		kategori   string
		mockFn     func(mock sqlmock.Sqlmock)
		wantErr    bool
	}{
		{
			name: "success - update all fields",
			id:   "se-123",
			req: dto.UpdateSERequest{
				NilaiInvestasi:                  stringPtr("A"),
				AnggaranOperasional:             stringPtr("A"),
				KepatuhanPeraturan:              stringPtr("A"),
				TeknikKriptografi:               stringPtr("A"),
				JumlahPengguna:                  stringPtr("A"),
				DataPribadi:                     stringPtr("A"),
				KlasifikasiData:                 stringPtr("A"),
				KekritisanProses:                stringPtr("A"),
				DampakKegagalan:                 stringPtr("A"),
				PotensiKerugiandanDampakNegatif: stringPtr("A"),
				IDPerusahaan:                    stringPtr("perusahaan-updated"),
				IDSubSektor:                     stringPtr("sub-updated"),
				IDCsirt:                         stringPtr("csirt-updated"),
				NamaSE:                          stringPtr("Updated System"),
				IpSE:                            stringPtr("192.168.1.200"),
				AsNumberSE:                      stringPtr("AS65001"),
				PengelolaSE:                     stringPtr("Updated Dept"),
				FiturSE:                         stringPtr("Updated Features"),
			},
			totalBobot: 100,
			kategori:   "Strategis",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE se SET").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "success - update partial fields (only kategori changed)",
			id:   "se-456",
			req: dto.UpdateSERequest{
				NamaSE: stringPtr("Partially Updated System"),
			},
			totalBobot: 70,
			kategori:   "Tinggi",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE se SET").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "success - update from B to C (degradation)",
			id:   "se-789",
			req: dto.UpdateSERequest{
				NilaiInvestasi:      stringPtr("C"),
				AnggaranOperasional: stringPtr("C"),
			},
			totalBobot: 40,
			kategori:   "Rendah",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE se SET").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name:       "error - database error on update",
			id:         "se-123",
			req:        dto.UpdateSERequest{},
			totalBobot: 70,
			kategori:   "Tinggi",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE se SET").
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

			repo := NewSERepository(db)

			if tt.mockFn != nil {
				tt.mockFn(mock)
			}

			err = repo.Update(tt.id, tt.req, tt.totalBobot, tt.kategori)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestSERepository_Delete(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		mockFn  func(mock sqlmock.Sqlmock, id string)
		wantErr bool
	}{
		{
			name: "success - delete SE",
			id:   "se-123",
			mockFn: func(mock sqlmock.Sqlmock, id string) {
				mock.ExpectExec("DELETE FROM se WHERE id = \\?").
					WithArgs(id).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name: "success - delete non-existent SE (no rows affected)",
			id:   "non-existent",
			mockFn: func(mock sqlmock.Sqlmock, id string) {
				mock.ExpectExec("DELETE FROM se WHERE id = \\?").
					WithArgs(id).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			wantErr: false,
		},
		{
			name: "error - database error on delete",
			id:   "se-123",
			mockFn: func(mock sqlmock.Sqlmock, id string) {
				mock.ExpectExec("DELETE FROM se WHERE id = \\?").
					WithArgs(id).
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

			repo := NewSERepository(db)

			if tt.mockFn != nil {
				tt.mockFn(mock, tt.id)
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

// Test for validating A/B/C values
func TestSERepository_ValidateABCValues(t *testing.T) {
	t.Run("verify all A/B/C combinations", func(t *testing.T) {
		combinations := []struct {
			values   []string
			kategori string
			bobot    int
		}{
			{[]string{"A", "A", "A", "A", "A", "A", "A", "A", "A", "A"}, "Strategis", 100},
			{[]string{"B", "B", "B", "B", "B", "B", "B", "B", "B", "B"}, "Tinggi", 70},
			{[]string{"C", "C", "C", "C", "C", "C", "C", "C", "C", "C"}, "Rendah", 30},
			{[]string{"A", "B", "C", "A", "B", "C", "A", "B", "C", "A"}, "Tinggi", 65},
		}

		for i, combo := range combinations {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)

			repo := NewSERepository(db)

			req := dto.CreateSERequest{
				IDPerusahaan:                    "perusahaan-1",
				NilaiInvestasi:                  combo.values[0],
				AnggaranOperasional:             combo.values[1],
				KepatuhanPeraturan:              combo.values[2],
				TeknikKriptografi:               combo.values[3],
				JumlahPengguna:                  combo.values[4],
				DataPribadi:                     combo.values[5],
				KlasifikasiData:                 combo.values[6],
				KekritisanProses:                combo.values[7],
				DampakKegagalan:                 combo.values[8],
				PotensiKerugiandanDampakNegatif: combo.values[9],
				NamaSE:                          "Test SE",
				IpSE:                            "192.168.1.1",
				AsNumberSE:                      "AS65000",
				PengelolaSE:                     "IT",
			}

			id := uuid.New().String()

			mock.ExpectExec("INSERT INTO se").
				WithArgs(id, req.IDPerusahaan, req.IDSubSektor, req.IDCsirt,
					req.NilaiInvestasi, req.AnggaranOperasional, req.KepatuhanPeraturan,
					req.TeknikKriptografi, req.JumlahPengguna, req.DataPribadi,
					req.KlasifikasiData, req.KekritisanProses, req.DampakKegagalan,
					req.PotensiKerugiandanDampakNegatif, req.NamaSE, req.IpSE,
					req.AsNumberSE, req.PengelolaSE, req.FiturSE, combo.bobot, combo.kategori).
				WillReturnResult(sqlmock.NewResult(1, 1))

			err = repo.Create(req, id, combo.bobot, combo.kategori)
			assert.NoError(t, err, "Combination %d should succeed", i)
			assert.NoError(t, mock.ExpectationsWereMet())

			db.Close()
		}
	})
}

func TestSERepository_GetByPerusahaan(t *testing.T) {
	cols := []string{
		"id", "id_perusahaan", "id_sub_sektor", "id_csirt",
		"nilai_investasi", "anggaran_operasional", "kepatuhan_peraturan", "teknik_kriptografi",
		"jumlah_pengguna", "data_pribadi", "klasifikasi_data", "kekritisan_proses",
		"dampak_kegagalan", "potensi_kerugian_dan_dampak_negatif",
		"nama_se", "ip_se", "as_number_se", "pengelola_se", "fitur_se",
		"total_bobot", "kategori_se", "created_at", "updated_at",
		"p_id", "nama_perusahaan",
		"ss_id", "nama_sub_sektor", "s_id", "nama_sektor",
		"c_id", "nama_csirt",
	}

	tests := []struct {
		name         string
		idPerusahaan string
		mockFn       func(mock sqlmock.Sqlmock)
		wantLen      int
		wantErr      bool
	}{
		{
			name:         "success - returns SE milik perusahaan",
			idPerusahaan: "perusahaan-1",
			mockFn: func(mock sqlmock.Sqlmock) {
				now := time.Now()
				rows := sqlmock.NewRows(cols).
					AddRow("se-1", "perusahaan-1", "sub-1", "csirt-1",
						"A", "A", "A", "A", "A", "A", "A", "A", "A", "A",
						"Core Banking", "192.168.1.1", "AS1", "IT", "Features",
						100, "Strategis", now, now,
						"perusahaan-1", "PT ABC",
						"sub-1", "Perbankan", "s-1", "Keuangan",
						"csirt-1", "CSIRT ABC").
					AddRow("se-2", "perusahaan-1", "", "",
						"B", "B", "B", "B", "B", "B", "B", "B", "B", "B",
						"Internal Portal", "10.0.0.1", "AS2", "IT", "",
						70, "Tinggi", now, now,
						"perusahaan-1", "PT ABC",
						"", "", "", "",
						"", "")

				mock.ExpectQuery("SELECT (.+) FROM se JOIN perusahaan p (.+) WHERE se.id_perusahaan = \\?").
					WithArgs("perusahaan-1").
					WillReturnRows(rows)
			},
			wantLen: 2,
			wantErr: false,
		},
		{
			name:         "success - perusahaan tidak punya SE (empty)",
			idPerusahaan: "perusahaan-kosong",
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(cols)
				mock.ExpectQuery("SELECT (.+) FROM se JOIN perusahaan p (.+) WHERE se.id_perusahaan = \\?").
					WithArgs("perusahaan-kosong").
					WillReturnRows(rows)
			},
			wantLen: 0,
			wantErr: false,
		},
		{
			name:         "error - database error",
			idPerusahaan: "perusahaan-1",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM se JOIN perusahaan p (.+) WHERE se.id_perusahaan = \\?").
					WithArgs("perusahaan-1").
					WillReturnError(sql.ErrConnDone)
			},
			wantLen: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := NewSERepository(db)
			tt.mockFn(mock)

			result, err := repo.GetByPerusahaan(tt.idPerusahaan)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.wantLen)
				for _, se := range result {
					assert.NotEmpty(t, se.ID)
					assert.NotNil(t, se.Perusahaan)
					assert.Equal(t, tt.idPerusahaan, se.IDPerusahaan)
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}