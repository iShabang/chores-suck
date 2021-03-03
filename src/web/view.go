package web

import (
	"chores-suck/core/types"
	"chores-suck/core/users"
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
type ViewService interface {
	Index(http.ResponseWriter, *http.Request) error
	BuildDashboard(http.ResponseWriter, *http.Request, uint64) error
	RegisterForm(http.ResponseWriter, *http.Request, *messages.RegisterMessage) error
	LoginForm(http.ResponseWriter, *http.Request) error
	NewGroupForm(http.ResponseWriter, *http.Request, uint64, *messages.CreateGroup) error
	EditGroupForm(http.ResponseWriter, *http.Request, *types.Group, *types.User) error
}

type viewService struct {
	store *sessions.Store
	users users.Service
}

func NewViewService(s *sessions.Store, u users.Service) ViewService {
	return &viewService{
		store: s,
		users: u,
	}
}

func (s *viewService) Index(wr http.ResponseWriter, req *http.Request) error {
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

func (s *viewService) BuildDashboard(wr http.ResponseWriter, req *http.Request, uid uint64) error {
	user := types.User{}
	user.ID = uid
	err := s.users.GetUserByID(&user)
	if err != nil {
		return internalError(err)
	}

	err = s.users.GetChores(&user)
	if err != nil {
		return internalError(err)
	}

	err = s.users.GetMemberships(&user)
	if err != nil {
		return internalError(err)
	}
	model := struct {
		User *types.User
	}{
		User: &user,
	}
	err = executeTemplate(wr, model, "../html/dashboard.html")
	if err != nil {
		return internalError(err)
	}
	return nil
}

func (s *viewService) RegisterForm(wr http.ResponseWriter, req *http.Request, msg *messages.RegisterMessage) error {
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

func (s *viewService) LoginForm(wr http.ResponseWriter, req *http.Request) error {
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

func (s *viewService) NewGroupForm(wr http.ResponseWriter, req *http.Request, uid uint64, msg *messages.CreateGroup) error {
	user := types.User{ID: uid}
	e := s.users.GetUserByID(&user)
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

func (s *viewService) EditGroupForm(wr http.ResponseWriter, req *http.Request, group *types.Group, user *types.User) error {
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
