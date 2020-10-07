package library

import (
	"log"
	"net/http"

	"github.com/haleyrc/cheevos-simple/api/lib/json"
	"github.com/haleyrc/cheevos-simple/api/pkg/core"
)

type HTTPAdapter struct {
	svc  *Service
	resp core.HTTPResponder
}

func (a *HTTPAdapter) CreateAuthor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateAuthorRequest
	if err := json.Decode(r.Body, &req); err != nil {
		a.resp.RespondWithError(w, err)
		return
	}

	resp, err := a.svc.CreateAuthor(ctx, req)
	if err != nil {
		a.resp.RespondWithError(w, err)
		return
	}
	log.Println("created author:", resp.Author)

	a.resp.RespondOK(w, resp)
}

func (a *HTTPAdapter) ListAuthors(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	resp, err := a.svc.ListAuthors(ctx, ListAuthorsRequest{})
	if err != nil {
		a.resp.RespondWithError(w, err)
		return
	}

	a.resp.RespondOK(w, resp)
}

func (a *HTTPAdapter) CreateBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateBookRequest
	if err := json.Decode(r.Body, &req); err != nil {
		a.resp.RespondWithError(w, err)
		return
	}

	resp, err := a.svc.CreateBook(ctx, req)
	if err != nil {
		a.resp.RespondWithError(w, err)
		return
	}
	log.Println("created book:", resp.Book)

	a.resp.RespondOK(w, resp)
}

func (a *HTTPAdapter) ListBooks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	resp, err := a.svc.ListBooks(ctx, ListBooksRequest{})
	if err != nil {
		a.resp.RespondWithError(w, err)
		return
	}

	a.resp.RespondOK(w, resp)
}

func (a *HTTPAdapter) CreateTag(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateTagRequest
	if err := json.Decode(r.Body, &req); err != nil {
		a.resp.RespondWithError(w, err)
		return
	}

	resp, err := a.svc.CreateTag(ctx, req)
	if err != nil {
		a.resp.RespondWithError(w, err)
		return
	}
	log.Println("created tag:", resp.Tag)

	a.resp.RespondOK(w, resp)
}

func (a *HTTPAdapter) ListTags(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	resp, err := a.svc.ListTags(ctx, ListTagsRequest{})
	if err != nil {
		a.resp.RespondWithError(w, err)
		return
	}

	a.resp.RespondOK(w, resp)
}
