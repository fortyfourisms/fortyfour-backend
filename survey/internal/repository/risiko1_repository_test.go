package repository

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"survey/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
)

// GET ALL RISIKO
func TestGetAllRisiko_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewRisikoRepository(db)

	rows := sqlmock.NewRows([]string{"id", "nama", "deskripsi"}).
		AddRow(1, "Risiko A", "Deskripsi A")

	mock.ExpectQuery("SELECT id, nama, deskripsi FROM risiko").
		WillReturnRows(rows)

	result, err := repo.GetAllRisiko()

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 1 {
		t.Errorf("expected 1, got %d", len(result))
	}
}

// GET BY ID
func TestGetByID_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewRisikoRepository(db)

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	mock.ExpectQuery("SELECT id, responden_id").
		WithArgs(1).
		WillReturnRows(rows)

	res, err := repo.GetByID(1)

	if err != nil {
		t.Fatal(err)
	}

	if res.ID != 1 {
		t.Error("invalid ID")
	}
}

// UPSERT ELIGIBILITY
func TestUpsertEligibility(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewRisikoRepository(db)

	mock.ExpectExec("INSERT INTO risiko_eligibility").
		WithArgs(1, 2, true).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.UpsertEligibility(models.RisikoEligibility{
		RespondenID:   1,
		RisikoID:      2,
		PernahTerjadi: true,
	})

	if err != nil {
		t.Error(err)
	}
}

// FIND BY RESPONDENT
func TestFindByRespondentID_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewRisikoRepository(db)

	rows := sqlmock.NewRows([]string{
		"pernah_terjadi", "alasan",
		"dampak_reputasi", "dampak_operasional",
		"dampak_finansial", "dampak_hukum",
		"frekuensi",
		"ada_pengendalian", "deskripsi_pengendalian",
	}).AddRow(
		true, "alasan",
		"tinggi", "sedang", "rendah", "tinggi",
		"sering",
		true, "kontrol",
	)

	mock.ExpectQuery("SELECT").
		WithArgs(1).
		WillReturnRows(rows)

	result, err := repo.FindByRespondentID(1)

	if err != nil {
		t.Fatal(err)
	}

	if result["pernah_terjadi"] != true {
		t.Error("invalid mapping")
	}
}

// FIND NOT FOUND
func TestFindByRespondentID_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewRisikoRepository(db)

	mock.ExpectQuery("SELECT").
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)

	_, err := repo.FindByRespondentID(1)

	if !errors.Is(err, ErrNotFound) {
		t.Error("expected ErrNotFound")
	}
}

// UPDATE PARTIAL
func TestUpdatePartial(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewRisikoRepository(db)

	mock.ExpectExec("UPDATE risiko SET").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.UpdatePartial(1, map[string]interface{}{
		"nama": "baru",
	})

	if err != nil {
		t.Error(err)
	}
}

// PROGRESS INSERT FLOW
func TestGetProgress_InsertDefault(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewRisikoRepository(db)

	// first query -> no rows
	mock.ExpectQuery("SELECT id, responden_id").
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)

	// insert default
	mock.ExpectExec("INSERT INTO survey_progress").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// second query (after insert)
	rows := sqlmock.NewRows([]string{
		"id", "responden_id", "risiko_id",
		"langkah_saat_ini", "selesai", "terakhir_update",
	}).AddRow(1, 1, nil, "eligibility", false, time.Now())

	mock.ExpectQuery("SELECT id, responden_id").
		WithArgs(1).
		WillReturnRows(rows)

	res, err := repo.GetProgress(1)

	if err != nil {
		t.Fatal(err)
	}

	if res.RespondenID != 1 {
		t.Error("invalid responden id")
	}
}

// EXISTS RISIKO
func TestExistsRisiko(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewRisikoRepository(db)

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	exists, err := repo.ExistsRisiko(1)

	if err != nil || !exists {
		t.Error("expected true")
	}
}

// INSERT CUSTOM RISIKO
func TestInsertCustomRisiko(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewRisikoRepository(db)

	mock.ExpectExec("INSERT INTO risiko_custom").
		WithArgs(1, "custom").
		WillReturnResult(sqlmock.NewResult(10, 1))

	id, err := repo.InsertCustomRisiko(1, "custom")

	if err != nil {
		t.Fatal(err)
	}

	if id != 10 {
		t.Error("invalid id")
	}
}
