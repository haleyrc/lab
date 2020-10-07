package web

import "net/http"

type route struct {
	Path    string
	Method  string
	Handler http.HandlerFunc
	Private bool
}
