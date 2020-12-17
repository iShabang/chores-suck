package rest

import (
	"net/http"
	"strconv"
	"time"

	"chores-suck/users"
)

// Login function that handles http login requests
func Login(service users.Service) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
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
