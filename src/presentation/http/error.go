package http

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator"
)

type httpError struct {
	Message *string `json:"message"`
	Type    *string `json:"type"`
	Field   *string `json:"field"`
}

type httpErrors struct {
	Errors []httpError `json:"errors"`
}

// ValidationErrors examples: https://github.com/go-playground/validator/blob/master/_examples/simple/main.go#L69
func newValidationErrors(errs validator.ValidationErrors) (httpErrs httpErrors) {
	for _, e := range errs {
		fmt.Println("-- INI --")
		fmt.Println("Field", e.Field())
		fmt.Println("ActualTag", e.ActualTag())
		fmt.Println("Value", e.Value())
		fmt.Println("Param", e.Param())
		fmt.Println("Kind", e.Kind())
		fmt.Println("-- END --")
		// TODO handle birthDate struct required and other fields

		var msg string
		switch e.ActualTag() {
		case "required":
			msg = "field is required"
		case "min", "max":
			{
				// TODO handle each type to properly formats it
				if e.Kind() == reflect.String {
					msg = fmt.Sprintf("value '%s' doesn't match with a %s of %s allowed length", e.Value(), e.ActualTag(), e.Param())
				} else {
					msg = fmt.Sprintf("value '%d' doesn't match with validation '%s=%s'", e.Value(), e.ActualTag(), e.Param())
				}
			}
		default:
			msg = fmt.Sprint(e)
		}

		tp := "validation"
		field := e.Field()

		httpErrs.Errors = append(httpErrs.Errors, httpError{
			Message: &msg,
			Type:    &tp,
			Field:   &field,
		})
	}

	return httpErrs
}

func newConversionError(err error) httpError {
	msg := err.Error()
	tp := "validation"

	return httpError{
		Message: &msg,
		Type:    &tp,
		Field:   nil,
	}
}
