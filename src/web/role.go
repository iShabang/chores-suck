package web

import (
	"chores-suck/core"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type RoleService interface {
	Update(wr http.ResponseWriter, req *http.Request, ps httprouter.Params, user *core.User, role *core.Role)
	RoleMW(handler func(wr http.ResponseWriter, req *http.Request,
		ps httprouter.Params, user *core.User, role *core.Role)) authParamHandle
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

func (s *roleService) Update(wr http.ResponseWriter, req *http.Request,
	ps httprouter.Params, user *core.User, role *core.Role) {
	if submit := req.PostFormValue("submit_1"); submit != "" {
		s.updateRole(wr, req, user, role)
	} else if submit = req.PostFormValue("submit_2"); submit != "" {
		s.delMember(wr, req, user, role)
	} else if submit = req.PostFormValue("submit_3"); submit != "" {
		s.addMember(wr, req, user, role)
	}
}

func (s *roleService) updateRole(wr http.ResponseWriter, req *http.Request,
	user *core.User, role *core.Role) {
	var msg string
	name := req.PostFormValue("rolename")
	editMem := req.PostFormValue("editmembers") == "true"
	editChores := req.PostFormValue("editchores") == "true"
	editGroup := req.PostFormValue("editgroup") == "true"
	editRoles := req.PostFormValue("editroles") == "true"
	getsChores := req.PostFormValue("getschores") == "true"
	if e := validateGroupName(name); e != nil {
		msg = e.Error()
	} else {
		newRole := core.Role{ID: role.ID, Group: role.Group}
		newRole.Name = name
		newRole.GetsChores = getsChores
		newRole.Set(core.EditMembers, editMem)
		newRole.Set(core.EditChores, editChores)
		newRole.Set(core.EditGroup, editGroup)
		newRole.Set(core.EditRoles, editRoles)
		if e := s.rs.Update(role, &newRole, user); e != nil {
			msg = e.Error()
		}
	}
	if msg != "" {
		SetFlash(wr, "genError", []byte(msg))
	}
	url := fmt.Sprintf("/roles/update/%v", role.ID)
	http.Redirect(wr, req, url, 302)
}

func (s *roleService) addMember(wr http.ResponseWriter, req *http.Request,
	user *core.User, role *core.Role) {
	var msg string
	username := req.PostFormValue("username")
	if e := s.rs.AddMember(role, username, user); e != nil {
		msg = e.Error()
	}
	if msg != "" {
		SetFlash(wr, "genError", []byte(msg))
	}
	url := fmt.Sprintf("/roles/update/%v", role.ID)
	http.Redirect(wr, req, url, 302)
}

func (s *roleService) delMember(wr http.ResponseWriter, req *http.Request,
	user *core.User, role *core.Role) {
	var msg string
	delID, e := strconv.ParseUint(req.PostFormValue("user_id"), 10, 64)
	if e != nil {
		msg = "Invalid request"
	}
	if e := s.rs.RemoveMember(role, delID, user); e != nil {
		msg = e.Error()
	}
	if msg != "" {
		SetFlash(wr, "genError", []byte(msg))
	}
	url := fmt.Sprintf("/roles/update/%v", role.ID)
	http.Redirect(wr, req, url, 302)
}

/***************************************************************
MIDDLEWARE
***************************************************************/
func (s *roleService) RoleMW(handler func(wr http.ResponseWriter, req *http.Request,
	ps httprouter.Params, user *core.User, role *core.Role)) authParamHandle {
	return func(wr http.ResponseWriter, req *http.Request, ps httprouter.Params, uid uint64) {
		//Get the role
		roleID, e := strconv.ParseUint(ps.ByName("roleID"), 10, 64)
		if e != nil {
			http.Error(wr, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		role := core.Role{ID: roleID}
		if e = s.rs.GetRole(&role); e != nil {
			http.Error(wr, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		} else if role.Name == "" {
			http.Error(wr, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		//Get the group
		if e = s.gs.GetGroup(role.Group); e != nil {
			http.Error(wr, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		//Get group members
		if e = s.gs.GetMemberships(role.Group); e != nil {
			http.Error(wr, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		//Check if user is a member
		mem := role.Group.FindMember(uid)
		if mem == nil {
			http.Error(wr, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		//Get the user
		user := core.User{ID: uid}
		if e = s.us.GetUserByID(&user); e != nil {
			http.Error(wr, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		//Get member roles
		if e = s.gs.GetRoles(mem); e != nil {
			http.Error(wr, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		//Check if member has permission to update roles
		if !mem.SuperRole.Can(core.EditRoles) {
			http.Error(wr, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		handler(wr, req, ps, &user, &role)
	}
}
