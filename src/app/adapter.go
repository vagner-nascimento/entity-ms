package app

import (
	"entity/src/apperrors"
	"entity/src/model"
)

// data adaption and control flow layer
type entAdapter struct {
	uc entUseCase
}

func (ea *entAdapter) Save(ent *model.Entity) *apperrors.Error {
	return ea.uc.save(ent)
}

func (ea *entAdapter) Get(id interface{}) (*model.Entity, *apperrors.Error) {
	return ea.uc.get(id)
}

func NewEnityAdapter() EntityAdapter {
	return &entAdapter{
		uc: newEntityUseCase(),
	}
}
