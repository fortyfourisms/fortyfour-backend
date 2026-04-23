package services

import (
	"errors"
	"testing"
	"time"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/models"
)

type fakeSEEditRequestRepo struct {
	createFn        func(req *models.SEEditRequest) error
	findByIDFn      func(id string) (*models.SEEditRequest, error)
	findAllPendingFn func() ([]models.SEEditRequest, error)
	findByUserFn    func(userID string) ([]models.SEEditRequest, error)
	updateStatusFn  func(id string, status models.SEEditRequestStatus, catatan *string) error
}

func (f *fakeSEEditRequestRepo) Create(req *models.SEEditRequest) error {
	if f.createFn != nil {
		return f.createFn(req)
	}
	return nil
}

func (f *fakeSEEditRequestRepo) FindByID(id string) (*models.SEEditRequest, error) {
	if f.findByIDFn != nil {
		return f.findByIDFn(id)
	}
	return nil, nil
}

func (f *fakeSEEditRequestRepo) FindPendingBySE(idSE string) ([]models.SEEditRequest, error) {
	return nil, nil
}

func (f *fakeSEEditRequestRepo) FindAllPending() ([]models.SEEditRequest, error) {
	if f.findAllPendingFn != nil {
		return f.findAllPendingFn()
	}
	return nil, nil
}

func (f *fakeSEEditRequestRepo) FindByUser(idUser string) ([]models.SEEditRequest, error) {
	if f.findByUserFn != nil {
		return f.findByUserFn(idUser)
	}
	return nil, nil
}

func (f *fakeSEEditRequestRepo) UpdateStatus(id string, status models.SEEditRequestStatus, catatan *string) error {
	if f.updateStatusFn != nil {
		return f.updateStatusFn(id, status, catatan)
	}
	return nil
}

type fakeSERepo struct {
	getByIDFn func(id string) (*dto.SEResponse, error)
}

func (f *fakeSERepo) Create(req dto.CreateSERequest, id string, totalBobot int, kategori string) error {
	return nil
}

func (f *fakeSERepo) GetAll() ([]dto.SEResponse, error) {
	return nil, nil
}

func (f *fakeSERepo) GetByID(id string) (*dto.SEResponse, error) {
	if f.getByIDFn != nil {
		return f.getByIDFn(id)
	}
	return nil, nil
}

func (f *fakeSERepo) GetByPerusahaan(idPerusahaan string) ([]dto.SEResponse, error) {
	return nil, nil
}

func (f *fakeSERepo) Update(id string, req dto.UpdateSERequest, totalBobot int, kategori string) error {
	return nil
}

func (f *fakeSERepo) Delete(id string) error {
	return nil
}

type fakeSEService struct {
	updateFn func(id string, req dto.UpdateSERequest) (*dto.SEResponse, error)
}

func (f *fakeSEService) Create(req dto.CreateSERequest) (*dto.SEResponse, error) {
	return nil, nil
}

func (f *fakeSEService) GetAll() ([]dto.SEResponse, error) {
	return nil, nil
}

func (f *fakeSEService) GetByID(id string) (*dto.SEResponse, error) {
	return nil, nil
}

func (f *fakeSEService) GetByPerusahaan(idPerusahaan string) ([]dto.SEResponse, error) {
	return nil, nil
}

func (f *fakeSEService) Update(id string, req dto.UpdateSERequest) (*dto.SEResponse, error) {
	if f.updateFn != nil {
		return f.updateFn(id, req)
	}
	return nil, nil
}

func (f *fakeSEService) Delete(id string) error {
	return nil
}

