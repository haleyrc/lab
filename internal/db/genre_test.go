package db

import (
	"context"
	"testing"

	"github.com/haleyrc/cheevos-simple/api/lib/check"
	"github.com/pborman/uuid"
)

func TestCreateGenre(t *testing.T) {
	ctx := context.Background()
	check := check.New(t)

	id := uuid.New()
	err := store.CreateGenre(ctx, &Genre{
		ID:   id,
		Name: "Test Genre",
	})
	check.OK(err).Fatal()

	genre, err := store.GetGenre(ctx, id)
	check.OK(err).Fatal()
	check.Equals(genre.ID, id)
	check.Equals(genre.Name, "Test Genre")
}
