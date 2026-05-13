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

// helper
func int64Ptr(v int64) *int64 { return &v }

// GET ALL RISIKO
func TestGetAllRisiko_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewRisikoRepository(db)

	rows := sqlmock.NewRows([]string{"id", "nama", "deskripsi"}).
		AddRow(1, "Risiko A", "Deskripsi A")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, nama, COALESCE(deskripsi, '')")).
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
func TestGetRisikoByID_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewRisikoRepository(db)

	rows := sqlmock.NewRows([]string{"id", "kode", "nama", "deskripsi", "urutan", "aktif", "created_at", "updated_at"}).
		AddRow(1, "R01", "Risiko A", nil, 1, true, time.Now(), time.Now())

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, kode, nama, deskripsi")).
		WithArgs(int64(1)).
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
		WithArgs(int64(1), int64Ptr(2), true).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.UpsertEligibility(models.RisikoEligibility{
		RespondenID:   1,
		RisikoID:      int64Ptr(2),
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
		"responden_id", "risiko_id",
		"pernah_terjadi", "alasan",
		"dampak_reputasi", "dampak_operasional",
		"dampak_finansial", "dampak_hukum",
		"frekuensi",
		"ada_pengendalian", "deskripsi_pengendalian",
	}).AddRow(
		int64(1), int64(2),
		true, "alasan",
		"tinggi", "sedang", "rendah", "tinggi",
		"sering",
		true, "kontrol",
	).AddRow(
		int64(1), int64(3),
		false, "",
		nil, nil, nil, nil,
		nil,
		false, nil,
	)

	mock.ExpectQuery("SELECT").
		WithArgs(int64(1)).
		WillReturnRows(rows)

	result, err := repo.FindByRespondentID(1)

	if err != nil {
		t.Fatal(err)
	}

	if result["pernah_terjadi"] != true {
		t.Error("invalid mapping")
	}
	if result["responden_id"] != int64(1) {
		t.Error("invalid responden_id mapping")
	}
	if result["risiko_id"] != int64(2) {
		t.Error("invalid risiko_id mapping")
	}
	items, ok := result["items"].([]map[string]interface{})
	if !ok {
		t.Fatal("expected items slice in result")
	}
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}
}

// FIND NOT FOUND
func TestFindByRespondentID_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewRisikoRepository(db)

	mock.ExpectQuery("SELECT").
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{
			"responden_id", "risiko_id",
			"pernah_terjadi", "alasan",
			"dampak_reputasi", "dampak_operasional",
			"dampak_finansial", "dampak_hukum",
			"frekuensi",
			"ada_pengendalian", "deskripsi_pengendalian",
		}))

	_, err := repo.FindByRespondentID(1)

	if !errors.Is(err, ErrNotFound) {
		t.Error("expected ErrNotFound")
	}
}

// PROGRESS INSERT FLOW
func TestGetProgress_InsertDefault(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewRisikoRepository(db)

	// first query -> no rows
	mock.ExpectQuery("SELECT id, responden_id").
		WithArgs(int64(1)).
		WillReturnError(sql.ErrNoRows)

	// insert default
	mock.ExpectExec("INSERT INTO survey_progress").
		WithArgs(int64(1)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// second query (after insert)
	rows := sqlmock.NewRows([]string{
		"id", "responden_id", "risiko_id",
		"langkah_saat_ini", "selesai", "status",
		"edit_request_reason", "edit_request_response",
		"submitted_at", "edit_requested_at",
		"edit_approved_at", "edit_approved_by",
		"edit_rejected_at", "edit_rejected_by",
		"terakhir_update",
	}).AddRow(1, 1, nil, "eligibility", false, "draft", nil, nil, nil, nil, nil, nil, nil, nil, time.Now())

	mock.ExpectQuery("SELECT id, responden_id").
		WithArgs(int64(1)).
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
		WithArgs(int64(1)).
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
		WithArgs(int64(1), "custom").
		WillReturnResult(sqlmock.NewResult(10, 1))

	id, err := repo.InsertCustomRisiko(1, "custom")

	if err != nil {
		t.Fatal(err)
	}

	if id != 10 {
		t.Error("invalid id")
	}
}
