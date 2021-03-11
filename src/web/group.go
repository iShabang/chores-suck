package web

import (
	"chores-suck/core"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

var (
	// ErrInvalidFormData occurs when invalid form data is received when submitting a post request for groups
	ErrInvalidFormData = errors.New("Invalid form input")
)

type GroupService interface {
	CreateGroup(wr http.ResponseWriter, req *http.Request, uid uint64)
	EditGroup(wr http.ResponseWriter, req *http.Request, ps httprouter.Params, user *core.User, group *core.Group)
	DeleteMember(wr http.ResponseWriter, req *http.Request, ps httprouter.Params, user *core.User, group *core.Group)
	AddMember(wr http.ResponseWriter, req *http.Request, ps httprouter.Params, user *core.User, group *core.Group)
	AddRole(wr http.ResponseWriter, req *http.Request, ps httprouter.Params, user *core.User, group *core.Group)
	UpdateRole(wr http.ResponseWriter, req *http.Request, ps httprouter.Params, user *core.User, group *core.Group)
	GroupAccess(handler func(wr http.ResponseWriter, req *http.Request,
		ps httprouter.Params, user *core.User, group *core.Group)) authParamHandle
}

type groupService struct {
	gs core.GroupService
	us core.UserService
	vs ViewService
}

func NewGroupService(g core.GroupService, u core.UserService, v ViewService) GroupService {
	return &groupService{
		gs: g,
		us: u,
		vs: v,
	}
}

func (s *groupService) CreateGroup(wr http.ResponseWriter, req *http.Request, uid uint64) {
	groupName := req.PostFormValue("groupname")
	user := core.User{ID: uid}
	e := s.us.GetUserByID(&user)
	if e != nil {
		handleError(internalError(e), wr)
		return
	}
	if e = validateGroupName(groupName); e != nil {
		SetFlash(wr, "nameError", []byte(e.Error()))
		http.Redirect(wr, req, "/new/group", 302)
		return
	}
	e = s.gs.CreateGroup(groupName, &user)
	if e != nil {
		handleError(internalError(e), wr)
		return
	}
	http.Redirect(wr, req, "/dashboard", 302)
}

func (s *groupService) EditGroup(wr http.ResponseWriter, req *http.Request,
	ps httprouter.Params, user *core.User, group *core.Group) {
	groupName := req.PostFormValue("groupname")
	if e := validateGroupName(groupName); e != nil {
		SetFlash(wr, "nameError", []byte(e.Error()))
	} else {
		group.Name = groupName
		e := s.gs.UpdateGroup(group, user)
		if e != nil {
			SetFlash(wr, "genError", []byte("An unexpected error occurred"))
		}
	}
	http.Redirect(wr, req, fmt.Sprintf("/groups/%v", group.ID), 302)
}

func (s *groupService) DeleteMember(wr http.ResponseWriter, req *http.Request,
	ps httprouter.Params, user *core.User, group *core.Group) {
	userID, e := strconv.ParseUint(ps.ByName("userID"), 10, 64)
	msg := ""
	if e != nil {
		msg = "Invalid request"
	} else {
		delUser := core.User{ID: userID}
		delMem := core.Membership{User: &delUser, Group: group}
		if e = s.gs.DeleteMember(&delMem, user); e != nil {
			msg = e.Error()
		}
	}
	if msg != "" {
		SetFlash(wr, "memError", []byte(msg))
	}
	http.Redirect(wr, req, fmt.Sprintf("/groups/%v", group.ID), 302)
}

func (s *groupService) AddMember(wr http.ResponseWriter, req *http.Request,
	ps httprouter.Params, user *core.User, group *core.Group) {
	msg := ""
	uname := req.PostFormValue("username")
	userNew := core.User{Username: uname}
	if e := s.us.GetUserByName(&userNew); e != nil {
		msg = "User not found"
	} else {
		memNew := core.Membership{User: &userNew, Group: group}
		if e := s.gs.AddMember(&memNew, user); e != nil {
			msg = e.Error()
		}
	}
	if msg != "" {
		SetFlash(wr, "memError", []byte(msg))
	}
	url := fmt.Sprintf("/groups/%v", group.ID)
	http.Redirect(wr, req, url, 302)
}

func (s *groupService) AddRole(wr http.ResponseWriter, req *http.Request,
	ps httprouter.Params, user *core.User, group *core.Group) {
	var msg string
	name := req.PostFormValue("name")
	editGroup := req.PostFormValue("editgroup")
	editChores := req.PostFormValue("editchores")
	editMembers := req.PostFormValue("editmembers")
	editRoles := req.PostFormValue("editroles")
	getsChores := req.PostFormValue("getsChores")
	if e := validateGroupName(name); e != nil {
		msg = e.Error()
	}
	if msg == "" {
		role := core.Role{Name: name, GetsChores: getsChores == "true", Group: group}
		role.Set(core.EditGroup, editGroup == "true")
		role.Set(core.EditChores, editChores == "true")
		role.Set(core.EditMembers, editMembers == "true")
		role.Set(core.EditRoles, editRoles == "true")
		if e := s.gs.AddRole(&role, user); e != nil {
			msg = e.Error()
		}
	}
	if msg != "" {
		SetFlash(wr, "genError", []byte(msg))
		url := fmt.Sprintf("/groups/%v/add/role", group.ID)
		http.Redirect(wr, req, url, 302)
	} else {
		url := fmt.Sprintf("/groups/%v", group.ID)
		http.Redirect(wr, req, url, 302)
	}
}

func (s *groupService) UpdateRole(wr http.ResponseWriter, req *http.Request,
	ps httprouter.Params, user *core.User, group *core.Group) {
	var msg string
	name := req.PostFormValue("rolename")
	editMem := req.PostFormValue("editmembers") == "true"
	editChores := req.PostFormValue("editchores") == "true"
	editGroup := req.PostFormValue("editgroup") == "true"
	editRoles := req.PostFormValue("editroles") == "true"
	getsChores := req.PostFormValue("getschores") == "true"
	roleID, e := strconv.ParseUint(ps.ByName("roleID"), 10, 64)
	if e != nil {
		msg = "Invalid request"
	} else if e := validateGroupName(name); e != nil {
		msg = e.Error()
	} else {
		role := core.Role{ID: roleID, Group: group}
		role.Name = name
		role.GetsChores = getsChores
		role.Set(core.EditMembers, editMem)
		role.Set(core.EditChores, editChores)
		role.Set(core.EditGroup, editGroup)
		role.Set(core.EditRoles, editRoles)
		if e := s.gs.UpdateRole(&role, user); e != nil {
			msg = e.Error()
		}
	}
	if msg != "" {
		SetFlash(wr, "genError", []byte(msg))
	}
	url := fmt.Sprintf("/groups/%v/update/role/%v", group.ID, roleID)
	http.Redirect(wr, req, url, 302)
}

func (s *groupService) GroupAccess(handler func(wr http.ResponseWriter, req *http.Request,
	ps httprouter.Params, user *core.User, group *core.Group)) authParamHandle {
	return func(wr http.ResponseWriter, req *http.Request, ps httprouter.Params, uid uint64) {
		groupID, e := strconv.ParseUint(ps.ByName("groupID"), 10, 64)
		if e != nil {
			http.Error(wr, e.Error(), http.StatusInternalServerError)
			return
		}
		group := core.Group{ID: groupID}
		e = s.gs.GetGroup(&group)
		if e != nil {
			http.Error(wr, e.Error(), http.StatusInternalServerError)
			return
		}
		user := core.User{ID: uid}
		e = s.us.GetUserByID(&user)
		if e != nil {
			http.Error(wr, e.Error(), http.StatusInternalServerError)
			return
		}
		var mem *core.Membership
		if e := s.gs.GetMemberships(&group); e != nil {
			log.Printf("Web.GroupAccess: %s", e.Error())
			http.Error(wr, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		mem = group.FindMember(user.ID)
		if e := s.gs.GetRoles(mem); e != nil {
			log.Printf("Web.GroupAccess: %s", e.Error())
			http.Error(wr, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if mem == nil || !mem.SuperRole.CanEdit() {
			var msg string
			if mem == nil {
				msg = "Member not found"
			} else {
				msg = "Member has insufficient privileges"
			}
			log.Printf("UserID: %v, GroupID: %v, GroupAccess: %s", user.ID, group.ID, msg)
			http.Error(wr, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		handler(wr, req, ps, &user, &group)
	}
}
