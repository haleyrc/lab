package library

import (
	"context"
	"testing"

	"github.com/haleyrc/cheevos-simple/api/internal/db"
	"github.com/haleyrc/cheevos-simple/api/lib/check"
	"github.com/pborman/uuid"
)

type mockBookCreator struct{}

func (mbc mockBookCreator) CreateBook(_ context.Context, _ *db.Book) error {
	return nil
}

type mockBookTagCreator struct {
	CreateBookTagsCalled bool
}

func (mbtc *mockBookTagCreator) CreateBookTags(_ context.Context, _ []*db.BookTag) error {
	mbtc.CreateBookTagsCalled = true
	return nil
}

func TestNewBook_CreatesAuthorIfProvided(t *testing.T) {
	ctx := context.Background()
	check := check.New(t)
	ac := &mockAuthorCreator{}
	bc := mockBookCreator{}
	btc := &mockBookTagCreator{}

	book, err := newBook(ctx, ac, bc, btc, NewBookParams{
		Title: "Dune",
		Author: NewAuthorParams{
			FirstName: "Frank",
			LastName:  "Herbert",
		},
		TagIDs: []string{"tag1", "tag2"},
	})
	check.OK(err).Fatal()
	check.NotEmpty(book.ID)
	check.Equals(book.Title, "Dune")
	check.NotEmpty(book.AuthorID)
	check.True(ac.CreateAuthorCalled)
	check.True(btc.CreateBookTagsCalled)
}

func TestNewBook_DoesntCreateAuthorIfIDProvided(t *testing.T) {
	ctx := context.Background()
	check := check.New(t)
	ac := &mockAuthorCreator{}
	bc := mockBookCreator{}
	btc := &mockBookTagCreator{}

	authorID := uuid.New()
	book, err := newBook(ctx, ac, bc, btc, NewBookParams{
		Title:    "Dune",
		AuthorID: authorID,
	})
	check.OK(err).Fatal()
	check.NotEmpty(book.ID)
	check.Equals(book.Title, "Dune")
	check.Equals(book.AuthorID, authorID)
	check.False(ac.CreateAuthorCalled)
}
