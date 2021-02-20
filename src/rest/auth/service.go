package auth

import (
	cerror "chores-suck/rest/errors"
	"chores-suck/rest/messages"
	"chores-suck/types"
	"errors"
	"log"
	"net/http"
	"os"

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

	// ErrValueName occurs when attempting to access an invalid session value
	ErrValueName = errors.New("invalid session value name")

	// Name of the session cookie that will be sent to clients
	SessionName = os.Getenv("SESSION_NAME")
)

// Service provides functionality for authentication and authorization
type Service interface {
	Login(http.ResponseWriter, *http.Request) error
	Logout(http.ResponseWriter, *http.Request) error
	Authorize(http.ResponseWriter, *http.Request) (uint64, error)
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
	ses, e := s.store.Get(req, SessionName)
	if e != nil {
		return internalError(e)
	} else if !ses.IsNew {
		var authorized bool
		e = getSessionValue("auth", authorized, ses)
		if e != nil {
			return internalError(e)
		} else if authorized {
			return nil
		}
	}

	n := req.FormValue("username")
	u := types.User{Username: n}
	e = s.repo.GetUserByName(&u)
	if e != nil {
		log.Print(e)
		return authError(ErrNotAuthorized)
	}

	p := req.FormValue("pword")
	r := checkpword(p, u.Password)

	if !r {
		log.Printf("Incorrect Password. Passed-In: %s, Database: %s", p, u.Password)
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
	ses, e := s.store.Get(req, SessionName)
	if e != nil {
		return internalError(e)
	} else if ses.IsNew {
		return nil
	}

	var authorized bool
	e = getSessionValue("auth", authorized, ses)
	if e != nil {
		return internalError(e)
	} else if !authorized {
		return authError(ErrNotAuthorized)
	}

	ses.Values["auth"] = false
	ses.Options.MaxAge = -1
	ses.Save(req, wr)
	return nil
}

func (s *service) Authorize(wr http.ResponseWriter, req *http.Request) (uint64, error) {
	ses, e := s.store.Get(req, SessionName)
	if e != nil {
		log.Printf("Get session: %s", e.Error())
		return 0, internalError(e)
	} else if ses.IsNew {
		log.Print("New session not authorized")
		return 0, authError(ErrNotAuthorized)
	}

	var authorized bool
	if e = getSessionValue("auth", &authorized, ses); e != nil {
		log.Printf("Auth value: %s", e.Error())
		return 0, internalError(e)
	}
	if !authorized {
		log.Print("session not authorized")
		return 0, authError(ErrNotAuthorized)
	}

	var uid uint64
	if e = getSessionValue("userid", &uid, ses); e != nil {
		log.Printf("UserID value: %s", e.Error())
		return 0, internalError(e)
	}
	if uid == 0 {
		log.Printf("Invalid user id: %v", uid)
		return 0, authError(ErrNotAuthorized)
	}

	return uid, nil
}

// TODO: Move this method to a new User service
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
	if err != nil {
		return internalError(err)
	}
	return nil
}

/////////////////////////////////////////////////////////////////
// Helper methods
/////////////////////////////////////////////////////////////////

func (s *service) isLoggedIn(r *http.Request) bool {
	ses, e := s.store.Get(r, SessionName)
	if e != nil {
		return false
	}
	return ses.Values["auth"] == "true"
}

func getSessionValue(name string, result interface{}, ses *sessions.Session) error {
	if name == "" {
		return ErrValueName
	}

	var ok bool
	switch v := result.(type) {
	case *uint64:
		*v, ok = ses.Values[name].(uint64)
	case *bool:
		*v, ok = ses.Values[name].(bool)
	}

	if !ok {
		return ErrSessionType
	}

	return nil
}

func internalError(e error) cerror.StatusError {
	return cerror.StatusError{Code: http.StatusInternalServerError, Err: e}
}

func authError(e error) cerror.StatusError {
	return cerror.StatusError{Code: http.StatusUnauthorized, Err: e}
}
