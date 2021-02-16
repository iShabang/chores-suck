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
	ro.HandlerFunc("POST", "/login", s.login)
	ro.HandlerFunc("POST", "/logout", s.requiresLogin(s.logout))
	ro.HandlerFunc("GET", "/dashboard", s.requiresLogin(s.dashboard))
	ro.HandlerFunc("GET", "/register", s.register)
	return ro
}

// Create middleware for authentication
func (s *Services) requiresLogin(handler func(wr http.ResponseWriter, req *http.Request, uid string)) http.HandlerFunc {
	return func(wr http.ResponseWriter, req *http.Request) {
		uid, err := s.auth.Authorize(wr, req)
		if err != nil {
			handleError(err, wr)
			http.Redirect(wr, req, "/", 302)
			return
		}

		handler(wr, req, uid)
	}
}

func (s *Services) dashboard(wr http.ResponseWriter, req *http.Request, uid string) {
	s.views.BuildDashboard(wr, req, uid)
}

func (s *Services) logout(wr http.ResponseWriter, req *http.Request, uid string) {
	err := s.auth.Logout(wr, req)
	handleError(err, wr)
	http.Redirect(wr, req, "/", 302)
}

func (s *Services) login(wr http.ResponseWriter, req *http.Request) {
	err := s.auth.Login(wr, req)
	if err != nil {
		handleError(err, wr)
		return
	}
	http.Redirect(wr, req, "/", 302)
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
