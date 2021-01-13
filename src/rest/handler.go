package rest

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"chores-suck/http/auth"
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

func login(service auth.Service) func(wr http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	return func(wr http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		err := service.Login(wr, req)
		if err != nil {
			switch e := err.(type) {
			case errors.Error:
				http.Error(wr, e.Error(), e.Status())
			default:
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}

		http.Redirect(wr, req, "/", 302)
	}
}
