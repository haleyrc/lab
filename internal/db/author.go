package db

import (
	"context"
)

type Author struct {
	ID         string `db:"id"`
	FirstName  string `db:"first_name"`
	MiddleName string `db:"middle_name"`
	LastName   string `db:"last_name"`
}

const createAuthorQuery = `
INSERT INTO authors (
	id,
	first_name,
	middle_name,
	last_name
) VALUES (
	$1,
	$2,
	$3,
	$4
);
`

func (s *Store) CreateAuthor(ctx context.Context, a *Author) error {
	if _, err := s.db.ExecContext(
		ctx,
		createAuthorQuery,
		a.ID,
		a.FirstName,
		a.MiddleName,
		a.LastName,
	); err != nil {
		return err
	}
	return nil
}

const getAuthorQuery = `
SELECT
	id,
	first_name,
	middle_name,
	last_name
FROM
	authors
WHERE
	id = $1;
`

func (s *Store) GetAuthor(ctx context.Context, id string) (*Author, error) {
	var a Author
	if err := s.db.GetContext(
		ctx,
		&a,
		getAuthorQuery,
		id,
	); err != nil {
		return nil, err
	}
	return &a, nil
}

const listAuthorsByIDQuery = `
SELECT
	id,
	first_name,
	middle_name,
	last_name
FROM
	authors
WHERE
	id in (?);
`

func (s *Store) ListAuthorsByIDs(ctx context.Context, ids []string) ([]*Author, error) {
	if len(ids) == 0 {
		return []*Author{}, nil
	}
	q, args, err := bindIDs(s.db, listAuthorsByIDQuery, ids)
	if err != nil {
		return nil, err
	}

	var as []*Author
	if err := s.db.SelectContext(ctx, &as, q, args...); err != nil {
		return nil, err
	}

	return as, nil
}

const listAuthorsQuery = `
SELECT
	id,
	first_name,
	middle_name,
	last_name
FROM
	authors;
`

func (s *Store) ListAuthors(ctx context.Context) ([]*Author, error) {
	var as []*Author
	if err := s.db.SelectContext(ctx, &as, listAuthorsQuery); err != nil {
		return nil, err
	}

	return as, nil
}
