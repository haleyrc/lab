package library

import (
	"context"
	"testing"

	"github.com/haleyrc/cheevos-simple/api/internal/db"
	"github.com/haleyrc/cheevos-simple/api/lib/check"
)

type mockTagCreator struct{}

func (mtc mockTagCreator) CreateTag(_ context.Context, _ *db.Tag) error {
	return nil
}

func TestNewTag(t *testing.T) {
	ctx := context.Background()
	check := check.New(t)
	tc := &mockTagCreator{}

	tag, err := newTag(ctx, tc, NewTagParams{
		Name:        "Test",
		Description: "A tag for testing.",
		Color:       "BADA55",
	})

	check.OK(err)
	check.NotEmpty(tag.ID)
	check.Equals(tag.Name, "Test")
	check.Equals(tag.Description, "A tag for testing.")
	check.Equals(tag.Color, "BADA55")
}
