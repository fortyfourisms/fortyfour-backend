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

// ─────────────────────────────────────────────────────────
// CREATE
// ─────────────────────────────────────────────────────────

func TestCsirtRepository_Create(t *testing.T) {
	db, mock, repo := setupTest(t)
	defer db.Close()

	t.Run("success — semua field terisi termasuk yang baru", func(t *testing.T) {
		req := dto.CreateCsirtRequest{
			IdPerusahaan:           "perusahaan-123",
			NamaCsirt:              "CSIRT Test",
			WebCsirt:               "https://csirt.test.com",
			TeleponCsirt:           "021-12345678",
			PhotoCsirt:             "photo.jpg",
			FileRFC2350:            "rfc2350.pdf",
			FilePublicKeyPGP:       "pgp.key",
			FileStr:                "uploads/str_csirt/str.pdf",
			TanggalRegistrasi:      "2024-01-15",
			TanggalKadaluarsa:      "2025-01-15",
			TanggalRegistrasiUlang: "2025-01-20",
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
				req.FileStr,                // nullable — terisi
				req.TanggalRegistrasi,      // nullable — terisi
				req.TanggalKadaluarsa,      // nullable — terisi
				req.TanggalRegistrasiUlang, // nullable — terisi
			).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Create(req, id)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success — field nullable kosong → NULL di DB", func(t *testing.T) {
		req := dto.CreateCsirtRequest{
			IdPerusahaan: "perusahaan-123",
			NamaCsirt:    "CSIRT Minimal",
			// FileStr, Tanggal* sengaja dikosongkan
		}
		id := "csirt-124"

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
				nil, // FileStr kosong → nullableStr → nil
				nil, // TanggalRegistrasi kosong → nil
				nil, // TanggalKadaluarsa kosong → nil
				nil, // TanggalRegistrasiUlang kosong → nil
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
			WillReturnError(errors.New("database error"))

		err := repo.Create(req, id)
		assert.Error(t, err)
		assert.Equal(t, "database error", err.Error())
	})
}

// ─────────────────────────────────────────────────────────
// GET ALL
// ─────────────────────────────────────────────────────────

