package web

import (
	"chores-suck/core"
	"chores-suck/web/messages"
	"chores-suck/web/sessions"
	"html/template"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type RegisterFormData struct {
	Username string
	Email    string
	Messages *messages.RegisterMessage
}

// Service provides functionality for generating views
type ViewService interface {
	Index(http.ResponseWriter, *http.Request)
	Dashboard(http.ResponseWriter, *http.Request, uint64)
	RegisterForm(http.ResponseWriter, *http.Request)
	LoginForm(http.ResponseWriter, *http.Request)
	NewGroupForm(http.ResponseWriter, *http.Request, uint64)
	EditGroupForm(http.ResponseWriter, *http.Request, httprouter.Params, *core.User, *core.Group)
	NewRoleForm(http.ResponseWriter, *http.Request, httprouter.Params, *core.User, *core.Group)
	UpdateRoleForm(http.ResponseWriter, *http.Request, httprouter.Params, *core.User, *core.Role)
}

type viewService struct {
	store  *sessions.Store
	users  core.UserService
	groups core.GroupService
	auth   AuthService
}

func NewViewService(s *sessions.Store, u core.UserService, a AuthService, g core.GroupService) ViewService {
	return &viewService{
		store:  s,
		users:  u,
		auth:   a,
		groups: g,
	}
}

func (s *viewService) Index(wr http.ResponseWriter, req *http.Request) {
	var t *template.Template
	t, err := template.ParseFiles("../html/index.html", "../html/navbar.html", "../html/head.html")
	if err != nil {
		handleError(internalError(err), wr)
		return
	}
	err = t.ExecuteTemplate(wr, "index", nil)
	if err != nil {
		handleError(internalError(err), wr)
		return
	}
}

func (s *viewService) Dashboard(wr http.ResponseWriter, req *http.Request, uid uint64) {
	user := core.User{}
	user.ID = uid
	err := s.users.GetUserByID(&user)
	if err != nil {
		handleError(internalError(err), wr)
		return
	}

	err = s.users.GetChores(&user)
	if err != nil {
		handleError(internalError(err), wr)
		return
	}

	err = s.users.GetMemberships(&user)
	if err != nil {
		handleError(internalError(err), wr)
		return
	}
	model := struct {
		User *core.User
	}{
		User: &user,
	}
	err = executeTemplate(wr, model, "../html/dashboard.html")
	if err != nil {
		handleError(internalError(err), wr)
		return
	}
}

func (s *viewService) RegisterForm(wr http.ResponseWriter, req *http.Request) {
	var nameErr string
	var emailErr string
	var passErr string
	if data, _ := GetFlash(wr, req, "nameError"); data != nil {
		nameErr = string(data)
	}
	if data, _ := GetFlash(wr, req, "emailError"); data != nil {
		emailErr = string(data)
	}
	if data, _ := GetFlash(wr, req, "passError"); data != nil {
		passErr = string(data)
	}
	model := struct {
		Username   string
		Email      string
		User       *core.User
		NameError  string
		EmailError string
		PassError  string
	}{
		Username:   req.FormValue("username"),
		Email:      req.FormValue("email"),
		User:       nil,
		NameError:  nameErr,
		EmailError: emailErr,
		PassError:  passErr,
	}
	e := executeTemplate(wr, model, "../html/register.html")
	if e != nil {
		handleError(internalError(e), wr)
	}
}

func (s *viewService) LoginForm(wr http.ResponseWriter, req *http.Request) {
	_, e := s.auth.Authorize(wr, req)
	if e == nil {
		http.Redirect(wr, req, "/dashboard", 302)
	}
	var err string
	if data, _ := GetFlash(wr, req, "genError"); data != nil {
		err = string(data)
	}
	model := struct {
		User  *core.User
		Error string
	}{
		User:  nil,
		Error: err,
	}
	e = executeTemplate(wr, model, "../html/login.html")
	if e != nil {
		handleError(internalError(e), wr)
		return
	}
}
func (s *viewService) NewGroupForm(wr http.ResponseWriter, req *http.Request, uid uint64) {
	user := core.User{ID: uid}
	e := s.users.GetUserByID(&user)
	if e != nil {
		handleError(internalError(e), wr)
		return
	}
	var genErr string
	var nameErr string
	if data, _ := GetFlash(wr, req, "genError"); data != nil {
		genErr = string(data)
	}
	if data, _ := GetFlash(wr, req, "nameError"); data != nil {
		nameErr = string(data)
	}
	model := struct {
		User      *core.User
		GenError  string
		NameError string
	}{
		User:      &user,
		GenError:  genErr,
		NameError: nameErr,
	}
	e = executeTemplate(wr, model, "../html/newgroup.html")
	if e != nil {
		handleError(internalError(e), wr)
		return
	}
}

func (s *viewService) EditGroupForm(wr http.ResponseWriter, req *http.Request,
	ps httprouter.Params, user *core.User, group *core.Group) {
	if e := s.groups.GetRoles(group); e != nil {
		log.Printf("UserID: %v, GroupID: %v, EditGroupForm: Error: %s", user.ID, group.ID, e.Error())
	}
	var nameErr string
	var memErr string
	if data, _ := GetFlash(wr, req, "nameError"); data != nil {
		nameErr = string(data)
	}
	if data, _ := GetFlash(wr, req, "memError"); data != nil {
		memErr = string(data)
	}
	model := struct {
		User      *core.User
		Group     *core.Group
		NameError string
		MemError  string
	}{
		User:      user,
		Group:     group,
		NameError: nameErr,
		MemError:  memErr,
	}
	err := executeTemplate(wr, model, "../html/editgroup.html")
	if err != nil {
		handleError(internalError(err), wr)
	}
}

func (s *viewService) NewRoleForm(wr http.ResponseWriter, req *http.Request,
	ps httprouter.Params, user *core.User, group *core.Group) {
	mem := group.FindMember(user.ID)
	if e := s.groups.GetRoles(mem); e != nil {
		http.Error(wr, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if !mem.SuperRole.Can(core.EditRoles) {
		http.Error(wr, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	var genErr string
	if data, _ := GetFlash(wr, req, "genError"); data != nil {
		genErr = string(data)
	}
	model := struct {
		User  *core.User
		Group *core.Group
		Error string
	}{
		User:  user,
		Group: group,
		Error: genErr,
	}
	executeTemplate(wr, model, "../html/addrole.html")
}

func (s *viewService) UpdateRoleForm(wr http.ResponseWriter, req *http.Request,
	ps httprouter.Params, user *core.User, role *core.Role) {
	s.groups.GetMemberships(role)
	var msg string
	if data, e := GetFlash(wr, req, "genError"); data != nil {
		msg = string(data)
	} else if e != nil {
		log.Printf("UpdateRoleForm: Failed to get flash message: %s", e.Error())
	}
	model := struct {
		User  *core.User
		Group *core.Group
		Role  *core.Role
		Error string
	}{
		User:  user,
		Group: role.Group,
		Role:  role,
		Error: msg,
	}
	executeTemplate(wr, model, "../html/editrole.html")
}

func executeTemplate(wr http.ResponseWriter, model interface{}, files ...string) error {
	common := []string{"../html/layout.html", "../html/navbar.html"}
	files = append(files, common...)
	var t *template.Template
	t, err := template.ParseFiles(files...)
	if err == nil {
		err = t.ExecuteTemplate(wr, "layout", model)
	}
	if err != nil {
		log.Print(err)
	}
	return err
}

func findMembership(user *core.User, groupID uint64) *core.Membership {
	for _, v := range user.Memberships {
		if v.Group.ID == groupID {
			return &v
		}
	}
	return nil
}

func findPermission(mem *core.Membership, action core.PermBit) bool {
	for _, v := range mem.Roles {
		if v.Can(action) {
			return true
		}
	}
	return false
}
