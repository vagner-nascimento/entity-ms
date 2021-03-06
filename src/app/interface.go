package app

import (
	"entity/src/apperrors"
	"entity/src/model"
)

type EntityAdapter interface {
	Save(ent *model.Entity) *apperrors.Error
	Get(id interface{}) (*model.Entity, *apperrors.Error)
	Update(id interface{}, newEnt model.Entity) (*model.Entity, *apperrors.Error)
	Delete(id interface{}) (*model.Entity, *apperrors.Error)
	Search(parms map[string][]string) ([]model.Entity, *apperrors.Error)
}
