package http

import (
	"entity/src/apperrors"
	"fmt"
)

type httpErrors struct {
	Errors []apperrors.Error `json:"errors"`
}

func (h *httpErrors) Error() (msg string) {
	for _, err := range h.Errors {
		msg += fmt.Sprintf("|%s", err.Error())
	}

	return msg
}
