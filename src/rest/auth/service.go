package auth

import (
	cerror "chores-suck/rest/errors"
	"chores-suck/types"
	"errors"
	"net/http"

	"github.com/gorilla/sessions"
)

var (
	// ErrSessionType describes a type assertion error for retrieving session values
	ErrSessionType = errors.New("session value type assertion failed")

	// ErrNotAuthorized occurs when the current session is not authorized
	ErrNotAuthorized = errors.New("authorization unsuccessfull")

	// ErrInvalidValue occurs when an invalid session value is retrieved from the current session
	ErrInvalidValue = errors.New("invalid session value")
)

// Service provides functionality for authentication and authorization
type Service interface {
	Login(http.ResponseWriter, *http.Request) error
	Logout(http.ResponseWriter, *http.Request) error
	Authorize(http.ResponseWriter, *http.Request) (string, error)
}

// Repository defines storage functionality for a service
type Repository interface {
	GetUserByName(string) (types.User, error)
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

func (s *service) Login(wr http.ResponseWriter, req *http.Request) error {
	ses, e := s.getSession(req)
	if e != nil {
		return e
	}
	if e = checkAuthValue(ses); e == nil {
		return nil
	}

	n := req.FormValue("username")
	u, e := s.repo.GetUserByName(n)
	if e != nil {
		return cerror.StatusError{Code: http.StatusInternalServerError, Err: e}
	} else if u.ID == 0 {
		return cerror.StatusError{Code: http.StatusUnauthorized, Err: ErrNotAuthorized}
	}

	p := req.FormValue("password")
	r := checkpword(p, u.Password)

	if !r {
		return cerror.StatusError{Code: http.StatusUnauthorized, Err: ErrNotAuthorized}
	}

	ses.Values["userid"] = u.ID
	ses.Values["auth"] = true
	if e = ses.Save(req, wr); e != nil {
		return cerror.StatusError{Code: http.StatusInternalServerError, Err: e}
	}

	return nil
}

func (s *service) Logout(wr http.ResponseWriter, req *http.Request) error {
	ses, e := s.getSession(req)
	if e != nil {
		return e
	}
	ses.Values["auth"] = false
	ses.Options.MaxAge = -1
	ses.Save(req, wr)
	return nil
}

func (s *service) Authorize(wr http.ResponseWriter, req *http.Request) (string, error) {
	ses, e := s.getSession(req)
	if e != nil {
		return "", e
	}

	if e = checkAuthValue(ses); e != nil {
		return "", e
	}

	var u string
	if u, e = getUserIdValue(ses); e != nil {
		return "", e
	}

	return u, nil
}

func (s *service) isLoggedIn(r *http.Request) bool {
	ses, e := s.store.Get(r, "session")
	if e != nil {
		return false
	}
	return ses.Values["auth"] == "true"
}

func (s *service) getSession(req *http.Request) (*sessions.Session, error) {
	ses, e := s.store.Get(req, "session")
	if e != nil {
		return nil, cerror.StatusError{Code: http.StatusInternalServerError, Err: e}
	}
	return ses, nil
}

func checkAuthValue(s *sessions.Session) error {
	auth, ok := s.Values["auth"].(bool)
	if !ok {
		return cerror.StatusError{Code: http.StatusInternalServerError, Err: ErrSessionType}
	} else if !auth {
		return cerror.StatusError{Code: http.StatusUnauthorized, Err: ErrNotAuthorized}
	}
	return nil
}

func getUserIdValue(s *sessions.Session) (string, error) {
	u, ok := s.Values["userid"].(string)
	if !ok {
		return "", cerror.StatusError{Code: http.StatusInternalServerError, Err: ErrSessionType}
	} else if u == "" {
		return "", cerror.StatusError{Code: http.StatusUnauthorized, Err: ErrNotAuthorized}
	}

	return u, nil
}
