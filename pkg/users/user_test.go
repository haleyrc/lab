package users

import (
	"context"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/haleyrc/cheevos-simple/api/internal/db"
	"github.com/haleyrc/cheevos-simple/api/lib/check"
)

type mockUserCreator struct{}

func (muc *mockUserCreator) CreateUser(_ context.Context, _ *db.User) error {
	return nil
}

type mockUserGetter struct {
	GetUserByIDCalled    bool
	GetUserByEmailCalled bool
}

func (mug *mockUserGetter) GetUserByID(_ context.Context, _ string) (*db.User, error) {
	mug.GetUserByIDCalled = true
	return &db.User{}, nil
}

func (mug *mockUserGetter) GetUserByEmail(_ context.Context, _ string) (*db.User, error) {
	mug.GetUserByEmailCalled = true
	return &db.User{}, nil
}

func TestUser_Authenticate(t *testing.T) {
	check := check.New(t)

	hash, err := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)
	check.OK(err).Fatal()

	u := User{PasswordHash: string(hash)}
	check.OK(u.Authenticate("test"))
	check.Error(u.Authenticate("oops"))
}

func TestGetUserByID_FailsWithBlankID(t *testing.T) {
	ctx := context.Background()
	check := check.New(t)
	ug := &mockUserGetter{}

	user, err := getUserByID(ctx, ug, "")

	check.Error(err)
	check.Nil(user)
	check.False(ug.GetUserByIDCalled)
}

func TestGetUserByEmail_FailsWithBlankEmail(t *testing.T) {
	ctx := context.Background()
	check := check.New(t)
	ug := &mockUserGetter{}

	user, err := getUserByEmail(ctx, ug, "")

	check.Error(err)
	check.Nil(user)
	check.False(ug.GetUserByEmailCalled)
}

func TestNewUser(t *testing.T) {
	ctx := context.Background()
	check := check.New(t)
	uc := &mockUserCreator{}

	user, err := newUser(ctx, uc, NewUserParams{
		Email:    "test@example.com",
		Password: "testtest",
	})

	check.OK(err)
	check.NotEmpty(user.ID)
	check.Equals(user.Email, "test@example.com")
	check.NotEquals(user.PasswordHash, "testtest")
}
