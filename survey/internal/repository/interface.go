package repository

import "survey/internal/models"

type EditRequestRepositoryInterface interface {
	Create(req *models.EditRequest) error
	FindByID(id string) (*models.EditRequest, error)
	FindAllPending() ([]models.EditRequest, error)
	FindByUser(userID string) ([]models.EditRequest, error)
	FindPendingByRespondenRisiko(respondenID, risikoID int) ([]models.EditRequest, error)
	UpdateStatus(id string, status models.EditRequestStatus, catatan *string) error
}

type RisikoRepositoryInterface interface {
	GetByID(id int) (*models.Risiko, error)
	UpdatePartial(id int, data map[string]interface{}) error
}
