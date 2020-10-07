package sdk

import (
	"context"
	"strings"
)

type Book struct {
	ID              string   `json:"id"`
	Title           string   `json:"title"`
	PublicationYear int64    `json:"pub_year"`
	GenreID         string   `json:"genre_id"`
	AuthorID        string   `json:"author_id"`
	TagIDs          []string `json:"tag_ids"`
}

func NewBookService(c *Client) *BookService {
	if c == nil {
		c = NewClient(Config{})
	}
	return &BookService{c: c}
}

type BookService struct {
	c *Client
}

type CreateBookRequest struct {
	Title           string              `json:"title"`
	PublicationYear int64               `json:"pub_year"`
	GenreID         string              `json:"genre_id"`
	AuthorID        string              `json:"author_id"`
	Author          CreateAuthorRequest `json:"author"`
	TagIDs          []string            `json:"tag_ids"`
}

func (cbr *CreateBookRequest) normalize() {
	cbr.Title = strings.TrimSpace(cbr.Title)
	cbr.Author.FirstName = strings.TrimSpace(cbr.Author.FirstName)
	cbr.Author.MiddleName = strings.TrimSpace(cbr.Author.MiddleName)
	cbr.Author.LastName = strings.TrimSpace(cbr.Author.LastName)
	for i, tag := range cbr.TagIDs {
		if tag == "" {
			cbr.TagIDs = append(cbr.TagIDs[:i], cbr.TagIDs[i+1:]...)
		} else {
			cbr.TagIDs[i] = strings.TrimSpace(tag)
		}
	}
}

func (cbr CreateBookRequest) validate() error {
	if cbr.Title == "" {
		return ValidationError{
			Message: "You must specify a title.",
		}
	}
	return nil
}

type CreateBookResponse struct {
	Book Book `json:"book"`
}

func (bs *BookService) Create(ctx context.Context, req CreateBookRequest) (*CreateBookResponse, error) {
	req.normalize()
	if err := req.validate(); err != nil {
		return nil, err
	}

	res, err := bs.c.post(ctx, `/books`, req)
	if err != nil {
		return nil, err
	}

	var resp CreateBookResponse
	if err := decode(res, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

type ListBooksResponse struct {
	Authors []Author `json:"authors"`
	Books   []Book   `json:"books"`
	Genres  []Genre  `json:"genres"`
	Tags    []Tag    `json:"tags"`
}
