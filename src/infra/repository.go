package infra

import (
	"entity/src/apperrors"
	"entity/src/infra/database"
	"entity/src/model"
	"time"
)

// data adaption and control flow layer
type entRepository struct {
	db database.DataBaseHandler
	tb string
}

func (er *entRepository) Save(ent *model.Entity) (err *apperrors.Error) {
	cdate := time.Now()
	ent.CreatedAt = &cdate
	ent.UpdatedAt = nil
	ent.DeletedAt = nil

	var id interface{}
	if id, err = er.db.Insert(ent, "entity"); err == nil {
		ent.Id = id
	}

	return err
}

func (er *entRepository) Get(id interface{}) (res *model.Entity, err *apperrors.Error) {
	var data []byte
	if data, err = er.db.Get(id, "entity"); err == nil {
		parseEntity(res, data)
	}

	return res, err
}

func NewEntityRepository() EntityDataAdapter {
	return &entRepository{
		db: database.NewDatabaseConnection(),
		tb: "entity",
	}
}

/*
 * Auxiliar Functions
 */
func parseEntity(ent *model.Entity, data []byte) {
	nent, _ := model.NewEntityFromBytes(data)
	ent = &nent
}
