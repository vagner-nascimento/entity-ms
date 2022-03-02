package http

import (
	"entity/src/infra/logger"
	"errors"
	"fmt"
	netHttp "net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi"
)

func StartHttpServer() <-chan error {
	ch := make(chan error)
	router := chi.NewRouter()

	router.Use(getMiddlewareList()...)
	router.Route("/", func(r chi.Router) {
		r.Mount("/entity", getEntityRoutes())
	})

	logger.Info("avaliable routes:")
	walkThroughRoutes := func(method string, route string, handler netHttp.Handler, middleware ...func(netHttp.Handler) netHttp.Handler) error {
		logger.Info(fmt.Sprintf("> %s %s", method, route))
		return nil
	}

	if err := chi.Walk(router, walkThroughRoutes); err != nil {
		logger.Error("error on walkThroughRoutes", err)
		ch <- errors.New("error on try to start http server")
	} else {
		port := 80

		if envPort := os.Getenv("LISTEN_PORT"); len(envPort) > 0 {
			port, _ = strconv.Atoi(envPort)
		}

		logger.Info(fmt.Sprintf("http server connected on port: %d", port))

		ch <- netHttp.ListenAndServe(fmt.Sprintf(":%d", port), router)
	}

	return ch
}
