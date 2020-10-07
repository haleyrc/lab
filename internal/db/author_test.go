package db

import (
	"context"
	"testing"

	"github.com/pborman/uuid"

	"github.com/haleyrc/cheevos-simple/api/lib/check"
)

func TestCreateAuthor(t *testing.T) {
	ctx := context.Background()
	check := check.New(t)

	id := uuid.New()
	err := store.CreateAuthor(ctx, &Author{
		ID:         id,
		FirstName:  "Frank",
		MiddleName: "Patrick",
		LastName:   "Herbert",
	})
	check.OK(err).Fatal()

	author, err := store.GetAuthor(ctx, id)
	check.OK(err).Fatal()
	check.Equals(author.ID, id)
	check.Equals(author.FirstName, "Frank")
	check.Equals(author.MiddleName, "Patrick")
	check.Equals(author.LastName, "Herbert")
}
