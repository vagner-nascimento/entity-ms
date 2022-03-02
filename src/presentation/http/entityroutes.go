package http

import (
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

}

func getEntityFromBody(b io.ReadCloser) (ent model.Entity, err error) {
	defer b.Close()

	var bys []byte
	if bys, err = io.ReadAll(b); err == nil {
		logger.Info("data received", string(bys))

		ent, err = model.NewEntityFromBytes(bys)
	}

	return ent, err
}
