package sdk

import (
	"context"
	"strings"

	"github.com/haleyrc/cheevos-simple/api/pkg/domain"
)

type Tag struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

func NewTagService(c *Client) *TagService {
	if c == nil {
		c = NewClient(Config{})
	}
	return &TagService{c: c}
}

type TagService struct {
	c *Client
}

type CreateTagRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

func (ctr *CreateTagRequest) normalize() {
	ctr.Name = strings.TrimSpace(ctr.Name)
	ctr.Description = strings.TrimSpace(ctr.Description)
	ctr.Color = strings.TrimSpace(ctr.Color)
	ctr.Color = strings.TrimPrefix(ctr.Color, "#")
	if ctr.Color == "" {
		ctr.Color = domain.DefaultTagColor
	}
	ctr.Color = "#" + ctr.Color
}

func (ctr CreateTagRequest) validate() error {
	if ctr.Name == "" {
		return ValidationError{Message: "You must specify a tag name."}
	}
	if err := validColorCode(ctr.Color); err != nil {
		return err
	}
	return nil
}

type CreateTagResponse struct {
	Tag Tag `json:"tag"`
}

func (ts *TagService) Create(ctx context.Context, req CreateTagRequest) (*CreateTagResponse, error) {
	req.normalize()
	if err := req.validate(); err != nil {
		return nil, err
	}

	res, err := ts.c.post(ctx, `/tags`, req)
	if err != nil {
		return nil, err
	}

	var resp CreateTagResponse
	if err := decode(res, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

type ListTagsResponse struct {
	Tags []Tag `json:"tags"`
}

// TODO (RCH): Implement this
func validColorCode(c string) error {
	return nil
}
