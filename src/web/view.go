package web

import (
	"chores-suck/core"
	"chores-suck/web/messages"
	"chores-suck/web/sessions"
	"html/template"
	"log"
	"net/http"
	"strconv"

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
	RegisterFail(http.ResponseWriter, *http.Request, *messages.RegisterMessage)
	LoginForm(http.ResponseWriter, *http.Request)
	NewGroupForm(http.ResponseWriter, *http.Request, uint64)
	NewGroupFail(http.ResponseWriter, *http.Request, *core.User, *messages.Group)
	EditGroupForm(http.ResponseWriter, *http.Request, httprouter.Params, *core.User, *core.Group)
	EditGroupFail(http.ResponseWriter, *http.Request, *core.User, *core.Group, *messages.Group)
	NewRoleForm(http.ResponseWriter, *http.Request, httprouter.Params, *core.User, *core.Group)
	NewRoleFail(http.ResponseWriter, *core.User, *core.Group, *messages.Group)
	UpdateRoleForm(http.ResponseWriter, *http.Request, httprouter.Params, *core.User, *core.Group)
	UpdateRoleFail(http.ResponseWriter, *core.User, *core.Group, *core.Role, *messages.Group)
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
	s.regFormInternal(wr, req, nil)
}

func (s *viewService) RegisterFail(wr http.ResponseWriter, req *http.Request, msg *messages.RegisterMessage) {
	s.regFormInternal(wr, req, msg)
}

func (s *viewService) regFormInternal(wr http.ResponseWriter, req *http.Request, msg *messages.RegisterMessage) {
	model := struct {
		Username string
		Email    string
		Messages *messages.RegisterMessage
		User     *core.User
	}{
		Username: req.FormValue("username"),
		Email:    req.FormValue("email"),
		Messages: msg,
		User:     nil,
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
	model := struct {
		User *core.User
	}{
		User: nil,
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
	s.newGroupInternal(wr, req, &user, nil)
}

func (s *viewService) NewGroupFail(wr http.ResponseWriter, req *http.Request,
	user *core.User, msg *messages.Group) {
	s.newGroupInternal(wr, req, user, msg)
}

func (s *viewService) newGroupInternal(wr http.ResponseWriter, req *http.Request, user *core.User, msg *messages.Group) {
	model := struct {
		User *core.User
		Msg  *messages.Group
	}{
		User: user,
		Msg:  msg,
	}
	e := executeTemplate(wr, model, "../html/newgroup.html")
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
	s.editGroupInternal(wr, req, user, group, nil)
}

func (s *viewService) EditGroupFail(wr http.ResponseWriter, req *http.Request,
	user *core.User, group *core.Group, msg *messages.Group) {
	s.editGroupInternal(wr, req, user, group, msg)
}

func (s *viewService) editGroupInternal(wr http.ResponseWriter, req *http.Request,
	user *core.User, group *core.Group, msg *messages.Group) {
	model := struct {
		User  *core.User
		Group *core.Group
		Msg   *messages.Group
	}{
		User:  user,
		Group: group,
		Msg:   msg,
	}
	err := executeTemplate(wr, model, "../html/editgroup.html")
	if err != nil {
		handleError(internalError(err), wr)
	}
}

func (s *viewService) NewRoleForm(wr http.ResponseWriter, req *http.Request,
	ps httprouter.Params, user *core.User, group *core.Group) {
	s.users.GetMemberships(user)
	mem := core.Membership{}
	for _, m := range user.Memberships {
		if m.Group.ID == group.ID {
			mem = m
			break
		}
	}
	s.groups.GetRoles(&mem)
	if !mem.SuperRole.Can(core.EditRoles) {
		http.Error(wr, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	s.newRoleInternal(wr, user, group, nil)
}

func (s *viewService) NewRoleFail(wr http.ResponseWriter, user *core.User, group *core.Group, msg *messages.Group) {
	s.newRoleInternal(wr, user, group, msg)
}

func (s *viewService) newRoleInternal(wr http.ResponseWriter, user *core.User,
	group *core.Group, msg *messages.Group) {
	model := struct {
		User  *core.User
		Group *core.Group
		Msg   *messages.Group
	}{
		User:  user,
		Group: group,
		Msg:   msg,
	}
	executeTemplate(wr, model, "../html/addrole.html")
}

func (s *viewService) UpdateRoleForm(wr http.ResponseWriter, req *http.Request,
	ps httprouter.Params, user *core.User, group *core.Group) {
	s.users.GetMemberships(user)
	s.groups.GetRoles(group)
	mem := findMembership(user, group.ID)
	s.groups.GetRoles(mem)
	if !mem.SuperRole.Can(core.EditRoles) {
		log.Printf("UserID: %v, GroupID: %v, UpdateRoleForm: Member cannot edit roles. Perm: %v",
			user.ID, group.ID, mem.SuperRole.Permissions)
		http.Error(wr, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	} else {
		roleID, e := strconv.ParseUint(ps.ByName("roleID"), 10, 64)
		if e != nil {
			log.Print(e.Error())
		}
		var role *core.Role
		for i := range group.Roles {
			if group.Roles[i].ID == roleID {
				role = &group.Roles[i]
				break
			}
		}
		if role == nil {
			log.Printf("UserID: %v, GroupID: %v, UpdateRoleForm: Role not found", user.ID, group.ID)
			http.Error(wr, "An unexpected error occurred", http.StatusInternalServerError)
			return
		}
		s.groups.GetMemberships(role)
		s.updateRoleInternal(wr, user, group, role, nil)
	}

}
func (s *viewService) UpdateRoleFail(wr http.ResponseWriter, user *core.User,
	group *core.Group, role *core.Role, msg *messages.Group) {
	s.updateRoleInternal(wr, user, group, role, msg)
}

func (s *viewService) updateRoleInternal(wr http.ResponseWriter, user *core.User,
	group *core.Group, role *core.Role, msg *messages.Group) {
	model := struct {
		User  *core.User
		Group *core.Group
		Role  *core.Role
		Msg   *messages.Group
	}{
		User:  user,
		Group: group,
		Role:  role,
		Msg:   msg,
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
