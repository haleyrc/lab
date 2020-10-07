package db

import (
	"context"
	"testing"

	"github.com/haleyrc/cheevos-simple/api/lib/check"
	"github.com/jmoiron/sqlx"
	"github.com/pborman/uuid"
)

func TestCreateBookTag(t *testing.T) {
	ctx := context.Background()
	check := check.New(t)

	bookID := uuid.New()
	err := store.CreateBook(ctx, &Book{ID: bookID})
	check.OK(err).Fatal()

	tagID := uuid.New()
	err = store.CreateTag(ctx, &Tag{ID: tagID})
	check.OK(err).Fatal()

	err = store.CreateBookTag(ctx, &BookTag{
		BookID: bookID,
		TagID:  tagID,
	})
	check.OK(err).Fatal()
	check.Equals(countBookTags(ctx, store.db, bookID), 1)
}

func countBookTags(ctx context.Context, db *sqlx.DB, bookID string) int {
	var count int
	err := db.GetContext(ctx, &count, `SELECT COUNT(*) FROM book_tags WHERE book_id = $1`, bookID)
	if err != nil {
		return -1
	}
	return count
}
