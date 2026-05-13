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
func stringPtr(v string) *string { return &v }

// GET ALL RISIKO
func TestGetAllRisiko_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewRisikoRepository(db)

	rows := sqlmock.NewRows([]string{"id", "nama", "deskripsi"}).
		AddRow("uuid-1", "Risiko A", "Deskripsi A")

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
		AddRow("uuid-1", "R01", "Risiko A", nil, 1, true, time.Now(), time.Now())

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, kode, nama, deskripsi")).
		WithArgs("uuid-1").
		WillReturnRows(rows)

	res, err := repo.GetByID("uuid-1")

	if err != nil {
		t.Fatal(err)
	}

	if res.ID != "uuid-1" {
		t.Error("invalid ID")
	}
}

// UPSERT ELIGIBILITY
func TestUpsertEligibility(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewRisikoRepository(db)

	mock.ExpectExec("INSERT INTO risiko_eligibility").
		WithArgs(sqlmock.AnyArg(), "uuid-resp", "uuid-risk", true).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.UpsertEligibility(models.RisikoEligibility{
		RespondenID:   "uuid-resp",
		RisikoID:      stringPtr("uuid-risk"),
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
		"uuid-resp", "uuid-risk-1",
		true, "alasan",
		"tinggi", "sedang", "rendah", "tinggi",
		"sering",
		true, "kontrol",
	).AddRow(
		"uuid-resp", "uuid-risk-2",
		false, "",
		nil, nil, nil, nil,
		nil,
		false, nil,
	)

	mock.ExpectQuery("SELECT").
		WithArgs("uuid-resp").
		WillReturnRows(rows)

	result, err := repo.FindByRespondentID("uuid-resp")

	if err != nil {
		t.Fatal(err)
	}

	if result["pernah_terjadi"] != true {
		t.Error("invalid mapping")
	}
	if result["responden_id"] != "uuid-resp" {
		t.Error("invalid responden_id mapping")
	}
	if result["risiko_id"] != "uuid-risk-1" {
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
		WithArgs("uuid-resp").
		WillReturnRows(sqlmock.NewRows([]string{
			"responden_id", "risiko_id",
			"pernah_terjadi", "alasan",
			"dampak_reputasi", "dampak_operasional",
			"dampak_finansial", "dampak_hukum",
			"frekuensi",
			"ada_pengendalian", "deskripsi_pengendalian",
		}))

	_, err := repo.FindByRespondentID("uuid-resp")

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
		WithArgs("uuid-resp").
		WillReturnError(sql.ErrNoRows)

	// insert default
	mock.ExpectExec("INSERT INTO survey_progress").
		WithArgs(sqlmock.AnyArg(), "uuid-resp").
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
	}).AddRow("uuid-prog", "uuid-resp", nil, "eligibility", false, "draft", nil, nil, nil, nil, nil, nil, nil, nil, time.Now())

	mock.ExpectQuery("SELECT id, responden_id").
		WithArgs("uuid-resp").
		WillReturnRows(rows)

	res, err := repo.GetProgress("uuid-resp")

	if err != nil {
		t.Fatal(err)
	}

	if res.RespondenID != "uuid-resp" {
		t.Error("invalid responden id")
	}
}

// EXISTS RISIKO
func TestExistsRisiko(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewRisikoRepository(db)

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs("uuid-risk").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	exists, err := repo.ExistsRisiko("uuid-risk")

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
		WithArgs(sqlmock.AnyArg(), "uuid-resp", "custom").
		WillReturnResult(sqlmock.NewResult(0, 1))

	id, err := repo.InsertCustomRisiko("uuid-resp", "custom")

	if err != nil {
		t.Fatal(err)
	}

	if id == "" {
		t.Error("invalid id")
	}
}
