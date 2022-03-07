package app

import (
	"entity/src/apperrors"
	"entity/src/infra"
	"entity/src/model"
)

// business rules and data transformation layer
type entEntity struct {
	repo infra.EntityDataAdapter
}

func (e *entEntity) save(ent *model.Entity) *apperrors.Error {
	return e.repo.Save(ent)
}

func (e *entEntity) get(id interface{}) (*model.Entity, *apperrors.Error) {
	return e.repo.Get(id)
}

func (e *entEntity) update(id interface{}, newEnt model.Entity) (*model.Entity, *apperrors.Error) {
	return e.repo.Update(id, newEnt)
}

func newEntEntity() entEntity {
	return entEntity{
		repo: infra.NewEntityRepository(),
	}
}
