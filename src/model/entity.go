package model

import (
	"encoding/json"
	"entity/src/apperrors"
	"entity/src/infra/logger"

	"github.com/go-playground/validator"
)

// TODO custom validation to avoid inform a future date
type BirthDate struct {
	Day   int16 `json:"day" validate:"required,min=1,max=31"`
	Month int16 `json:"month" validate:"required,min=1,max=12"`
	Year  int32 `json:"year" validate:"required,min=1900"`
}

type Entity struct {
	Id        *int64    `json:"id"`
	Name      string    `json:"name" validate:"required,min=4,max=100"`
	BirthDate BirthDate `json:"birthDate" validate:"required"`
	Weight    float32   `json:"weight" validate:"min=1.5,max=599.99"`
}

func (e *Entity) Validate() (valid bool, errs []apperrors.Error) {
	v := validator.New()

	if err := v.Struct(*e); err != nil {
		errs = apperrors.NewValidationErrors(err.(validator.ValidationErrors))
	}

	return len(errs) == 0, errs
}

func (e *Entity) String() (str string) {
	if bys, err := json.Marshal(e); err == nil {
		str = string(bys)
	} else {
		logger.Error("Entity.String error", err)
	}

	return str
}

func NewEntityFromBytes(bys []byte) (ent Entity, err error) {
	err = json.Unmarshal(bys, &ent)

	return ent, err
}