func TestCsirtRepository_GetAll(t *testing.T) {
	db, mock, repo := setupTest(t)
	defer db.Close()

	newCols := []string{
		"id", "id_perusahaan", "nama_csirt", "web_csirt", "telepon_csirt",
		"photo_csirt", "file_rfc2350", "file_public_key_pgp",
		"file_str", "tanggal_registrasi", "tanggal_kadaluarsa", "tanggal_registrasi_ulang",
	}

	t.Run("success — field baru terisi", func(t *testing.T) {
		rows := sqlmock.NewRows(newCols).
			AddRow("csirt-1", "perusahaan-1", "CSIRT 1", "https://csirt1.com", "021-111",
				"photo1.jpg", "rfc1.pdf", "pgp1.key",
				"uploads/str.pdf", "2024-01-15", "2025-01-15", "2025-01-20").
			AddRow("csirt-2", "perusahaan-2", "CSIRT 2", "https://csirt2.com", "021-222",
				"photo2.jpg", "rfc2.pdf", "pgp2.key",
				nil, nil, nil, nil) // field nullable kosong

		mock.ExpectQuery("SELECT (.+) FROM csirt").WillReturnRows(rows)

		result, err := repo.GetAll()
		assert.NoError(t, err)
		assert.Len(t, result, 2)

		// csirt-1: field baru terisi
		assert.NotNil(t, result[0].FileStr)
		assert.Equal(t, "uploads/str.pdf", *result[0].FileStr)
		assert.NotNil(t, result[0].TanggalRegistrasi)
		assert.Equal(t, "2024-01-15", *result[0].TanggalRegistrasi)
		assert.NotNil(t, result[0].TanggalKadaluarsa)
		assert.Equal(t, "2025-01-15", *result[0].TanggalKadaluarsa)
		assert.NotNil(t, result[0].TanggalRegistrasiUlang)
		assert.Equal(t, "2025-01-20", *result[0].TanggalRegistrasiUlang)

		// csirt-2: field nullable → nil
		assert.Nil(t, result[1].FileStr)
		assert.Nil(t, result[1].TanggalRegistrasi)
		assert.Nil(t, result[1].TanggalKadaluarsa)
		assert.Nil(t, result[1].TanggalRegistrasiUlang)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM csirt").
			WillReturnError(errors.New("query error"))

		result, err := repo.GetAll()
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("scan error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id"}).AddRow("csirt-1")
		mock.ExpectQuery("SELECT (.+) FROM csirt").WillReturnRows(rows)

		result, err := repo.GetAll()
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("empty result", func(t *testing.T) {
		rows := sqlmock.NewRows(newCols)
		mock.ExpectQuery("SELECT (.+) FROM csirt").WillReturnRows(rows)

		result, err := repo.GetAll()
		assert.NoError(t, err)
		assert.Empty(t, result)
	})
}

// ─────────────────────────────────────────────────────────
// GET BY ID
// ─────────────────────────────────────────────────────────

func TestCsirtRepository_GetByID(t *testing.T) {
	db, mock, repo := setupTest(t)
	defer db.Close()

	newCols := []string{
		"id", "id_perusahaan", "nama_csirt", "web_csirt", "telepon_csirt",
		"photo_csirt", "file_rfc2350", "file_public_key_pgp",
		"file_str", "tanggal_registrasi", "tanggal_kadaluarsa", "tanggal_registrasi_ulang",
	}

	t.Run("success — field baru terisi", func(t *testing.T) {
		id := "csirt-123"
		row := sqlmock.NewRows(newCols).
			AddRow(id, "perusahaan-1", "CSIRT Test", "https://csirt.test.com", "021-12345",
				"photo.jpg", "rfc.pdf", "pgp.key",
				"uploads/str_csirt/str.pdf", "2024-03-01", "2025-03-01", "2025-03-05")

		mock.ExpectQuery("SELECT (.+) FROM csirt WHERE id = ?").
			WithArgs(id).
			WillReturnRows(row)

		result, err := repo.GetByID(id)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, id, result.ID)
		assert.NotNil(t, result.FileStr)
		assert.Equal(t, "uploads/str_csirt/str.pdf", *result.FileStr)
		assert.NotNil(t, result.TanggalRegistrasi)
		assert.Equal(t, "2024-03-01", *result.TanggalRegistrasi)
		assert.NotNil(t, result.TanggalKadaluarsa)
		assert.Equal(t, "2025-03-01", *result.TanggalKadaluarsa)
		assert.NotNil(t, result.TanggalRegistrasiUlang)
		assert.Equal(t, "2025-03-05", *result.TanggalRegistrasiUlang)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success — field nullable nil", func(t *testing.T) {
		id := "csirt-456"
		row := sqlmock.NewRows(newCols).
			AddRow(id, "perusahaan-1", "CSIRT Null", "https://csirt.com", nil,
				nil, nil, nil,
				nil, nil, nil, nil)

		mock.ExpectQuery("SELECT (.+) FROM csirt WHERE id = ?").
			WithArgs(id).
			WillReturnRows(row)

		result, err := repo.GetByID(id)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Nil(t, result.FileStr)
		assert.Nil(t, result.TanggalRegistrasi)
		assert.Nil(t, result.TanggalKadaluarsa)
		assert.Nil(t, result.TanggalRegistrasiUlang)
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
	})
}

// ─────────────────────────────────────────────────────────
// GET ALL WITH PERUSAHAAN (JOIN)
// ─────────────────────────────────────────────────────────

