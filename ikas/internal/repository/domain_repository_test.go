package repository

import (
	"database/sql"
	"errors"
	"ikas/internal/dto"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestDomainRepository(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	defer db.Close()

	repo := NewDomainRepository(db)

	t.Run("Create_Success", func(t *testing.T) {
		req := dto.CreateDomainRequest{NamaDomain: "Test Domain"}
		mock.ExpectExec("INSERT INTO domain").
			WithArgs(req.NamaDomain).
			WillReturnResult(sqlmock.NewResult(1, 1))

		id, err := repo.Create(req)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), id)
	})

	t.Run("Create_Error", func(t *testing.T) {
		req := dto.CreateDomainRequest{NamaDomain: "Test Domain"}
		mock.ExpectExec("INSERT INTO domain").
			WithArgs(req.NamaDomain).
			WillReturnError(errors.New("db error"))

		id, err := repo.Create(req)
		assert.Error(t, err)
		assert.Equal(t, int64(0), id)
	})

	t.Run("GetAll_Success", func(t *testing.T) {
		now := time.Now()
		rows := sqlmock.NewRows([]string{"id", "nama_domain", "created_at", "updated_at"}).
			AddRow(1, "Domain A", now, now).
			AddRow(2, "Domain B", now, now)

		mock.ExpectQuery("SELECT id, nama_domain, created_at, updated_at FROM domain").
			WillReturnRows(rows)

		result, err := repo.GetAll()
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "Domain A", result[0].NamaDomain)
	})

	t.Run("GetAll_QueryError", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, nama_domain, created_at, updated_at FROM domain").
			WillReturnError(errors.New("query error"))

		result, err := repo.GetAll()
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("GetAll_ScanError", func(t *testing.T) {
		// Mock a row that will fail to scan (wrong type for ID)
		rows := sqlmock.NewRows([]string{"id", "nama_domain", "created_at", "updated_at"}).
			AddRow("not-an-int", "Domain A", time.Now(), time.Now())

		mock.ExpectQuery("SELECT id, nama_domain, created_at, updated_at FROM domain").
			WillReturnRows(rows)

		result, err := repo.GetAll()
		assert.NoError(t, err)
		assert.Len(t, result, 0) // It should continue on scan error
	})

	t.Run("GetByID_Success", func(t *testing.T) {
		now := time.Now()
		rows := sqlmock.NewRows([]string{"id", "nama_domain", "created_at", "updated_at"}).
			AddRow(1, "Domain A", now, now)

		mock.ExpectQuery("SELECT id, nama_domain, created_at, updated_at FROM domain WHERE id=?").
			WithArgs(1).
			WillReturnRows(rows)

		result, err := repo.GetByID(1)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Domain A", result.NamaDomain)
	})

	t.Run("GetByID_Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, nama_domain, created_at, updated_at FROM domain WHERE id=?").
			WithArgs(1).
			WillReturnError(sql.ErrNoRows)

		result, err := repo.GetByID(1)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("Update_Success", func(t *testing.T) {
		name := "Updated Domain"
		req := dto.UpdateDomainRequest{NamaDomain: &name}
		mock.ExpectExec("UPDATE domain SET nama_domain=\\? WHERE id=\\?").
			WithArgs(name, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Update(1, req)
		assert.NoError(t, err)
	})

	t.Run("Update_NoChanges", func(t *testing.T) {
		req := dto.UpdateDomainRequest{NamaDomain: nil}
		err := repo.Update(1, req)
		assert.NoError(t, err)
		// DB shouldn't be touched
	})

	t.Run("Update_Error", func(t *testing.T) {
		name := "Updated Domain"
		req := dto.UpdateDomainRequest{NamaDomain: &name}
		mock.ExpectExec("UPDATE domain SET nama_domain=\\? WHERE id=\\?").
			WithArgs(name, 1).
			WillReturnError(errors.New("update error"))

		err := repo.Update(1, req)
		assert.Error(t, err)
	})

	t.Run("Delete_Success", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM domain WHERE id=\\?").
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Delete(1)
		assert.NoError(t, err)
	})

	t.Run("Delete_Error", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM domain WHERE id=\\?").
			WithArgs(1).
			WillReturnError(errors.New("delete error"))

		err := repo.Delete(1)
		assert.Error(t, err)
	})

	t.Run("CheckDuplicateName_Found", func(t *testing.T) {
		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM domain").
			WithArgs("test-domain").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		found, err := repo.CheckDuplicateName("test-domain", 0)
		assert.NoError(t, err)
		assert.True(t, found)
	})

	t.Run("CheckDuplicateName_NotFound", func(t *testing.T) {
		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM domain").
			WithArgs("test-domain", 1).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		found, err := repo.CheckDuplicateName("test-domain", 1)
		assert.NoError(t, err)
		assert.False(t, found)
	})

	t.Run("CheckDuplicateName_Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM domain").
			WithArgs("test-domain").
			WillReturnError(errors.New("count error"))

		found, err := repo.CheckDuplicateName("test-domain", 0)
		assert.Error(t, err)
		assert.False(t, found)
	})
}
