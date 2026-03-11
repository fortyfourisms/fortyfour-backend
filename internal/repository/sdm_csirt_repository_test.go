package repository

import (
	"database/sql"
	"errors"
	"fortyfour-backend/internal/dto"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupSdmCsirtTest(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *SdmCsirtRepository) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	repo := NewSdmCsirtRepository(db)
	return db, mock, repo
}

func TestNewSdmCsirtRepository(t *testing.T) {
	db, _, _ := setupSdmCsirtTest(t)
	defer db.Close()

	repo := NewSdmCsirtRepository(db)
	assert.NotNil(t, repo)
	assert.Equal(t, db, repo.db)
}

func TestSdmCsirtRepository_Create(t *testing.T) {
	db, mock, repo := setupSdmCsirtTest(t)
	defer db.Close()

	t.Run("success with all fields", func(t *testing.T) {
		idCsirt := "csirt001"
		namaPersonel := "John Doe"
		jabatanCsirt := "Security Analyst"
		jabatanPerusahaan := "IT Manager"
		skill := "Cybersecurity, Network Security"
		sertifikasi := "CISSP, CEH"

		req := dto.CreateSdmCsirtRequest{
			IdCsirt:           &idCsirt,
			NamaPersonel:      &namaPersonel,
			JabatanCsirt:      &jabatanCsirt,
			JabatanPerusahaan: &jabatanPerusahaan,
			Skill:             &skill,
			Sertifikasi:       &sertifikasi,
		}
		id := "sdm001"

		mock.ExpectExec("INSERT INTO sdm_csirt").
			WithArgs(id, idCsirt, namaPersonel, jabatanCsirt, jabatanPerusahaan, skill, sertifikasi).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Create(req, id)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success with empty optional fields", func(t *testing.T) {
		req := dto.CreateSdmCsirtRequest{
			IdCsirt:           nil,
			NamaPersonel:      nil,
			JabatanCsirt:      nil,
			JabatanPerusahaan: nil,
			Skill:             nil,
			Sertifikasi:       nil,
		}
		id := "sdm001"

		// utils.ValueOrEmpty returns empty string for nil pointer
		// But we need to pass the actual pointer value (nil) to the function
		// The function will handle it and convert to empty string
		mock.ExpectExec("INSERT INTO sdm_csirt").
			WithArgs(id, nil, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Create(req, id)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		idCsirt := "csirt001"
		namaPersonel := "John Doe"
		req := dto.CreateSdmCsirtRequest{
			IdCsirt:      &idCsirt,
			NamaPersonel: &namaPersonel,
		}
		id := "sdm001"

		mock.ExpectExec("INSERT INTO sdm_csirt").
			WithArgs(id, &idCsirt, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnError(errors.New("database error"))

		err := repo.Create(req, id)
		assert.Error(t, err)
		assert.Equal(t, "database error", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("foreign key constraint error", func(t *testing.T) {
		idCsirt := "csirt-invalid"
		namaPersonel := "John Doe"
		req := dto.CreateSdmCsirtRequest{
			IdCsirt:      &idCsirt,
			NamaPersonel: &namaPersonel,
		}
		id := "sdm001"

		mock.ExpectExec("INSERT INTO sdm_csirt").
			WithArgs(id, &idCsirt, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnError(errors.New("foreign key constraint fails"))

		err := repo.Create(req, id)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "foreign key")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSdmCsirtRepository_GetAll(t *testing.T) {
	db, mock, repo := setupSdmCsirtTest(t)
	defer db.Close()

	t.Run("success with csirt", func(t *testing.T) {
		createdAt := "2024-01-01 10:00:00"
		updatedAt := "2024-01-02 15:30:00"

		rows := sqlmock.NewRows([]string{
			"s.id", "s.nama_personel", "s.jabatan_csirt", "s.jabatan_perusahaan",
			"s.skill", "s.sertifikasi", "s.created_at", "s.updated_at",
			"c.id", "c.nama_csirt", "c.web_csirt", "c.telepon_csirt",
		}).
			AddRow(
				"sdm001", "John Doe", "Security Analyst", "IT Manager",
				"Cybersecurity", "CISSP", createdAt, updatedAt,
				"csirt-1", "CSIRT ABC", "https://csirt-abc.com", "021-111",
			).
			AddRow(
				"sdm002", "Jane Smith", "Incident Responder", "Security Lead",
				"Forensics", "CEH", createdAt, updatedAt,
				"csirt-2", "CSIRT XYZ", "https://csirt-xyz.com", "021-222",
			)

		mock.ExpectQuery("SELECT (.+) FROM sdm_csirt s LEFT JOIN csirt c").
			WillReturnRows(rows)

		result, err := repo.GetAll()
		assert.NoError(t, err)
		require.Len(t, result, 2)

		// Check first record
		assert.Equal(t, "sdm001", result[0].ID)
		assert.Equal(t, "John Doe", result[0].NamaPersonel)
		assert.Equal(t, "Security Analyst", result[0].JabatanCsirt)
		assert.Equal(t, "IT Manager", result[0].JabatanPerusahaan)
		assert.Equal(t, "Cybersecurity", result[0].Skill)
		assert.Equal(t, "CISSP", result[0].Sertifikasi)
		assert.Equal(t, createdAt, result[0].CreatedAt)
		assert.Equal(t, updatedAt, result[0].UpdatedAt)
		require.NotNil(t, result[0].Csirt)
		assert.Equal(t, "csirt-1", result[0].Csirt.ID)
		assert.Equal(t, "CSIRT ABC", result[0].Csirt.NamaCsirt)
		require.NotNil(t, result[0].Csirt.WebCsirt)
		assert.Equal(t, "https://csirt-abc.com", *result[0].Csirt.WebCsirt)
		require.NotNil(t, result[0].Csirt.TeleponCsirt)
		assert.Equal(t, "021-111", *result[0].Csirt.TeleponCsirt)

		// Check second record
		assert.Equal(t, "sdm002", result[1].ID)
		assert.Equal(t, "Jane Smith", result[1].NamaPersonel)
		require.NotNil(t, result[1].Csirt)
		assert.Equal(t, "csirt-2", result[1].Csirt.ID)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success without csirt (null)", func(t *testing.T) {
		createdAt := "2024-01-01 10:00:00"
		updatedAt := "2024-01-02 15:30:00"

		rows := sqlmock.NewRows([]string{
			"s.id", "s.nama_personel", "s.jabatan_csirt", "s.jabatan_perusahaan",
			"s.skill", "s.sertifikasi", "s.created_at", "s.updated_at",
			"c.id", "c.nama_csirt", "c.web_csirt", "c.telepon_csirt",
		}).
			AddRow(
				"sdm001", "John Doe", "Security Analyst", "IT Manager",
				"Cybersecurity", "CISSP", createdAt, updatedAt,
				sql.NullString{Valid: false}, // c.id - this is the key one to check
				"",                           // c.nama_csirt - empty string for NULL
				"",                           // c.web_csirt - empty string for NULL
				"",                           // c.telepon_csirt - empty string for NULL
			)

		mock.ExpectQuery("SELECT (.+) FROM sdm_csirt s LEFT JOIN csirt c").
			WillReturnRows(rows)

		result, err := repo.GetAll()
		assert.NoError(t, err)
		require.Len(t, result, 1)

		assert.Equal(t, "sdm001", result[0].ID)
		assert.Equal(t, "John Doe", result[0].NamaPersonel)
		assert.Nil(t, result[0].Csirt) // Should be nil when csirtID is not valid

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM sdm_csirt s LEFT JOIN csirt c").
			WillReturnError(errors.New("query error"))

		result, err := repo.GetAll()
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "query error", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("empty result", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"s.id", "s.nama_personel", "s.jabatan_csirt", "s.jabatan_perusahaan",
			"s.skill", "s.sertifikasi", "s.created_at", "s.updated_at",
			"c.id", "c.nama_csirt", "c.web_csirt", "c.telepon_csirt",
		})

		mock.ExpectQuery("SELECT (.+) FROM sdm_csirt s LEFT JOIN csirt c").
			WillReturnRows(rows)

		result, err := repo.GetAll()
		assert.NoError(t, err)
		assert.Empty(t, result) // Should be empty slice, not nil
		assert.Len(t, result, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id"}).
			AddRow("sdm-1")

		mock.ExpectQuery("SELECT (.+) FROM sdm_csirt s LEFT JOIN csirt c").
			WillReturnRows(rows)

		result, err := repo.GetAll()
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSdmCsirtRepository_GetByID(t *testing.T) {
	db, mock, repo := setupSdmCsirtTest(t)
	defer db.Close()

	t.Run("success with csirt", func(t *testing.T) {
		id := "sdm001"
		createdAt := "2024-01-01 10:00:00"
		updatedAt := "2024-01-02 15:30:00"

		row := sqlmock.NewRows([]string{
			"s.id", "s.nama_personel", "s.jabatan_csirt", "s.jabatan_perusahaan",
			"s.skill", "s.sertifikasi", "s.created_at", "s.updated_at",
			"c.id", "c.nama_csirt", "c.web_csirt", "c.telepon_csirt",
		}).
			AddRow(
				id, "John Doe", "Security Analyst", "IT Manager",
				"Cybersecurity", "CISSP", createdAt, updatedAt,
				"csirt-1", "CSIRT ABC", "https://csirt-abc.com", "021-111",
			)

		mock.ExpectQuery("SELECT (.+) FROM sdm_csirt s LEFT JOIN csirt c (.+) WHERE s.id = ?").
			WithArgs(id).
			WillReturnRows(row)

		result, err := repo.GetByID(id)
		assert.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, id, result.ID)
		assert.Equal(t, "John Doe", result.NamaPersonel)
		assert.Equal(t, "Security Analyst", result.JabatanCsirt)
		assert.Equal(t, "IT Manager", result.JabatanPerusahaan)
		assert.Equal(t, "Cybersecurity", result.Skill)
		assert.Equal(t, "CISSP", result.Sertifikasi)
		assert.Equal(t, createdAt, result.CreatedAt)
		assert.Equal(t, updatedAt, result.UpdatedAt)
		require.NotNil(t, result.Csirt)
		assert.Equal(t, "csirt-1", result.Csirt.ID)
		assert.Equal(t, "CSIRT ABC", result.Csirt.NamaCsirt)
		require.NotNil(t, result.Csirt.WebCsirt)
		assert.Equal(t, "https://csirt-abc.com", *result.Csirt.WebCsirt)
		require.NotNil(t, result.Csirt.TeleponCsirt)
		assert.Equal(t, "021-111", *result.Csirt.TeleponCsirt)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success without csirt (null)", func(t *testing.T) {
		id := "sdm001"
		createdAt := "2024-01-01 10:00:00"
		updatedAt := "2024-01-02 15:30:00"

		row := sqlmock.NewRows([]string{
			"s.id", "s.nama_personel", "s.jabatan_csirt", "s.jabatan_perusahaan",
			"s.skill", "s.sertifikasi", "s.created_at", "s.updated_at",
			"c.id", "c.nama_csirt", "c.web_csirt", "c.telepon_csirt",
		}).
			AddRow(
				id, "John Doe", "Security Analyst", "IT Manager",
				"Cybersecurity", "CISSP", createdAt, updatedAt,
				sql.NullString{Valid: false}, // c.id - this determines if csirt is nil
				"",                           // c.nama_csirt - empty string for NULL
				"",                           // c.web_csirt - empty string for NULL
				"",                           // c.telepon_csirt - empty string for NULL
			)

		mock.ExpectQuery("SELECT (.+) FROM sdm_csirt s LEFT JOIN csirt c (.+) WHERE s.id = ?").
			WithArgs(id).
			WillReturnRows(row)

		result, err := repo.GetByID(id)
		assert.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, id, result.ID)
		assert.Equal(t, "John Doe", result.NamaPersonel)
		assert.Nil(t, result.Csirt) // Should be nil when csirtID is not valid
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		id := "sdm-nonexistent"

		mock.ExpectQuery("SELECT (.+) FROM sdm_csirt s LEFT JOIN csirt c (.+) WHERE s.id = ?").
			WithArgs(id).
			WillReturnError(sql.ErrNoRows)

		result, err := repo.GetByID(id)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, sql.ErrNoRows, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error", func(t *testing.T) {
		id := "sdm001"
		row := sqlmock.NewRows([]string{"id"}).
			AddRow(id)

		mock.ExpectQuery("SELECT (.+) FROM sdm_csirt s LEFT JOIN csirt c (.+) WHERE s.id = ?").
			WithArgs(id).
			WillReturnRows(row)

		result, err := repo.GetByID(id)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		id := "sdm001"

		mock.ExpectQuery("SELECT (.+) FROM sdm_csirt s LEFT JOIN csirt c (.+) WHERE s.id = ?").
			WithArgs(id).
			WillReturnError(errors.New("database connection error"))

		result, err := repo.GetByID(id)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "database connection error", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSdmCsirtRepository_Update(t *testing.T) {
	db, mock, repo := setupSdmCsirtTest(t)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		id := "sdm001"
		createdAt := "2024-01-01 10:00:00"
		updatedAt := "2024-01-02 15:30:00"

		sdm := dto.SdmCsirtResponse{
			ID:                id,
			NamaPersonel:      "John Doe Updated",
			JabatanCsirt:      "Senior Security Analyst",
			JabatanPerusahaan: "Security Manager",
			Skill:             "Advanced Cybersecurity",
			Sertifikasi:       "CISSP, CISM",
			CreatedAt:         createdAt,
			UpdatedAt:         updatedAt,
		}

		mock.ExpectExec("UPDATE sdm_csirt SET").
			WithArgs(
				sdm.NamaPersonel,
				sdm.JabatanCsirt,
				sdm.JabatanPerusahaan,
				sdm.Skill,
				sdm.Sertifikasi,
				id,
			).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Update(id, sdm)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		id := "sdm001"
		sdm := dto.SdmCsirtResponse{
			NamaPersonel: "John Doe Updated",
		}

		mock.ExpectExec("UPDATE sdm_csirt SET").
			WithArgs(
				sdm.NamaPersonel,
				sdm.JabatanCsirt,
				sdm.JabatanPerusahaan,
				sdm.Skill,
				sdm.Sertifikasi,
				id,
			).
			WillReturnError(errors.New("update error"))

		err := repo.Update(id, sdm)
		assert.Error(t, err)
		assert.Equal(t, "update error", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("no rows affected - not checked by function", func(t *testing.T) {
		id := "sdm-nonexistent"
		sdm := dto.SdmCsirtResponse{
			NamaPersonel: "John Doe",
		}

		mock.ExpectExec("UPDATE sdm_csirt SET").
			WithArgs(
				sdm.NamaPersonel,
				sdm.JabatanCsirt,
				sdm.JabatanPerusahaan,
				sdm.Skill,
				sdm.Sertifikasi,
				id,
			).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.Update(id, sdm)
		// Function doesn't check rows affected, so no error
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSdmCsirtRepository_Delete(t *testing.T) {
	db, mock, repo := setupSdmCsirtTest(t)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		id := "sdm001"

		mock.ExpectExec("DELETE FROM sdm_csirt WHERE id=\\?").
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Delete(id)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		id := "sdm001"

		mock.ExpectExec("DELETE FROM sdm_csirt WHERE id=\\?").
			WithArgs(id).
			WillReturnError(errors.New("delete error"))

		err := repo.Delete(id)
		assert.Error(t, err)
		assert.Equal(t, "delete error", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("no rows affected - not checked by function", func(t *testing.T) {
		id := "sdm-nonexistent"

		mock.ExpectExec("DELETE FROM sdm_csirt WHERE id=\\?").
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.Delete(id)
		// Function doesn't check rows affected, so no error
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("multiple rows deleted - not checked by function", func(t *testing.T) {
		id := "sdm001"

		mock.ExpectExec("DELETE FROM sdm_csirt WHERE id=\\?").
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 5))

		err := repo.Delete(id)
		// Function doesn't check if more than 1 row was deleted
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSdmCsirtRepository_GetByCsirt(t *testing.T) {
	sdmCols := []string{
		"s.id", "s.nama_personel", "s.jabatan_csirt", "s.jabatan_perusahaan",
		"s.skill", "s.sertifikasi", "s.created_at", "s.updated_at",
		"c.id", "c.nama_csirt", "c.web_csirt", "c.telepon_csirt",
	}

	tests := []struct {
		name    string
		idCsirt string
		mockFn  func(mock sqlmock.Sqlmock)
		wantLen int
		wantErr bool
	}{
		{
			name:    "success - returns SDM milik CSIRT",
			idCsirt: "csirt-1",
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(sdmCols).
					AddRow(
						"sdm-1", "John Doe", "Security Analyst", "IT Manager",
						"Cybersecurity", "CISSP", "2024-01-01 00:00:00", "2024-01-01 00:00:00",
						"csirt-1", "CSIRT ABC", "https://csirt.abc.com", "021-111",
					).
					AddRow(
						"sdm-2", "Jane Smith", "Incident Responder", "Security Lead",
						"Forensics", "CEH", "2024-01-01 00:00:00", "2024-01-01 00:00:00",
						"csirt-1", "CSIRT ABC", "https://csirt.abc.com", "021-111",
					)

				mock.ExpectQuery("SELECT (.+) FROM sdm_csirt s (.+) WHERE s.id_csirt = \\?").
					WithArgs("csirt-1").
					WillReturnRows(rows)
			},
			wantLen: 2,
			wantErr: false,
		},
		{
			name:    "success - CSIRT tidak punya SDM (empty)",
			idCsirt: "csirt-kosong",
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(sdmCols)
				mock.ExpectQuery("SELECT (.+) FROM sdm_csirt s (.+) WHERE s.id_csirt = \\?").
					WithArgs("csirt-kosong").
					WillReturnRows(rows)
			},
			wantLen: 0,
			wantErr: false,
		},
		{
			name:    "success - SDM tanpa CSIRT yang terhubung (LEFT JOIN null)",
			idCsirt: "csirt-2",
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(sdmCols).
					AddRow(
						"sdm-3", "Budi", "Analyst", "Staff IT",
						"Networking", "", "2024-01-01 00:00:00", "2024-01-01 00:00:00",
						sql.NullString{Valid: false}, "", "", "",
					)

				mock.ExpectQuery("SELECT (.+) FROM sdm_csirt s (.+) WHERE s.id_csirt = \\?").
					WithArgs("csirt-2").
					WillReturnRows(rows)
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name:    "error - database error",
			idCsirt: "csirt-1",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM sdm_csirt s (.+) WHERE s.id_csirt = \\?").
					WithArgs("csirt-1").
					WillReturnError(sql.ErrConnDone)
			},
			wantLen: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, repo := setupSdmCsirtTest(t)
			defer db.Close()

			tt.mockFn(mock)

			result, err := repo.GetByCsirt(tt.idCsirt)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.wantLen)
				for _, sdm := range result {
					assert.NotEmpty(t, sdm.ID)
					assert.NotEmpty(t, sdm.NamaPersonel)
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}