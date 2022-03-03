package infra

import (
	"entity/src/apperrors"
	"entity/src/model"
)

type EntityDataAdapter interface {
	Save(ent *model.Entity) *apperrors.Error
}
