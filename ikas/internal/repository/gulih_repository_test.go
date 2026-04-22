package repository

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGulihRepository(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	defer db.Close()

	repo := NewGulihRepository(db)

	columns := []string{"id", "ikas_id", "nilai_gulih", "nilai_subdomain1", "nilai_subdomain2", "nilai_subdomain3", "nilai_subdomain4"}

	t.Run("GetAll_Success", func(t *testing.T) {
		rows := sqlmock.NewRows(columns).
			AddRow("1", "ikas-1", 4.5, 4.0, 4.5, 5.0, 4.2).
			AddRow("2", "ikas-2", 3.5, 3.0, 3.5, 4.0, 3.8)

		mock.ExpectQuery("SELECT id, ikas_id, nilai_gulih, nilai_subdomain1, nilai_subdomain2, nilai_subdomain3, nilai_subdomain4 FROM gulih").
			WillReturnRows(rows)

		result, err := repo.GetAll()
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "ikas-1", result[0].IkasID)
	})

	t.Run("GetAll_Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, ikas_id, nilai_gulih, nilai_subdomain1, nilai_subdomain2, nilai_subdomain3, nilai_subdomain4 FROM gulih").
			WillReturnError(errors.New("query error"))

		result, err := repo.GetAll()
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("GetByIkasID_Success", func(t *testing.T) {
		rows := sqlmock.NewRows(columns).
			AddRow("1", "ikas-1", 4.5, 4.0, 4.5, 5.0, 4.2)

		mock.ExpectQuery("SELECT id, ikas_id, nilai_gulih, nilai_subdomain1, nilai_subdomain2, nilai_subdomain3, nilai_subdomain4 FROM gulih WHERE ikas_id = ?").
			WithArgs("ikas-1").
			WillReturnRows(rows)

		result, err := repo.GetByIkasID("ikas-1")
		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "1", result[0].ID)
	})

	t.Run("GetByIkasID_Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, ikas_id, nilai_gulih, nilai_subdomain1, nilai_subdomain2, nilai_subdomain3, nilai_subdomain4 FROM gulih WHERE ikas_id = ?").
			WithArgs("ikas-1").
			WillReturnError(errors.New("query error"))

		result, err := repo.GetByIkasID("ikas-1")
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("GetByID_Success", func(t *testing.T) {
		rows := sqlmock.NewRows(columns).
			AddRow("1", "ikas-1", 4.5, 4.0, 4.5, 5.0, 4.2)

		mock.ExpectQuery("SELECT id, ikas_id, nilai_gulih, nilai_subdomain1, nilai_subdomain2, nilai_subdomain3, nilai_subdomain4 FROM gulih WHERE id = ?").
			WithArgs("1").
			WillReturnRows(rows)

		result, err := repo.GetByID("1")
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "1", result.ID)
	})

	t.Run("GetByID_NotFound", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, ikas_id, nilai_gulih, nilai_subdomain1, nilai_subdomain2, nilai_subdomain3, nilai_subdomain4 FROM gulih WHERE id = ?").
			WithArgs("non-existent").
			WillReturnError(sql.ErrNoRows)

		result, err := repo.GetByID("non-existent")
		assert.Error(t, err)
		assert.True(t, errors.Is(err, sql.ErrNoRows))
		assert.Nil(t, result)
	})

	t.Run("GetByPerusahaanID_Success", func(t *testing.T) {
		rows := sqlmock.NewRows(columns).
			AddRow("1", "ikas-1", 4.5, 4.0, 4.5, 5.0, 4.2)

		mock.ExpectQuery("SELECT t.id, t.ikas_id, t.nilai_gulih, t.nilai_subdomain1, t.nilai_subdomain2, t.nilai_subdomain3, t.nilai_subdomain4").
			WithArgs("comp-1").
			WillReturnRows(rows)

		result, err := repo.GetByPerusahaanID("comp-1")
		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "ikas-1", result[0].IkasID)
	})

	t.Run("GetByPerusahaanID_Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT t.id, t.ikas_id, t.nilai_gulih, t.nilai_subdomain1, t.nilai_subdomain2, t.nilai_subdomain3, t.nilai_subdomain4").
			WithArgs("comp-1").
			WillReturnError(errors.New("db error"))

		result, err := repo.GetByPerusahaanID("comp-1")
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("CloneByIkasID_Success", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO gulih").
			WithArgs("target-ikas", "source-ikas").
			WillReturnResult(sqlmock.NewResult(123, 1))

		id, err := repo.CloneByIkasID("source-ikas", "target-ikas")
		assert.NoError(t, err)
		assert.Equal(t, "123", id)
	})

	t.Run("CloneByIkasID_ExecError", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO gulih").
			WithArgs("target-ikas", "source-ikas").
			WillReturnError(errors.New("exec error"))

		id, err := repo.CloneByIkasID("source-ikas", "target-ikas")
		assert.Error(t, err)
		assert.Equal(t, "", id)
	})

	t.Run("CloneByIkasID_LastInsertIdError", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO gulih").
			WithArgs("target-ikas", "source-ikas").
			WillReturnResult(sqlmock.NewErrorResult(errors.New("last insert id error")))

		id, err := repo.CloneByIkasID("source-ikas", "target-ikas")
		assert.Error(t, err)
		assert.Equal(t, "", id)
	})
}