func TestCsirtRepository_GetAllWithPerusahaan(t *testing.T) {
	db, mock, repo := setupTest(t)
	defer db.Close()

	joinCols := []string{
		"c.id", "c.nama_csirt", "c.web_csirt", "c.telepon_csirt",
		"c.photo_csirt", "c.file_rfc2350", "c.file_public_key_pgp",
		"c.file_str", "c.tanggal_registrasi", "c.tanggal_kadaluarsa", "c.tanggal_registrasi_ulang",
		"p.id", "p.photo", "p.nama_perusahaan",
		"p.alamat", "p.telepon", "p.email", "p.website",
		"p.created_at", "p.updated_at",
		"ss.id", "ss.nama_sub_sektor", "ss.id_sektor", "ss.created_at", "ss.updated_at",
		"s.nama_sektor",
	}

	t.Run("success — field baru terisi + ada sub_sektor", func(t *testing.T) {
		now := time.Now()
		rows := sqlmock.NewRows(joinCols).
			AddRow(
				"csirt-1", "CSIRT 1", "https://csirt1.com", "021-111",
				"photo1.jpg", "rfc1.pdf", "pgp1.key",
				"uploads/str_csirt/abc.pdf", "2024-01-15", "2025-01-15", "2025-01-20",
				"perusahaan-1", "logo1.png", "PT Test 1",
				"Jl. Test 1", "021-1111", "test1@test.com", "https://test1.com",
				now, now,
				"sub-1", "Sub Sektor 1", "sektor-1", now.Format(time.RFC3339), now.Format(time.RFC3339),
				"Sektor 1",
			)

		mock.ExpectQuery("SELECT (.+) FROM csirt c JOIN perusahaan p").
			WillReturnRows(rows)

		result, err := repo.GetAllWithPerusahaan()
		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "uploads/str_csirt/abc.pdf", result[0].FileStr)
		assert.Equal(t, "2024-01-15", result[0].TanggalRegistrasi)
		assert.Equal(t, "2025-01-15", result[0].TanggalKadaluarsa)
		assert.Equal(t, "2025-01-20", result[0].TanggalRegistrasiUlang)
		assert.NotNil(t, result[0].Perusahaan.SubSektor)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success — field baru NULL + tanpa sub_sektor", func(t *testing.T) {
		now := time.Now()
		rows := sqlmock.NewRows(joinCols).
			AddRow(
				"csirt-2", "CSIRT 2", "https://csirt2.com", "021-222",
				"photo2.jpg", "rfc2.pdf", "pgp2.key",
				nil, nil, nil, nil, // field baru semua NULL
				"perusahaan-2", "logo2.png", "PT Test 2",
				"Jl. Test 2", "021-2222", "test2@test.com", "https://test2.com",
				now, now,
				nil, nil, nil, nil, nil, // sub_sektor NULL
				nil,
			)

		mock.ExpectQuery("SELECT (.+) FROM csirt c JOIN perusahaan p").
			WillReturnRows(rows)

		result, err := repo.GetAllWithPerusahaan()
		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "", result[0].FileStr)
		assert.Equal(t, "", result[0].TanggalRegistrasi)
		assert.Equal(t, "", result[0].TanggalKadaluarsa)
		assert.Equal(t, "", result[0].TanggalRegistrasiUlang)
		assert.Nil(t, result[0].Perusahaan.SubSektor)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM csirt c JOIN perusahaan p").
			WillReturnError(errors.New("query error"))

		result, err := repo.GetAllWithPerusahaan()
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("scan error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id"}).AddRow("csirt-1")
		mock.ExpectQuery("SELECT (.+) FROM csirt c JOIN perusahaan p").
			WillReturnRows(rows)

		result, err := repo.GetAllWithPerusahaan()
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

// ─────────────────────────────────────────────────────────
// GET BY ID WITH PERUSAHAAN (JOIN)
// ─────────────────────────────────────────────────────────

