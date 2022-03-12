package http

import (
	"entity/src/app"
	"entity/src/apperrors"
	"entity/src/infra/logger"
	"entity/src/model"
	"fmt"
	"io"
	netHttp "net/http"

	"github.com/go-chi/chi"
)

// TODO list
/*
 * - Return a bad request response when coming unallowed fields
 */
func getEntityRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Post("/", postEntity)
	router.Get("/{id}", getEntity)
	router.Get("/", getEntities)
	router.Put("/{id}", putEntity)
	router.Delete("/{id}", deleteEntity)
	router.Patch("/{id}/name", patchEntityName)
	router.Patch("/{id}/weight", patchEntityWeight)
	router.Patch("/{id}/birthdate", patchEntityBirthDate)

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
	if id := getIdFromPath(r.URL.Path, 1, w); id != "" {
		if ent, err := app.NewEnityAdapter().Get(id); err == nil {
			writeSuccessResponse(w, ent)
		} else {
			writeErrorResponse(w, *err)
		}
	}
}

func getEntities(w netHttp.ResponseWriter, r *netHttp.Request) {
	if errs := validateEntityQuery(r.URL.Query()); errs == nil {
		res, err := app.NewEnityAdapter().Search(r.URL.Query())
		fmt.Println(res, err) // TODO handle res and error
	} else {
		writeBadRequestResponse(w, httpErrors{Errors: errs})
	}
}

func putEntity(w netHttp.ResponseWriter, r *netHttp.Request) {
	if id := getIdFromPath(r.URL.Path, 1, w); id != "" {
		if ent := getValidatedEntity(r.Body, w); ent != nil {
			if newEnt, err := app.NewEnityAdapter().Update(id, *ent); err == nil {
				writeSuccessResponse(w, newEnt)
			} else {
				writeErrorResponse(w, *err)
			}
		}
	}
}

func patchEntityName(w netHttp.ResponseWriter, r *netHttp.Request) {
	if id := getIdFromPath(r.URL.Path, 2, w); id != "" {
		if ent := getEntityData(r.Body, w); ent != nil {
			if valid, verr := ent.ValidateName(); valid {
				ent.NilAllButName()
				if res, aerr := app.NewEnityAdapter().Update(id, *ent); aerr == nil {
					writeSuccessResponse(w, res)
				} else {
					writeErrorResponse(w, *aerr)
				}
			} else {
				writeBadRequestResponse(w, httpErrors{Errors: verr})
			}
		}
	}
}

func patchEntityWeight(w netHttp.ResponseWriter, r *netHttp.Request) {
	if id := getIdFromPath(r.URL.Path, 2, w); id != "" {
		if ent := getEntityData(r.Body, w); ent != nil {
			if valid, verr := ent.ValidateWeigth(); valid {
				ent.NilAllButWeight()
				if res, aerr := app.NewEnityAdapter().Update(id, *ent); aerr == nil {
					writeSuccessResponse(w, res)
				} else {
					writeErrorResponse(w, *aerr)
				}
			} else {
				writeBadRequestResponse(w, httpErrors{Errors: verr})
			}
		}
	}
}

func patchEntityBirthDate(w netHttp.ResponseWriter, r *netHttp.Request) {
	if id := getIdFromPath(r.URL.Path, 2, w); id != "" {
		if ent := getEntityData(r.Body, w); ent != nil {
			if valid, verr := ent.ValidateBirthDate(); valid {
				ent.NilAllButBithDate()
				if res, aerr := app.NewEnityAdapter().Update(id, *ent); aerr == nil {
					writeSuccessResponse(w, res)
				} else {
					writeErrorResponse(w, *aerr)
				}
			} else {
				writeBadRequestResponse(w, httpErrors{Errors: verr})
			}
		}
	}
}

func deleteEntity(w netHttp.ResponseWriter, r *netHttp.Request) {
	if id := getIdFromPath(r.URL.Path, 1, w); id != "" {
		if ent, err := app.NewEnityAdapter().Delete(id); err == nil {
			writeSuccessResponse(w, ent)
		} else {
			writeErrorResponse(w, *err)
		}
	}
}

/*
 * Auxiliar functions
 */
// Get entity from a reader, that is the request body, and validates the received data.
// If is valid, return an Entity struct filled with data. If invalid, write a bad request response with details.
func getValidatedEntity(reader io.ReadCloser, w netHttp.ResponseWriter) (res *model.Entity) {
	if res = getEntityData(reader, w); res != nil {
		if isValid, errs := res.Validate(); !isValid {
			writeBadRequestResponse(w, httpErrors{Errors: errs})
		}
	}

	return
}

// Get entity from a reader, that is the request body. If parse fails, wirte a bad request response with parse fail message
func getEntityData(reader io.ReadCloser, w netHttp.ResponseWriter) (res *model.Entity) {
	if ent, err := getEntityFromBody(reader); err == nil {
		res = &ent
	} else {
		errs := []apperrors.Error{apperrors.NewValidationError(err.Error(), nil, nil)}
		writeBadRequestResponse(w, httpErrors{Errors: errs})
	}

	return
}

func getEntityFromBody(reader io.ReadCloser) (ent model.Entity, err error) {
	defer reader.Close()

	var bys []byte
	if bys, err = io.ReadAll(reader); err == nil {
		logger.Info("data received", string(bys))

		ent, err = model.NewEntityFromBytes(bys)
	}

	return
}

// Get id from path. If id was not found, writes a bad request response
func getIdFromPath(path string, skip int, w netHttp.ResponseWriter) (id string) {
	if id = getPathParam(path, skip); id == "" {
		fild := "path '/{id}'"
		writeErrorResponse(w, apperrors.NewValidationError("id must be informed", &fild, nil))
	}

	return
}
