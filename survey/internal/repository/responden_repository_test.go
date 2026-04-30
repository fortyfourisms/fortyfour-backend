package repository

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"survey/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
)

//
// =====================
// CREATE
// =====================
func TestCreate_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewRespondenRepository(db)

	mock.ExpectExec("INSERT INTO responden").
		WithArgs("user1", "08123", "IT", nil, "yes").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Create(models.Responden{
		UserID:              "user1",
		NoTelepon:           "08123",
		Sektor:              "IT",
		SektorLainnya:       nil,
		SertifikatTraining:  "yes",
	})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

//
// =====================
// GET ALL
// =====================
func TestGetAllDetail_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewRespondenRepository(db)

	rows := sqlmock.NewRows([]string{
		"id", "user_id", "display_name", "email",
		"jabatan", "perusahaan", "perusahaan_id",
		"no_telepon", "sektor", "sektor_lainnya",
		"sertifikat_training", "created_at", "updated_at",
	}).AddRow(
		1, "user1", "Nama", "email@mail.com",
		"Manager", "PT A", "10",
		"08123", "IT", "lainnya",
		"yes", time.Now(), time.Now(),
	)

	mock.ExpectQuery(regexp.QuoteMeta(baseDetailQuery)).
		WillReturnRows(rows)

	result, err := repo.GetAllDetail()

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(result) != 1 {
		t.Errorf("expected 1 result, got %d", len(result))
	}
}

//
// =====================
// GET BY ID
// =====================
func TestGetDetailByID_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewRespondenRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(baseDetailQuery + " WHERE r.id = ?")).
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)

	_, err := repo.GetDetailByID(1)

	if err == nil || err.Error() != "data tidak ditemukan" {
		t.Errorf("expected not found error")
	}
}

//
// =====================
// UPDATE SUCCESS
// =====================
func TestUpdate_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewRespondenRepository(db)

	mock.ExpectExec("UPDATE responden").
		WithArgs("08123", "IT", nil, "yes", 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Update(1, models.Responden{
		NoTelepon:          "08123",
		Sektor:             "IT",
		SektorLainnya:      nil,
		SertifikatTraining: "yes",
	})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

//
// =====================
// UPDATE NOT FOUND
// =====================
func TestUpdate_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewRespondenRepository(db)

	mock.ExpectExec("UPDATE responden").
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.Update(1, models.Responden{})

	if err == nil || err.Error() != "data tidak ditemukan" {
		t.Errorf("expected not found error")
	}
}

//
// =====================
// EXISTS
// =====================
func TestExists_True(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewRespondenRepository(db)

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	exists, err := repo.Exists(1)

	if err != nil {
		t.Errorf("unexpected error")
	}

	if !exists {
		t.Errorf("expected true")
	}
}

//
// =====================
// EXISTS ERROR
// =====================
func TestExists_Error(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewRespondenRepository(db)

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(1).
		WillReturnError(errors.New("db error"))

	_, err := repo.Exists(1)

	if err == nil {
		t.Errorf("expected error")
	}
}