func TestCsirtRepository_GetByIDWithPerusahaan(t *testing.T) {
	db, mock, repo := setupTest(t)
	defer db.Close()

	joinCols := []string{
		"c.id", "c.nama_csirt", "c.web_csirt", "c.telepon_csirt",
		"c.photo_csirt", "c.file_rfc2350", "c.file_public_key_pgp",
		"c.file_str", "c.tanggal_registrasi", "c.tanggal_kadaluarsa", "c.tanggal_registrasi_ulang",
		"p.id", "p.photo", "p.nama_perusahaan",
		"p.alamat", "p.telepon", "p.email", "p.website",
		"p.created_at", "p.updated_at",
		"ss.id", "ss.nama_sub_sektor", "ss.id_sektor", "ss.created_at", "ss.updated_at",
		"s.nama_sektor",
	}

	t.Run("success — dengan field baru", func(t *testing.T) {
		id := "csirt-123"
		now := time.Now()
		row := sqlmock.NewRows(joinCols).
			AddRow(
				id, "CSIRT Test", "https://csirt.test.com", "021-12345",
				"photo.jpg", "rfc.pdf", "pgp.key",
				"uploads/str_csirt/str.pdf", "2024-03-01", "2025-03-01", "2025-03-05",
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
		assert.Equal(t, "uploads/str_csirt/str.pdf", result.FileStr)
		assert.Equal(t, "2024-03-01", result.TanggalRegistrasi)
		assert.Equal(t, "2025-03-01", result.TanggalKadaluarsa)
		assert.Equal(t, "2025-03-05", result.TanggalRegistrasiUlang)
		assert.NotNil(t, result.Perusahaan.SubSektor)
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
	})
}

// ─────────────────────────────────────────────────────────
// UPDATE
// ─────────────────────────────────────────────────────────

func TestCsirtRepository_Update(t *testing.T) {
	db, mock, repo := setupTest(t)
	defer db.Close()

	t.Run("success — semua field termasuk yang baru", func(t *testing.T) {
		id := "csirt-123"
		telepon := "021-98763"
		photo := "updated.jpg"
		rfc := "updated_rfc.pdf"
		pgp := "updated_pgp.key"
		fileStr := "uploads/str_csirt/new_str.pdf"
		tglReg := "2024-06-01"
		tglKad := "2025-06-01"
		tglRegU := "2025-06-10"

		csirt := models.Csirt{
			NamaCsirt:              "CSIRT Updated",
			WebCsirt:               "https://csirt.updated.com",
			TeleponCsirt:           &telepon,
			PhotoCsirt:             &photo,
			FileRFC2350:            &rfc,
			FilePublicKeyPGP:       &pgp,
			FileStr:                &fileStr,
			TanggalRegistrasi:      &tglReg,
			TanggalKadaluarsa:      &tglKad,
			TanggalRegistrasiUlang: &tglRegU,
		}

		mock.ExpectExec("UPDATE csirt SET").
			WithArgs(
				csirt.NamaCsirt,
				csirt.WebCsirt,
				csirt.TeleponCsirt,
				csirt.PhotoCsirt,
				csirt.FileRFC2350,
				csirt.FilePublicKeyPGP,
				csirt.FileStr,
				csirt.TanggalRegistrasi,
				csirt.TanggalKadaluarsa,
				csirt.TanggalRegistrasiUlang,
				id,
			).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Update(id, csirt)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success — field nullable nil", func(t *testing.T) {
		id := "csirt-456"
		csirt := models.Csirt{
			NamaCsirt: "CSIRT Null Fields",
			// FileStr, Tanggal* semua nil
		}

		mock.ExpectExec("UPDATE csirt SET").
			WithArgs(
				csirt.NamaCsirt,
				csirt.WebCsirt,
				csirt.TeleponCsirt,
				csirt.PhotoCsirt,
				csirt.FileRFC2350,
				csirt.FilePublicKeyPGP,
				csirt.FileStr,                // nil
				csirt.TanggalRegistrasi,      // nil
				csirt.TanggalKadaluarsa,      // nil
				csirt.TanggalRegistrasiUlang, // nil
				id,
			).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Update(id, csirt)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		id := "csirt-123"
		csirt := models.Csirt{NamaCsirt: "CSIRT Updated"}

		mock.ExpectExec("UPDATE csirt SET").
			WillReturnError(errors.New("update error"))

		err := repo.Update(id, csirt)
		assert.Error(t, err)
		assert.Equal(t, "update error", err.Error())
	})

	t.Run("no rows affected (not found)", func(t *testing.T) {
		id := "csirt-nonexistent"
		csirt := models.Csirt{NamaCsirt: "CSIRT Updated"}

		mock.ExpectExec("UPDATE csirt SET").
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.Update(id, csirt)
		assert.NoError(t, err) // fungsi tidak cek rows affected
	})
}

// ─────────────────────────────────────────────────────────
// DELETE
// ─────────────────────────────────────────────────────────

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
	})

	t.Run("no rows affected", func(t *testing.T) {
		id := "csirt-nonexistent"
		mock.ExpectExec("DELETE FROM csirt WHERE id = ?").
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.Delete(id)
		assert.NoError(t, err) // fungsi tidak cek rows affected
	})
}

// ─────────────────────────────────────────────────────────
// GET BY PERUSAHAAN
// ─────────────────────────────────────────────────────────

