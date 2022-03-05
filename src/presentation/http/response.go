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

func writeSuccessResponse(w http.ResponseWriter, data interface{}) {
	jsonRes, _ := json.Marshal(data)

	w.WriteHeader(http.StatusOK)
	w.Write(jsonRes)
}

// If err.Type == ERR_TP_VALIDATION or ERR_TP_DATA, the status is BadRequest with details
// If err.Type == ERR_TP_NOT_FOUND, the status is NoContent without details
// Else the status is InternalServerError with details
func writeErrorResponse(w http.ResponseWriter, err apperrors.Error) {
	sts, write := func() (int, bool) {
		switch *err.Type {
		case apperrors.ERR_TP_VALIDATION, apperrors.ERR_TP_DATA:
			return http.StatusBadRequest, true
		case apperrors.ERR_TP_NOT_FOUND:
			return http.StatusNoContent, false
		default:
			return http.StatusInternalServerError, true
		}
	}()

	w.WriteHeader(sts)

	if write {
		errs := httpErrors{Errors: []apperrors.Error{err}}
		jsonErr, _ := json.Marshal(errs)

		w.Write(jsonErr)
	}
}

func writeBadRequestResponse(w http.ResponseWriter, httpErr httpErrors) {
	jsonErr, _ := json.Marshal(httpErr)

	w.WriteHeader(http.StatusBadRequest)
	w.Write(jsonErr)
}
