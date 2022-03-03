package infra

import (
	"entity/src/apperrors"
	"entity/src/model"
)

type entRepository struct {
}

// TODO implements some SQL database
func (er *entRepository) Save(ent *model.Entity) *apperrors.Error {
	var id int64
	id = 69
	ent.Id = &id

	return nil
}

func NewEntityRepository() EntityDataAdapter {
	return &entRepository{}
}
