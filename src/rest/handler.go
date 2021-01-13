package rest

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"chores-suck/rest/auth"
	"chores-suck/rest/errors"
)

//Handler creates and returns a new http.Handler with the request handlers and functions pre-registered/routed
func Handler(a auth.Service) http.Handler {
	ro := httprouter.New()
	ro.POST("/login", login(a))
	return ro
}

// Create middleware for authentication
func requiresLogin(handler func(wr http.ResponseWriter, req *http.Request, ps httprouter.Params)) func(wr http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	return func(wr http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		// Do Login logic

		handler(wr, req, ps)
	}
}

func logout(service auth.Service) func(wr http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	return func(wr http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		err := service.Logout(wr, req)
		handleError(err, wr)
		http.Redirect(wr, req, "/", 302)
	}
}

func login(service auth.Service) func(wr http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	return func(wr http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		err := service.Login(wr, req)
		handleError(err, wr)
		http.Redirect(wr, req, "/", 302)
	}
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
