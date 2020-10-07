package library

import "time"

type Author struct {
	ID         string `json:"id"`
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name"`
}

type Book struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	PublicationYear int       `json:"pub_year"`
	Genre           string    `json:"genre"`
	AuthorID        *string   `json:"author_id"`
	TagIDs          []string  `json:"tag_ids"`
	Readers         []*Reader `json:"readers"`
	Read            *bool     `json:"read,omitempty"`
	Rating          *int      `json:"rating,omitempty"`
	AverageRating   int       `json:"average_rating"`
	CoverURL        *string   `json:"cover_url"`
}

type Genre string

func (g Genre) String() string {
	return string(g)
}

func (g Genre) Valid() bool {
	for _, genre := range Genres {
		if g == genre {
			return true
		}
	}
	return false
}

const (
	Adventure      Genre = "Adventure"
	Fantasy        Genre = "Fantasy"
	Reference      Genre = "Reference"
	ScienceFiction Genre = "Science Fiction"
	Textbook       Genre = "Textbook"
)

var Genres = []Genre{
	Adventure,
	Fantasy,
	Reference,
	ScienceFiction,
	Textbook,
}

type Reader struct {
	ID        string `json:"id"`
	AvatarURL string `json:"avatar_url"`
}

type ReadThrough struct {
	ID       string     `json:"id"`
	BookID   string     `json:"book_id"`
	ReaderID string     `json:"reader_id"`
	Notes    string     `json:"notes"`
	Started  time.Time  `json:"started"`
	Finished *time.Time `json:"finished"`
}

const DefaultTagColor = "dddddd"

type Tag struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
}
