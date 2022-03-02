package start

import "entity/src/presentation/http"

func StartApplication() <-chan error {
	return http.StartHttpServer()
}
