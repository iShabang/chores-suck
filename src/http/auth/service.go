package auth

import (
	"net/http"

	"github.com/gorilla/sessions"
)

// Service provides functionality for authentication and authorization
type Service interface {
	Authenticate(http.ResponseWriter, *http.Request) (bool, error)
	Authorize(http.ResponseWriter, *http.Request) bool
}

// Repository defines storage functionality for a service
type Repository interface {
	//AddUser(User) error
	//UpdateUser(User, string, string) error
	GetUser(string) (User, error)
}

type service struct {
	repo  Repository
	store sessions.Store
}

// NewService creates and returns a new auth Service
func NewService(rep Repository, ses sessions.Store) Service {
	return &service{
		repo:  rep,
		store: ses,
	}
}

func (s service) Authenticate(wr http.ResponseWriter, req *http.Request) (bool, error) {
	if s.isLoggedIn(req) {
		return true, nil
	}

	n := req.FormValue("username")
	u, e := s.repo.GetUser(n)

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
