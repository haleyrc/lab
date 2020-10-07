package web

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/haleyrc/cheevos-simple/api/lib/json"
	"github.com/haleyrc/cheevos-simple/api/lib/token"
	"github.com/haleyrc/cheevos-simple/api/pkg/domain"
)

type Domain interface {
	NewUser(context.Context, domain.NewUserParams) (*domain.User, error)
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)

	NewGenre(context.Context, domain.NewGenreParams) (*domain.Genre, error)
	ListGenres(context.Context, domain.GenreQuery) ([]*domain.Genre, error)

	NewAuthor(context.Context, domain.NewAuthorParams) (*domain.Author, error)
	ListAuthors(context.Context, domain.AuthorQuery) ([]*domain.Author, error)

	NewTag(context.Context, domain.NewTagParams) (*domain.Tag, error)
	ListTags(context.Context, domain.TagQuery) ([]*domain.Tag, error)

	NewBook(context.Context, domain.NewBookParams) (*domain.Book, error)
	ListBooks(context.Context, domain.BookQuery) ([]*domain.Book, error)
}

type Config struct {
	Domain     Domain
	Router     *mux.Router
	SigningKey string
}

func NewServer(cfg Config) *Server {
	if cfg.Router == nil {
		cfg.Router = mux.NewRouter()
	}
	srv := &Server{
		dom:    cfg.Domain,
		router: cfg.Router,
		tokens: token.New(cfg.SigningKey),
	}
	srv.routes()
	return srv
}

type Server struct {
	dom    Domain
	router *mux.Router
	tokens *token.ParserGenerator
}

func (s *Server) registerRoutes(router *mux.Router, routes ...route) {
	for _, route := range routes {
		h := route.Handler
		if route.Private {
			h = private(h)
		}
		h = s.authenticate(h)
		router.Path(route.Path).Methods(route.Method).HandlerFunc(h)
	}
}

func (s *Server) routes() {
	routes := []route{
		{
			Path:    `/UserService.SignUp`,
			Method:  "POST",
			Handler: s.signUp(),
		},
		{
			Path:    `/UserService.SignIn`,
			Method:  "POST",
			Handler: s.signIn(),
		},
		{
			Path:    `/UserService.SignOut`,
			Method:  "GET",
			Handler: s.signOut(),
		},
		{
			Path:    `/GenreService.Create`,
			Method:  "POST",
			Handler: s.createGenre(),
			Private: true,
		},
		{
			Path:    `/GenreService.List`,
			Method:  "GET",
			Handler: s.listGenres(),
		},
		{
			Path:    `/AuthorService.Create`,
			Method:  "POST",
			Handler: s.createAuthor(),
			Private: true,
		},
		{
			Path:    `/AuthorService.List`,
			Method:  "GET",
			Handler: s.listAuthors(),
		},
		{
			Path:    `/TagService.Create`,
			Method:  "POST",
			Handler: s.createTag(),
			Private: true,
		},
		{
			Path:    `/TagService.List`,
			Method:  "GET",
			Handler: s.listTags(),
		},
		{
			Path:    `/BookService.Create`,
			Method:  "POST",
			Handler: s.createBook(),
			Private: true,
		},
		{
			Path:    "/BookService.List",
			Method:  "GET",
			Handler: s.listBooks(),
		},
		{
			Path:    "/UserService.Me",
			Method:  "GET",
			Handler: s.me(),
			Private: true,
		},
	}
	s.registerRoutes(s.router, routes...)
	s.router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.Error(w, 404, fmt.Errorf("path not found: %s", r.URL.Path))
	})
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
