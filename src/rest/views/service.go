package views

import (
	cerror "chores-suck/rest/errors"
	"chores-suck/rest/messages"
	"chores-suck/rest/sessions"
	"chores-suck/types"
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
	BuildDashboard(http.ResponseWriter, *http.Request, string) error
	RegisterForm(http.ResponseWriter, *http.Request, *messages.RegisterMessage) error
}

// Repository describes the interface necessary for grabbing data necessary for views
type Repository interface {
	GetUserByID(string) (types.User, error)
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

func (s *service) BuildDashboard(wr http.ResponseWriter, req *http.Request, uid string) error {
	// Get User
	// Get Memberships
	// Get Group Data
	// Populate template with data
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
	t.ExecuteTemplate(wr, "register.html", data)
	return nil
}

func internalError(e error) cerror.StatusError {
	return cerror.StatusError{Code: http.StatusInternalServerError, Err: e}
}
