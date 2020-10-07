package db

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type Book struct {
	ID              string         `db:"id"`
	Title           string         `db:"title"`
	PublicationYear int            `db:"pub_year"`
	Genre           string         `db:"genre"`
	AuthorID        sql.NullString `db:"author_id"`
	TagIDs          pq.StringArray `db:"tag_ids"`
}

const createBookQuery = `
INSERT INTO books (
	id,
	title,
	pub_year,
	genre,
	author_id
) VALUES (
	$1,
	$2,
	$3,
	$4,
	$5
);
`

func (s *Store) CreateBook(ctx context.Context, b *Book) error {
	if _, err := s.db.ExecContext(
		ctx,
		createBookQuery,
		b.ID,
		b.Title,
		b.PublicationYear,
		b.Genre,
		b.AuthorID,
	); err != nil {
		return err
	}
	return nil
}

const getBookQuery = `
SELECT
	id,
	title,
	pub_year,
	genre,
	author_id
FROM
	books
WHERE
	id = $1;
`

func (s *Store) GetBook(ctx context.Context, id string) (*Book, error) {
	var book Book
	if err := s.db.GetContext(
		ctx,
		&book,
		getBookQuery,
		id,
	); err != nil {
		return nil, err
	}
	return &book, nil
}

const listBooksQuery = `
SELECT
	id,
	title,
	pub_year,
	genre,
	author_id,
	ARRAY(
		SELECT
			tag_id
		FROM
			book_tags bt
		WHERE
			bt.book_id = books.id
	) as tag_ids
FROM
	books;
`

func (s *Store) ListBooks(ctx context.Context) ([]*Book, error) {
	books := []*Book{}
	if err := s.db.SelectContext(
		ctx,
		&books,
		listBooksQuery,
	); err != nil {
		return nil, err
	}
	return books, nil
}

const readBookQuery = `
INSERT INTO user_books (
	user_id,
	book_id
) VALUES (
	$1,
	$2
);
`

func (s *Store) MarkBookRead(ctx context.Context, userID, bookID string) error {
	if _, err := s.db.ExecContext(
		ctx,
		readBookQuery,
		userID,
		bookID,
	); err != nil {
		return err
	}
	return nil
}

const unreadBookQuery = `
DELETE FROM
	user_books
WHERE
	user_id = $1 AND
	book_id = $2;
`

func (s *Store) MarkBookUnread(ctx context.Context, userID, bookID string) error {
	if _, err := s.db.ExecContext(
		ctx,
		unreadBookQuery,
		userID,
		bookID,
	); err != nil {
		return err
	}
	return nil
}
