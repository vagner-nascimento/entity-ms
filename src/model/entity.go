package model

import (
	"encoding/json"

	"github.com/go-playground/validator"
)

// TODO custom validation to avoid inform a future date
type Birthdate struct {
	Day   int16 `json:"day" validate:"required,min=1,max=31"`
	Month int16 `json:"month" validate:"required,min=1,max=12"`
	Year  int32 `json:"year" validate:"required,min=1900"`
}

type Entity struct {
	Id        int64     `json:"id" validate:"required,min=1"`
	Name      string    `json:"name" validate:"required,min=4"`
	Birthdate Birthdate `json:"birthDate" validate:"required"`
	Weight    float32   `json:"weight" validate:"min=1.5,max=599.99"`
}

func (p *Entity) Validate() (errs validator.ValidationErrors) {
	v := validator.New()

	if err := v.Struct(*p); err != nil {
		errs = err.(validator.ValidationErrors)
	}

	return errs
}

func NewEntityFromBytes(bys []byte) (ent Entity, err error) {
	err = json.Unmarshal(bys, &ent)

	return ent, err
}
