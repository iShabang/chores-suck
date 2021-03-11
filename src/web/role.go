package web

import (
	"chores-suck/core"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type RoleService interface {
	RemoveMember(wr http.ResponseWriter, req *http.Request,
		ps httprouter.Params, user *core.User, group *core.Group)
	AddMember(wr http.ResponseWriter, req *http.Request,
		ps httprouter.Params, user *core.User, group *core.Group)
}

type roleService struct {
	gs core.GroupService
	rs core.RoleService
	us core.UserService
	vs ViewService
}

func NewRoleService(g core.GroupService, r core.RoleService, u core.UserService, v ViewService) RoleService {
	return &roleService{
		gs: g,
		rs: r,
		us: u,
		vs: v,
	}
}

func (s *roleService) RemoveMember(wr http.ResponseWriter, req *http.Request,
	ps httprouter.Params, user *core.User, group *core.Group) {
	var msg string
	roleID, e := strconv.ParseUint(ps.ByName("roleID"), 10, 64)
	if e != nil {
		msg = "Invalid request"
	}
	delID, e := strconv.ParseUint(ps.ByName("userID"), 10, 64)
	if e != nil {
		msg = "Invalid request"
	}
	role := core.Role{ID: roleID, Group: group}
	if e := s.rs.RemoveMember(&role, delID, user); e != nil {
		msg = e.Error()
	}
	if msg != "" {
		SetFlash(wr, "genError", []byte(msg))
	}
	url := fmt.Sprintf("/groups/%v/update/role/%v", group.ID, role.ID)
	http.Redirect(wr, req, url, 302)
}

func (s *roleService) AddMember(wr http.ResponseWriter, req *http.Request,
	ps httprouter.Params, user *core.User, group *core.Group) {
	var msg string
	roleID, e := strconv.ParseUint(ps.ByName("roleID"), 10, 64)
	if e != nil {
		msg = "Invalid request"
	}
	role := core.Role{ID: roleID, Group: group}
	username := req.PostFormValue("username")
	if e := s.rs.AddMember(&role, username, user); e != nil {
		msg = e.Error()
	}
	if msg != "" {
		SetFlash(wr, "genError", []byte(msg))
	}
	url := fmt.Sprintf("/groups/%v/update/role/%v", group.ID, roleID)
	http.Redirect(wr, req, url, 302)
}
