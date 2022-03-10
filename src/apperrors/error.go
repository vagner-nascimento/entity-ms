package apperrors

import (
	"entity/src/utils"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator"
)

type Error struct {
	Message *string     `json:"message"`
	Type    *string     `json:"type"`
	Field   *string     `json:"field"`
	Value   interface{} `json:"value"`
}

func (e *Error) Error() string {
	return *e.Message
}

const ERR_TP_VALIDATION string = "validation"
const ERR_TP_DATA string = "data"
const ERR_TP_INFRA string = "infrastructure"
const ERR_TP_NOT_FOUND string = "not_found"

func newError(msg *string, tp string, val interface{}, fil *string) Error {
	return Error{
		Message: msg,
		Type:    &tp,
		Value:   val,
		Field:   fil,
	}
}

func NewValidationError(msg string, field *string, val interface{}) Error {
	return newError(&msg, ERR_TP_VALIDATION, val, field)
}

func NewDataError(msg string, val interface{}) Error {
	return newError(&msg, ERR_TP_DATA, val, nil)
}

func NewInfraError(msg string, val interface{}) Error {
	return newError(&msg, ERR_TP_INFRA, val, nil)
}

func NewNotFoundError(msg string, val interface{}) Error {
	return newError(&msg, ERR_TP_NOT_FOUND, val, nil)
}

// NOTE: only required, min and max validation tags are implemented
func NewValidationErrors(vErs validator.ValidationErrors) (errs []Error) {
	for _, e := range vErs {
		valTag := e.ActualTag()
		dtTyp := e.Kind()
		parm := e.Param()
		var msg string

		switch valTag {
		case "required":
			msg = "field is required"
		case "min", "max":
			switch dtTyp {
			case reflect.String:
				msg = fmt.Sprintf("value must have %s length equals %s ", valTag, parm)
			default:
				msg = fmt.Sprintf("value doesn't match with validation %s %s", valTag, parm)
			}
		default:
			msg = fmt.Sprint(e)
		}

		tp := ERR_TP_VALIDATION
		indx := strings.Index(e.Namespace(), ".") + 1
		fld := utils.LowerFirst(e.Namespace()[indx:])

		errs = append(errs, Error{
			Message: &msg,
			Type:    &tp,
			Field:   &fld,
			Value:   e.Value(),
		})
	}

	return
}
