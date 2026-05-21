package services

import (
	"context"
	"database/sql"
	"testing"

	"survey/internal/dto"
	"survey/internal/models"
)

// MOCK CACHE
type mockCache struct{}

func (m *mockCache) Get(ctx context.Context, key string) (string, bool, error) {
	return "", false, nil
}

func (m *mockCache) Set(ctx context.Context, key string, value string, ttlSeconds int) error {
	return nil
}

func (m *mockCache) Del(ctx context.Context, key string) error {
	return nil
}

// MOCK REPOSITORY
type mockRisikoRepo struct {
	getAllRisikoFn            func() ([]models.RisikoResponse, error)
	existsRespondenFn         func(id string) (bool, error)
	existsRisikoFn            func(id string) (bool, error)
	existsCustomFn            func(id string) (bool, error)
	getRisikoIDByUrutanFn     func(urutan int) (string, error)
	getUrutanByRisikoIDFn     func(id string) (int, error)
	upsertEligibilityFn       func(m models.RisikoEligibility) error
	upsertAlasanFn            func(m models.RisikoAlasan) error
	upsertDampakFn            func(m models.RisikoDampak) error
	upsertPengendalianFn      func(m models.RisikoPengendalian) error
	findByRespondentFn        func(id string) (map[string]interface{}, error)
	getProgressFn             func(id string) (*models.SurveyProgress, error)
	upsertProgressFn          func(p models.SurveyProgress) error
	getRespondentIDByUserIDFn func(userID string) (string, error)
	insertCustomFn            func(id string, nama string) (int, error)
}

func (m *mockRisikoRepo) GetAllRisiko() ([]models.RisikoResponse, error) {
	return m.getAllRisikoFn()
}
func (m *mockRisikoRepo) ExistsResponden(id string) (bool, error) {
	return m.existsRespondenFn(id)
}
func (m *mockRisikoRepo) ExistsRisiko(id string) (bool, error) {
	return m.existsRisikoFn(id)
}
func (m *mockRisikoRepo) GetRisikoIDByUrutan(urutan int) (string, error) {
	return m.getRisikoIDByUrutanFn(urutan)
}
func (m *mockRisikoRepo) GetUrutanByRisikoID(id string) (int, error) {
	return m.getUrutanByRisikoIDFn(id)
}
func (m *mockRisikoRepo) ExistsCustomRisiko(id string) (bool, error) {
	return m.existsCustomFn(id)
}
func (m *mockRisikoRepo) UpsertEligibility(r models.RisikoEligibility) error {
	return m.upsertEligibilityFn(r)
}
func (m *mockRisikoRepo) UpsertAlasan(r models.RisikoAlasan) error {
	return m.upsertAlasanFn(r)
}
func (m *mockRisikoRepo) UpsertDampak(r models.RisikoDampak) error {
	return m.upsertDampakFn(r)
}
func (m *mockRisikoRepo) UpsertPengendalian(r models.RisikoPengendalian) error {
	return m.upsertPengendalianFn(r)
}
func (m *mockRisikoRepo) FindByRespondentID(id string) (map[string]interface{}, error) {
	return m.findByRespondentFn(id)
}
func (m *mockRisikoRepo) GetProgress(id string) (*models.SurveyProgress, error) {
	return m.getProgressFn(id)
}
func (m *mockRisikoRepo) UpsertProgress(p models.SurveyProgress) error {
	return m.upsertProgressFn(p)
}
func (m *mockRisikoRepo) GetRespondentIDByUserID(userID string) (string, error) {
	return m.getRespondentIDByUserIDFn(userID)
}
func (m *mockRisikoRepo) InsertCustomRisiko(id string, nama string) (int, error) {
	return m.insertCustomFn(id, nama)
}

// HELPER
func ptrStr(v string) *string {
	return &v
}

func newService(mock *mockRisikoRepo) *RisikoService {
	if mock.getAllRisikoFn == nil {
		mock.getAllRisikoFn = func() ([]models.RisikoResponse, error) {
			return nil, nil
		}
	}
	if mock.getRespondentIDByUserIDFn == nil {
		mock.getRespondentIDByUserIDFn = func(userID string) (string, error) {
			return "1", nil
		}
	}
	if mock.getProgressFn == nil {
		mock.getProgressFn = func(id string) (*models.SurveyProgress, error) {
			return &models.SurveyProgress{
				RespondenID: id,
				Status:      SurveyStatusDraft,
			}, nil
		}
	}
	if mock.getRisikoIDByUrutanFn == nil {
		mock.getRisikoIDByUrutanFn = func(urutan int) (string, error) {
			return "1", nil
		}
	}
	if mock.getUrutanByRisikoIDFn == nil {
		mock.getUrutanByRisikoIDFn = func(id string) (int, error) {
			return 1, nil
		}
	}
	if mock.upsertProgressFn == nil {
		mock.upsertProgressFn = func(p models.SurveyProgress) error {
			return nil
		}
	}
	return NewRisikoService(mock, &mockCache{})
}

