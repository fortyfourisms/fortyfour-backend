package services

import (
	"survey/internal/dto"
	"survey/internal/utils"
)

type DefaultValidator struct{}

func (d DefaultValidator) ValidateCreate(req dto.CreateRespondenRequest) error {
	return utils.ValidateCreateResponden(req)
}

func (d DefaultValidator) ValidateUpdate(req dto.UpdateRespondenRequest) error {
	return utils.ValidateUpdateResponden(req)
}