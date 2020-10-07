package db

import "context"

type BookTag struct {
	BookID string `db:"book_id"`
	TagID  string `db:"tag_id"`
}

func (s *Store) CreateBookTags(ctx context.Context, bookTags []*BookTag) error {
	for _, bt := range bookTags {
		if err := s.CreateBookTag(ctx, bt); err != nil {
			return err
		}
	}
	return nil
}

const createBookTagQuery = `
INSERT INTO book_tags (
	book_id,
	tag_id
) VALUES (
	$1,
	$2
);
`

func (s *Store) CreateBookTag(ctx context.Context, bookTag *BookTag) error {
	if _, err := s.db.ExecContext(
		ctx,
		createBookTagQuery,
		bookTag.BookID,
		bookTag.TagID,
	); err != nil {
		return err
	}
	return nil
}
