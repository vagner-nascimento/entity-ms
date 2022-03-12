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

func (ea *entAdapter) Update(id interface{}, newEnt model.Entity) (*model.Entity, *apperrors.Error) {
	return ea.uc.update(id, newEnt)
}

func (ea *entAdapter) Delete(id interface{}) (*model.Entity, *apperrors.Error) {
	return ea.uc.delete(id)
}

func (ea *entAdapter) Search(parms map[string][]string) ([]model.Entity, *apperrors.Error) {
	return ea.uc.search(parms)
}

func NewEnityAdapter() EntityAdapter {
	return &entAdapter{
		uc: newEntityUseCase(),
	}
}
