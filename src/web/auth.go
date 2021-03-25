package web

import (
	"log"
	"net/http"
	"os"

	"chores-suck/core"

	"github.com/gorilla/sessions"
)

var (
	// Name of the session cookie that will be sent to clients
	SessionName = os.Getenv("SESSION_NAME")
)

// Service provides functionality for authentication and authorization
type AuthService interface {
	Login(http.ResponseWriter, *http.Request)
	Logout(http.ResponseWriter, *http.Request)
	Authorize(http.ResponseWriter, *http.Request) (uint64, error)
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

func (s *authService) Login(wr http.ResponseWriter, req *http.Request) {
	authorized := false
	ses, e := s.store.Get(req, SessionName)
	if e != nil {
		handleError(internalError(e), wr)
		return
	} else if !ses.IsNew {
		log.Print("Session not new")
		e = getSessionValue("auth", &authorized, ses)
		if e != nil {
			log.Print("Existing session failure during login")
			handleError(internalError(e), wr)
			return
		}
	}

	if !authorized {
		n := req.FormValue("username")
		p := req.FormValue("pword")
		u := core.User{Username: n, Password: p}
		e = s.checkCredentials(&u)
		if e == ErrNotAuthorized {
			SetFlash(wr, "genError", []byte("Invalid username/password"))
			http.Redirect(wr, req, "/login", 302)
			return
		} else if e != nil {
			handleError(internalError(e), wr)
		}
		ses.Values["userid"] = u.ID
		ses.Values["auth"] = true
		if e = ses.Save(req, wr); e != nil {
			handleError(internalError(e), wr)
			return
		}
	}
	// TODO: Check cookie for a redirect url
	http.Redirect(wr, req, "/dashboard", 302)
}

func (s *authService) checkCredentials(u *core.User) error {
	p := u.Password
	e := s.users.GetUserByName(u)
	if e != nil {
		log.Print(e)
		return ErrNotAuthorized
	}
	r := checkpword(p, u.Password)
	if !r {
		return ErrNotAuthorized
	}
	return nil
}

func (s *authService) Logout(wr http.ResponseWriter, req *http.Request) {
	ses, e := s.store.Get(req, SessionName)
	if e != nil {
		handleError(internalError(e), wr)
	} else if !ses.IsNew {
		var authorized bool
		e = getSessionValue("auth", &authorized, ses)
		if e != nil {
			handleError(internalError(e), wr)
		} else if !authorized {
			handleError(authError(ErrNotAuthorized), wr)
		}

		ses.Values["auth"] = false
		ses.Options.MaxAge = -1
		ses.Save(req, wr)
	}

	http.Redirect(wr, req, "/", 302)
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
