package app

import (
	"entity/src/apperrors"
	"entity/src/model"
)

type EntityAdapter interface {
	Save(ent *model.Entity) *apperrors.Error
}
