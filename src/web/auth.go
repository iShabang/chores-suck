package web

import (
	"log"
	"net/http"
	"os"

	"chores-suck/core"
	"chores-suck/web/messages"

	"github.com/gorilla/sessions"
)

var (
	// Name of the session cookie that will be sent to clients
	SessionName = os.Getenv("SESSION_NAME")
)

// Service provides functionality for authentication and authorization
type AuthService interface {
	Login(http.ResponseWriter, *http.Request) error
	Logout(http.ResponseWriter, *http.Request) error
	Authorize(http.ResponseWriter, *http.Request) (uint64, error)
	Create(http.ResponseWriter, *http.Request, *messages.RegisterMessage) error
}

type authService struct {
	users core.UserService
	store sessions.Store
}

// NewService creates and returns a new auth Service
func NewAuthService(us core.UserService, ses sessions.Store) AuthService {
	return &authService{
		users: us,
		store: ses,
	}
}

func (s *authService) Login(wr http.ResponseWriter, req *http.Request) error {
	ses, e := s.store.Get(req, SessionName)
	if e != nil {
		return internalError(e)
	} else if !ses.IsNew {
		log.Print("Session not new")
		var authorized bool
		e = getSessionValue("auth", &authorized, ses)
		if e != nil {
			log.Print("Existing session failure during login")
			return internalError(e)
		} else if authorized {
			return nil
		}
	}

	n := req.FormValue("username")
	u := core.User{Username: n}
	e = s.users.GetUserByName(&u)
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

func (s *authService) Logout(wr http.ResponseWriter, req *http.Request) error {
	ses, e := s.store.Get(req, SessionName)
	if e != nil {
		return internalError(e)
	} else if ses.IsNew {
		return nil
	}

	var authorized bool
	e = getSessionValue("auth", &authorized, ses)
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

func (s *authService) Authorize(wr http.ResponseWriter, req *http.Request) (uint64, error) {
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
func (s *authService) Create(wr http.ResponseWriter, req *http.Request, msg *messages.RegisterMessage) error {
	username := req.FormValue("username")
	email := req.FormValue("email")
	password := req.FormValue("pword")
	password2 := req.FormValue("pwordConf")

	if !validateRegisterInput(username, password, password2, email, msg) {
		return ErrInvalidInput
	}

	user := core.User{Username: username, Email: email}
	var err error
	user.Password, err = hashPassword(password)
	if err != nil {
		return internalError(err)
	}

	err = s.users.CreateUser(&user)
	if err != nil {
		switch err {
		case core.ErrEmailExists:
			msg.Email = "Email already registered"
			return ErrInvalidInput
		case core.ErrNameExists:
			msg.Username = "Username already taken"
			return ErrInvalidInput
		default:
			return internalError(err)
		}
	}

	return nil
}

/////////////////////////////////////////////////////////////////
// Helper methods
/////////////////////////////////////////////////////////////////

func (s *authService) isLoggedIn(r *http.Request) bool {
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
