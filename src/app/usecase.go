package app

import (
	"entity/src/apperrors"
	"entity/src/model"
)

// application rules layer
type entUseCase struct {
	ent entEntity
}

func (eu *entUseCase) save(ent *model.Entity) *apperrors.Error {
	if ent.Id != nil {
		fildName := "Id"
		err := apperrors.NewValidationError("invalid field", &fildName, &ent.Id)

		return &err
	}

	return eu.ent.save(ent)
}

func newEntityUseCase() entUseCase {
	return entUseCase{
		ent: newEntEntity(),
	}
}
