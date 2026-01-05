package services

import "fortyfour-backend/internal/dto"

type mockSdmCsirtRepo struct {
	CreateFn  func(req dto.CreateSdmCsirtRequest, id string) error
	GetAllFn  func() ([]dto.SdmCsirtResponse, error)
	GetByIDFn func(id string) (*dto.SdmCsirtResponse, error)
	UpdateFn  func(id string, req dto.SdmCsirtResponse) error
	DeleteFn  func(id string) error
}

func (m *mockSdmCsirtRepo) Create(req dto.CreateSdmCsirtRequest, id string) error {
	return m.CreateFn(req, id)
}
func (m *mockSdmCsirtRepo) GetAll() ([]dto.SdmCsirtResponse, error) {
	return m.GetAllFn()
}
func (m *mockSdmCsirtRepo) GetByID(id string) (*dto.SdmCsirtResponse, error) {
	return m.GetByIDFn(id)
}
func (m *mockSdmCsirtRepo) Update(id string, req dto.SdmCsirtResponse) error {
	return m.UpdateFn(id, req)
}
func (m *mockSdmCsirtRepo) Delete(id string) error {
	return m.DeleteFn(id)
}
