package db

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Tag struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

const createTagQuery = `
INSERT INTO tags (
	id,
	name,
	description,
	color
) VALUES (
	$1,
	$2,
	$3,
	$4
);
`

func (s *Store) CreateTag(ctx context.Context, t *Tag) error {
	if _, err := s.db.ExecContext(
		ctx,
		createTagQuery,
		t.ID,
		t.Name,
		t.Description,
		t.Color,
	); err != nil {
		return err
	}
	return nil
}

const getTagQuery = `
SELECT
	id,
	name,
	description,
	color
FROM
	tags
WHERE
	id = $1;
`

func (s *Store) GetTag(ctx context.Context, id string) (*Tag, error) {
	var tag Tag
	if err := s.db.GetContext(
		ctx,
		&tag,
		getTagQuery,
		id,
	); err != nil {
		return nil, err
	}
	return &tag, nil
}

const listTagsByIDsQuery = `
SELECT
	id,
	name,
	description,
	color
FROM
	tags
WHERE
	id in (?);
`

func (s *Store) ListTagsByIDs(ctx context.Context, ids []string) ([]*Tag, error) {
	if len(ids) == 0 {
		return []*Tag{}, nil
	}
	q, args, err := bindIDs(s.db, listTagsByIDsQuery, ids)
	if err != nil {
		return nil, err
	}

	tags := []*Tag{}
	if err := s.db.SelectContext(ctx, &tags, q, args...); err != nil {
		return nil, err
	}
	return tags, nil
}

const listTagsQuery = `
SELECT
	id,
	name,
	description,
	color
FROM
	tags;
`

func (s *Store) ListTags(ctx context.Context) ([]*Tag, error) {
	tags := []*Tag{}
	if err := s.db.SelectContext(ctx, &tags, listTagsQuery); err != nil {
		return nil, err
	}
	return tags, nil
}

func bindIDs(db *sqlx.DB, query string, ids []string) (string, []interface{}, error) {
	query, args, err := sqlx.In(query, ids)
	if err != nil {
		return "", nil, err
	}

	query = db.Rebind(query)

	return query, args, err
}
