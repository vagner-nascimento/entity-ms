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

func newEntEntity() entEntity {
	return entEntity{
		repo: infra.NewEntityRepository(),
	}
}
