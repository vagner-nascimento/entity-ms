package infra

import (
	"entity/src/apperrors"
	"entity/src/infra/database"
	"entity/src/model"
	"time"
)

// data adaption and control flow layer
type entRepository struct {
	db       database.DataBaseHandler
	tb       string
	activFil map[string]interface{}
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

	return
}

func (er *entRepository) Get(id interface{}) (*model.Entity, *apperrors.Error) {
	var ent model.Entity
	err := er.db.Get(id, "entity", &ent, er.activFil)
	ent.Id = id

	return &ent, err
}

func (er *entRepository) Update(id interface{}, newEnt model.Entity) (res *model.Entity, err *apperrors.Error) {
	updt := time.Now()
	newEnt.UpdatedAt = &updt
	newEnt.DeletedAt = nil
	newEnt.CreatedAt = nil

	if err = er.db.Update(id, newEnt, "entity", &res, er.activFil); err == nil {
		res.Id = id
	}

	return
}

func (er *entRepository) Delete(id interface{}) (res *model.Entity, err *apperrors.Error) {
	dldt := time.Now()
	del := model.Entity{DeletedAt: &dldt}

	if err = er.db.Update(id, del, "entity", &res, er.activFil); err == nil {
		res.Id = id
	}

	return
}

func (er *entRepository) Search(parms map[string][]string) (res []model.Entity, err *apperrors.Error) {
	e := apperrors.NewDataError("not implemented yet", nil)
	err = &e
	// TODO convert res to model.Entity and return the properly response
	er.db.Search(parms, "entity")

	return
}

func NewEntityRepository() EntityDataAdapter {
	return &entRepository{
		db:       database.NewDatabaseConnection(),
		tb:       "entity",
		activFil: map[string]interface{}{"deletedAt": nil},
	}
}

/*
 * Auxiliar Functions
 */
func parseEntity(ent *model.Entity, data []byte) {
	nent, _ := model.NewEntityFromBytes(data)
	ent = &nent
}
