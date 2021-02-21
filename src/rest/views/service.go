package views

import (
	"chores-suck/core/types"
	cerror "chores-suck/rest/errors"
	"chores-suck/rest/messages"
	"chores-suck/rest/sessions"
	"html/template"
	"net/http"
)

type RegisterFormData struct {
	Username string
	Email    string
	Messages *messages.RegisterMessage
}

// Service provides functionality for generating views
type Service interface {
	Index(http.ResponseWriter, *http.Request) error
	BuildDashboard(http.ResponseWriter, *http.Request, uint64) error
	RegisterForm(http.ResponseWriter, *http.Request, *messages.RegisterMessage) error
	LoginForm(http.ResponseWriter, *http.Request) error
}

// Repository describes the interface necessary for grabbing data necessary for views
type Repository interface {
	GetUserByID(*types.User) error
	GetUserChoreList(*types.User) ([]types.ChoreListItem, error)
	GetUserMemberships(*types.User) error
}

type service struct {
	store *sessions.Store
	repo  Repository
}

func NewService(s *sessions.Store, r Repository) Service {
	return &service{
		store: s,
		repo:  r,
	}
}

func (s *service) Index(wr http.ResponseWriter, req *http.Request) error {
	var t *template.Template
	t, err := template.ParseFiles("../html/index.html", "../html/navbar.html", "../html/head.html")
	if err != nil {
		return internalError(err)
	}
	err = t.ExecuteTemplate(wr, "index", nil)
	if err != nil {
		return internalError(err)
	}
	return nil
}

func (s *service) BuildDashboard(wr http.ResponseWriter, req *http.Request, uid uint64) error {
	user := types.User{}

	user.ID = uid
	err := s.repo.GetUserByID(&user)
	if err != nil {
		return internalError(err)
	}

	var t *template.Template
	t, err = template.ParseFiles("../html/dashboard.html", "../html/sessionNav.html", "../html/head.html")
	if err != nil {
		return internalError(err)
	}

	chores, err := s.repo.GetUserChoreList(&user)
	if err != nil {
		return internalError(err)
	}

	err = s.repo.GetUserMemberships(&user)
	if err != nil {
		return internalError(err)
	}

	dm := DashboardModel{
		User:   &user,
		Chores: chores,
	}
	err = t.ExecuteTemplate(wr, "index", dm)
	if err != nil {
		return internalError(err)
	}
	return nil
}

func (s *service) RegisterForm(wr http.ResponseWriter, req *http.Request, msg *messages.RegisterMessage) error {
	var t *template.Template
	t, err := template.ParseFiles("../html/register.html", "../html/navbar.html", "../html/head.html")
	if err != nil {
		return internalError(err)
	}
	data := RegisterFormData{
		Username: req.FormValue("username"),
		Email:    req.FormValue("email"),
		Messages: msg,
	}
	err = t.ExecuteTemplate(wr, "register.html", data)
	if err != nil {
		return internalError(err)
	}
	return nil
}

func (s *service) LoginForm(wr http.ResponseWriter, req *http.Request) error {
	var t *template.Template
	t, err := template.ParseFiles("../html/login.html", "../html/head.html", "../html/navbar.html")
	if err != nil {
		return internalError(err)
	}
	err = t.ExecuteTemplate(wr, "login.html", nil)
	if err != nil {
		return internalError(err)
	}
	return nil
}

func internalError(e error) cerror.StatusError {
	return cerror.StatusError{Code: http.StatusInternalServerError, Err: e}
}
