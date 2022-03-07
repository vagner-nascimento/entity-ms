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

// TODO list
/*
 * - DELETE
 * - PATCH
 */
func getEntityRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Post("/", postEntity)
	router.Get("/{id}", getEntity)
	router.Put("/{id}", putEntity)

	return router
}

func postEntity(w netHttp.ResponseWriter, r *netHttp.Request) {
	if ent := getValidatedEntity(r.Body, w); ent != nil {
		if err := app.NewEnityAdapter().Save(ent); err == nil {
			writeCreatedResponse(w, ent)
		} else {
			writeErrorResponse(w, *err)
		}
	}
}

func getEntity(w netHttp.ResponseWriter, r *netHttp.Request) {
	if id := getIdFromPath(r.URL.Path, w); id != "" {
		if ent, err := app.NewEnityAdapter().Get(id); err == nil {
			writeSuccessResponse(w, ent)
		} else {
			writeErrorResponse(w, *err)
		}
	}
}

func putEntity(w netHttp.ResponseWriter, r *netHttp.Request) {
	if id := getIdFromPath(r.URL.Path, w); id != "" {
		if ent := getValidatedEntity(r.Body, w); ent != nil {
			if newEnt, err := app.NewEnityAdapter().Update(id, *ent); err == nil {
				writeSuccessResponse(w, newEnt)
			} else {
				writeErrorResponse(w, *err)
			}
		}
	}
}

func deleteEntity(w netHttp.ResponseWriter, r *netHttp.Request) {
	if id := getIdFromPath(r.URL.Path, w); id != "" {

	}
}

/*
 * Auxiliar functions
 */
// Get entity from a reader, that is the request body, and validates the received data.
// If is valid, return an Entity struct filled with data. If invalid, write a bad request response with details.
func getValidatedEntity(reader io.ReadCloser, w netHttp.ResponseWriter) (res *model.Entity) {
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

	if len(resErr.Errors) > 0 {
		writeBadRequestResponse(w, resErr)
	}

	return res
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

// Get id from path. If id was not found, writes a bad request response
func getIdFromPath(path string, w netHttp.ResponseWriter) (id string) {
	if id = getPathParam(path, 1); id == "" {
		fild := "path '/{id}'"
		writeErrorResponse(w, apperrors.NewValidationError("id must be informed", &fild, nil))
	}

	return id
}
