package library

import (
	"context"

	"github.com/haleyrc/cheevos-simple/api/internal/db"
)

type Repository interface {
	CreateBook(ctx context.Context, b *db.Book) error
	GetBook(ctx context.Context, id string) (*db.Book, error)
	ListBooks(ctx context.Context) ([]*db.Book, error)
	MarkBookRead(ctx context.Context, userID, bookID string) error
	MarkBookUnread(ctx context.Context, userID, bookID string) error

	// TODO (RCH): Convert this to add tag to book and replace the private method
	CreateBookTags(ctx context.Context, bookTags []*db.BookTag) error

	CreateAuthor(ctx context.Context, a *db.Author) error
	GetAuthor(ctx context.Context, id string) (*db.Author, error)
	ListAuthors(ctx context.Context) ([]*db.Author, error)
	ListAuthorsByIDs(ctx context.Context, ids []string) ([]*db.Author, error)

	CreateTag(ctx context.Context, t *db.Tag) error
	ListTags(ctx context.Context) ([]*db.Tag, error)
	ListTagsByIDs(ctx context.Context, ids []string) ([]*db.Tag, error)
}

type Service struct {
	repo Repository
}
