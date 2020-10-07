package db

import "context"

type User struct {
	ID           string `db:"id"`
	Email        string `db:"email"`
	PasswordHash string `db:"password_hash"`
}

const createUserQuery = `
INSERT INTO users (
	id,
	email,
	password_hash
) VALUES (
	$1,
	$2,
	$3
);
`

func (s *Store) CreateUser(ctx context.Context, u *User) error {
	if _, err := s.db.ExecContext(
		ctx,
		createUserQuery,
		u.ID,
		u.Email,
		u.PasswordHash,
	); err != nil {
		return err
	}
	return nil
}

const getUserByEmailQuery = `
SELECT
	id,
	email,
	password_hash
FROM
	users
WHERE
	email = $1;
`

func (s *Store) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var u User
	if err := s.db.GetContext(
		ctx,
		&u,
		getUserByEmailQuery,
		email,
	); err != nil {
		return nil, err
	}
	return &u, nil
}

const getUserByIDQuery = `
SELECT
	id,
	email,
	password_hash
FROM
	users
WHERE
	id = $1;
`

func (s *Store) GetUserByID(ctx context.Context, id string) (*User, error) {
	var u User
	if err := s.db.GetContext(
		ctx,
		&u,
		getUserByIDQuery,
		id,
	); err != nil {
		return nil, err
	}
	return &u, nil
}
