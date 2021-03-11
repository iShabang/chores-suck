package web

import (
	"net/http"

	"chores-suck/core"
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
	e := validateUsername(username)
	if e != nil {
		SetFlash(wr, "nameError", []byte(e.Error()))
		ok = false
	}
	e = validateEmail(email)
	if e != nil {
		SetFlash(wr, "emailError", []byte(e.Error()))
		ok = false
	}
	e = validatePassword(password, password2)
	if e != nil {
		SetFlash(wr, "passError", []byte(e.Error()))
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
				SetFlash(wr, "emailError", []byte("Email already registered"))
				ok = false
			case core.ErrNameExists:
				SetFlash(wr, "nameError", []byte("Username already taken"))
				ok = false
			default:
				handleError(internalError(err), wr)
				return
			}
		}
	}

	if !ok {
		http.Redirect(wr, req, "/new/user", 302)
		return
	}

	http.Redirect(wr, req, "/login", 302)
}