func TestSEEditRequestService_CreateRequest_Success(t *testing.T) {
	var created *models.SEEditRequest
	catatan := "tolong ubah nama"
	namaSE := "SE Baru"

	svc := NewSEEditRequestService(
		&fakeSEEditRequestRepo{
			createFn: func(req *models.SEEditRequest) error {
				created = req
				return nil
			},
		},
		&fakeSERepo{
			getByIDFn: func(id string) (*dto.SEResponse, error) {
				return &dto.SEResponse{ID: id}, nil
			},
		},
		&fakeSEService{},
	)

	resp, err := svc.CreateRequest("user-1", "se-1", dto.CreateSEEditRequestDTO{
		Catatan: &catatan,
		DataPerubahan: dto.UpdateSERequest{
			NamaSE: &namaSE,
		},
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp == nil {
		t.Fatal("expected response, got nil")
	}
	if created == nil {
		t.Fatal("expected repo.Create to be called")
	}
	if created.IDSE != "se-1" || created.IDUser != "user-1" {
		t.Fatalf("unexpected created request: %+v", created)
	}
	if created.Status != models.SEEditRequestPending {
		t.Fatalf("expected pending status, got %s", created.Status)
	}
	if resp.DataPerubahan == nil || resp.DataPerubahan.NamaSE == nil || *resp.DataPerubahan.NamaSE != namaSE {
		t.Fatal("expected data_perubahan to be mapped back to response")
	}
}

func TestSEEditRequestService_CreateRequest_SENotFound(t *testing.T) {
	svc := NewSEEditRequestService(
		&fakeSEEditRequestRepo{},
		&fakeSERepo{
			getByIDFn: func(id string) (*dto.SEResponse, error) {
				return nil, errors.New("not found")
			},
		},
		&fakeSEService{},
	)

	resp, err := svc.CreateRequest("user-1", "se-x", dto.CreateSEEditRequestDTO{})
	if err == nil {
		t.Fatal("expected error when SE is missing")
	}
	if resp != nil {
		t.Fatal("expected nil response on error")
	}
}

func TestSEEditRequestService_GetPending_Success(t *testing.T) {
	now := time.Now()
	namaSE := "SE Baru"
	data := `{"nama_se":"SE Baru"}`

	svc := NewSEEditRequestService(
		&fakeSEEditRequestRepo{
			findAllPendingFn: func() ([]models.SEEditRequest, error) {
				return []models.SEEditRequest{
					{
						ID:            "req-1",
						IDSE:          "se-1",
						IDUser:        "user-1",
						Status:        models.SEEditRequestPending,
						DataPerubahan: data,
						CreatedAt:     now,
						UpdatedAt:     now,
					},
				}, nil
			},
		},
		&fakeSERepo{},
		&fakeSEService{},
	)

	resp, err := svc.GetPending()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(resp) != 1 {
		t.Fatalf("expected 1 response, got %d", len(resp))
	}
	if resp[0].DataPerubahan == nil || resp[0].DataPerubahan.NamaSE == nil || *resp[0].DataPerubahan.NamaSE != namaSE {
		t.Fatal("expected pending item to contain parsed update payload")
	}
}

func TestSEEditRequestService_GetByUser_RepoError(t *testing.T) {
	svc := NewSEEditRequestService(
		&fakeSEEditRequestRepo{
			findByUserFn: func(userID string) ([]models.SEEditRequest, error) {
				return nil, errors.New("db error")
			},
		},
		&fakeSERepo{},
		&fakeSEService{},
	)

	resp, err := svc.GetByUser("user-1")
	if err == nil {
		t.Fatal("expected repo error")
	}
	if resp != nil {
		t.Fatal("expected nil response on error")
	}
}

func TestSEEditRequestService_Review_ApproveSuccess(t *testing.T) {
	namaSE := "SE Approved"
	catatan := "approve"
	updateCalled := false
	updateStatusCalled := false

	svc := NewSEEditRequestService(
		&fakeSEEditRequestRepo{
			findByIDFn: func(id string) (*models.SEEditRequest, error) {
				return &models.SEEditRequest{
					ID:            "req-1",
					IDSE:          "se-1",
					IDUser:        "user-1",
					Status:        models.SEEditRequestPending,
					DataPerubahan: `{"nama_se":"SE Approved"}`,
					CreatedAt:     time.Now(),
					UpdatedAt:     time.Now(),
				}, nil
			},
			updateStatusFn: func(id string, status models.SEEditRequestStatus, catatanArg *string) error {
				updateStatusCalled = true
				if id != "req-1" || status != models.SEEditRequestApproved {
					t.Fatalf("unexpected update status call: %s %s", id, status)
				}
				if catatanArg == nil || *catatanArg != catatan {
					t.Fatal("expected review note to be forwarded")
				}
				return nil
			},
		},
		&fakeSERepo{},
		&fakeSEService{
			updateFn: func(id string, req dto.UpdateSERequest) (*dto.SEResponse, error) {
				updateCalled = true
				if id != "se-1" {
					t.Fatalf("expected SE update on se-1, got %s", id)
				}
				if req.NamaSE == nil || *req.NamaSE != namaSE {
					t.Fatal("expected update payload to be parsed from request data")
				}
				return &dto.SEResponse{ID: id}, nil
			},
		},
	)

	resp, err := svc.Review("req-1", dto.ReviewSEEditRequestDTO{
		Status:  "approved",
		Catatan: &catatan,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !updateCalled || !updateStatusCalled {
		t.Fatal("expected approve review to update SE and request status")
	}
	if resp.Status != models.SEEditRequestApproved {
		t.Fatalf("expected approved status, got %s", resp.Status)
	}
}

func TestSEEditRequestService_Review_RejectSuccess(t *testing.T) {
	catatan := "ditolak"
	updateCalled := false

	svc := NewSEEditRequestService(
		&fakeSEEditRequestRepo{
			findByIDFn: func(id string) (*models.SEEditRequest, error) {
				return &models.SEEditRequest{
					ID:            "req-1",
					IDSE:          "se-1",
					IDUser:        "user-1",
					Status:        models.SEEditRequestPending,
					DataPerubahan: `{}`,
					CreatedAt:     time.Now(),
					UpdatedAt:     time.Now(),
				}, nil
			},
			updateStatusFn: func(id string, status models.SEEditRequestStatus, catatanArg *string) error {
				if status != models.SEEditRequestRejected {
					t.Fatalf("expected rejected status, got %s", status)
				}
				return nil
			},
		},
		&fakeSERepo{},
		&fakeSEService{
			updateFn: func(id string, req dto.UpdateSERequest) (*dto.SEResponse, error) {
				updateCalled = true
				return nil, nil
			},
		},
	)

	resp, err := svc.Review("req-1", dto.ReviewSEEditRequestDTO{
		Status:  "rejected",
		Catatan: &catatan,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if updateCalled {
		t.Fatal("reject review must not update SE data")
	}
	if resp.Status != models.SEEditRequestRejected {
		t.Fatalf("expected rejected status, got %s", resp.Status)
	}
}

func TestSEEditRequestService_Review_InvalidStatus(t *testing.T) {
	svc := NewSEEditRequestService(
		&fakeSEEditRequestRepo{
			findByIDFn: func(id string) (*models.SEEditRequest, error) {
				return &models.SEEditRequest{
					ID:            "req-1",
					Status:        models.SEEditRequestPending,
					DataPerubahan: `{}`,
				}, nil
			},
		},
		&fakeSERepo{},
		&fakeSEService{},
	)

	resp, err := svc.Review("req-1", dto.ReviewSEEditRequestDTO{Status: "draft"})
	if err == nil {
		t.Fatal("expected invalid status error")
	}
	if resp != nil {
		t.Fatal("expected nil response on error")
	}
}

func TestSEEditRequestService_Review_AlreadyProcessed(t *testing.T) {
	svc := NewSEEditRequestService(
		&fakeSEEditRequestRepo{
			findByIDFn: func(id string) (*models.SEEditRequest, error) {
				return &models.SEEditRequest{
					ID:            "req-1",
					Status:        models.SEEditRequestApproved,
					DataPerubahan: `{}`,
				}, nil
			},
		},
		&fakeSERepo{},
		&fakeSEService{},
	)

	resp, err := svc.Review("req-1", dto.ReviewSEEditRequestDTO{Status: "approved"})
	if err == nil {
		t.Fatal("expected already processed error")
	}
	if resp != nil {
		t.Fatal("expected nil response on error")
	}
}

func TestSEEditRequestService_Review_InvalidStoredPayload(t *testing.T) {
	svc := NewSEEditRequestService(
		&fakeSEEditRequestRepo{
			findByIDFn: func(id string) (*models.SEEditRequest, error) {
				return &models.SEEditRequest{
					ID:            "req-1",
					IDSE:          "se-1",
					Status:        models.SEEditRequestPending,
					DataPerubahan: `{invalid-json}`,
				}, nil
			},
		},
		&fakeSERepo{},
		&fakeSEService{},
	)

	resp, err := svc.Review("req-1", dto.ReviewSEEditRequestDTO{Status: "approved"})
	if err == nil {
		t.Fatal("expected invalid payload error")
	}
	if resp != nil {
		t.Fatal("expected nil response on error")
	}
}
