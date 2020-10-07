package sdk

import (
	"context"
	"strings"
)

type Author struct {
	ID         string `json:"id"`
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name"`
}

func NewAuthorService(c *Client) *AuthorService {
	if c == nil {
		c = NewClient(Config{})
	}
	return &AuthorService{c: c}
}

type AuthorService struct {
	c *Client
}

type CreateAuthorRequest struct {
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name"`
}

func (car CreateAuthorRequest) isEmpty() bool {
	return car.FirstName == "" && car.MiddleName == "" && car.LastName == ""
}

func (car *CreateAuthorRequest) normalize() {
	car.FirstName = strings.TrimSpace(car.FirstName)
	car.MiddleName = strings.TrimSpace(car.MiddleName)
	car.LastName = strings.TrimSpace(car.LastName)
}

func (car CreateAuthorRequest) validate() error {
	if car.isEmpty() {
		return ValidationError{
			Message: "You must specify at least one of the following: first name, middle name, or last name.",
		}
	}
	return nil
}

type CreateAuthorResponse struct {
	Author Author `json:"author"`
}

func (as *AuthorService) Create(ctx context.Context, req CreateAuthorRequest) (*CreateAuthorResponse, error) {
	req.normalize()
	if err := req.validate(); err != nil {
		return nil, err
	}

	res, err := as.c.post(ctx, `/authors`, req)
	if err != nil {
		return nil, err
	}

	var resp CreateAuthorResponse
	if err := decode(res, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

type ListAuthorsResponse struct {
	Authors []Author `json:"authors"`
}
