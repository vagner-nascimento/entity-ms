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

func (e *entEntity) delete(id interface{}) (*model.Entity, *apperrors.Error) {
	return e.repo.Delete(id)
}

func (e *entEntity) search(parms map[string][]string) ([]model.Entity, *apperrors.Error) {
	return e.repo.Search(parms)
}

func newEntEntity() entEntity {
	return entEntity{
		repo: infra.NewEntityRepository(),
	}
}
