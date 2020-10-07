package token

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const DefaultExpiration = time.Hour

var (
	ErrInvalidToken         = errors.New("invalid token")
	ErrInvalidSigningMethod = errors.New("invalid signing method")
)

type Claims struct {
	Sub string
}

func (c Claims) Valid() error {
	if c.Sub == "" {
		return ErrInvalidToken
	}
	return nil
}

func New(key string) *ParserGenerator {
	return &ParserGenerator{
		exp: DefaultExpiration,
		key: []byte(key),
	}
}

type ParserGenerator struct {
	exp time.Duration
	key []byte
}

func (pg ParserGenerator) Generate(sub string) (string, error) {
	claims := jwt.MapClaims{
		"sub": sub,
		"exp": int64(time.Now().Add(pg.exp).Unix()),
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tok, err := t.SignedString(pg.key)
	if err != nil {
		return "", err
	}

	return tok, nil
}

func (pg ParserGenerator) Parse(tok string) (*Claims, error) {
	t, err := jwt.ParseWithClaims(tok, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSigningMethod
		}

		return pg.key, nil
	})

	claims, ok := t.Claims.(*Claims)
	if !ok || !t.Valid {
		return nil, &InvalidTokenError{err}
	}

	return claims, nil
}

type InvalidTokenError struct {
	Err error
}

func (ite InvalidTokenError) Error() string {
	return ite.Err.Error()
}
