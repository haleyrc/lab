package library

import (
	"context"
	"testing"

	"github.com/haleyrc/cheevos-simple/api/internal/db"
	"github.com/haleyrc/cheevos-simple/api/lib/check"
)

type mockAuthorCreator struct {
	CreateAuthorCalled bool
}

func (mac *mockAuthorCreator) CreateAuthor(_ context.Context, _ *db.Author) error {
	mac.CreateAuthorCalled = true
	return nil
}

func TestNewAuthor(t *testing.T) {
	ctx := context.Background()
	check := check.New(t)
	ac := &mockAuthorCreator{}

	author, err := newAuthor(ctx, ac, NewAuthorParams{
		FirstName: "Frank",
		LastName:  "Herbert",
	})

	check.OK(err)
	check.NotEmpty(author.ID)
	check.Equals(author.FirstName, "Frank")
	check.Equals(author.MiddleName, "")
	check.Equals(author.LastName, "Herbert")
}
