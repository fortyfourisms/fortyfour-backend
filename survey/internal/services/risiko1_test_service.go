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
	existsRespondenFn    func(id int) (bool, error)
	existsRisikoFn       func(id int) (bool, error)
	existsCustomFn       func(id int) (bool, error)
	upsertEligibilityFn  func(m models.RisikoEligibility) error
	upsertAlasanFn       func(m models.RisikoAlasan) error
	upsertDampakFn       func(m models.RisikoDampak) error
	upsertPengendalianFn func(m models.RisikoPengendalian) error
	findByRespondentFn   func(id int) (map[string]interface{}, error)
	getProgressFn        func(id int) (*models.SurveyProgress, error)
	upsertProgressFn     func(p models.SurveyProgress) error
	insertCustomFn       func(id int, nama string) (int, error)
}

func (m *mockRisikoRepo) ExistsResponden(id int) (bool, error) {
	return m.existsRespondenFn(id)
}
func (m *mockRisikoRepo) ExistsRisiko(id int) (bool, error) {
	return m.existsRisikoFn(id)
}
func (m *mockRisikoRepo) ExistsCustomRisiko(id int) (bool, error) {
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
func (m *mockRisikoRepo) FindByRespondentID(id int) (map[string]interface{}, error) {
	return m.findByRespondentFn(id)
}
func (m *mockRisikoRepo) GetProgress(id int) (*models.SurveyProgress, error) {
	return m.getProgressFn(id)
}
func (m *mockRisikoRepo) UpsertProgress(p models.SurveyProgress) error {
	return m.upsertProgressFn(p)
}
func (m *mockRisikoRepo) InsertCustomRisiko(id int, nama string) (int, error) {
	return m.insertCustomFn(id, nama)
}

// HELPER 
func newService(mock *mockRisikoRepo) *RisikoService {
	return NewRisikoService(mock, &mockCache{})
}

// TEST ELIGIBILITY SUCCESS 
func TestProcessEligibility_Success(t *testing.T) {
	mock := &mockRisikoRepo{
		existsRespondenFn: func(id int) (bool, error) { return true, nil },
		existsRisikoFn:    func(id int) (bool, error) { return true, nil },
		upsertEligibilityFn: func(m models.RisikoEligibility) error {
			return nil
		},
	}

	svc := newService(mock)

	res, err := svc.ProcessEligibility(dto.EligibilityRequest{
		RespondenID:   1,
		RisikoID:      1,
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
		existsRespondenFn: func(id int) (bool, error) { return false, nil },
	}

	svc := newService(mock)

	_, err := svc.ProcessEligibility(dto.EligibilityRequest{
		RespondenID: 1,
		RisikoID:    1,
	})

	if err == nil {
		t.Error("expected error")
	}
}

// TEST ALASAN SUCCESS 
func TestProcessAlasan_Success(t *testing.T) {
	mock := &mockRisikoRepo{
		existsRespondenFn: func(id int) (bool, error) { return true, nil },
		existsRisikoFn:    func(id int) (bool, error) { return true, nil },
		upsertAlasanFn: func(m models.RisikoAlasan) error {
			return nil
		},
	}

	svc := newService(mock)

	res, err := svc.ProcessAlasan(dto.AlasanRequest{
		RespondenID: 1,
		RisikoID:    1,
		Alasan:      "test",
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
		existsRespondenFn: func(id int) (bool, error) { return true, nil },
		existsRisikoFn:    func(id int) (bool, error) { return true, nil },
	}

	svc := newService(mock)

	_, err := svc.ProcessDampak(dto.DampakRequest{
		RespondenID: 1,
		RisikoID:    1,
	})

	if err == nil {
		t.Error("expected validation error")
	}
}

// TEST PENGENDALIAN SUCCESS 
func TestProcessPengendalian_Success(t *testing.T) {
	mock := &mockRisikoRepo{
		existsRespondenFn: func(id int) (bool, error) { return true, nil },
		existsRisikoFn:    func(id int) (bool, error) { return true, nil },
		upsertPengendalianFn: func(m models.RisikoPengendalian) error {
			return nil
		},
	}

	svc := newService(mock)

	res, err := svc.ProcessPengendalian(dto.PengendalianRequest{
		RespondenID:           1,
		RisikoID:              1,
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

// TEST NAVIGATE NEXT 
func TestNavigate_Next(t *testing.T) {
	mock := &mockRisikoRepo{
		getProgressFn: func(id int) (*models.SurveyProgress, error) {
			return &models.SurveyProgress{
				RespondenID: id,
				RisikoID:    sql.NullInt64{Int64: 1, Valid: true},
			}, nil
		},
		upsertProgressFn: func(p models.SurveyProgress) error {
			return nil
		},
	}

	svc := newService(mock)

	res, err := svc.Navigate(dto.NavigateRequest{
		RespondenID: 1,
		Direction:   "next",
	})

	if err != nil {
		t.Fatal(err)
	}

	if *res.RisikoID != 2 {
		t.Error("expected next risk")
	}
}

// TEST SAVE PROGRESS 
func TestSaveProgress_Success(t *testing.T) {
	mock := &mockRisikoRepo{
		getProgressFn: func(id int) (*models.SurveyProgress, error) {
			return &models.SurveyProgress{
				RespondenID: id,
			}, nil
		},
		upsertProgressFn: func(p models.SurveyProgress) error {
			return nil
		},
	}

	svc := newService(mock)

	res, err := svc.SaveProgress(dto.NavigateRequest{
		RespondenID: 1,
		CurrentRisk: 5,
	})

	if err != nil {
		t.Fatal(err)
	}

	if *res.RisikoID != 5 {
		t.Error("expected saved risk")
	}
}

// TEST CREATE CUSTOM RISIKO 
func TestCreateCustomRisiko_Success(t *testing.T) {
	mock := &mockRisikoRepo{
		insertCustomFn: func(id int, nama string) (int, error) {
			return 99, nil
		},
	}

	svc := newService(mock)

	id, err := svc.CreateCustomRisiko(dto.CustomRisikoRequest{
		RespondenID: 1,
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
		getProgressFn: func(id int) (*models.SurveyProgress, error) {
			return &models.SurveyProgress{
				RespondenID: id,
			}, nil
		},
		upsertProgressFn: func(p models.SurveyProgress) error {
			return nil
		},
	}

	svc := newService(mock)

	err := svc.FinishSurvey(1)

	if err != nil {
		t.Fatal(err)
	}
}
