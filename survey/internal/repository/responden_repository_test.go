package repository

import (
	"database/sql"
	"regexp"
	"testing"
	"time"

	"survey/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
)

// helper
func strPtr(s string) *string { return &s }

// CREATE
func TestCreate_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewRespondenRepository(db)

	mock.ExpectExec("INSERT INTO responden").
		WithArgs("user1", "perusahaan1", "Nama Lengkap", "Manager", "email@mail.com", "08123", strPtr("yes")).
		WillReturnResult(sqlmock.NewResult(1, 1))

	id, err := repo.Create(models.Responden{
		UserID:             "user1",
		IdPerusahaan:       "perusahaan1",
		NamaLengkap:        "Nama Lengkap",
		Jabatan:            "Manager",
		Email:              "email@mail.com",
		NoTelepon:          "08123",
		SertifikatTraining: strPtr("yes"),
	})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if id != 1 {
		t.Errorf("expected insert ID 1, got %d", id)
	}
}

// GET ALL
func TestGetAllDetail_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewRespondenRepository(db)

	rows := sqlmock.NewRows([]string{
		"id", "user_id", "id_perusahaan", "nama_lengkap", "jabatan", "email",
		"no_telepon", "sertifikat_training",
		"nama_perusahaan", "nama_sub_sektor", "nama_sektor",
		"created_at", "updated_at",
	}).AddRow(
		1, "user1", "perusahaan1", "Nama", "Manager", "email@mail.com",
		"08123", "yes",
		"PT A", "Sub Sektor 1", "Sektor 1",
		time.Now(), time.Now(),
	).AddRow(
		2, "user2", "perusahaan2", "Nama 2", "Manager 2", "email2@mail.com",
		"081234", nil, // test scanning null string
		nil, nil, nil, // test other null fields
		time.Now(), time.Now(),
	)

	mock.ExpectQuery(regexp.QuoteMeta(baseDetailQuery)).
		WillReturnRows(rows)

	result, err := repo.GetAllDetail()

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("expected 2 results, got %d", len(result))
	}
	
	if result[1].SertifikatTraining != nil {
		t.Errorf("expected SertifikatTraining to be nil, got %v", *result[1].SertifikatTraining)
	}
}

// GET BY ID
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

// GET BY USER ID
func TestGetByUserID_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewRespondenRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(baseDetailQuery + " WHERE r.user_id = ?")).
		WithArgs("nonexistent").
		WillReturnError(sql.ErrNoRows)

	_, err := repo.GetByUserID("nonexistent")

	if err == nil || err.Error() != "data tidak ditemukan" {
		t.Errorf("expected not found error")
	}
}
