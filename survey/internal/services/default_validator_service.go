package services

import (
	"survey/internal/dto"
	"survey/internal/utils"
)

type DefaultValidator struct{}

func (d DefaultValidator) ValidateCreate(req dto.CreateRespondenRequest) error {
	return utils.ValidateCreateResponden(req)
}

func (v *DefaultValidator) ValidateUpdate(req dto.UpdateRespondenRequest) error {
	return utils.ValidateCreateResponden(dto.CreateRespondenRequest{
		IdPerusahaan:       req.IdPerusahaan,
		NamaLengkap:        req.NamaLengkap,
		Jabatan:            req.Jabatan,
		Email:              req.Email,
		NoTelepon:          req.NoTelepon,
		SertifikatTraining: req.SertifikatTraining,
	})
}
