package web

import (
	"fmt"
	"net/http"

	"github.com/haleyrc/cheevos-simple/api/lib/json"
)

var ErrNotAuthorized = fmt.Errorf("not authorized")

func private(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getCurrentUserFromContext(r.Context())
		if userID == "" {
			json.Error(w, 401, ErrNotAuthorized)
			return
		}
		next(w, r)
	}
}

func (s *Server) authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		cookie, err := r.Cookie("access_token")
		if err == nil {
			claims, err := s.tokens.Parse(cookie.Value)
			if err != nil {
				fmt.Println("authenticate:", err)
				json.Error(w, 401, ErrNotAuthorized)
				return
			}
			ctx = setCurrentUserOnContext(ctx, claims.Sub)
		}
		// TODO (RCH): This feels gross
		if err != nil && err != http.ErrNoCookie {
			fmt.Println("authenticate: error getting cookie:", err)
		}

		next(w, r.WithContext(ctx))
	}
}
