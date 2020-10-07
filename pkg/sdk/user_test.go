package sdk

import (
	"context"
	"testing"

	"github.com/haleyrc/cheevos-simple/api/lib/check"
)

func TestSignUp(t *testing.T) {
	t.Skip("TODO: Setup a test server")

	ctx := context.Background()
	check := check.New(t)
	userService := NewUserService(nil)

	err := userService.SignUp(ctx, SignUpRequest{
		Email:    "paul.atreides@example.com",
		Password: "gomjabbar",
	})

	check.OK(err)

	resp, err := userService.SignIn(ctx, SignInRequest{
		Email:    "paul.atreides@example.com",
		Password: "gomjabbar",
	})

	check.OK(err)
	check.NotEmpty(resp.User.ID)
	check.Equals(resp.User.Email, "paul.atreides@example.com")
	check.NotEmpty(resp.Token)
}
