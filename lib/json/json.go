package json

import (
	"encoding/json"
	"io"
	"net/http"
)

type ErrorHandlerFunc func(err error)

var ErrorHandler ErrorHandlerFunc = func(err error) {}

type e struct {
	Message string `json:"message"`
}

type d struct {
	D interface{} `json:"data"`
}

func Decode(r io.Reader, dest interface{}) error {
	return json.NewDecoder(r).Decode(dest)
}

// TODO (RCH): Get code for errors
func Error(w http.ResponseWriter, code int, err error) {
	encode(w, code, map[string]interface{}{
		"error": e{Message: err.Error()},
	})
}

func encode(w http.ResponseWriter, code int, body interface{}) {
	w.WriteHeader(code)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")
	err := enc.Encode(body)
	if ErrorHandler != nil {
		ErrorHandler(err)
	}
}

func Created(w http.ResponseWriter, data interface{}) {
	encode(w, 201, d{data})
}

func OK(w http.ResponseWriter, data interface{}) {
	encode(w, 200, d{data})
}