func TestCsirtRepository_GetByPerusahaan(t *testing.T) {
	joinCols := []string{
		"c.id", "c.nama_csirt", "c.web_csirt", "c.telepon_csirt",
		"c.photo_csirt", "c.file_rfc2350", "c.file_public_key_pgp",
		"c.file_str", "c.tanggal_registrasi", "c.tanggal_kadaluarsa", "c.tanggal_registrasi_ulang",
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
		checkFn      func(t *testing.T, result []dto.CsirtResponse)
	}{
		{
			name:         "success — field baru terisi dengan sub_sektor",
			idPerusahaan: "perusahaan-1",
			mockFn: func(mock sqlmock.Sqlmock) {
				now := time.Now()
				rows := sqlmock.NewRows(joinCols).
					AddRow(
						"csirt-1", "CSIRT ABC", "https://csirt.abc.com", "021-111",
						"photo.jpg", "rfc.pdf", "pgp.key",
						"uploads/str_csirt/str.pdf", "2024-01-15", "2025-01-15", "2025-01-20",
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
			checkFn: func(t *testing.T, result []dto.CsirtResponse) {
				assert.Equal(t, "uploads/str_csirt/str.pdf", result[0].FileStr)
				assert.Equal(t, "2024-01-15", result[0].TanggalRegistrasi)
				assert.Equal(t, "2025-01-15", result[0].TanggalKadaluarsa)
				assert.Equal(t, "2025-01-20", result[0].TanggalRegistrasiUlang)
				assert.NotNil(t, result[0].Perusahaan.SubSektor)
			},
		},
		{
			name:         "success — field baru NULL tanpa sub_sektor",
			idPerusahaan: "perusahaan-2",
			mockFn: func(mock sqlmock.Sqlmock) {
				now := time.Now()
				rows := sqlmock.NewRows(joinCols).
					AddRow(
						"csirt-2", "CSIRT XYZ", "https://csirt.xyz.com", "021-222",
						"photo2.jpg", "rfc2.pdf", "pgp2.key",
						nil, nil, nil, nil, // field baru NULL
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
			checkFn: func(t *testing.T, result []dto.CsirtResponse) {
				assert.Equal(t, "", result[0].FileStr)
				assert.Equal(t, "", result[0].TanggalRegistrasi)
				assert.Nil(t, result[0].Perusahaan.SubSektor)
			},
		},
		{
			name:         "success — perusahaan tidak punya CSIRT (empty)",
			idPerusahaan: "perusahaan-kosong",
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(joinCols)
				mock.ExpectQuery("SELECT (.+) FROM csirt c JOIN perusahaan p (.+) WHERE c.id_perusahaan = \\?").
					WithArgs("perusahaan-kosong").
					WillReturnRows(rows)
			},
			wantLen: 0,
		},
		{
			name:         "error — database error",
			idPerusahaan: "perusahaan-1",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM csirt c JOIN perusahaan p (.+) WHERE c.id_perusahaan = \\?").
					WithArgs("perusahaan-1").
					WillReturnError(sql.ErrConnDone)
			},
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
				if tt.checkFn != nil {
					tt.checkFn(t, result)
				}
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

// ─────────────────────────────────────────────────────────
// EXISTS BY PERUSAHAAN
// ─────────────────────────────────────────────────────────

func TestCsirtRepository_ExistsByPerusahaan(t *testing.T) {
	tests := []struct {
		name         string
		idPerusahaan string
		mockFn       func(mock sqlmock.Sqlmock)
		wantExists   bool
		wantErr      bool
	}{
		{
			name:         "returns true when csirt exists for perusahaan",
			idPerusahaan: "perusahaan-1",
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(2)
				mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM csirt WHERE id_perusahaan = \\?").
					WithArgs("perusahaan-1").
					WillReturnRows(rows)
			},
			wantExists: true,
			wantErr:    false,
		},
		{
			name:         "returns false when no csirt for perusahaan",
			idPerusahaan: "perusahaan-kosong",
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0)
				mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM csirt WHERE id_perusahaan = \\?").
					WithArgs("perusahaan-kosong").
					WillReturnRows(rows)
			},
			wantExists: false,
			wantErr:    false,
		},
		{
			name:         "returns false and error on db failure",
			idPerusahaan: "perusahaan-1",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM csirt WHERE id_perusahaan = \\?").
					WithArgs("perusahaan-1").
					WillReturnError(sql.ErrConnDone)
			},
			wantExists: false,
			wantErr:    true,
		},
		{
			name:         "returns true when exactly one csirt exists",
			idPerusahaan: "perusahaan-2",
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1)
				mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM csirt WHERE id_perusahaan = \\?").
					WithArgs("perusahaan-2").
					WillReturnRows(rows)
			},
			wantExists: true,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, repo := setupTest(t)
			defer db.Close()

			tt.mockFn(mock)

			exists, err := repo.ExistsByPerusahaan(tt.idPerusahaan)

			if tt.wantErr {
				assert.Error(t, err)
				assert.False(t, exists)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantExists, exists)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
