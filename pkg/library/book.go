package library

import (
	"context"
	"database/sql"

	"github.com/pborman/uuid"

	"github.com/haleyrc/cheevos-simple/api/internal/db"
)

type CreateBookRequest struct {
	Title           string   `json:"title"`
	PublicationYear int      `json:"pub_year"`
	Genre           Genre    `json:"genre"`
	AuthorID        string   `json:"author_id"`
	TagIDs          []string `json:"tag_ids"`
}

type CreateBookResponse struct {
	Book *Book `json:"book"`
}

func (s *Service) CreateBook(ctx context.Context, req CreateBookRequest) (*CreateBookResponse, error) {
	book := &db.Book{
		ID:              uuid.New(),
		Title:           req.Title,
		PublicationYear: req.PublicationYear,
		Genre:           string(req.Genre),
		AuthorID: sql.NullString{
			String: req.AuthorID,
			Valid:  req.AuthorID != "",
		},
	}
	if err := s.repo.CreateBook(ctx, book); err != nil {
		return nil, err
	}

	if err := s.addTagsToBook(ctx, book.ID, req.TagIDs); err != nil {
		return nil, err
	}

	return &CreateBookResponse{Book: newBookFromDB(book)}, nil
}

type GetBookRequest struct {
	ID string `json:"id"`
}

type GetBookResponse struct {
	Book *Book `json:"book"`
}

func (s *Service) GetBook(ctx context.Context, req GetBookRequest) (*GetBookResponse, error) {
	book, err := s.repo.GetBook(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return &GetBookResponse{Book: newBookFromDB(book)}, nil
}

type ListBooksRequest struct{}

type ListBooksResponse struct {
	Books []*Book `json:"books"`
}

func (s *Service) ListBooks(ctx context.Context, req ListBooksRequest) (*ListBooksResponse, error) {
	books, err := s.repo.ListBooks(ctx)
	if err != nil {
		return nil, err
	}
	return &ListBooksResponse{Books: newBooksFromDB(books)}, nil
}

type MarkBookReadRequest struct {
	UserID string `json:"user_id"`
	BookID string `json:"book_id"`
}

type MarkBookReadResponse struct{}

func (s *Service) MarkBookRead(ctx context.Context, req MarkBookReadRequest) (*MarkBookReadResponse, error) {
	return &MarkBookReadResponse{}, s.repo.MarkBookRead(ctx, req.UserID, req.BookID)
}

type MarkBookUnreadRequest struct {
	UserID string `json:"user_id"`
	BookID string `json:"book_id"`
}

type MarkBookUnreadResponse struct{}

func (s *Service) MarkBookUnread(ctx context.Context, req MarkBookUnreadRequest) (*MarkBookUnreadResponse, error) {
	return &MarkBookUnreadResponse{}, s.repo.MarkBookUnread(ctx, req.UserID, req.BookID)
}

func (s *Service) addTagsToBook(ctx context.Context, bookID string, tagIDs []string) error {
	bts := make([]*db.BookTag, 0, len(tagIDs))
	for _, tagID := range tagIDs {
		bts = append(bts, &db.BookTag{
			BookID: bookID,
			TagID:  tagID,
		})
	}
	if err := s.repo.CreateBookTags(ctx, bts); err != nil {
		return err
	}
	return nil
}

// TODO (RCH): Handle tags somehow
func newBookFromDB(dbb *db.Book) *Book {
	return &Book{
		ID:              dbb.ID,
		Title:           dbb.Title,
		PublicationYear: dbb.PublicationYear,
		Genre:           dbb.Genre,
		AuthorID:        dbb.AuthorID.String,
		TagIDs:          dbb.TagIDs,
	}
}

func newBooksFromDB(dbbs []*db.Book) []*Book {
	books := make([]*Book, 0, len(dbbs))
	for _, dbb := range dbbs {
		books = append(books, newBookFromDB(dbb))
	}
	return books
}
