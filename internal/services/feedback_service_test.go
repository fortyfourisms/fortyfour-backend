package services_test

import (
	"errors"
	"testing"
	"time"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/models"
	"fortyfour-backend/internal/services"
)

// ═══════════════════════════════════════════════════════════════════════════
// MOCK: FeedbackRepository
// ═══════════════════════════════════════════════════════════════════════════

type MockFeedbackRepository struct {
	UpsertFunc            func(feedback *models.Feedback) error
	FindByUserAndMateriFn func(idUser, idMateri string) (*models.Feedback, error)
	FindByMateriFn        func(idMateri string) ([]dto.FeedbackListItem, error)
	DeleteFunc            func(id string) error
}

func (m *MockFeedbackRepository) Upsert(feedback *models.Feedback) error {
	if m.UpsertFunc != nil {
		return m.UpsertFunc(feedback)
	}
	return nil
}

func (m *MockFeedbackRepository) FindByUserAndMateri(idUser, idMateri string) (*models.Feedback, error) {
	if m.FindByUserAndMateriFn != nil {
		return m.FindByUserAndMateriFn(idUser, idMateri)
	}
	return nil, nil
}

func (m *MockFeedbackRepository) FindByMateri(idMateri string) ([]dto.FeedbackListItem, error) {
	if m.FindByMateriFn != nil {
		return m.FindByMateriFn(idMateri)
	}
	return nil, nil
}

