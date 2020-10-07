package db

import "github.com/jmoiron/sqlx"

type Config struct {
	DB *sqlx.DB
}

func New(cfg Config) *Store {
	return &Store{db: cfg.DB}
}

type Store struct {
	db *sqlx.DB
}
