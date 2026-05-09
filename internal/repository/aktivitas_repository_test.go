package repository_test

import (
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/repository"

	"github.com/DATA-DOG/go-sqlmock"
)

// ═══════════════════════════════════════════════════════════════════════════
// TEST: Create
// ═══════════════════════════════════════════════════════════════════════════

func TestAktivitasRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock error: %v", err)
	}
	defer db.Close()

	repo := repository.NewAktivitasRepository(db)

	req := dto.CreateAktivitasRequest{
		PerusahaanID:   "uuid-1",
		Judul:          "Rapat",
		Deskripsi:      "Deskripsi rapat",
		TanggalMulai:   "2024-01-01",
		TanggalSelesai: "2024-01-02",
		JenisAktivitas: []string{"dinas", "rapat koordinasi"},
	}

	jenisJSON, _ := json.Marshal(req.JenisAktivitas)

	mock.ExpectExec("INSERT INTO aktivitas").
		WithArgs(req.PerusahaanID, req.Judul, req.Deskripsi, req.TanggalMulai, req.TanggalSelesai, string(jenisJSON)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	id, err := repo.Create(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if id != 1 {
		t.Errorf("expected ID 1, got %d", id)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestAktivitasRepository_Create_ExecError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock error: %v", err)
	}
	defer db.Close()

	repo := repository.NewAktivitasRepository(db)

	req := dto.CreateAktivitasRequest{
		PerusahaanID:   "uuid-1",
		Judul:          "Rapat",
		TanggalMulai:   "2024-01-01",
		TanggalSelesai: "2024-01-02",
		JenisAktivitas: []string{"dinas"},
	}

	mock.ExpectExec("INSERT INTO aktivitas").
		WillReturnError(sql.ErrConnDone)

	_, err = repo.Create(req)
	if err == nil {
		t.Error("expected error")
	}
}

// ═══════════════════════════════════════════════════════════════════════════
// TEST: GetAll
// ═══════════════════════════════════════════════════════════════════════════

func TestAktivitasRepository_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock error: %v", err)
	}
	defer db.Close()

	repo := repository.NewAktivitasRepository(db)
	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "perusahaan_id", "judul", "deskripsi", "tanggal_mulai", "tanggal_selesai", "jenis_aktivitas", "created_at", "updated_at"}).
		AddRow(1, "uuid-1", "Rapat", "Desc", "2024-01-01", "2024-01-02", `["dinas"]`, now, now).
		AddRow(2, "uuid-2", "Workshop", "Desc 2", "2024-02-01", "2024-02-02", `["workshop","seminar"]`, now, now)

	mock.ExpectQuery("SELECT id, perusahaan_id, judul, deskripsi, tanggal_mulai, tanggal_selesai, jenis_aktivitas, created_at, updated_at FROM aktivitas").
		WillReturnRows(rows)

	result, err := repo.GetAll()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 results, got %d", len(result))
	}
	if result[0].Judul != "Rapat" {
		t.Errorf("expected 'Rapat', got '%s'", result[0].Judul)
	}
	if len(result[1].JenisAktivitas) != 2 {
		t.Errorf("expected 2 jenis for second row, got %d", len(result[1].JenisAktivitas))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestAktivitasRepository_GetAll_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock error: %v", err)
	}
	defer db.Close()

	repo := repository.NewAktivitasRepository(db)

	mock.ExpectQuery("SELECT id").WillReturnError(sql.ErrConnDone)

	_, err = repo.GetAll()
	if err == nil {
		t.Error("expected error")
	}
}

// ═══════════════════════════════════════════════════════════════════════════
// TEST: GetByID
// ═══════════════════════════════════════════════════════════════════════════

func TestAktivitasRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock error: %v", err)
	}
	defer db.Close()

	repo := repository.NewAktivitasRepository(db)
	now := time.Now()

	row := sqlmock.NewRows([]string{"id", "perusahaan_id", "judul", "deskripsi", "tanggal_mulai", "tanggal_selesai", "jenis_aktivitas", "created_at", "updated_at"}).
		AddRow(1, "uuid-1", "Rapat", "Desc", "2024-01-01", "2024-01-02", `["dinas"]`, now, now)

	mock.ExpectQuery("SELECT id, perusahaan_id, judul, deskripsi, tanggal_mulai, tanggal_selesai, jenis_aktivitas, created_at, updated_at FROM aktivitas WHERE id").
		WithArgs(1).
		WillReturnRows(row)

	result, err := repo.GetByID(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.ID != 1 {
		t.Errorf("expected ID 1, got %d", result.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestAktivitasRepository_GetByID_NoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock error: %v", err)
	}
	defer db.Close()

	repo := repository.NewAktivitasRepository(db)

	mock.ExpectQuery("SELECT id").WithArgs(999).WillReturnError(sql.ErrNoRows)

	_, err = repo.GetByID(999)
	if err == nil {
		t.Error("expected error")
	}
}

// ═══════════════════════════════════════════════════════════════════════════
// TEST: GetByPerusahaanID
// ═══════════════════════════════════════════════════════════════════════════

func TestAktivitasRepository_GetByPerusahaanID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock error: %v", err)
	}
	defer db.Close()

	repo := repository.NewAktivitasRepository(db)
	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "perusahaan_id", "judul", "deskripsi", "tanggal_mulai", "tanggal_selesai", "jenis_aktivitas", "created_at", "updated_at"}).
		AddRow(1, "uuid-1", "Rapat", "Desc", "2024-01-01", "2024-01-02", `["dinas"]`, now, now)

	mock.ExpectQuery("SELECT id, perusahaan_id, judul, deskripsi").
		WithArgs("uuid-1").
		WillReturnRows(rows)

	result, err := repo.GetByPerusahaanID("uuid-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 result, got %d", len(result))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

// ═══════════════════════════════════════════════════════════════════════════
// TEST: Update
// ═══════════════════════════════════════════════════════════════════════════

func TestAktivitasRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock error: %v", err)
	}
	defer db.Close()

	repo := repository.NewAktivitasRepository(db)

	judul := "Updated"
	req := dto.UpdateAktivitasRequest{Judul: &judul}

	mock.ExpectExec("UPDATE aktivitas").
		WithArgs("Updated", 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Update(1, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestAktivitasRepository_Update_NoFields(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock error: %v", err)
	}
	defer db.Close()

	repo := repository.NewAktivitasRepository(db)

	// Empty update — should be a no-op
	req := dto.UpdateAktivitasRequest{}
	err = repo.Update(1, req)
	if err != nil {
		t.Fatalf("expected no error for empty update, got %v", err)
	}
}

// ═══════════════════════════════════════════════════════════════════════════
// TEST: Delete
// ═══════════════════════════════════════════════════════════════════════════

func TestAktivitasRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock error: %v", err)
	}
	defer db.Close()

	repo := repository.NewAktivitasRepository(db)

	mock.ExpectExec("DELETE FROM aktivitas").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Delete(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestAktivitasRepository_Delete_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock error: %v", err)
	}
	defer db.Close()

	repo := repository.NewAktivitasRepository(db)

	mock.ExpectExec("DELETE FROM aktivitas").
		WithArgs(1).
		WillReturnError(sql.ErrConnDone)

	err = repo.Delete(1)
	if err == nil {
		t.Error("expected error")
	}
}
