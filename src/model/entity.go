package model

import (
	"encoding/json"
	"entity/src/apperrors"
	"entity/src/infra/logger"
	"time"

	"github.com/go-playground/validator"
)

// TODO custom validation to avoid inform a future date
type BirthDate struct {
	Day   int16 `json:"day" validate:"required,min=1,max=31" bson:"day"`
	Month int16 `json:"month" validate:"required,min=1,max=12" bson:"month"`
	Year  int32 `json:"year" validate:"required,min=1900" bson:"year"`
}

type Entity struct {
	Id        interface{} `json:"id" bson:"-"`
	Name      string      `json:"name" validate:"required,min=4,max=100" bson:"name"`
	BirthDate BirthDate   `json:"birthDate" validate:"required" bson:"birthDate"`
	Weight    *float32    `json:"weight" validate:"min=1.5,max=599.99" bson:"weight"`
	CreatedAt *time.Time  `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt *time.Time  `json:"updatedAt" bson:"updatedAt,omitempty"`
	DeletedAt *time.Time  `json:"deletedAt" bson:"deletedAt,omitempty"`
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
