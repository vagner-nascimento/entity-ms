package model

import (
	"encoding/json"
	"entity/src/apperrors"
	"entity/src/infra/logger"
	"time"

	"github.com/go-playground/validator"
)

/*
 * Entity Definitions.
 * If add new fields, change or delete a existing one, remember to handle it on the functions.
 */
type BirthDate struct {
	Day   int16 `json:"day" validate:"required,min=1,max=31" bson:"day,omitempty"`
	Month int16 `json:"month" validate:"required,min=1,max=12" bson:"month,omitempty"`
	Year  int32 `json:"year" validate:"required,min=1900" bson:"year,omitempty"`
}

// TODO custom validation to avoid inform a future date
type Entity struct {
	Id        interface{} `json:"id" bson:"-"`
	Name      string      `json:"name" validate:"required,min=4,max=100" bson:"name,omitempty"`
	Weight    *float32    `json:"weight" validate:"min=1.5,max=599.99" bson:"weight,omitempty"`
	BirthDate *BirthDate  `json:"birthDate" validate:"required" bson:"birthDate,omitempty"`
	CreatedAt *time.Time  `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt *time.Time  `json:"updatedAt" bson:"updatedAt,omitempty"`
	DeletedAt *time.Time  `json:"deletedAt" bson:"deletedAt,omitempty"`
}

func GetEntityBsonFields() (fils []string) {
	fils = append(fils, "name")
	fils = append(fils, "weight")
	fils = append(fils, "birthDate")
	fils = append(fils, "birthDate.day")
	fils = append(fils, "birthDate.month")
	fils = append(fils, "birthDate.year")
	fils = append(fils, "createdAt")
	fils = append(fils, "updatedAt")
	fils = append(fils, "deletedAt")

	return
}

func (e *Entity) Validate() (valid bool, errs []apperrors.Error) {
	v := validator.New()

	if err := v.Struct(*e); err != nil {
		errs = apperrors.NewValidationErrors(err.(validator.ValidationErrors))
	}

	valid = len(errs) == 0

	return
}

func (e *Entity) ValidateName() (valid bool, errs []apperrors.Error) {
	var w float32 = 500.1
	valEnt := Entity{
		Name:   e.Name,
		Weight: &w,
		BirthDate: &BirthDate{
			Day:   int16(time.Now().Day()) - 1,
			Month: int16(time.Now().Month()),
			Year:  int32(time.Now().Year()),
		},
	}

	return valEnt.Validate()
}

func (e *Entity) ValidateWeigth() (valid bool, errs []apperrors.Error) {
	valEnt := Entity{
		Name:   "AAAAAAAAA",
		Weight: e.Weight,
		BirthDate: &BirthDate{
			Day:   int16(time.Now().Day()) - 1,
			Month: int16(time.Now().Month()),
			Year:  int32(time.Now().Year()),
		},
	}

	return valEnt.Validate()
}

func (e *Entity) ValidateBirthDate() (valid bool, errs []apperrors.Error) {
	var w float32 = 500.1
	valEnt := Entity{
		Name:      "AAAAAAAAA",
		Weight:    &w,
		BirthDate: e.BirthDate,
	}

	return valEnt.Validate()
}

func (e *Entity) NilAllButBithDate() {
	e.Id = nil
	e.Weight = nil
	e.Name = ""
	e.CreatedAt = nil
	e.UpdatedAt = nil
	e.DeletedAt = nil
}

func (e *Entity) NilAllButWeight() {
	e.Id = nil
	e.BirthDate = nil
	e.Name = ""
	e.CreatedAt = nil
	e.UpdatedAt = nil
	e.DeletedAt = nil
}

func (e *Entity) NilAllButName() {
	e.Id = nil
	e.BirthDate = nil
	e.Weight = nil
	e.CreatedAt = nil
	e.UpdatedAt = nil
	e.DeletedAt = nil
}

func (e *Entity) String() (str string) {
	if bys, err := json.Marshal(e); err == nil {
		str = string(bys)
	} else {
		logger.Error("Entity.String error", err)
	}

	return
}

func NewEntityFromBytes(bys []byte) (ent Entity, err error) {
	err = json.Unmarshal(bys, &ent)

	return
}
