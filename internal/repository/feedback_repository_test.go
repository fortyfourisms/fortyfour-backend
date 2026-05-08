package repository_test

import (
	"database/sql"
	"testing"
	"time"

	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/repository"

	"github.com/DATA-DOG/go-sqlmock"
)

// ═══════════════════════════════════════════════════════════════════════════
// TEST: Upsert
// ═══════════════════════════════════════════════════════════════════════════

func TestFeedbackRepository_Upsert(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock error: %v", err)
	}
	defer db.Close()

	repo := repository.NewFeedbackRepository(db)

	f := &models.Feedback{
		ID:       "fb-1",
		IDMateri: "materi-1",
		IDUser:   "user-1",
		Konten:   "Great content!",
	}

	mock.ExpectExec("INSERT INTO catatan_pribadi").
		WithArgs(f.ID, f.IDMateri, f.IDUser, f.Konten).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Upsert(f)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestFeedbackRepository_Upsert_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock error: %v", err)
	}
	defer db.Close()

	repo := repository.NewFeedbackRepository(db)

	f := &models.Feedback{
		ID:       "fb-1",
		IDMateri: "materi-1",
		IDUser:   "user-1",
		Konten:   "Test",
	}

	mock.ExpectExec("INSERT INTO catatan_pribadi").
		WillReturnError(sql.ErrConnDone)

	err = repo.Upsert(f)
	if err == nil {
		t.Error("expected error")
	}
}

// ═══════════════════════════════════════════════════════════════════════════
// TEST: FindByUserAndMateri
// ═══════════════════════════════════════════════════════════════════════════

func TestFeedbackRepository_FindByUserAndMateri(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock error: %v", err)
	}
	defer db.Close()

	repo := repository.NewFeedbackRepository(db)
	now := time.Now()

	row := sqlmock.NewRows([]string{"id", "id_materi", "id_user", "konten", "created_at", "updated_at"}).
		AddRow("fb-1", "materi-1", "user-1", "Nice content", now, now)

	mock.ExpectQuery("SELECT id, id_materi, id_user, konten, created_at, updated_at FROM catatan_pribadi").
		WithArgs("user-1", "materi-1").
		WillReturnRows(row)

	result, err := repo.FindByUserAndMateri("user-1", "materi-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.ID != "fb-1" {
		t.Errorf("expected ID 'fb-1', got '%s'", result.ID)
	}
	if result.Konten != "Nice content" {
		t.Errorf("expected konten 'Nice content', got '%s'", result.Konten)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestFeedbackRepository_FindByUserAndMateri_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock error: %v", err)
	}
	defer db.Close()

	repo := repository.NewFeedbackRepository(db)

	mock.ExpectQuery("SELECT id, id_materi, id_user, konten, created_at, updated_at FROM catatan_pribadi").
		WithArgs("user-1", "materi-1").
		WillReturnError(sql.ErrNoRows)

	_, err = repo.FindByUserAndMateri("user-1", "materi-1")
	if err == nil {
		t.Error("expected error for not found")
	}
}

// ═══════════════════════════════════════════════════════════════════════════
// TEST: FindByMateri
// ═══════════════════════════════════════════════════════════════════════════

func TestFeedbackRepository_FindByMateri(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock error: %v", err)
	}
	defer db.Close()

	repo := repository.NewFeedbackRepository(db)
	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "id_materi", "id_user", "username", "konten", "created_at", "updated_at"}).
		AddRow("fb-1", "materi-1", "user-1", "john", "Good", now, now).
		AddRow("fb-2", "materi-1", "user-2", "jane", "Great", now, now)

	mock.ExpectQuery("SELECT cp.id, cp.id_materi, cp.id_user, .* FROM catatan_pribadi").
		WithArgs("materi-1").
		WillReturnRows(rows)

	result, err := repo.FindByMateri("materi-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 results, got %d", len(result))
	}
	if result[0].Username != "john" {
		t.Errorf("expected username 'john', got '%s'", result[0].Username)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestFeedbackRepository_FindByMateri_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock error: %v", err)
	}
	defer db.Close()

	repo := repository.NewFeedbackRepository(db)

	mock.ExpectQuery("SELECT cp.id").WillReturnError(sql.ErrConnDone)

	_, err = repo.FindByMateri("materi-1")
	if err == nil {
		t.Error("expected error")
	}
}

func TestFeedbackRepository_FindByMateri_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock error: %v", err)
	}
	defer db.Close()

	repo := repository.NewFeedbackRepository(db)

	rows := sqlmock.NewRows([]string{"id", "id_materi", "id_user", "username", "konten", "created_at", "updated_at"})

	mock.ExpectQuery("SELECT cp.id").
		WithArgs("materi-empty").
		WillReturnRows(rows)

	result, err := repo.FindByMateri("materi-empty")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result != nil && len(result) > 0 {
		t.Errorf("expected nil or empty slice, got %d items", len(result))
	}
}

// ═══════════════════════════════════════════════════════════════════════════
// TEST: Delete
// ═══════════════════════════════════════════════════════════════════════════

func TestFeedbackRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock error: %v", err)
	}
	defer db.Close()

	repo := repository.NewFeedbackRepository(db)

	mock.ExpectExec("DELETE FROM catatan_pribadi").
		WithArgs("fb-1").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Delete("fb-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestFeedbackRepository_Delete_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock error: %v", err)
	}
	defer db.Close()

	repo := repository.NewFeedbackRepository(db)

	mock.ExpectExec("DELETE FROM catatan_pribadi").
		WithArgs("fb-1").
		WillReturnError(sql.ErrConnDone)

	err = repo.Delete("fb-1")
	if err == nil {
		t.Error("expected error")
	}
}
