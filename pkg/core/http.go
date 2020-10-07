package core

import "net/http"

type HTTPResponder interface {
	RespondOK(http.ResponseWriter, interface{})
	RespondWithError(http.ResponseWriter, error)
}
