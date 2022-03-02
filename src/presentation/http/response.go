package http

import (
	"encoding/json"
	"net/http"
)

func writeBadRequestResponse(w http.ResponseWriter, httpErr httpErrors) {
	jsonErr, _ := json.Marshal(httpErr)

	w.WriteHeader(http.StatusBadRequest)
	w.Write(jsonErr)
}
