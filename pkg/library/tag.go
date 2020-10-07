package library

import (
	"context"

	"github.com/pborman/uuid"

	"github.com/haleyrc/cheevos-simple/api/internal/db"
)

type CreateTagRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

type CreateTagResponse struct {
	Tag *Tag `json:"tag"`
}

func (s *Service) CreateTag(ctx context.Context, req CreateTagRequest) (*CreateTagResponse, error) {
	tag := &db.Tag{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		Color:       req.Color,
	}
	if err := s.repo.CreateTag(ctx, tag); err != nil {
		return nil, err
	}
	return &CreateTagResponse{Tag: newTagFromDB(tag)}, nil
}

type ListTagsRequest struct {
	IDs []string `json:"ids"`
}

type ListTagsResponse struct {
	Tags []*Tag `json:"tags"`
}

func (s *Service) ListTags(ctx context.Context, req ListTagsRequest) (*ListTagsResponse, error) {
	if len(req.IDs) > 0 {
		tags, err := s.repo.ListTagsByIDs(ctx, req.IDs)
		if err != nil {
			return nil, err
		}
		return &ListTagsResponse{Tags: newTagsFromDB(tags)}, nil
	}

	tags, err := s.repo.ListTags(ctx)
	if err != nil {
		return nil, err
	}
	return &ListTagsResponse{Tags: newTagsFromDB(tags)}, nil
}

func newTagsFromDB(dbts []*db.Tag) []*Tag {
	tags := make([]*Tag, len(dbts))
	for i := range dbts {
		tags[i] = newTagFromDB(dbts[i])
	}
	return tags
}

func newTagFromDB(dbt *db.Tag) *Tag {
	return &Tag{
		ID:          dbt.ID,
		Name:        dbt.Name,
		Description: dbt.Description,
		Color:       dbt.Color,
	}
}
