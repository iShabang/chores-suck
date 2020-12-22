package auth

import (
	"chores-suck/users"
	"net/http"

	"github.com/gorilla/sessions"
)

// Service provides functionality for user types
type Service interface {
	Authenticate(http.ResponseWriter, *http.Request) (bool, error)
	Authorize(http.ResponseWriter, *http.Request) bool
}

// NewService creates and returns a new auth Service
func NewService(ur users.Repository, ses sessions.Store) Service {
	return &service{
		userRepo: ur,
		store:    ses,
	}
}

type service struct {
	userRepo users.Repository
	store    sessions.Store
}

func (s service) Authenticate(wr http.ResponseWriter, req *http.Request) (bool, error) {
	if s.isLoggedIn(req) {
		return true, nil
	}

	n := req.FormValue("username")
	u, e := s.userRepo.Get(n)

	if e != nil {
		return false, nil
	}

	p := req.FormValue("password")
	r := checkpword(p, u.Password)

	if !r {
		return false, nil
	}

	ses, e := s.store.New(req, "session")
	ses.Values["username"] = n
	e = ses.Save(req, wr)

	if e != nil {
		return false, e
	}

	return true, nil
}

func (s service) Authorize(wr http.ResponseWriter, req *http.Request) bool {
	return s.isLoggedIn(req)
}

func (s service) isLoggedIn(r *http.Request) bool {
	ses, e := s.store.Get(r, "session")
	if e != nil {
		return false
	}
	return ses.Values["loggedin"] == "true"
}
