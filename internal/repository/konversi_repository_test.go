package repository_test

import (
	"database/sql"
	"testing"

	"fortyfour-backend/internal/repository"

	"github.com/DATA-DOG/go-sqlmock"
)

// ═══════════════════════════════════════════════════════════════════════════
// TEST: GetAllKonversi — tanpa filter
// ═══════════════════════════════════════════════════════════════════════════

func TestKonversiRepository_GetAllKonversi_All(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock error: %v", err)
	}
	defer db.Close()

	repo := repository.NewKonversiRepository(db)

	rows := sqlmock.NewRows([]string{"perusahaan_id", "nama_perusahaan", "has_ikas", "has_kse", "has_survey", "has_csirt"}).
		AddRow("uuid-1", "PT A", true, true, false, true).
		AddRow("uuid-2", "PT B", false, false, true, false)

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	results, err := repo.GetAllKonversi("")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}

	// PT A: ikas=1, kse=1, survey=0, csirt=1
	if results[0].PoinIkas != 1 {
		t.Errorf("PT A: expected PoinIkas 1, got %d", results[0].PoinIkas)
	}
	if results[0].PoinKse != 1 {
		t.Errorf("PT A: expected PoinKse 1, got %d", results[0].PoinKse)
	}
	if results[0].PoinSurvey != 0 {
		t.Errorf("PT A: expected PoinSurvey 0, got %d", results[0].PoinSurvey)
	}
	if results[0].PoinCsirt != 1 {
		t.Errorf("PT A: expected PoinCsirt 1, got %d", results[0].PoinCsirt)
	}

	// PT B: ikas=0, kse=0, survey=1, csirt=0
	if results[1].PoinSurvey != 1 {
		t.Errorf("PT B: expected PoinSurvey 1, got %d", results[1].PoinSurvey)
	}
	if results[1].PoinIkas != 0 {
		t.Errorf("PT B: expected PoinIkas 0, got %d", results[1].PoinIkas)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

// ═══════════════════════════════════════════════════════════════════════════
// TEST: GetAllKonversi — dengan filter perusahaan_id
// ═══════════════════════════════════════════════════════════════════════════

func TestKonversiRepository_GetAllKonversi_WithFilter(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock error: %v", err)
	}
	defer db.Close()

	repo := repository.NewKonversiRepository(db)

	rows := sqlmock.NewRows([]string{"perusahaan_id", "nama_perusahaan", "has_ikas", "has_kse", "has_survey", "has_csirt"}).
		AddRow("uuid-1", "PT A", true, true, true, true)

	mock.ExpectQuery("SELECT").
		WithArgs("uuid-1").
		WillReturnRows(rows)

	results, err := repo.GetAllKonversi("uuid-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
	if results[0].PerusahaanID != "uuid-1" {
		t.Errorf("expected perusahaan_id 'uuid-1', got '%s'", results[0].PerusahaanID)
	}
	// All true → all poin = 1
	if results[0].PoinIkas != 1 || results[0].PoinKse != 1 || results[0].PoinSurvey != 1 || results[0].PoinCsirt != 1 {
		t.Error("expected all poin = 1 when all have data")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

// ═══════════════════════════════════════════════════════════════════════════
// TEST: GetAllKonversi — query error
// ═══════════════════════════════════════════════════════════════════════════

func TestKonversiRepository_GetAllKonversi_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock error: %v", err)
	}
	defer db.Close()

	repo := repository.NewKonversiRepository(db)

	mock.ExpectQuery("SELECT").WillReturnError(sql.ErrConnDone)

	_, err = repo.GetAllKonversi("")
	if err == nil {
		t.Error("expected error")
	}
}

// ═══════════════════════════════════════════════════════════════════════════
// TEST: GetAllKonversi — empty result
// ═══════════════════════════════════════════════════════════════════════════

func TestKonversiRepository_GetAllKonversi_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock error: %v", err)
	}
	defer db.Close()

	repo := repository.NewKonversiRepository(db)

	rows := sqlmock.NewRows([]string{"perusahaan_id", "nama_perusahaan", "has_ikas", "has_kse", "has_survey", "has_csirt"})

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	results, err := repo.GetAllKonversi("")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if results != nil {
		t.Errorf("expected nil or empty slice, got %d items", len(results))
	}
}

// ═══════════════════════════════════════════════════════════════════════════
// TEST: GetAllKonversi — all false (zero poin)
// ═══════════════════════════════════════════════════════════════════════════

func TestKonversiRepository_GetAllKonversi_AllFalse(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock error: %v", err)
	}
	defer db.Close()

	repo := repository.NewKonversiRepository(db)

	rows := sqlmock.NewRows([]string{"perusahaan_id", "nama_perusahaan", "has_ikas", "has_kse", "has_survey", "has_csirt"}).
		AddRow("uuid-3", "PT C", false, false, false, false)

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	results, err := repo.GetAllKonversi("")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if results[0].PoinIkas != 0 || results[0].PoinKse != 0 || results[0].PoinSurvey != 0 || results[0].PoinCsirt != 0 {
		t.Error("expected all poin = 0 when none have data")
	}
}