// TEST ELIGIBILITY SUCCESS
func TestProcessEligibility_Success(t *testing.T) {
	mock := &mockRisikoRepo{
		existsRespondenFn: func(id string) (bool, error) { return true, nil },
		existsRisikoFn:    func(id string) (bool, error) { return true, nil },
		upsertEligibilityFn: func(m models.RisikoEligibility) error {
			return nil
		},
	}

	svc := newService(mock)

	res, err := svc.ProcessEligibility("user-1", dto.EligibilityRequest{
		RisikoID:      ptrStr("1"),
		PernahTerjadi: true,
	})

	if err != nil {
		t.Fatal(err)
	}

	if res["next_step"] != "dampak" {
		t.Error("expected dampak")
	}
}

// TEST ELIGIBILITY INVALID RESPONDEN
func TestProcessEligibility_InvalidResponden(t *testing.T) {
	mock := &mockRisikoRepo{
		getRespondentIDByUserIDFn: func(userID string) (string, error) {
			return "", sql.ErrNoRows
		},
	}

	svc := newService(mock)

	_, err := svc.ProcessEligibility("user-1", dto.EligibilityRequest{
		RisikoID: ptrStr("1"),
	})

	if err == nil {
		t.Error("expected error")
	}
}

// TEST ALASAN SUCCESS
func TestProcessAlasan_Success(t *testing.T) {
	mock := &mockRisikoRepo{
		existsRespondenFn: func(id string) (bool, error) { return true, nil },
		existsRisikoFn:    func(id string) (bool, error) { return true, nil },
		upsertAlasanFn: func(m models.RisikoAlasan) error {
			return nil
		},
	}

	svc := newService(mock)

	res, err := svc.ProcessAlasan("user-1", dto.AlasanRequest{
		RisikoID: ptrStr("1"),
		Alasan:   "test",
	})

	if err != nil {
		t.Fatal(err)
	}

	if res["next_step"] != "finish" {
		t.Error("expected finish")
	}
}

// TEST DAMPAK INVALID
func TestProcessDampak_InvalidImpact(t *testing.T) {
	mock := &mockRisikoRepo{
		existsRespondenFn: func(id string) (bool, error) { return true, nil },
		existsRisikoFn:    func(id string) (bool, error) { return true, nil },
	}

	svc := newService(mock)

	_, err := svc.ProcessDampak("user-1", dto.DampakRequest{
		RisikoID: ptrStr("1"),
	})

	if err == nil {
		t.Error("expected validation error")
	}
}

// TEST PENGENDALIAN SUCCESS
func TestProcessPengendalian_Success(t *testing.T) {
	mock := &mockRisikoRepo{
		existsRespondenFn: func(id string) (bool, error) { return true, nil },
		existsRisikoFn:    func(id string) (bool, error) { return true, nil },
		upsertPengendalianFn: func(m models.RisikoPengendalian) error {
			return nil
		},
	}

	svc := newService(mock)

	res, err := svc.ProcessPengendalian("user-1", dto.PengendalianRequest{
		RisikoID:              ptrStr("1"),
		AdaPengendalian:       false,
		DeskripsiPengendalian: "",
	})

	if err != nil {
		t.Fatal(err)
	}

	if res["next_step"] != "finish" {
		t.Error("expected finish")
	}
}

func TestProcessDampak_MapsEnumValuesForDatabase(t *testing.T) {
	var saved models.RisikoDampak
	mock := &mockRisikoRepo{
		upsertDampakFn: func(m models.RisikoDampak) error {
			saved = m
			return nil
		},
	}

	svc := newService(mock)

	_, err := svc.ProcessDampak("user-1", dto.DampakRequest{
		RisikoID:          ptrStr("7"),
		DampakReputasi:    models.ImpactVerySignificant,
		DampakOperasional: models.ImpactSignificant,
		DampakFinansial:   models.ImpactFairlySignificant,
		DampakHukum:       models.ImpactNotSignificant,
		Frekuensi:         models.FrequencyLarge,
	})
	if err != nil {
		t.Fatal(err)
	}

	if saved.DampakReputasi != "Sangat Signifikan" {
		t.Fatalf("expected mapped reputasi, got %q", saved.DampakReputasi)
	}
	if saved.Frekuensi != "Besar" {
		t.Fatalf("expected mapped frekuensi, got %q", saved.Frekuensi)
	}
}

// TEST NAVIGATE NEXT
func TestNavigate_Next(t *testing.T) {
	mock := &mockRisikoRepo{
		getProgressFn: func(id string) (*models.SurveyProgress, error) {
			return &models.SurveyProgress{
				RespondenID: id,
				RisikoID:    sql.NullString{String: "1", Valid: true},
			}, nil
		},
		upsertProgressFn: func(p models.SurveyProgress) error {
			return nil
		},
	}

	svc := newService(mock)

	res, err := svc.Navigate("user-1", dto.NavigateRequest{
		Direction: "next",
	})

	if err != nil {
		t.Fatal(err)
	}

	if *res.RisikoID != "2" {
		t.Error("expected next risk")
	}
}

