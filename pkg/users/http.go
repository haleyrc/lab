package users

import (
	"log"
	"net/http"
	"time"

	"github.com/haleyrc/cheevos-simple/api/lib/json"
)

type HTTPAdapter struct {
	svc  *Service
	resp core.HTTPResponder
}

func (a *HTTPAdapter) Me(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := getCurrentUserFromContext(ctx)

	resp, err := a.svc.GetUser(ctx, GetUserRequest{
		ID: userID,
	})
	if err != nil {
		a.resp.RespondWithError(w, err)
		return
	}

	a.resp.RespondOK(w, resp)
}

func (a *HTTPAdapter) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateUserRequest
	if err := json.Decode(r.Body, &req); err != nil {
		a.resp.RespondWithError(w, err)
		return
	}

	resp, err := a.svc.CreateUser(ctx, req)
	if err != nil {
		a.resp.RespondWithError(w, err)
		return
	}

	token, err := a.tokens.Generate(user.ID)
	if err != nil {
		a.resp.RespondWithError(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    token,
		MaxAge:   int(time.Hour.Seconds()),
		HttpOnly: true,
	})

	log.Println("created user:", user)

	a.resp.RespondOK(w, resp)
}

func (a *HTTPAdapter) SignIn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req SignInRequest
	if err := json.Decode(r.Body, &req); err != nil {
		a.resp.RespondWithError(w, err)
		return
	}

	resp, err := a.svc.SignIn(ctx, req)
	if err != nil {
		a.resp.RespondWithError(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    resp.Token,
		MaxAge:   int(time.Hour.Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	log.Println("signed in user:", user)

	a.resp.RespondOK(w, resp)
}

func (a *HTTPAdapter) SignOut(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
	})
	a.resp.RespondOK(w, struct{}{})
}
