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
	existsRespondenFn         func(id int64) (bool, error)
	existsRisikoFn            func(id int64) (bool, error)
	existsCustomFn            func(id int64) (bool, error)
	getAllRisikoFn             func() ([]models.RisikoResponse, error)
	upsertEligibilityFn       func(m models.RisikoEligibility) error
	upsertAlasanFn            func(m models.RisikoAlasan) error
	upsertDampakFn            func(m models.RisikoDampak) error
	upsertPengendalianFn      func(m models.RisikoPengendalian) error
	findByRespondentFn        func(id int64) (map[string]interface{}, error)
	getProgressFn             func(id int64) (*models.SurveyProgress, error)
	upsertProgressFn          func(p models.SurveyProgress) error
	getRespondentIDByUserIDFn func(userID string) (int64, error)
	insertCustomFn            func(id int64, nama string) (int, error)
	getAllEditRequestsFn      func() ([]models.EditRequestItem, error)
	getEditRequestByUserIDFn  func(userID string) (*models.EditRequestItem, error)
}

func (m *mockRisikoRepo) ExistsResponden(id int64) (bool, error) {
	return m.existsRespondenFn(id)
}
func (m *mockRisikoRepo) ExistsRisiko(id int64) (bool, error) {
	return m.existsRisikoFn(id)
}
func (m *mockRisikoRepo) ExistsCustomRisiko(id int64) (bool, error) {
	return m.existsCustomFn(id)
}
func (m *mockRisikoRepo) GetAllRisiko() ([]models.RisikoResponse, error) {
	if m.getAllRisikoFn != nil {
		return m.getAllRisikoFn()
	}
	return nil, nil
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
func (m *mockRisikoRepo) FindByRespondentID(id int64) (map[string]interface{}, error) {
	return m.findByRespondentFn(id)
}
func (m *mockRisikoRepo) GetProgress(id int64) (*models.SurveyProgress, error) {
	return m.getProgressFn(id)
}
func (m *mockRisikoRepo) UpsertProgress(p models.SurveyProgress) error {
	return m.upsertProgressFn(p)
}
func (m *mockRisikoRepo) GetRespondentIDByUserID(userID string) (int64, error) {
	return m.getRespondentIDByUserIDFn(userID)
}
func (m *mockRisikoRepo) InsertCustomRisiko(id int64, nama string) (int, error) {
	return m.insertCustomFn(id, nama)
}
func (m *mockRisikoRepo) GetAllEditRequests() ([]models.EditRequestItem, error) {
	if m.getAllEditRequestsFn != nil {
		return m.getAllEditRequestsFn()
	}
	return nil, nil
}
func (m *mockRisikoRepo) GetEditRequestByUserID(userID string) (*models.EditRequestItem, error) {
	if m.getEditRequestByUserIDFn != nil {
		return m.getEditRequestByUserIDFn(userID)
	}
	return nil, nil
}

// HELPER
func ptrInt(v int) *int {
	return &v
}

func newService(mock *mockRisikoRepo) *RisikoService {
	if mock.getRespondentIDByUserIDFn == nil {
		mock.getRespondentIDByUserIDFn = func(userID string) (int64, error) {
			return 1, nil
		}
	}
	if mock.getProgressFn == nil {
		mock.getProgressFn = func(id int64) (*models.SurveyProgress, error) {
			return &models.SurveyProgress{
				RespondenID: id,
				Status:      SurveyStatusDraft,
			}, nil
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
		existsRespondenFn: func(id int64) (bool, error) { return true, nil },
		existsRisikoFn:    func(id int64) (bool, error) { return true, nil },
		upsertEligibilityFn: func(m models.RisikoEligibility) error {
			return nil
		},
	}

	svc := newService(mock)

	res, err := svc.ProcessEligibility("user-1", dto.EligibilityRequest{
		RespondenID:   1,
		RisikoID:      ptrInt(1),
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
		getRespondentIDByUserIDFn: func(userID string) (int64, error) {
			return 0, sql.ErrNoRows
		},
	}

	svc := newService(mock)

	_, err := svc.ProcessEligibility("user-1", dto.EligibilityRequest{
		RespondenID: 1,
		RisikoID:    ptrInt(1),
	})

	if err == nil {
		t.Error("expected error")
	}
}

// TEST ALASAN SUCCESS
func TestProcessAlasan_Success(t *testing.T) {
	mock := &mockRisikoRepo{
		existsRespondenFn: func(id int64) (bool, error) { return true, nil },
		existsRisikoFn:    func(id int64) (bool, error) { return true, nil },
		upsertAlasanFn: func(m models.RisikoAlasan) error {
			return nil
		},
	}

	svc := newService(mock)

	res, err := svc.ProcessAlasan("user-1", dto.AlasanRequest{
		RespondenID: 1,
		RisikoID:    ptrInt(1),
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
		existsRespondenFn: func(id int64) (bool, error) { return true, nil },
		existsRisikoFn:    func(id int64) (bool, error) { return true, nil },
	}

	svc := newService(mock)

	_, err := svc.ProcessDampak("user-1", dto.DampakRequest{
		RespondenID: 1,
		RisikoID:    ptrInt(1),
	})

	if err == nil {
		t.Error("expected validation error")
	}
}

// TEST PENGENDALIAN SUCCESS
func TestProcessPengendalian_Success(t *testing.T) {
	mock := &mockRisikoRepo{
		existsRespondenFn: func(id int64) (bool, error) { return true, nil },
		existsRisikoFn:    func(id int64) (bool, error) { return true, nil },
		upsertPengendalianFn: func(m models.RisikoPengendalian) error {
			return nil
		},
	}

	svc := newService(mock)

	res, err := svc.ProcessPengendalian("user-1", dto.PengendalianRequest{
		RespondenID:           1,
		RisikoID:              ptrInt(1),
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
		getProgressFn: func(id int64) (*models.SurveyProgress, error) {
			return &models.SurveyProgress{
				RespondenID: int64(id),
				RisikoID:    sql.NullInt64{Int64: 1, Valid: true},
			}, nil
		},
		upsertProgressFn: func(p models.SurveyProgress) error {
			return nil
		},
	}

	svc := newService(mock)

	res, err := svc.Navigate("user-1", dto.NavigateRequest{
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
		getProgressFn: func(id int64) (*models.SurveyProgress, error) {
			return &models.SurveyProgress{
				RespondenID: int64(id),
			}, nil
		},
		upsertProgressFn: func(p models.SurveyProgress) error {
			return nil
		},
	}

	svc := newService(mock)

	res, err := svc.SaveProgress("user-1", dto.NavigateRequest{
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

func TestNavigate_SubmittedSurveyRejected(t *testing.T) {
	mock := &mockRisikoRepo{
		getProgressFn: func(id int64) (*models.SurveyProgress, error) {
			return &models.SurveyProgress{
				RespondenID: int64(id),
				Selesai:     true,
				Status:      SurveyStatusSubmitted,
			}, nil
		},
	}

	svc := newService(mock)

	_, err := svc.Navigate("user-1", dto.NavigateRequest{
		RespondenID: 1,
		Direction:   "next",
	})

	if err == nil {
		t.Fatal("expected error for submitted survey")
	}
}

func TestNavigate_EditApprovedSurveyAllowed(t *testing.T) {
	mock := &mockRisikoRepo{
		getProgressFn: func(id int64) (*models.SurveyProgress, error) {
			return &models.SurveyProgress{
				RespondenID: int64(id),
				RisikoID:    sql.NullInt64{Int64: 2, Valid: true},
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
		RespondenID: 1,
		Direction:   "next",
	})

	if err != nil {
		t.Fatal(err)
	}

	if *res.RisikoID != 3 {
		t.Error("expected next risk for edit-approved survey")
	}
}

// TEST CREATE CUSTOM RISIKO
func TestCreateCustomRisiko_Success(t *testing.T) {
	mock := &mockRisikoRepo{
		insertCustomFn: func(id int64, nama string) (int, error) {
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
		getProgressFn: func(id int64) (*models.SurveyProgress, error) {
			return &models.SurveyProgress{
				RespondenID: int64(id),
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
