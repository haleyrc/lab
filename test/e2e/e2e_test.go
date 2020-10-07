package e2e

import (
	"context"
	"fmt"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gorilla/mux"
	"github.com/haleyrc/cheevos-simple/api/internal/db"
	"github.com/haleyrc/cheevos-simple/api/internal/migrate"
	"github.com/haleyrc/cheevos-simple/api/internal/testutil"
	"github.com/haleyrc/cheevos-simple/api/internal/web"
	"github.com/haleyrc/cheevos-simple/api/lib/check"
	"github.com/haleyrc/cheevos-simple/api/pkg/domain"
	"github.com/haleyrc/cheevos-simple/api/pkg/sdk"
)

var url string

func TestMain(m *testing.M) {
	pgdb, cleanup, err := testutil.ConnectToTestDatabase()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dir := filepath.Join("..", "..", "migrations")
	fmt.Println("dir:", dir)
	err = migrate.UpFrom(context.Background(), pgdb, dir)
	if err != nil {
		fmt.Println(err)
		cleanup()
		os.Exit(1)
	}

	store := db.New(db.Config{
		DB: pgdb,
	})
	dom := domain.New(domain.Config{
		Database: store,
	})
	h := web.NewServer(web.Config{
		Domain: dom,
		Router: mux.NewRouter().PathPrefix("/api").Subrouter(),
	})

	srv := httptest.NewServer(h)
	url = srv.URL + "/api"

	code := m.Run()

	srv.Close()
	cleanup()

	os.Exit(code)
}

func TestFlow(t *testing.T) {
	ctx := context.Background()
	check := check.New(t)

	client := sdk.NewClient(sdk.Config{
		Debug: true,
		URL:   url,
	})
	authorService := sdk.NewAuthorService(client)
	bookService := sdk.NewBookService(client)
	genreService := sdk.NewGenreService(client)
	tagService := sdk.NewTagService(client)
	userService := sdk.NewUserService(client)

	err := userService.SignUp(ctx, sdk.SignUpRequest{
		Email:    "ryan@haleylab.com",
		Password: "testtest",
	})
	check.OK(err).Fatal()

	signInResponse, err := userService.SignIn(ctx, sdk.SignInRequest{
		Email:    "ryan@haleylab.com",
		Password: "testtest",
	})
	check.OK(err).Fatal()
	check.NotEmpty(signInResponse.Token)
	check.NotEmpty(signInResponse.User.ID)
	check.Equals(signInResponse.User.Email, "ryan@haleylab.com")

	createGenreResponse, err := genreService.Create(ctx, sdk.CreateGenreRequest{
		Name: "Science Fiction",
	})
	check.OK(err).Then(func() {
		check.NotEmpty(createGenreResponse.Genre.ID)
		check.Equals(createGenreResponse.Genre.Name, "Science Fiction")
	})

	createAuthorResponse, err := authorService.Create(ctx, sdk.CreateAuthorRequest{
		FirstName:  "Frank",
		MiddleName: "Patrick",
		LastName:   "Herbert",
	})
	check.OK(err).Then(func() {
		check.NotEmpty(createAuthorResponse.Author.ID)
		check.Equals(createAuthorResponse.Author.FirstName, "Frank")
		check.Equals(createAuthorResponse.Author.MiddleName, "Patrick")
		check.Equals(createAuthorResponse.Author.LastName, "Herbert")
	})

	createTagResponse, err := tagService.Create(ctx, sdk.CreateTagRequest{
		Name:        "Favorite",
		Description: "The best of the best.",
		Color:       "FFD100",
	})
	check.OK(err).Then(func() {
		check.NotEmpty(createTagResponse.Tag.ID)
		check.Equals(createTagResponse.Tag.Name, "Favorite")
		check.Equals(createTagResponse.Tag.Description, "The best of the best.")
		check.Equals(createTagResponse.Tag.Color, "#FFD100")
	})

	createBookResponse, err := bookService.Create(ctx, sdk.CreateBookRequest{
		Title:           "Dune",
		PublicationYear: 1965,
		GenreID:         createGenreResponse.Genre.ID,
		AuthorID:        createAuthorResponse.Author.ID,
		TagIDs:          []string{createTagResponse.Tag.ID},
	})
	check.OK(err).Then(func() {
		check.NotEmpty(createBookResponse.Book.ID)
		check.Equals(createBookResponse.Book.Title, "Dune")
		check.Equals(createBookResponse.Book.PublicationYear, int64(1965))
		check.Equals(createBookResponse.Book.GenreID, createGenreResponse.Genre.ID)
		check.Equals(createBookResponse.Book.AuthorID, createAuthorResponse.Author.ID)
		check.Equals(len(createBookResponse.Book.TagIDs), 1).Then(func() {
			check.Equals(createBookResponse.Book.TagIDs[0], createTagResponse.Tag.ID)
		})
	})

	createBook2Response, err := bookService.Create(ctx, sdk.CreateBookRequest{
		Title:           "2001: A Space Odyssey",
		PublicationYear: 1968,
		GenreID:         createGenreResponse.Genre.ID,
		Author: sdk.CreateAuthorRequest{
			FirstName:  "Arthur",
			MiddleName: "C.",
			LastName:   "Clarke",
		},
	})
	check.OK(err).Then(func() {
		check.NotEmpty(createBook2Response.Book.ID)
		check.Equals(createBook2Response.Book.Title, "2001: A Space Odyssey")
		check.Equals(createBook2Response.Book.PublicationYear, int64(1968))
		check.Equals(createBook2Response.Book.GenreID, createGenreResponse.Genre.ID)
		check.NotEmpty(createBook2Response.Book.AuthorID)
		check.Equals(len(createBook2Response.Book.TagIDs), 0)
	})
}
