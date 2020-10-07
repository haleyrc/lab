package db

import (
	"context"
)

type Genre struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

const createGenreQuery = `
INSERT INTO genres (
	id,
	name
) VALUES (
	$1,
	$2
);
`

func (s *Store) CreateGenre(ctx context.Context, g *Genre) error {
	if _, err := s.db.ExecContext(
		ctx,
		createGenreQuery,
		g.ID,
		g.Name,
	); err != nil {
		return err
	}
	return nil
}

const getGenreQuery = `
SELECT
	id,
	name
FROM
	genres
WHERE
	id = $1;
`

func (s *Store) GetGenre(ctx context.Context, id string) (*Genre, error) {
	var genre Genre
	if err := s.db.GetContext(
		ctx,
		&genre,
		getGenreQuery,
		id,
	); err != nil {
		return nil, err
	}
	return &genre, nil
}

const listGenresByIDsQuery = `
SELECT
	id,
	name
FROM
	genres
WHERE
	id IN (?);
`

func (s *Store) ListGenresByIDs(ctx context.Context, ids []string) ([]*Genre, error) {
	if len(ids) == 0 {
		return []*Genre{}, nil
	}
	q, args, err := bindIDs(s.db, listGenresByIDsQuery, ids)
	if err != nil {
		return nil, err
	}

	var genres []*Genre
	if err := s.db.SelectContext(ctx, &genres, q, args...); err != nil {
		return nil, err
	}

	return genres, nil
}

const listGenresQuery = `
SELECT
	id,
	name
FROM
	genres;
`

func (s *Store) ListGenres(ctx context.Context) ([]*Genre, error) {
	var genres []*Genre
	if err := s.db.SelectContext(ctx, &genres, listGenresQuery); err != nil {
		return nil, err
	}

	return genres, nil
}
