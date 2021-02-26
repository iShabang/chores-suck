package views

import (
	"chores-suck/core/types"
	cerror "chores-suck/web/errors"
	"chores-suck/web/messages"
	"chores-suck/web/sessions"
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
	NewGroupForm(http.ResponseWriter, *http.Request, uint64, *messages.CreateGroup) error
	EditGroupForm(http.ResponseWriter, *http.Request, *types.Group, *types.User) error
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

	chores, err := s.repo.GetUserChoreList(&user)
	if err != nil {
		return internalError(err)
	}

	err = s.repo.GetUserMemberships(&user)
	if err != nil {
		return internalError(err)
	}
	model := struct {
		User   *types.User
		Chores []types.ChoreListItem
	}{
		User:   &user,
		Chores: chores,
	}
	err = executeTemplate(wr, model, "../html/dashboard.html")
	if err != nil {
		return internalError(err)
	}
	return nil
}

func (s *service) RegisterForm(wr http.ResponseWriter, req *http.Request, msg *messages.RegisterMessage) error {
	model := struct {
		Username string
		Email    string
		Messages *messages.RegisterMessage
		User     *types.User
	}{
		Username: req.FormValue("username"),
		Email:    req.FormValue("email"),
		Messages: msg,
		User:     nil,
	}
	e := executeTemplate(wr, model, "../html/register.html")
	if e != nil {
		return internalError(e)
	}
	return nil
}

func (s *service) LoginForm(wr http.ResponseWriter, req *http.Request) error {
	model := struct {
		User *types.User
	}{
		User: nil,
	}
	e := executeTemplate(wr, model, "../html/login.html")
	if e != nil {
		return internalError(e)
	}
	return nil
}

func (s *service) NewGroupForm(wr http.ResponseWriter, req *http.Request, uid uint64, msg *messages.CreateGroup) error {
	user := types.User{ID: uid}
	e := s.repo.GetUserByID(&user)
	if e != nil {
		return internalError(e)
	}
	model := struct {
		User *types.User
		Msg  *messages.CreateGroup
	}{
		User: &user,
		Msg:  msg,
	}
	e = executeTemplate(wr, model, "../html/newgroup.html")
	if e != nil {
		return internalError(e)
	}
	return nil
}

func (s *service) EditGroupForm(wr http.ResponseWriter, req *http.Request, group *types.Group, user *types.User) error {
	model := struct {
		User  *types.User
		Group *types.Group
	}{
		User:  user,
		Group: group,
	}
	err := executeTemplate(wr, model, "../html/editgroup.html")
	if err != nil {
		return internalError(err)
	}
	return nil
}

func internalError(e error) cerror.StatusError {
	return cerror.StatusError{Code: http.StatusInternalServerError, Err: e}
}

func executeTemplate(wr http.ResponseWriter, model interface{}, files ...string) error {
	common := []string{"../html/layout.html", "../html/navbar.html"}
	files = append(files, common...)
	var t *template.Template
	t, err := template.ParseFiles(files...)
	if err != nil {
		return err
	}
	err = t.ExecuteTemplate(wr, "layout", model)
	return err
}
