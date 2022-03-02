package http

import netHttp "net/http"

func responseHeaders(next netHttp.Handler) netHttp.Handler {
	return netHttp.HandlerFunc(func(w netHttp.ResponseWriter, r *netHttp.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func getMiddlewareList() (middleware []func(netHttp.Handler) netHttp.Handler) {
	middleware = append(middleware, responseHeaders)

	return middleware
}
