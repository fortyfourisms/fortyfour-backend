package repository_test

import (
	"database/sql"
	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/repository"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestEventRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewEventRepository(db)

	now := time.Now()
	event := &models.Event{
		Judul:     "Test",
		Deskripsi: "Desc",
		Lokasi:    "Loc",
		Tanggal:   now,
	}

	mock.ExpectExec("INSERT INTO events").
		WithArgs("Test", "Desc", "Loc", now).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(event)

	if err != nil {
		t.Errorf("error was not expected while inserting: %s", err)
	}
	if event.ID != 1 {
		t.Errorf("expected ID 1, got %d", event.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestEventRepository_Create_ExecError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewEventRepository(db)

	now := time.Now()
	event := &models.Event{
		Judul:     "Test",
		Deskripsi: "Desc",
		Lokasi:    "Loc",
		Tanggal:   now,
	}

	mock.ExpectExec("INSERT INTO events").
		WithArgs("Test", "Desc", "Loc", now).
		WillReturnError(sql.ErrConnDone)

	err = repo.Create(event)

	if err == nil {
		t.Error("error was expected")
	}
}

func TestEventRepository_Create_LastInsertIdError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewEventRepository(db)

	now := time.Now()
	event := &models.Event{
		Judul:     "Test",
		Deskripsi: "Desc",
		Lokasi:    "Loc",
		Tanggal:   now,
	}

	mock.ExpectExec("INSERT INTO events").
		WithArgs("Test", "Desc", "Loc", now).
		WillReturnResult(sqlmock.NewErrorResult(sql.ErrConnDone))

	err = repo.Create(event)

	if err == nil {
		t.Error("error was expected")
	}
}

func TestEventRepository_FindAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewEventRepository(db)
	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "judul", "deskripsi", "lokasi", "tanggal", "created_at", "updated_at"}).
		AddRow(1, "Test", "Desc", "Loc", now, now, now).
		AddRow(2, "Test 2", "Desc 2", "Loc 2", now, now, now)

	mock.ExpectQuery("SELECT id, judul, deskripsi, lokasi, tanggal, created_at, updated_at").
		WillReturnRows(rows)

	res, err := repo.FindAll()

	if err != nil {
		t.Errorf("error was not expected: %s", err)
	}
	if len(res) != 2 {
		t.Errorf("expected 2 records, got %d", len(res))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestEventRepository_FindAll_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewEventRepository(db)

	mock.ExpectQuery("SELECT id").WillReturnError(sql.ErrConnDone)

	_, err = repo.FindAll()

	if err == nil {
		t.Error("error was expected")
	}
}

func TestEventRepository_FindAll_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewEventRepository(db)

	// Missing columns
	rows := sqlmock.NewRows([]string{"id", "judul"}).
		AddRow(1, "Test")

	mock.ExpectQuery("SELECT id").WillReturnRows(rows)

	_, err = repo.FindAll()

	if err == nil {
		t.Error("error was expected")
	}
}

func TestEventRepository_FindByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewEventRepository(db)
	now := time.Now()

	row := sqlmock.NewRows([]string{"id", "judul", "deskripsi", "lokasi", "tanggal", "created_at", "updated_at"}).
		AddRow(1, "Test", "Desc", "Loc", now, now, now)

	mock.ExpectQuery("SELECT id, judul, deskripsi, lokasi, tanggal, created_at, updated_at").
		WithArgs(1).
		WillReturnRows(row)

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

func TestEventRepository_FindByID_NoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewEventRepository(db)

	mock.ExpectQuery("SELECT id").WithArgs(1).WillReturnError(sql.ErrNoRows)

	res, err := repo.FindByID(1)

	if err != nil {
		t.Errorf("error was not expected: %s", err)
	}
	if res != nil {
		t.Errorf("expected nil result, got %v", res)
	}
}

func TestEventRepository_FindByID_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewEventRepository(db)

	mock.ExpectQuery("SELECT id").WithArgs(1).WillReturnError(sql.ErrConnDone)

	_, err = repo.FindByID(1)

	if err == nil {
		t.Error("error was expected")
	}
}

func TestEventRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewEventRepository(db)

	now := time.Now()
	event := &models.Event{
		ID:        1,
		Judul:     "Test Updated",
		Deskripsi: "Desc Updated",
		Lokasi:    "Loc Updated",
		Tanggal:   now,
	}

	mock.ExpectExec("UPDATE events").
		WithArgs("Test Updated", "Desc Updated", "Loc Updated", now, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Update(event)

	if err != nil {
		t.Errorf("error was not expected while updating: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestEventRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewEventRepository(db)

	mock.ExpectExec("DELETE FROM events").
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
