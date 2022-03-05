package http

import (
	"entity/src/app"
	"entity/src/apperrors"
	"entity/src/infra/logger"
	"entity/src/model"
	"io"
	netHttp "net/http"

	"github.com/go-chi/chi"
)

/*
 * TODO list
 * - PUT
 * - DELETE
 * - PATCH
 */
func getEntityRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Post("/", postEntity)
	router.Get("/{id}", getEntity)

	return router
}

func postEntity(w netHttp.ResponseWriter, r *netHttp.Request) {
	if ent, err := getValidatedEntity(r.Body); len(err.Errors) == 0 {
		if err := app.NewEnityAdapter().Save(ent); err == nil {
			writeCreatedResponse(w, ent)
		} else {
			writeErrorResponse(w, *err)
		}
	} else {
		writeBadRequestResponse(w, err)
	}
}

func getEntity(w netHttp.ResponseWriter, r *netHttp.Request) {
	if id := getPathParam(r.URL.Path, 1); id != "" {
		if ent, err := app.NewEnityAdapter().Get(id); err == nil {
			logger.Info("ent", ent)
			writeSuccessResponse(w, ent)
		} else {
			writeErrorResponse(w, *err)
		}
	} else {
		fild := "path '/{id}'"
		writeErrorResponse(w, apperrors.NewValidationError("id must be informed", &fild, nil))
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
		resErr.Errors = append(resErr.Errors, apperrors.NewValidationError(err.Error(), nil, nil))
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
