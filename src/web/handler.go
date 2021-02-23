package web

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"chores-suck/web/auth"
	"chores-suck/web/errors"
	"chores-suck/web/groups"
	"chores-suck/web/messages"
	"chores-suck/web/views"
)

// Services holds references to services that handlers utilize to carry out requests
type Services struct {
	auth   auth.Service
	views  views.Service
	groups groups.Service
}

// NewServices creates a new Services object
func NewServices(a auth.Service, v views.Service, g groups.Service) *Services {
	return &Services{
		auth:   a,
		views:  v,
		groups: g,
	}
}

//Handler creates and returns a new http.Handler with the request handlers and functions pre-registered/routed
func Handler(s *Services) http.Handler {
	ro := httprouter.New()
	ro.HandlerFunc("GET", "/login", s.loginGet)
	ro.HandlerFunc("POST", "/login", s.loginPost)
	ro.HandlerFunc("GET", "/logout", s.logout)
	ro.HandlerFunc("POST", "/createuser", s.createUser)
	ro.HandlerFunc("GET", "/dashboard", s.authorize(s.dashboard))
	ro.HandlerFunc("GET", "/register", s.register)
	ro.HandlerFunc("GET", "/creategroup", s.authorize(s.createGroupGet))
	ro.HandlerFunc("POST", "/creategroup", s.authorize(s.createGroupPost))
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
		// TODO: Resend login form with error message
		handleError(err, wr)
		return
	}
	// TODO: Check cookie for a redirect url
	http.Redirect(wr, req, "/dashboard", 302)
}

func (s *Services) loginGet(wr http.ResponseWriter, req *http.Request) {
	_, err := s.auth.Authorize(wr, req)
	if err == nil {
		http.Redirect(wr, req, "/dashboard", 302)
	} else {
		err = s.views.LoginForm(wr, req)
		if err != nil {
			handleError(err, wr)
		}
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
		if err == auth.ErrInvalidInput {
			s.views.RegisterForm(wr, req, &msg)
		} else {
			handleError(err, wr)
		}
		return
	}
	http.Redirect(wr, req, "/login", 302)
}

func (s *Services) createGroupGet(wr http.ResponseWriter, req *http.Request, uid uint64) {
	e := s.views.NewGroupForm(wr, req, uid, &messages.CreateGroup{})
	if e != nil {
		log.Printf("handler.createGroupGet: %s", e.Error())
		handleError(e, wr)
	}
}

func (s *Services) createGroupPost(wr http.ResponseWriter, req *http.Request, uid uint64) {
	msg := messages.CreateGroup{}
	e := s.groups.CreateGroup(wr, req, uid, &msg)
	if e != nil {
		log.Printf("handler.createGroup: %s", e.Error())
		msg.General = "There was an error creating a new group"
		s.views.NewGroupForm(wr, req, uid, &msg)
		return
	}

	// TODO: if successfull, redirect to the group page
	http.Redirect(wr, req, "/dashboard", 302)
}

func (s *Services) editGroupGet(wr http.ResponseWriter, req *http.Request, ps httprouter.Params) {

}

/////////////////////////////////////////////////////////////////
// Middleware methods
/////////////////////////////////////////////////////////////////
func (s *Services) authorize(handler func(wr http.ResponseWriter, req *http.Request, uid uint64)) http.HandlerFunc {
	return func(wr http.ResponseWriter, req *http.Request) {
		// TODO: Save the requested url in a cookie that can be redirected to after logging in successfully
		uid, err := s.auth.Authorize(wr, req)
		if err != nil {
			log.Print(err)
			http.Redirect(wr, req, "/login", 302)
			return
		}

		handler(wr, req, uid)
	}
}

func (s *Services) authorizeParam(handler func(wr http.ResponseWriter, req *http.Request, ps httprouter.Params, uid uint64)) httprouter.Handle {
	return func(wr http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		uid, err := s.auth.Authorize(wr, req)
		if err != nil {
			log.Print(err)
			http.Redirect(wr, req, "/login", 302)
			return
		}

		handler(wr, req, ps, uid)
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
