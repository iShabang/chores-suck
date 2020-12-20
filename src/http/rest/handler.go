package rest

import (
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"

	"chores-suck/auth"
)

//Handler creates and returns a new http.Handler with the request handlers and functions pre-registered/routed
func Handler(a auth.Service) http.Handler {
	ro := httprouter.New()
	ro.POST("/login", login(a))
	return ro
}

func login(service auth.Service) func(writer http.ResponseWriter, request *http.Request, ps httprouter.Params) {
	return func(writer http.ResponseWriter, request *http.Request, ps httprouter.Params) {
		if request.Method != "POST" {
			http.Error(writer, "Invlalid login command", http.StatusMethodNotAllowed)
			return
		}

		n := request.FormValue("username")
		p := request.FormValue("password")

		u, e := service.Authenticate(n, p)

		// TODO: error check
		if e != nil {
			return
		}

		s, e := service.StartSession(u)

		// TODO: error check
		if e != nil {
			return
		}

		c := http.Cookie{
			Name:     "session",
			Value:    strconv.FormatUint(s.ID, 10),
			Expires:  time.Unix(s.ExpireTime, 0),
			Secure:   false,
			HttpOnly: true,
		}
		http.SetCookie(writer, &c)

		writer.WriteHeader(http.StatusOK)
	}
}
