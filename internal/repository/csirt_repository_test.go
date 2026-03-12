package repository

import (
	"database/sql"
	"errors"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/models"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTest(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *CsirtRepository) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	repo := NewCsirtRepository(db)
	return db, mock, repo
}

func TestNewCsirtRepository(t *testing.T) {
	db, _, _ := setupTest(t)
	defer db.Close()

	repo := NewCsirtRepository(db)
	assert.NotNil(t, repo)
	assert.Equal(t, db, repo.db)
}

func TestCsirtRepository_Create(t *testing.T) {
	db, mock, repo := setupTest(t)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		req := dto.CreateCsirtRequest{
			IdPerusahaan:     "perusahaan-123",
			NamaCsirt:        "CSIRT Test",
			WebCsirt:         "https://csirt.test.com",
			TeleponCsirt:     "021-12345678",
			PhotoCsirt:       "photo.jpg",
			FileRFC2350:      "rfc2350.pdf",
			FilePublicKeyPGP: "pgp.key",
		}
		id := "csirt-123"

		mock.ExpectExec("INSERT INTO csirt").
			WithArgs(
				id,
				req.IdPerusahaan,
				req.NamaCsirt,
				req.WebCsirt,
				req.TeleponCsirt,
				req.PhotoCsirt,
				req.FileRFC2350,
				req.FilePublicKeyPGP,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Create(req, id)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		req := dto.CreateCsirtRequest{
			IdPerusahaan: "perusahaan-123",
			NamaCsirt:    "CSIRT Test",
		}
		id := "csirt-123"

		mock.ExpectExec("INSERT INTO csirt").
			WithArgs(
				id,
				req.IdPerusahaan,
				req.NamaCsirt,
				req.WebCsirt,
				req.TeleponCsirt,
				req.PhotoCsirt,
				req.FileRFC2350,
				req.FilePublicKeyPGP,
			).
			WillReturnError(errors.New("database error"))

		err := repo.Create(req, id)
		assert.Error(t, err)
		assert.Equal(t, "database error", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCsirtRepository_GetAll(t *testing.T) {
	db, mock, repo := setupTest(t)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"id", "id_perusahaan", "nama_csirt", "web_csirt", "telepon_csirt",
			"photo_csirt", "file_rfc2350", "file_public_key_pgp",
		}).
			AddRow("csirt-1", "perusahaan-1", "CSIRT 1", "https://csirt1.com", "021-111",
				"photo1.jpg", "rfc1.pdf", "pgp1.key").
			AddRow("csirt-2", "perusahaan-2", "CSIRT 2", "https://csirt2.com", "021-222",
				"photo2.jpg", "rfc2.pdf", "pgp2.key")

		mock.ExpectQuery("SELECT (.+) FROM csirt").WillReturnRows(rows)

		result, err := repo.GetAll()
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "csirt-1", result[0].ID)
		assert.Equal(t, "CSIRT 1", result[0].NamaCsirt)
		assert.Equal(t, "csirt-2", result[1].ID)
		assert.Equal(t, "CSIRT 2", result[1].NamaCsirt)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM csirt").
			WillReturnError(errors.New("query error"))

		result, err := repo.GetAll()
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "query error", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id"}).
			AddRow("csirt-1")

		mock.ExpectQuery("SELECT (.+) FROM csirt").WillReturnRows(rows)

		result, err := repo.GetAll()
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("empty result", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"id", "id_perusahaan", "nama_csirt", "web_csirt", "telepon_csirt",
			"photo_csirt", "file_rfc2350", "file_public_key_pgp",
		})

		mock.ExpectQuery("SELECT (.+) FROM csirt").WillReturnRows(rows)

		result, err := repo.GetAll()
		assert.NoError(t, err)
		assert.Empty(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCsirtRepository_GetByID(t *testing.T) {
	db, mock, repo := setupTest(t)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		id := "csirt-123"
		row := sqlmock.NewRows([]string{
			"id", "id_perusahaan", "nama_csirt", "web_csirt", "telepon_csirt",
			"photo_csirt", "file_rfc2350", "file_public_key_pgp",
		}).
			AddRow(id, "perusahaan-1", "CSIRT Test", "https://csirt.test.com", "021-12345",
				"photo.jpg", "rfc.pdf", "pgp.key")

		mock.ExpectQuery("SELECT (.+) FROM csirt WHERE id = ?").
			WithArgs(id).
			WillReturnRows(row)

		result, err := repo.GetByID(id)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, id, result.ID)
		assert.Equal(t, "CSIRT Test", result.NamaCsirt)
		assert.Equal(t, "perusahaan-1", result.IdPerusahaan)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		id := "csirt-nonexistent"

		mock.ExpectQuery("SELECT (.+) FROM csirt WHERE id = ?").
			WithArgs(id).
			WillReturnError(sql.ErrNoRows)

		result, err := repo.GetByID(id)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, sql.ErrNoRows, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error", func(t *testing.T) {
		id := "csirt-123"
		row := sqlmock.NewRows([]string{"id"}).AddRow(id)

		mock.ExpectQuery("SELECT (.+) FROM csirt WHERE id = ?").
			WithArgs(id).
			WillReturnRows(row)

		result, err := repo.GetByID(id)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCsirtRepository_GetAllWithPerusahaan(t *testing.T) {
	db, mock, repo := setupTest(t)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		now := time.Now()
		rows := sqlmock.NewRows([]string{
			"c.id", "c.nama_csirt", "c.web_csirt", "c.telepon_csirt",
			"c.photo_csirt", "c.file_rfc2350", "c.file_public_key_pgp",
			"p.id", "p.photo", "p.nama_perusahaan",
			"p.alamat", "p.telepon", "p.email", "p.website",
			"p.created_at", "p.updated_at",
			"ss.id", "ss.nama_sub_sektor", "ss.id_sektor", "ss.created_at", "ss.updated_at",
			"s.nama_sektor",
		}).
			AddRow(
				"csirt-1", "CSIRT 1", "https://csirt1.com", "021-111",
				"photo1.jpg", "rfc1.pdf", "pgp1.key",
				"perusahaan-1", "logo1.png", "PT Test 1",
				"Jl. Test 1", "021-1111", "test1@test.com", "https://test1.com",
				now, now,
				"sub-1", "Sub Sektor 1", "sektor-1", now.Format(time.RFC3339), now.Format(time.RFC3339),
				"Sektor 1",
			).
			AddRow(
				"csirt-2", "CSIRT 2", "https://csirt2.com", "021-222",
				"photo2.jpg", "rfc2.pdf", "pgp2.key",
				"perusahaan-2", "logo2.png", "PT Test 2",
				"Jl. Test 2", "021-2222", "test2@test.com", "https://test2.com",
				now, now,
				nil, nil, nil, nil, nil,
				nil,
			)

		mock.ExpectQuery("SELECT (.+) FROM csirt c JOIN perusahaan p").
			WillReturnRows(rows)

		result, err := repo.GetAllWithPerusahaan()
		assert.NoError(t, err)
		assert.Len(t, result, 2)

		assert.Equal(t, "csirt-1", result[0].ID)
		assert.Equal(t, "CSIRT 1", result[0].NamaCsirt)
		assert.Equal(t, "perusahaan-1", result[0].Perusahaan.ID)
		assert.Equal(t, "PT Test 1", result[0].Perusahaan.NamaPerusahaan)
		assert.NotNil(t, result[0].Perusahaan.SubSektor)
		assert.Equal(t, "Sub Sektor 1", result[0].Perusahaan.SubSektor.NamaSubSektor)

		assert.Equal(t, "csirt-2", result[1].ID)
		assert.Equal(t, "CSIRT 2", result[1].NamaCsirt)
		assert.Equal(t, "perusahaan-2", result[1].Perusahaan.ID)
		assert.Equal(t, "PT Test 2", result[1].Perusahaan.NamaPerusahaan)
		assert.Nil(t, result[1].Perusahaan.SubSektor)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM csirt c JOIN perusahaan p").
			WillReturnError(errors.New("query error"))

		result, err := repo.GetAllWithPerusahaan()
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id"}).AddRow("csirt-1")

		mock.ExpectQuery("SELECT (.+) FROM csirt c JOIN perusahaan p").
			WillReturnRows(rows)

		result, err := repo.GetAllWithPerusahaan()
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCsirtRepository_GetByIDWithPerusahaan(t *testing.T) {
	db, mock, repo := setupTest(t)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		id := "csirt-123"
		now := time.Now()
		row := sqlmock.NewRows([]string{
			"c.id", "c.nama_csirt", "c.web_csirt", "c.telepon_csirt",
			"c.photo_csirt", "c.file_rfc2350", "c.file_public_key_pgp",
			"p.id", "p.photo", "p.nama_perusahaan",
			"p.alamat", "p.telepon", "p.email", "p.website",
			"p.created_at", "p.updated_at",
			"ss.id", "ss.nama_sub_sektor", "ss.id_sektor", "ss.created_at", "ss.updated_at",
			"s.nama_sektor",
		}).
			AddRow(
				id, "CSIRT Test", "https://csirt.test.com", "021-12345",
				"photo.jpg", "rfc.pdf", "pgp.key",
				"perusahaan-1", "logo.png", "PT Test",
				"Jl. Test", "021-1111", "test@test.com", "https://test.com",
				now, now,
				"sub-1", "Sub Sektor 1", "sektor-1", now.Format(time.RFC3339), now.Format(time.RFC3339),
				"Sektor 1",
			)

		mock.ExpectQuery("SELECT (.+) FROM csirt c JOIN perusahaan p (.+) WHERE c.id = ?").
			WithArgs(id).
			WillReturnRows(row)

		result, err := repo.GetByIDWithPerusahaan(id)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, id, result.ID)
		assert.Equal(t, "CSIRT Test", result.NamaCsirt)
		assert.Equal(t, "perusahaan-1", result.Perusahaan.ID)
		assert.Equal(t, "PT Test", result.Perusahaan.NamaPerusahaan)
		assert.NotNil(t, result.Perusahaan.SubSektor)
		assert.Equal(t, "Sub Sektor 1", result.Perusahaan.SubSektor.NamaSubSektor)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		id := "csirt-nonexistent"

		mock.ExpectQuery("SELECT (.+) FROM csirt c JOIN perusahaan p (.+) WHERE c.id = ?").
			WithArgs(id).
			WillReturnError(sql.ErrNoRows)

		result, err := repo.GetByIDWithPerusahaan(id)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, sql.ErrNoRows, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error", func(t *testing.T) {
		id := "csirt-123"
		row := sqlmock.NewRows([]string{"id"}).AddRow(id)

		mock.ExpectQuery("SELECT (.+) FROM csirt c JOIN perusahaan p (.+) WHERE c.id = ?").
			WithArgs(id).
			WillReturnRows(row)

		result, err := repo.GetByIDWithPerusahaan(id)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCsirtRepository_Update(t *testing.T) {
	db, mock, repo := setupTest(t)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		id := "csirt-123"
		telepon := "021-98763"
		photo := "updated.jpg"
		rfc := "updated_rfc.pdf"
		pgp := "updated_pgp.key"
		csirt := models.Csirt{
			NamaCsirt:        "CSIRT Updated",
			WebCsirt:         "https://csirt.updated.com",
			TeleponCsirt:     &telepon,
			PhotoCsirt:       &photo,
			FileRFC2350:      &rfc,
			FilePublicKeyPGP: &pgp,
		}

		mock.ExpectExec("UPDATE csirt SET").
			WithArgs(
				csirt.NamaCsirt,
				csirt.WebCsirt,
				csirt.TeleponCsirt,
				csirt.PhotoCsirt,
				csirt.FileRFC2350,
				csirt.FilePublicKeyPGP,
				id,
			).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Update(id, csirt)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		id := "csirt-123"
		csirt := models.Csirt{
			NamaCsirt: "CSIRT Updated",
		}

		mock.ExpectExec("UPDATE csirt SET").
			WithArgs(
				csirt.NamaCsirt,
				csirt.WebCsirt,
				csirt.TeleponCsirt,
				csirt.PhotoCsirt,
				csirt.FileRFC2350,
				csirt.FilePublicKeyPGP,
				id,
			).
			WillReturnError(errors.New("update error"))

		err := repo.Update(id, csirt)
		assert.Error(t, err)
		assert.Equal(t, "update error", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("no rows affected", func(t *testing.T) {
		id := "csirt-nonexistent"
		csirt := models.Csirt{
			NamaCsirt: "CSIRT Updated",
		}

		mock.ExpectExec("UPDATE csirt SET").
			WithArgs(
				csirt.NamaCsirt,
				csirt.WebCsirt,
				csirt.TeleponCsirt,
				csirt.PhotoCsirt,
				csirt.FileRFC2350,
				csirt.FilePublicKeyPGP,
				id,
			).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.Update(id, csirt)
		assert.NoError(t, err) // Function doesn't check rows affected
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCsirtRepository_Delete(t *testing.T) {
	db, mock, repo := setupTest(t)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		id := "csirt-123"

		mock.ExpectExec("DELETE FROM csirt WHERE id = ?").
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Delete(id)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		id := "csirt-123"

		mock.ExpectExec("DELETE FROM csirt WHERE id = ?").
			WithArgs(id).
			WillReturnError(errors.New("delete error"))

		err := repo.Delete(id)
		assert.Error(t, err)
		assert.Equal(t, "delete error", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("no rows affected", func(t *testing.T) {
		id := "csirt-nonexistent"

		mock.ExpectExec("DELETE FROM csirt WHERE id = ?").
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.Delete(id)
		assert.NoError(t, err) // Function doesn't check rows affected
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCsirtRepository_GetByPerusahaan(t *testing.T) {
	csirtCols := []string{
		"c.id", "c.nama_csirt", "c.web_csirt", "c.telepon_csirt",
		"c.photo_csirt", "c.file_rfc2350", "c.file_public_key_pgp",
		"p.id", "p.photo", "p.nama_perusahaan",
		"p.alamat", "p.telepon", "p.email", "p.website",
		"p.created_at", "p.updated_at",
		"ss.id", "ss.nama_sub_sektor", "ss.id_sektor", "ss.created_at", "ss.updated_at",
		"s.nama_sektor",
	}

	tests := []struct {
		name         string
		idPerusahaan string
		mockFn       func(mock sqlmock.Sqlmock)
		wantLen      int
		wantErr      bool
	}{
		{
			name:         "success - returns CSIRT milik perusahaan dengan sub sektor",
			idPerusahaan: "perusahaan-1",
			mockFn: func(mock sqlmock.Sqlmock) {
				now := time.Now()
				rows := sqlmock.NewRows(csirtCols).
					AddRow(
						"csirt-1", "CSIRT ABC", "https://csirt.abc.com", "021-111",
						"photo.jpg", "rfc.pdf", "pgp.key",
						"perusahaan-1", "logo.png", "PT ABC",
						"Jl. Test 1", "021-1111", "info@abc.com", "https://abc.com",
						now, now,
						"sub-1", "Perbankan", "sektor-1", now.Format(time.RFC3339), now.Format(time.RFC3339),
						"Keuangan",
					)

				mock.ExpectQuery("SELECT (.+) FROM csirt c JOIN perusahaan p (.+) WHERE c.id_perusahaan = \\?").
					WithArgs("perusahaan-1").
					WillReturnRows(rows)
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name:         "success - returns CSIRT tanpa sub sektor (nil)",
			idPerusahaan: "perusahaan-2",
			mockFn: func(mock sqlmock.Sqlmock) {
				now := time.Now()
				rows := sqlmock.NewRows(csirtCols).
					AddRow(
						"csirt-2", "CSIRT XYZ", "https://csirt.xyz.com", "021-222",
						"photo2.jpg", "rfc2.pdf", "pgp2.key",
						"perusahaan-2", "logo2.png", "PT XYZ",
						"Jl. Test 2", "021-2222", "info@xyz.com", "https://xyz.com",
						now, now,
						nil, nil, nil, nil, nil,
						nil,
					)

				mock.ExpectQuery("SELECT (.+) FROM csirt c JOIN perusahaan p (.+) WHERE c.id_perusahaan = \\?").
					WithArgs("perusahaan-2").
					WillReturnRows(rows)
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name:         "success - perusahaan tidak punya CSIRT (empty)",
			idPerusahaan: "perusahaan-kosong",
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(csirtCols)
				mock.ExpectQuery("SELECT (.+) FROM csirt c JOIN perusahaan p (.+) WHERE c.id_perusahaan = \\?").
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
				mock.ExpectQuery("SELECT (.+) FROM csirt c JOIN perusahaan p (.+) WHERE c.id_perusahaan = \\?").
					WithArgs("perusahaan-1").
					WillReturnError(sql.ErrConnDone)
			},
			wantLen: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, repo := setupTest(t)
			defer db.Close()

			tt.mockFn(mock)

			result, err := repo.GetByPerusahaan(tt.idPerusahaan)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.wantLen)
				for _, csirt := range result {
					assert.NotEmpty(t, csirt.ID)
					assert.NotEmpty(t, csirt.NamaCsirt)
					assert.NotNil(t, csirt.Perusahaan)
					assert.Equal(t, tt.idPerusahaan, csirt.Perusahaan.ID)
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}