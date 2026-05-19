package repository_test

import (
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
		ID:        "uuid-1",
		Slug:      "test-slug",
		Judul:     "Test",
		Deskripsi: "Desc",
		Lokasi:    "Loc",
		Tanggal:   now,
	}

	mock.ExpectExec("INSERT INTO events").
		WithArgs("uuid-1", "test-slug", "Test", "Desc", "Loc", now).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(event)

	if err != nil {
		t.Errorf("error was not expected while inserting: %s", err)
	}
	if event.ID != "uuid-1" {
		t.Errorf("expected ID uuid-1, got %s", event.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
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

	rows := sqlmock.NewRows([]string{"id", "slug", "judul", "deskripsi", "lokasi", "tanggal", "created_at", "updated_at"}).
		AddRow("uuid-1", "test-1", "Test", "Desc", "Loc", now, now, now).
		AddRow("uuid-2", "test-2", "Test 2", "Desc 2", "Loc 2", now, now, now)

	mock.ExpectQuery("SELECT id, slug, judul, deskripsi, lokasi, tanggal, created_at, updated_at").
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

func TestEventRepository_FindByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewEventRepository(db)
	now := time.Now()

	row := sqlmock.NewRows([]string{"id", "slug", "judul", "deskripsi", "lokasi", "tanggal", "created_at", "updated_at"}).
		AddRow("uuid-1", "test-1", "Test", "Desc", "Loc", now, now, now)

	mock.ExpectQuery("SELECT id, slug, judul, deskripsi, lokasi, tanggal, created_at, updated_at").
		WithArgs("uuid-1").
		WillReturnRows(row)

	res, err := repo.FindByID("uuid-1")

	if err != nil {
		t.Errorf("error was not expected: %s", err)
	}
	if res.ID != "uuid-1" {
		t.Errorf("expected ID uuid-1, got %s", res.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
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
		ID:        "uuid-1",
		Slug:      "test-updated",
		Judul:     "Test Updated",
		Deskripsi: "Desc Updated",
		Lokasi:    "Loc Updated",
		Tanggal:   now,
	}

	mock.ExpectExec("UPDATE events").
		WithArgs("test-updated", "Test Updated", "Desc Updated", "Loc Updated", now, "uuid-1").
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
		WithArgs("uuid-1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Delete("uuid-1")

	if err != nil {
		t.Errorf("error was not expected while deleting: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
