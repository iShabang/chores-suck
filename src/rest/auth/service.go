package auth

import (
	cerror "chores-suck/rest/errors"
	"chores-suck/rest/messages"
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

	// ErrInvalidFormData occurs when a form is submitted but the data is not valid
	ErrInvalidFormData = errors.New("one or more form values was invalid")
)

// Service provides functionality for authentication and authorization
type Service interface {
	Login(http.ResponseWriter, *http.Request) error
	Logout(http.ResponseWriter, *http.Request) error
	Authorize(http.ResponseWriter, *http.Request) (string, error)
	Create(http.ResponseWriter, *http.Request, *messages.RegisterMessage) error
}

// Repository defines storage functionality for a service
type Repository interface {
	GetUserByName(*types.User) error
	GetUserByEmail(*types.User) error
	CreateUser(*types.User) error
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
	u := types.User{Username: n}
	e = s.repo.GetUserByName(&u)
	if e != nil {
		return authError(ErrNotAuthorized)
	}

	p := req.FormValue("password")
	r := checkpword(p, u.Password)

	if !r {
		return authError(ErrNotAuthorized)
	}

	ses.Values["userid"] = u.ID
	ses.Values["auth"] = true
	if e = ses.Save(req, wr); e != nil {
		return internalError(e)
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

func (s *service) Create(wr http.ResponseWriter, req *http.Request, msg *messages.RegisterMessage) error {
	username := req.FormValue("username")
	email := req.FormValue("email")
	password := req.FormValue("pword")
	password2 := req.FormValue("pwordConf")

	if !validateInput(username, password, password2, email, msg) {
		return ErrInvalidFormData
	}

	user := types.User{Username: username, Email: email}
	err := s.repo.GetUserByName(&user)
	if err == nil {
		msg.Username = "Username already taken"
		return ErrInvalidFormData
	}

	err = s.repo.GetUserByEmail(&user)
	if err == nil {
		msg.Email = "Email already registered"
		return ErrInvalidFormData
	}

	user.Password, err = hashPassword(password)
	if err != nil {
		return internalError(err)
	}

	err = s.repo.CreateUser(&user)
	return err
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
		return nil, internalError(e)
	}
	return ses, nil
}

func checkAuthValue(s *sessions.Session) error {
	auth, ok := s.Values["auth"].(bool)
	if !ok {
		return internalError(ErrSessionType)
	} else if !auth {
		return authError(ErrNotAuthorized)
	}
	return nil
}

func getUserIdValue(s *sessions.Session) (string, error) {
	u, ok := s.Values["userid"].(string)
	if !ok {
		return "", internalError(ErrSessionType)
	} else if u == "" {
		return "", authError(ErrNotAuthorized)
	}

	return u, nil
}

func internalError(e error) cerror.StatusError {
	return cerror.StatusError{Code: http.StatusInternalServerError, Err: e}
}

func authError(e error) cerror.StatusError {
	return cerror.StatusError{Code: http.StatusUnauthorized, Err: e}
}
