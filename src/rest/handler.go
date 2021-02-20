package rest

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"chores-suck/rest/auth"
	"chores-suck/rest/errors"
	"chores-suck/rest/messages"
	"chores-suck/rest/views"
)

// Services holds references to services that handlers utilize to carry out requests
type Services struct {
	auth  auth.Service
	views views.Service
}

// NewServices creates a new Services object
func NewServices(a auth.Service, v views.Service) *Services {
	return &Services{
		auth:  a,
		views: v,
	}
}

//Handler creates and returns a new http.Handler with the request handlers and functions pre-registered/routed
func Handler(s *Services) http.Handler {
	ro := httprouter.New()
	ro.HandlerFunc("GET", "/login", s.loginGet)
	ro.HandlerFunc("POST", "/login", s.loginPost)
	ro.HandlerFunc("GET", "/logout", s.logout)
	ro.HandlerFunc("POST", "/createuser", s.createUser)
	ro.HandlerFunc("GET", "/dashboard", s.requiresLogin(s.dashboard))
	ro.HandlerFunc("GET", "/register", s.register)
	ro.HandlerFunc("GET", "/", s.index)
	return ro
}

func (s *Services) index(wr http.ResponseWriter, req *http.Request) {
	err := s.views.Index(wr, req)
	if err != nil {
		handleError(err, wr)
	}
}

func (s *Services) dashboard(wr http.ResponseWriter, req *http.Request, uid uint64) {
	err := s.views.BuildDashboard(wr, req, uid)
	if err != nil {
		handleError(err, wr)
	}
}

func (s *Services) logout(wr http.ResponseWriter, req *http.Request) {
	err := s.auth.Logout(wr, req)
	if err != nil {
		handleError(err, wr)
	}
	http.Redirect(wr, req, "/", 302)
}

func (s *Services) loginPost(wr http.ResponseWriter, req *http.Request) {
	err := s.auth.Login(wr, req)
	if err != nil {
		handleError(err, wr)
		return
	}
	// TODO: Check cookie for a redirect url
	http.Redirect(wr, req, "/dashboard", 302)
}

func (s *Services) loginGet(wr http.ResponseWriter, req *http.Request) {
	err := s.views.LoginForm(wr, req)
	if err != nil {
		handleError(err, wr)
	}
}

func (s *Services) register(wr http.ResponseWriter, req *http.Request) {
	msg := messages.RegisterMessage{}
	err := s.views.RegisterForm(wr, req, &msg)
	if err != nil {
		handleError(err, wr)
	}
}

func (s *Services) createUser(wr http.ResponseWriter, req *http.Request) {
	msg := messages.RegisterMessage{}
	err := s.auth.Create(wr, req, &msg)
	if err != nil {
		if err == auth.ErrInvalidFormData {
			s.views.RegisterForm(wr, req, &msg)
		} else {
			handleError(err, wr)
		}
		return
	}
	http.Redirect(wr, req, "/login", 302)
}

/////////////////////////////////////////////////////////////////
// Middleware methods
/////////////////////////////////////////////////////////////////
func (s *Services) requiresLogin(handler func(wr http.ResponseWriter, req *http.Request, uid uint64)) http.HandlerFunc {
	return func(wr http.ResponseWriter, req *http.Request) {
		// TODO: Save the requested url in a cookie that can be redirected to after logging in successfully
		uid, err := s.auth.Authorize(wr, req)
		if err != nil {
			http.Redirect(wr, req, "/login", 302)
			return
		}

		handler(wr, req, uid)
	}
}

/////////////////////////////////////////////////////////////////
// Helper methods
/////////////////////////////////////////////////////////////////

func handleError(err error, wr http.ResponseWriter) {
	if err != nil {
		switch e := err.(type) {
		case errors.Error:
			http.Error(wr, e.Error(), e.Status())
		default:
			http.Error(wr, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}