func (m *MockFeedbackRepository) Delete(id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

// ═══════════════════════════════════════════════════════════════════════════
// TESTS: Upsert
// ═══════════════════════════════════════════════════════════════════════════

func TestFeedbackService_Upsert_Success(t *testing.T) {
	now := time.Now()
	mockRepo := &MockFeedbackRepository{
		UpsertFunc: func(feedback *models.Feedback) error {
			return nil
		},
		FindByUserAndMateriFn: func(idUser, idMateri string) (*models.Feedback, error) {
			return &models.Feedback{
				ID:        "fb-1",
				IDMateri:  idMateri,
				IDUser:    idUser,
				Konten:    "Great content!",
				CreatedAt: now,
				UpdatedAt: now,
			}, nil
		},
	}
	svc := services.NewFeedbackService(mockRepo)

	req := dto.UpsertFeedbackRequest{Konten: "Great content!"}
	resp, err := svc.Upsert("materi-1", "user-1", req)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if resp.IDMateri != "materi-1" {
		t.Errorf("expected id_materi 'materi-1', got '%s'", resp.IDMateri)
	}
	if resp.Konten != "Great content!" {
		t.Errorf("expected konten 'Great content!', got '%s'", resp.Konten)
	}
}

func TestFeedbackService_Upsert_EmptyKonten(t *testing.T) {
	svc := services.NewFeedbackService(&MockFeedbackRepository{})

	req := dto.UpsertFeedbackRequest{Konten: ""}
	_, err := svc.Upsert("materi-1", "user-1", req)

	if err == nil {
		t.Error("expected error for empty konten")
	}
	if err.Error() != "konten tidak boleh kosong" {
		t.Errorf("expected 'konten tidak boleh kosong', got '%s'", err.Error())
	}
}

func TestFeedbackService_Upsert_RepoError(t *testing.T) {
	mockRepo := &MockFeedbackRepository{
		UpsertFunc: func(feedback *models.Feedback) error {
			return errors.New("db error")
		},
	}
	svc := services.NewFeedbackService(mockRepo)

	req := dto.UpsertFeedbackRequest{Konten: "test"}
	_, err := svc.Upsert("materi-1", "user-1", req)

	if err == nil {
		t.Error("expected error")
	}
}

func TestFeedbackService_Upsert_FindAfterUpsertError(t *testing.T) {
	mockRepo := &MockFeedbackRepository{
		UpsertFunc: func(feedback *models.Feedback) error {
			return nil
		},
		FindByUserAndMateriFn: func(idUser, idMateri string) (*models.Feedback, error) {
			return nil, errors.New("find error")
		},
	}
	svc := services.NewFeedbackService(mockRepo)

	req := dto.UpsertFeedbackRequest{Konten: "test"}
	_, err := svc.Upsert("materi-1", "user-1", req)

	if err == nil {
		t.Error("expected error from FindByUserAndMateri")
	}
}

// ═══════════════════════════════════════════════════════════════════════════
// TESTS: GetByUserAndMateri
// ═══════════════════════════════════════════════════════════════════════════

func TestFeedbackService_GetByUserAndMateri_Success(t *testing.T) {
	now := time.Now()
	mockRepo := &MockFeedbackRepository{
		FindByUserAndMateriFn: func(idUser, idMateri string) (*models.Feedback, error) {
			return &models.Feedback{
				ID:        "fb-1",
				IDMateri:  idMateri,
				IDUser:    idUser,
				Konten:    "Nice",
				CreatedAt: now,
				UpdatedAt: now,
			}, nil
		},
	}
	svc := services.NewFeedbackService(mockRepo)

	resp, err := svc.GetByUserAndMateri("user-1", "materi-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.ID != "fb-1" {
		t.Errorf("expected ID 'fb-1', got '%s'", resp.ID)
	}
}

func TestFeedbackService_GetByUserAndMateri_NotFound(t *testing.T) {
	mockRepo := &MockFeedbackRepository{
		FindByUserAndMateriFn: func(idUser, idMateri string) (*models.Feedback, error) {
			return nil, errors.New("not found")
		},
	}
	svc := services.NewFeedbackService(mockRepo)

	_, err := svc.GetByUserAndMateri("user-1", "materi-1")
	if err == nil {
		t.Error("expected error")
	}
	if err.Error() != "feedback tidak ditemukan" {
		t.Errorf("expected 'feedback tidak ditemukan', got '%s'", err.Error())
	}
}

// ═══════════════════════════════════════════════════════════════════════════
// TESTS: GetAllByMateri
// ═══════════════════════════════════════════════════════════════════════════

func TestFeedbackService_GetAllByMateri_Success(t *testing.T) {
	mockRepo := &MockFeedbackRepository{
		FindByMateriFn: func(idMateri string) ([]dto.FeedbackListItem, error) {
			return []dto.FeedbackListItem{
				{ID: "fb-1", IDMateri: idMateri, Username: "user1", Konten: "Good"},
				{ID: "fb-2", IDMateri: idMateri, Username: "user2", Konten: "Great"},
			}, nil
		},
	}
	svc := services.NewFeedbackService(mockRepo)

	items, err := svc.GetAllByMateri("materi-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(items) != 2 {
		t.Errorf("expected 2 items, got %d", len(items))
	}
}

func TestFeedbackService_GetAllByMateri_EmptyResult(t *testing.T) {
	mockRepo := &MockFeedbackRepository{
		FindByMateriFn: func(idMateri string) ([]dto.FeedbackListItem, error) {
			return nil, nil // nil slice from repo
		},
	}
	svc := services.NewFeedbackService(mockRepo)

	items, err := svc.GetAllByMateri("materi-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	// Service should return empty slice, not nil
	if items == nil {
		t.Error("expected non-nil empty slice")
	}
	if len(items) != 0 {
		t.Errorf("expected 0 items, got %d", len(items))
	}
}

func TestFeedbackService_GetAllByMateri_Error(t *testing.T) {
	mockRepo := &MockFeedbackRepository{
		FindByMateriFn: func(idMateri string) ([]dto.FeedbackListItem, error) {
			return nil, errors.New("db error")
		},
	}
	svc := services.NewFeedbackService(mockRepo)

	_, err := svc.GetAllByMateri("materi-1")
	if err == nil {
		t.Error("expected error")
	}
}
