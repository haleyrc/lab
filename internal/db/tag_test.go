package db

import (
	"context"
	"testing"

	"github.com/haleyrc/cheevos-simple/api/lib/check"
	"github.com/pborman/uuid"
)

func TestTags(t *testing.T) {
	ctx := context.Background()
	check := check.New(t)

	id := uuid.New()
	err := store.CreateTag(ctx, &Tag{
		ID:          id,
		Name:        "Test Tag",
		Description: "This is a test tag.",
		Color:       "BADA55",
	})
	check.OK(err).Fatal()

	tag, err := store.GetTag(ctx, id)
	check.OK(err).Fatal()
	check.Equals(tag.ID, id)
	check.Equals(tag.Name, "Test Tag")
	check.Equals(tag.Description, "This is a test tag.")
	check.Equals(tag.Color, "BADA55")

	tags, err := store.ListTags(ctx, []string{id})
	check.OK(err).Fatal()
	check.Equals(len(tags), 1)
	check.Equals(tags[0].ID, id)
	check.Equals(tags[0].Name, "Test Tag")
	check.Equals(tags[0].Description, "This is a test tag.")
	check.Equals(tags[0].Color, "BADA55")
}
