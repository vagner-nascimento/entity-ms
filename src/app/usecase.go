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
	if err := validateEntity(*ent); err != nil {
		return err
	}

	return eu.ent.save(ent)
}

func (eu *entUseCase) get(id interface{}) (*model.Entity, *apperrors.Error) {
	return eu.ent.get(id)
}

func (eu *entUseCase) update(id interface{}, newEnt model.Entity) (ent *model.Entity, err *apperrors.Error) {
	if err = validateEntity(newEnt); err == nil {
		ent, err = eu.ent.update(id, newEnt)
	}

	return ent, err
}

func validateEntity(ent model.Entity) *apperrors.Error {
	if ent.Id != nil {
		fildName := "Id"
		err := apperrors.NewValidationError("invalid field", &fildName, &ent.Id)

		return &err
	}

	return nil
}

func newEntityUseCase() entUseCase {
	return entUseCase{
		ent: newEntEntity(),
	}
}
