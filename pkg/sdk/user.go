package sdk

import (
	"context"
	"encoding/json"
	"strings"
)

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

func NewUserService(c *Client) *UserService {
	if c == nil {
		c = NewClient(Config{})
	}
	return &UserService{c: c}
}

type UserService struct {
	c *Client
}

type SignUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (sur *SignUpRequest) normalize() {
	sur.Email = strings.ToLower(strings.TrimSpace(sur.Email))
	sur.Password = strings.TrimSpace(sur.Password)
}

func (sur SignUpRequest) validate() error {
	if sur.Email == "" {
		return ValidationError{Message: "You must specify an email address."}
	}
	if sur.Password == "" {
		return ValidationError{Message: "You must specify a password."}
	}
	if err := validatePassword(sur.Password); err != nil {
		return err
	}
	return nil
}

type SignUpResponse struct {
	User        User   `json:"user"`
	AccessToken string `json:"access_token"`
}

func (us *UserService) SignUp(ctx context.Context, req SignUpRequest) (*SignUpResponse, error) {
	req.normalize()
	if err := req.validate(); err != nil {
		return nil, err
	}

	res, err := us.c.post(ctx, `/signup`, req)
	if err != nil {
		return nil, err
	}

	var resp SignUpResponse
	if err := decode(res, &resp); err != nil {
		return nil, err
	}
	us.c.bearerToken = resp.AccessToken

	return &resp, nil
}

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (sir *SignInRequest) normalize() {
	sir.Email = strings.ToLower(strings.TrimSpace(sir.Email))
	sir.Password = strings.TrimSpace(sir.Password)
}

func (sir SignInRequest) validate() error {
	if sir.Email == "" {
		return ValidationError{Message: "Email is required."}
	}
	if sir.Password == "" {
		return ValidationError{Message: "Password is required."}
	}
	return nil
}

type SignInResponse struct {
	User        User   `json:"user"`
	AccessToken string `json:"access_token"`
}

func (us *UserService) SignIn(ctx context.Context, req SignInRequest) (*SignInResponse, error) {
	req.normalize()
	if err := req.validate(); err != nil {
		return nil, err
	}

	res, err := us.c.post(ctx, `/signin`, req)
	if err != nil {
		return nil, err
	}

	var resp SignInResponse
	if err := decode(res, &resp); err != nil {
		return nil, err
	}
	us.c.bearerToken = resp.AccessToken

	return &resp, nil
}

type MeResponse struct {
	User User `json:"user"`
}

func decode(data json.RawMessage, dest interface{}) error {
	return json.Unmarshal(data, dest)
}

// TODO (RCH): Add some actual validation
func validatePassword(p string) error {
	return nil
}
