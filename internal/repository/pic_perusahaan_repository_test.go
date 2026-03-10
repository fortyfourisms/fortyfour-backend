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

func setupPICTest(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *PICRepository) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	repo := NewPICRepository(db)
	return db, mock, repo
}

func TestNewPICRepository(t *testing.T) {
	db, _, _ := setupPICTest(t)
	defer db.Close()

	repo := NewPICRepository(db)
	assert.NotNil(t, repo)
	assert.Equal(t, db, repo.db)
}

func TestPICRepository_Create(t *testing.T) {
	db, mock, repo := setupPICTest(t)
	defer db.Close()

	t.Run("success with all fields", func(t *testing.T) {
		nama := "John Doe"
		telepon := "081234567890"
		idPerusahaan := "perusahaan-123"

		req := dto.CreatePICRequest{
			Nama:         &nama,
			Telepon:      &telepon,
			IDPerusahaan: &idPerusahaan,
		}
		id := "pic-123"

		mock.ExpectExec("INSERT INTO pic_perusahaan").
			WithArgs(id, nama, telepon, idPerusahaan).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Create(req, id)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success with null fields", func(t *testing.T) {
		req := dto.CreatePICRequest{
			Nama:         nil,
			Telepon:      nil,
			IDPerusahaan: nil,
		}
		id := "pic-123"

		mock.ExpectExec("INSERT INTO pic_perusahaan").
			WithArgs(id, nil, nil, nil).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Create(req, id)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		nama := "John Doe"
		req := dto.CreatePICRequest{
			Nama: &nama,
		}
		id := "pic-123"

		mock.ExpectExec("INSERT INTO pic_perusahaan").
			WithArgs(id, nama, nil, nil).
			WillReturnError(errors.New("database error"))

		err := repo.Create(req, id)
		assert.Error(t, err)
		assert.Equal(t, "database error", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("duplicate key error", func(t *testing.T) {
		nama := "John Doe"
		req := dto.CreatePICRequest{
			Nama: &nama,
		}
		id := "pic-123"

		mock.ExpectExec("INSERT INTO pic_perusahaan").
			WithArgs(id, nama, nil, nil).
			WillReturnError(errors.New("duplicate entry"))

		err := repo.Create(req, id)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "duplicate")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPICRepository_GetAll(t *testing.T) {
	db, mock, repo := setupPICTest(t)
	defer db.Close()

	t.Run("success with perusahaan", func(t *testing.T) {
		createdAt := "2024-01-01 10:00:00"
		updatedAt := "2024-01-02 15:30:00"

		rows := sqlmock.NewRows([]string{
			"p.id", "p.nama", "p.telepon", "p.created_at", "p.updated_at",
			"per.id", "per.nama_perusahaan",
		}).
			AddRow("pic-1", "John Doe", "081234567890", createdAt, updatedAt, "perusahaan-1", "PT ABC").
			AddRow("pic-2", "Jane Smith", "081987654321", createdAt, updatedAt, "perusahaan-2", "PT XYZ")

		mock.ExpectQuery("SELECT (.+) FROM pic_perusahaan p LEFT JOIN perusahaan per").
			WillReturnRows(rows)

		result, err := repo.GetAll()
		assert.NoError(t, err)
		require.Len(t, result, 2)

		// Check first record
		assert.Equal(t, "pic-1", result[0].ID)
		assert.Equal(t, "John Doe", result[0].Nama)
		assert.Equal(t, "081234567890", result[0].Telepon)
		assert.Equal(t, createdAt, result[0].CreatedAt)
		assert.Equal(t, updatedAt, result[0].UpdatedAt)
		require.NotNil(t, result[0].Perusahaan)
		assert.Equal(t, "perusahaan-1", result[0].Perusahaan.ID)
		assert.Equal(t, "PT ABC", result[0].Perusahaan.NamaPerusahaan)

		// Check second record
		assert.Equal(t, "pic-2", result[1].ID)
		assert.Equal(t, "Jane Smith", result[1].Nama)
		require.NotNil(t, result[1].Perusahaan)
		assert.Equal(t, "perusahaan-2", result[1].Perusahaan.ID)
		assert.Equal(t, "PT XYZ", result[1].Perusahaan.NamaPerusahaan)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success without perusahaan (null)", func(t *testing.T) {
		createdAt := "2024-01-01 10:00:00"
		updatedAt := "2024-01-02 15:30:00"

		rows := sqlmock.NewRows([]string{
			"p.id", "p.nama", "p.telepon", "p.created_at", "p.updated_at",
			"per.id", "per.nama_perusahaan",
		}).
			AddRow("pic-1", "John Doe", "081234567890", createdAt, updatedAt, nil, nil)

		mock.ExpectQuery("SELECT (.+) FROM pic_perusahaan p LEFT JOIN perusahaan per").
			WillReturnRows(rows)

		result, err := repo.GetAll()
		assert.NoError(t, err)
		require.Len(t, result, 1)

		assert.Equal(t, "pic-1", result[0].ID)
		assert.Equal(t, "John Doe", result[0].Nama)
		assert.Nil(t, result[0].Perusahaan) // Should be nil when no perusahaan

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM pic_perusahaan p LEFT JOIN perusahaan per").
			WillReturnError(errors.New("query error"))

		result, err := repo.GetAll()
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "query error", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("empty result", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"p.id", "p.nama", "p.telepon", "p.created_at", "p.updated_at",
			"per.id", "per.nama_perusahaan",
		})

		mock.ExpectQuery("SELECT (.+) FROM pic_perusahaan p LEFT JOIN perusahaan per").
			WillReturnRows(rows)

		result, err := repo.GetAll()
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error - continues to next row", func(t *testing.T) {
		createdAt := "2024-01-01 10:00:00"
		updatedAt := "2024-01-02 15:30:00"

		rows := sqlmock.NewRows([]string{
			"p.id", "p.nama", "p.telepon", "p.created_at", "p.updated_at",
			"per.id", "per.nama_perusahaan",
		}).
			AddRow("pic-1", nil, nil, nil, nil, nil, nil). // Will cause scan error
			AddRow("pic-2", "Jane Smith", "081987654321", createdAt, updatedAt, "perusahaan-2", "PT XYZ")

		mock.ExpectQuery("SELECT (.+) FROM pic_perusahaan p LEFT JOIN perusahaan per").
			WillReturnRows(rows)

		result, err := repo.GetAll()
		// No error returned because scan error is caught and continues
		assert.NoError(t, err)
		assert.NotNil(t, result)
		// Should only have 1 record (second one) because first one failed scan
		assert.Len(t, result, 1)
		assert.Equal(t, "pic-2", result[0].ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPICRepository_GetByID(t *testing.T) {
	db, mock, repo := setupPICTest(t)
	defer db.Close()

	t.Run("success with perusahaan", func(t *testing.T) {
		id := "pic-123"
		createdAt := "2024-01-01 10:00:00"
		updatedAt := "2024-01-02 15:30:00"

		row := sqlmock.NewRows([]string{
			"p.id", "p.nama", "p.telepon", "p.created_at", "p.updated_at",
			"per.id", "per.nama_perusahaan",
		}).
			AddRow(id, "John Doe", "081234567890", createdAt, updatedAt, "perusahaan-1", "PT ABC")

		mock.ExpectQuery("SELECT (.+) FROM pic_perusahaan p LEFT JOIN perusahaan per (.+) WHERE p.id = ?").
			WithArgs(id).
			WillReturnRows(row)

		result, err := repo.GetByID(id)
		assert.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, id, result.ID)
		assert.Equal(t, "John Doe", result.Nama)
		assert.Equal(t, "081234567890", result.Telepon)
		assert.Equal(t, createdAt, result.CreatedAt)
		assert.Equal(t, updatedAt, result.UpdatedAt)
		require.NotNil(t, result.Perusahaan)
		assert.Equal(t, "perusahaan-1", result.Perusahaan.ID)
		assert.Equal(t, "PT ABC", result.Perusahaan.NamaPerusahaan)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success without perusahaan (null)", func(t *testing.T) {
		id := "pic-123"
		createdAt := "2024-01-01 10:00:00"
		updatedAt := "2024-01-02 15:30:00"

		row := sqlmock.NewRows([]string{
			"p.id", "p.nama", "p.telepon", "p.created_at", "p.updated_at",
			"per.id", "per.nama_perusahaan",
		}).
			AddRow(id, "John Doe", "081234567890", createdAt, updatedAt, nil, nil)

		mock.ExpectQuery("SELECT (.+) FROM pic_perusahaan p LEFT JOIN perusahaan per (.+) WHERE p.id = ?").
			WithArgs(id).
			WillReturnRows(row)

		result, err := repo.GetByID(id)
		assert.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, id, result.ID)
		assert.Equal(t, "John Doe", result.Nama)
		assert.Nil(t, result.Perusahaan) // Should be nil when no perusahaan
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		id := "pic-nonexistent"

		mock.ExpectQuery("SELECT (.+) FROM pic_perusahaan p LEFT JOIN perusahaan per (.+) WHERE p.id = ?").
			WithArgs(id).
			WillReturnError(sql.ErrNoRows)

		result, err := repo.GetByID(id)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, sql.ErrNoRows, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error", func(t *testing.T) {
		id := "pic-123"
		row := sqlmock.NewRows([]string{"id"}).
			AddRow(id)

		mock.ExpectQuery("SELECT (.+) FROM pic_perusahaan p LEFT JOIN perusahaan per (.+) WHERE p.id = ?").
			WithArgs(id).
			WillReturnRows(row)

		result, err := repo.GetByID(id)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		id := "pic-123"

		mock.ExpectQuery("SELECT (.+) FROM pic_perusahaan p LEFT JOIN perusahaan per (.+) WHERE p.id = ?").
			WithArgs(id).
			WillReturnError(errors.New("database connection error"))

		result, err := repo.GetByID(id)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "database connection error", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPICRepository_Update(t *testing.T) {
	db, mock, repo := setupPICTest(t)
	defer db.Close()

	t.Run("success update all fields", func(t *testing.T) {
		id := "pic-123"
		nama := "John Doe Updated"
		telepon := "081999999999"
		idPerusahaan := "perusahaan-456"

		req := dto.UpdatePICRequest{
			Nama:         &nama,
			Telepon:      &telepon,
			IDPerusahaan: &idPerusahaan,
		}

		mock.ExpectExec("UPDATE pic_perusahaan SET nama=\\?, telepon=\\?, id_perusahaan=\\? WHERE id=\\?").
			WithArgs(nama, telepon, idPerusahaan, id).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Update(id, req)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success update only nama", func(t *testing.T) {
		id := "pic-123"
		nama := "John Doe Updated"

		req := dto.UpdatePICRequest{
			Nama:         &nama,
			Telepon:      nil,
			IDPerusahaan: nil,
		}

		mock.ExpectExec("UPDATE pic_perusahaan SET nama=\\? WHERE id=\\?").
			WithArgs(nama, id).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Update(id, req)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success update only telepon", func(t *testing.T) {
		id := "pic-123"
		telepon := "081999999999"

		req := dto.UpdatePICRequest{
			Nama:         nil,
			Telepon:      &telepon,
			IDPerusahaan: nil,
		}

		mock.ExpectExec("UPDATE pic_perusahaan SET telepon=\\? WHERE id=\\?").
			WithArgs(telepon, id).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Update(id, req)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success update nama and telepon", func(t *testing.T) {
		id := "pic-123"
		nama := "John Doe Updated"
		telepon := "081999999999"

		req := dto.UpdatePICRequest{
			Nama:         &nama,
			Telepon:      &telepon,
			IDPerusahaan: nil,
		}

		mock.ExpectExec("UPDATE pic_perusahaan SET nama=\\?, telepon=\\? WHERE id=\\?").
			WithArgs(nama, telepon, id).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Update(id, req)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("no fields to update - returns nil without query", func(t *testing.T) {
		id := "pic-123"

		req := dto.UpdatePICRequest{
			Nama:         nil,
			Telepon:      nil,
			IDPerusahaan: nil,
		}

		// No fields to update: function akan COUNT(*) untuk cek keberadaan data
		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM pic_perusahaan WHERE id=\\?").
			WithArgs(id).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		err := repo.Update(id, req)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		id := "pic-123"
		nama := "John Doe Updated"

		req := dto.UpdatePICRequest{
			Nama: &nama,
		}

		mock.ExpectExec("UPDATE pic_perusahaan SET nama=\\? WHERE id=\\?").
			WithArgs(nama, id).
			WillReturnError(errors.New("update error"))

		err := repo.Update(id, req)
		assert.Error(t, err)
		assert.Equal(t, "update error", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("no rows affected - returns error", func(t *testing.T) {
		id := "pic-nonexistent"
		nama := "John Doe"

		req := dto.UpdatePICRequest{
			Nama: &nama,
		}

		mock.ExpectExec("UPDATE pic_perusahaan SET nama=\\? WHERE id=\\?").
			WithArgs(nama, id).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.Update(id, req)
		// Function sekarang cek rows affected, return error jika 0
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPICRepository_Delete(t *testing.T) {
	db, mock, repo := setupPICTest(t)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		id := "pic-123"

		mock.ExpectExec("DELETE FROM pic_perusahaan WHERE id=\\?").
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Delete(id)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		id := "pic-123"

		mock.ExpectExec("DELETE FROM pic_perusahaan WHERE id=\\?").
			WithArgs(id).
			WillReturnError(errors.New("delete error"))

		err := repo.Delete(id)
		assert.Error(t, err)
		assert.Equal(t, "delete error", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("no rows affected - not checked by function", func(t *testing.T) {
		id := "pic-nonexistent"

		mock.ExpectExec("DELETE FROM pic_perusahaan WHERE id=\\?").
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.Delete(id)
		// Function doesn't check rows affected, so no error
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("multiple rows deleted - not checked by function", func(t *testing.T) {
		id := "pic-123"

		mock.ExpectExec("DELETE FROM pic_perusahaan WHERE id=\\?").
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 5))

		err := repo.Delete(id)
		// Function doesn't check if more than 1 row was deleted
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}