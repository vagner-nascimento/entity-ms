package http

import "strings"

func getPathParam(path string, skip int) string {
	params := strings.Split(path, "/")
	return params[len(params)-skip]
}
