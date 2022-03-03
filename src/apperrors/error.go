package apperrors

import (
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

func NewValidationError(msg string, field *string, val interface{}) Error {
	tp := ERR_TP_VALIDATION
	return Error{
		Message: &msg,
		Type:    &tp,
		Field:   field,
		Value:   val,
	}
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
		fld := e.Namespace()[indx:]

		errs = append(errs, Error{
			Message: &msg,
			Type:    &tp,
			Field:   &fld,
			Value:   e.Value(),
		})
	}

	return errs
}
