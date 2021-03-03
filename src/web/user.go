package web

import (
	"net/http"

	"chores-suck/core"
	"chores-suck/web/messages"
)

type UserService interface {
	CreateUser(wr http.ResponseWriter, req *http.Request)
}

type userService struct {
	users core.UserService
	views ViewService
}

func NewUserService(u core.UserService, v ViewService) UserService {
	return &userService{
		users: u,
		views: v,
	}
}

func (s *userService) CreateUser(wr http.ResponseWriter, req *http.Request) {
	username := req.FormValue("username")
	email := req.FormValue("email")
	password := req.FormValue("pword")
	password2 := req.FormValue("pwordConf")

	ok := true
	msg := messages.RegisterMessage{}
	if !validateRegisterInput(username, password, password2, email, &msg) {
		ok = false
	}

	if ok {
		user := core.User{Username: username, Email: email}
		var err error
		user.Password, err = hashPassword(password)
		if err != nil {
			handleError(internalError(err), wr)
			return
		}
		err = s.users.CreateUser(&user)
		if err != nil {
			switch err {
			case core.ErrEmailExists:
				msg.Email = "Email already registered"
				ok = false
			case core.ErrNameExists:
				msg.Username = "Username already taken"
				ok = false
			default:
				handleError(internalError(err), wr)
				return
			}
		}
	}

	if !ok {
		s.views.RegisterFail(wr, req, &msg)
		return
	}

	http.Redirect(wr, req, "/login", 302)
}
