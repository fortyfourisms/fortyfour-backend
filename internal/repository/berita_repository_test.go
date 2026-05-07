package repository_test

import (
	"database/sql"
	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/repository"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestBeritaRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewBeritaRepository(db)

	berita := &models.Berita{
		Judul:     "Test",
		Deskripsi: "Desc",
		AuthorID:  "author1",
	}

	mock.ExpectExec("INSERT INTO berita").
		WithArgs("Test", "Desc", "", "author1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(berita)

	if err != nil {
		t.Errorf("error was not expected while inserting: %s", err)
	}
	if berita.ID != 1 {
		t.Errorf("expected ID 1, got %d", berita.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestBeritaRepository_Create_ExecError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewBeritaRepository(db)

	berita := &models.Berita{
		Judul:     "Test",
		Deskripsi: "Desc",
		AuthorID:  "author1",
	}

	mock.ExpectExec("INSERT INTO berita").
		WithArgs("Test", "Desc", "", "author1").
		WillReturnError(sql.ErrConnDone)

	err = repo.Create(berita)

	if err == nil {
		t.Error("error was expected")
	}
}

func TestBeritaRepository_Create_LastInsertIdError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewBeritaRepository(db)

	berita := &models.Berita{
		Judul:     "Test",
		Deskripsi: "Desc",
		AuthorID:  "author1",
	}

	mock.ExpectExec("INSERT INTO berita").
		WithArgs("Test", "Desc", "", "author1").
		WillReturnResult(sqlmock.NewErrorResult(sql.ErrConnDone))

	err = repo.Create(berita)

	if err == nil {
		t.Error("error was expected")
	}
}

func TestBeritaRepository_FindAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewBeritaRepository(db)
	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "judul", "deskripsi", "tags", "author_id", "created_at", "updated_at", "username", "display_name"}).
		AddRow(1, "Test", "Desc", "[]", "author1", now, now, "user1", "User One").
		AddRow(2, "Test 2", "Desc 2", "[\"tag1\"]", "author2", now, now, "user2", nil)

	mock.ExpectQuery("SELECT b.id, b.judul, b.deskripsi, b.tags, b.author_id, b.created_at, b.updated_at").
		WillReturnRows(rows)

	res, err := repo.FindAll()

	if err != nil {
		t.Errorf("error was not expected: %s", err)
	}
	if len(res) != 2 {
		t.Errorf("expected 2 records, got %d", len(res))
	}
	if res[0].Author.DisplayName == nil || *res[0].Author.DisplayName != "User One" {
		t.Errorf("expected display name User One")
	}
	if res[1].Author.DisplayName != nil {
		t.Errorf("expected nil display name")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestBeritaRepository_FindAll_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewBeritaRepository(db)

	mock.ExpectQuery("SELECT b.id").WillReturnError(sql.ErrConnDone)

	_, err = repo.FindAll()

	if err == nil {
		t.Error("error was expected")
	}
}

func TestBeritaRepository_FindAll_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewBeritaRepository(db)

	// Missing columns
	rows := sqlmock.NewRows([]string{"id", "judul"}).
		AddRow(1, "Test")

	mock.ExpectQuery("SELECT b.id").WillReturnRows(rows)

	_, err = repo.FindAll()

	if err == nil {
		t.Error("error was expected")
	}
}

func TestBeritaRepository_FindByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewBeritaRepository(db)
	now := time.Now()

	row := sqlmock.NewRows([]string{"id", "judul", "deskripsi", "tags", "author_id", "created_at", "updated_at", "username", "display_name"}).
		AddRow(1, "Test", "Desc", "[]", "author1", now, now, "user1", "User One")

	mock.ExpectQuery("SELECT b.id").WithArgs(1).WillReturnRows(row)

	res, err := repo.FindByID(1)

	if err != nil {
		t.Errorf("error was not expected: %s", err)
	}
	if res.ID != 1 {
		t.Errorf("expected ID 1, got %d", res.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestBeritaRepository_FindByID_NoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewBeritaRepository(db)

	mock.ExpectQuery("SELECT b.id").WithArgs(1).WillReturnError(sql.ErrNoRows)

	res, err := repo.FindByID(1)

	if err != nil {
		t.Errorf("error was not expected: %s", err)
	}
	if res != nil {
		t.Errorf("expected nil result, got %v", res)
	}
}

func TestBeritaRepository_FindByID_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewBeritaRepository(db)

	mock.ExpectQuery("SELECT b.id").WithArgs(1).WillReturnError(sql.ErrConnDone)

	_, err = repo.FindByID(1)

	if err == nil {
		t.Error("error was expected")
	}
}

func TestBeritaRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewBeritaRepository(db)

	berita := &models.Berita{
		ID:        1,
		Judul:     "Test Updated",
		Deskripsi: "Desc Updated",
	}

	mock.ExpectExec("UPDATE berita").
		WithArgs("Test Updated", "Desc Updated", "", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Update(berita)

	if err != nil {
		t.Errorf("error was not expected while updating: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestBeritaRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewBeritaRepository(db)

	mock.ExpectExec("DELETE FROM berita").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Delete(1)

	if err != nil {
		t.Errorf("error was not expected while deleting: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
