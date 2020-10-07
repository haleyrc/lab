package sdk

import (
	"context"
	"strings"
)

type Genre struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewGenreService(c *Client) *GenreService {
	if c == nil {
		c = NewClient(Config{})
	}
	return &GenreService{c: c}
}

type GenreService struct {
	c *Client
}

type CreateGenreRequest struct {
	Name string `json:"name"`
}

func (cgr *CreateGenreRequest) normalize() {
	cgr.Name = strings.TrimSpace(cgr.Name)
}

func (cgr *CreateGenreRequest) validate() error {
	if cgr.Name == "" {
		return ValidationError{Message: "You must specify a name."}
	}
	return nil
}

type CreateGenreResponse struct {
	Genre Genre `json:"genre"`
}

func (gs *GenreService) Create(ctx context.Context, req CreateGenreRequest) (*CreateGenreResponse, error) {
	req.normalize()
	if err := req.validate(); err != nil {
		return nil, err
	}

	res, err := gs.c.post(ctx, `/genres`, req)
	if err != nil {
		return nil, err
	}

	var resp CreateGenreResponse
	if err := decode(res, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

type ListGenresResponse struct {
	Genres []Genre `json:"genres"`
}
