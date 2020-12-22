package rest

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"chores-suck/http/auth"
)

//Handler creates and returns a new http.Handler with the request handlers and functions pre-registered/routed
func Handler(a auth.Service) http.Handler {
	ro := httprouter.New()
	ro.POST("/login", login(a))
	return ro
}

func login(service auth.Service) func(wr http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	return func(wr http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		a, e := service.Authenticate(wr, req)

		if !a {
			// TODO: Inform client that login failed (wrong username/password)
			if e == nil {
				http.Error(wr, e.Error(), http.StatusUnauthorized)
			} else {
				wr.WriteHeader(http.StatusUnauthorized)
			}
		}

		http.Redirect(wr, req, "/", 302)
	}
}
