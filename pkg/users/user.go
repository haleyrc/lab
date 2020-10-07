package users

import (
	"context"

	"github.com/pborman/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/haleyrc/cheevos-simple/api/internal/db"
)

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

/*
func (u User) Authenticate(password string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(u.PasswordHash),
		[]byte(password),
	)
}
*/

type Repository interface {
	CreateUser(ctx context.Context, u *db.User) error
	GetUserByEmail(ctx context.Context, email string) (*db.User, error)
	GetUserByID(ctx context.Context, id string) (*db.User, error)
}

type Service struct {
	repo Repository
}

type GetUserRequest struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type GetUserResponse struct {
	User *User `json:"user"`
}

func (s *Service) GetUser(ctx context.Context, req GetUserRequest) (*GetUserResponse, error) {
	var user *db.User
	var err error
	if req.ID != "" {
		user, err := s.repo.GetUserByID(ctx, req.ID)
	} else {
		user, err := s.repo.GetUserByEmail(ctx, req.Email)
	}
	if err != nil {
		return nil, err
	}
	return &GetUserResponse{User: newUserFromDB(user)}, nil
}

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserResponse struct {
	User *User `json:"user"`
}

func (s *Service) CreateUser(ctx context.Context, req CreateUserRequest) (*CreateUserResponse, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &db.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: string(hash),
	}
	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return &CreateUserResponse{User: newUserFromDB(user)}, nil
}

func newUserFromDB(dbu *db.User) *User {
	return &User{
		ID:    dbu.ID,
		Email: dbu.Email,
	}
}
