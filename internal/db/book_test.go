package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/haleyrc/cheevos-simple/api/lib/check"
	"github.com/jmoiron/sqlx"
	"github.com/pborman/uuid"
)

func TestCreateBook(t *testing.T) {
	ctx := context.Background()
	check := check.New(t)

	genreID := uuid.New()
	err := store.CreateGenre(ctx, &Genre{ID: genreID})
	check.OK(err).Fatal()

	authorID := uuid.New()
	err = store.CreateAuthor(ctx, &Author{ID: authorID})
	check.OK(err).Fatal()

	id := uuid.New()
	err = store.CreateBook(ctx, &Book{
		ID:              id,
		Title:           "Test Title",
		PublicationYear: 1965,
		GenreID: sql.NullString{
			String: genreID,
			Valid:  true,
		},
		AuthorID: sql.NullString{
			String: authorID,
			Valid:  true,
		},
	})
	check.OK(err).Fatal()

	book, err := store.GetBook(ctx, id)
	check.OK(err).Fatal()
	check.Equals(book.ID, id)
	check.Equals(book.Title, "Test Title")
	check.Equals(book.PublicationYear, int64(1965))
	check.True(book.GenreID.Valid)
	check.Equals(book.GenreID.String, genreID)
	check.True(book.AuthorID.Valid)
	check.Equals(book.AuthorID.String, authorID)

	books, err := store.ListBooks(ctx)
	check.OK(err).Fatal()
	check.Equals(len(books), 1).Fatal()
	check.Equals(books[0].ID, id)
	check.Equals(books[0].Title, "Test Title")
	check.Equals(books[0].PublicationYear, int64(1965))
	check.True(books[0].GenreID.Valid)
	check.Equals(books[0].GenreID.String, genreID)
	check.True(books[0].AuthorID.Valid)
	check.Equals(books[0].AuthorID.String, authorID)
}

func TestMarkBookRead(t *testing.T) {
	ctx := context.Background()
	check := check.New(t)

	userID := uuid.New()
	err := store.CreateUser(ctx, &User{ID: userID})
	check.OK(err).Fatal()

	bookID := uuid.New()
	err = store.CreateBook(ctx, &Book{ID: bookID})
	check.OK(err).Fatal()

	err = store.MarkBookRead(ctx, userID, bookID)
	check.OK(err).Fatal()
	check.Equals(countReadBooks(ctx, store.db, userID), 1)

	err = store.MarkBookUnread(ctx, userID, bookID)
	check.OK(err).Fatal()
	check.Equals(countReadBooks(ctx, store.db, userID), 0)
}

func countReadBooks(ctx context.Context, db *sqlx.DB, userID string) int {
	var count int
	err := db.GetContext(ctx, &count, `SELECT COUNT(*) FROM user_books WHERE user_id = $1`, userID)
	if err != nil {
		return -1
	}
	return count
}