// TEST SAVE PROGRESS
func TestSaveProgress_Success(t *testing.T) {
	var savedProgress models.SurveyProgress
	mock := &mockRisikoRepo{
		getRisikoIDByUrutanFn: func(urutan int) (string, error) {
			if urutan == 5 {
				return "42", nil
			}
			return "", sql.ErrNoRows
		},
		getProgressFn: func(id string) (*models.SurveyProgress, error) {
			return &models.SurveyProgress{
				RespondenID: id,
			}, nil
		},
		upsertProgressFn: func(p models.SurveyProgress) error {
			savedProgress = p
			return nil
		},
	}

	svc := newService(mock)

	res, err := svc.SaveProgress("user-1", dto.NavigateRequest{
		CurrentRisk: 5,
	})

	if err != nil {
		t.Fatal(err)
	}

	if savedProgress.RisikoID.String != "42" {
		t.Errorf("expected persisted risiko_id 42, got %s", savedProgress.RisikoID.String)
	}

	if *res.RisikoID != "42" {
		t.Error("expected response risk to use mapped risiko_id")
	}
}

func TestSaveProgress_UnknownRiskStoresNull(t *testing.T) {
	var savedProgress models.SurveyProgress
	mock := &mockRisikoRepo{
		getRisikoIDByUrutanFn: func(urutan int) (string, error) {
			return "", sql.ErrNoRows
		},
		getProgressFn: func(id string) (*models.SurveyProgress, error) {
			return &models.SurveyProgress{
				RespondenID: id,
			}, nil
		},
		upsertProgressFn: func(p models.SurveyProgress) error {
			savedProgress = p
			return nil
		},
	}

	svc := newService(mock)

	res, err := svc.SaveProgress("user-1", dto.NavigateRequest{
		CurrentRisk: 999,
	})

	if err != nil {
		t.Fatal(err)
	}

	if savedProgress.RisikoID.Valid {
		t.Fatal("expected persisted risiko_id to be NULL for unknown current_risk")
	}

	if res.RisikoID != nil {
		t.Fatal("expected response risiko_id to be nil for unknown current_risk")
	}
}

func TestNavigate_SubmittedSurveyRejected(t *testing.T) {
	mock := &mockRisikoRepo{
		getProgressFn: func(id string) (*models.SurveyProgress, error) {
			return &models.SurveyProgress{
				RespondenID: id,
				Selesai:     true,
				Status:      SurveyStatusSubmitted,
			}, nil
		},
	}

	svc := newService(mock)

	_, err := svc.Navigate("user-1", dto.NavigateRequest{
		Direction: "next",
	})

	if err == nil {
		t.Fatal("expected error for submitted survey")
	}
}

func TestNavigate_EditApprovedSurveyAllowed(t *testing.T) {
	mock := &mockRisikoRepo{
		getProgressFn: func(id string) (*models.SurveyProgress, error) {
			return &models.SurveyProgress{
				RespondenID: id,
				RisikoID:    sql.NullString{String: "2", Valid: true},
				Selesai:     true,
				Status:      SurveyStatusEditApproved,
			}, nil
		},
		upsertProgressFn: func(p models.SurveyProgress) error {
			return nil
		},
	}

	svc := newService(mock)

	res, err := svc.Navigate("user-1", dto.NavigateRequest{
		Direction: "next",
	})

	if err != nil {
		t.Fatal(err)
	}

	if *res.RisikoID != "3" {
		t.Error("expected next risk for edit-approved survey")
	}
}

// TEST CREATE CUSTOM RISIKO
func TestCreateCustomRisiko_Success(t *testing.T) {
	mock := &mockRisikoRepo{
		insertCustomFn: func(id string, nama string) (int, error) {
			return 99, nil
		},
	}

	svc := newService(mock)

	id, err := svc.CreateCustomRisiko(dto.CustomRisikoRequest{
		RespondenID: "1",
		NamaRisiko:  "custom",
	})

	if err != nil {
		t.Fatal(err)
	}

	if id != 99 {
		t.Error("invalid ID")
	}
}

// TEST FINISH SURVEY
func TestFinishSurvey_Success(t *testing.T) {
	mock := &mockRisikoRepo{
		getProgressFn: func(id string) (*models.SurveyProgress, error) {
			return &models.SurveyProgress{
				RespondenID: id,
			}, nil
		},
		upsertProgressFn: func(p models.SurveyProgress) error {
			return nil
		},
	}

	svc := newService(mock)

	err := svc.FinishSurvey("user-1")

	if err != nil {
		t.Fatal(err)
	}
}
