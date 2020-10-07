package db

import (
	"context"
	"testing"

	"github.com/haleyrc/cheevos-simple/api/lib/check"
	"github.com/pborman/uuid"
)

func TestUsers(t *testing.T) {
	ctx := context.Background()
	check := check.New(t)

	id := uuid.New()
	err := store.CreateUser(ctx, &User{
		ID:           id,
		Email:        "test@example.com",
		PasswordHash: "notarealhash",
	})
	check.OK(err).Fatal()

	{
		user, err := store.GetUserByID(ctx, id)
		check.OK(err).Fatal()
		check.Equals(user.ID, id)
		check.Equals(user.Email, "test@example.com")
		check.Equals(user.PasswordHash, "notarealhash")
	}

	{
		user, err := store.GetUserByEmail(ctx, "test@example.com")
		check.OK(err).Fatal()
		check.Equals(user.ID, id)
		check.Equals(user.Email, "test@example.com")
		check.Equals(user.PasswordHash, "notarealhash")
	}
}
