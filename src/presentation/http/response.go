package http

import (
	"encoding/json"
	"entity/src/apperrors"
	"net/http"
)

func writeCreatedResponse(w http.ResponseWriter, data interface{}) {
	jsonRes, _ := json.Marshal(data)

	w.WriteHeader(http.StatusCreated)
	w.Write(jsonRes)
}

// If err.Type == ERR_TP_VALIDATION, then status is BadRequest, otherwise is InternalServerError
func writeErrorResponse(w http.ResponseWriter, err apperrors.Error) {
	sts := func() int {
		if err.Type != nil && *err.Type == apperrors.ERR_TP_VALIDATION {
			return http.StatusBadRequest
		} else {
			return http.StatusInternalServerError
		}
	}()

	errs := httpErrors{
		Errors: []apperrors.Error{err},
	}

	jsonErr, _ := json.Marshal(errs)

	w.WriteHeader(sts)
	w.Write(jsonErr)
}

func writeBadRequestResponse(w http.ResponseWriter, httpErr httpErrors) {
	jsonErr, _ := json.Marshal(httpErr)

	w.WriteHeader(http.StatusBadRequest)
	w.Write(jsonErr)
}
