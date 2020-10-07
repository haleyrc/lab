package library

import (
	"context"

	"github.com/pborman/uuid"

	"github.com/haleyrc/cheevos-simple/api/internal/db"
)

type CreateAuthorRequest struct {
	FirstName  string `json:"firstName"`
	MiddleName string `json:"middleName"`
	LastName   string `json:"lastName"`
}

func (req CreateAuthorRequest) isEmpty() bool {
	return req.FirstName == "" && req.MiddleName == "" && req.LastName == ""
}

type CreateAuthorResponse struct {
	Author *Author `json:"author"`
}

func (s *Service) CreateAuthor(ctx context.Context, req CreateAuthorRequest) (*CreateAuthorResponse, error) {
	author := &db.Author{
		ID:         uuid.New(),
		FirstName:  req.FirstName,
		MiddleName: req.MiddleName,
		LastName:   req.LastName,
	}
	if err := s.repo.CreateAuthor(ctx, author); err != nil {
		return nil, err
	}

	return &CreateAuthorResponse{Author: newAuthorFromDB(author)}, nil
}

type ListAuthorsRequest struct{}

type ListAuthorsResponse struct {
	Authors []*Author `json:"authors"`
}

func (s *Service) ListAuthors(ctx context.Context, req ListAuthorsRequest) (*ListAuthorsResponse, error) {
	authors, err := s.repo.ListAuthors(ctx)
	if err != nil {
		return nil, err
	}
	return &ListAuthorsResponse{Authors: newAuthorsFromDB(authors)}, nil
}

func newAuthorsFromDB(dbas []*db.Author) []*Author {
	authors := make([]*Author, len(dbas))
	for i := range dbas {
		authors[i] = newAuthorFromDB(dbas[i])
	}
	return authors
}

func newAuthorFromDB(dba *db.Author) *Author {
	return &Author{
		ID:         dba.ID,
		FirstName:  dba.FirstName,
		MiddleName: dba.MiddleName,
		LastName:   dba.LastName,
	}
}
