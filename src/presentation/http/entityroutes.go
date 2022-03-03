package http

import (
	"entity/src/apperrors"
	"entity/src/infra/logger"
	"entity/src/model"
	"io"
	netHttp "net/http"

	"github.com/go-chi/chi"
)

func getEntityRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Post("/", postEntity)

	return router
}

func postEntity(w netHttp.ResponseWriter, r *netHttp.Request) {
	if ent, err := getValidatedEntity(r.Body); len(err.Errors) == 0 {
		//TODO proccess data
		logger.Info("http.postEntity success", ent.String())
	} else {
		logger.Error("http.postEntity error", &err)
		writeBadRequestResponse(w, err)
	}
}

func getValidatedEntity(reader io.ReadCloser) (*model.Entity, httpErrors) {
	var res *model.Entity
	var resErr httpErrors

	if ent, err := getEntityFromBody(reader); err == nil {
		if isValid, errs := ent.Validate(); isValid {
			res = &ent
		} else {
			resErr.Errors = errs
		}
	} else {
		resErr.Errors = append(resErr.Errors, apperrors.NewConversionValidationError(err.Error()))
	}

	return res, resErr
}

func getEntityFromBody(reader io.ReadCloser) (ent model.Entity, err error) {
	defer reader.Close()

	var bys []byte
	if bys, err = io.ReadAll(reader); err == nil {
		logger.Info("data received", string(bys))

		ent, err = model.NewEntityFromBytes(bys)
	}

	return ent, err
}
