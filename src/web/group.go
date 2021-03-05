package web

import (
	"chores-suck/core"
	"chores-suck/web/messages"
	"errors"
	"fmt"
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
	GetGroup(group *core.Group) error
	UpdateGroup(group *core.Group) error
	EditGroup(wr http.ResponseWriter, req *http.Request, ps httprouter.Params, user *core.User, group *core.Group)
	DeleteMember(wr http.ResponseWriter, req *http.Request, ps httprouter.Params, user *core.User, group *core.Group)
	AddMember(wr http.ResponseWriter, req *http.Request, ps httprouter.Params, user *core.User, group *core.Group)
	AddRole(wr http.ResponseWriter, req *http.Request, ps httprouter.Params, user *core.User, group *core.Group)
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
	msg := messages.Group{}
	if !validateGroupName(groupName, &msg) {
		s.vs.NewGroupFail(wr, req, &user, &msg)
		return
	}
	e = s.gs.CreateGroup(groupName, &user)
	if e != nil {
		handleError(internalError(e), wr)
		return
	}
	http.Redirect(wr, req, "/dashboard", 302)
}

func (s *groupService) GetGroup(group *core.Group) error {
	e := s.gs.GetGroup(group)
	if e != nil {
		return internalError(e)
	}

	e = s.gs.GetMemberships(group)
	if e != nil {
		return internalError(e)
	}

	e = s.gs.GetRoles(group)
	if e != nil {
		return internalError(e)
	}

	return nil
}

func (s *groupService) UpdateGroup(group *core.Group) error {
	e := s.gs.UpdateGroup(group)
	return StatusError{Code: http.StatusInternalServerError, Err: e}
}

func (s *groupService) EditGroup(wr http.ResponseWriter, req *http.Request,
	ps httprouter.Params, user *core.User, group *core.Group) {
	canEdit := false
	for _, role := range user.Memberships[0].Roles {
		if role.Can(core.EditGroup) {
			canEdit = true
			break
		}
	}
	if !canEdit {
		http.Error(wr, ErrNotAuthorized.Error(), http.StatusUnauthorized)
		return
	}
	groupName := req.PostFormValue("groupname")
	group.Name = groupName
	ok := true
	msg := messages.Group{}
	if validateGroupName(groupName, &msg) {
		e := s.gs.UpdateGroup(group)
		if e != nil {
			msg.General = "There was an unexpected error while updating the name"
			ok = false
		}
	} else {
		ok = false
	}

	if !ok {
		s.vs.EditGroupFail(wr, req, user, group, &msg)
	} else {
		url := fmt.Sprintf("/groups/%v", group.ID)
		http.Redirect(wr, req, url, 302)
	}
}

func (s *groupService) DeleteMember(wr http.ResponseWriter, req *http.Request,
	ps httprouter.Params, user *core.User, group *core.Group) {
	canDelete := false
	for _, v := range user.Memberships[0].Roles {
		if v.Can(core.EditMembers) {
			canDelete = true
			break
		}
	}
	if !canDelete {
		http.Error(wr, ErrNotAuthorized.Error(), http.StatusUnauthorized)
		return
	}
	userID, e := strconv.ParseUint(ps.ByName("userID"), 10, 64)
	delUser := core.User{ID: userID}
	e = s.us.GetUserByID(&delUser)
	if e != nil {
		handleError(internalError(e), wr)
		return
	}
	delMem := core.Membership{User: &delUser, Group: group}
	e = s.gs.GetRoles(&delMem)
	if e != nil {
		handleError(internalError(e), wr)
		return
	}
	ok := true
	msg := messages.Group{}
	e = s.gs.DeleteMember(&delMem)
	if e == nil {
		for i, m := range group.Memberships {
			if m.User.ID == delUser.ID {
				end := len(group.Memberships) - 1
				group.Memberships[i] = group.Memberships[end]
				group.Memberships = group.Memberships[:end]
			}
		}
	} else {
		ok = false
		msg.Member = e.Error()
	}
	if !ok {
		s.vs.EditGroupFail(wr, req, user, group, &msg)
	} else {
		url := fmt.Sprintf("/groups/%v", group.ID)
		http.Redirect(wr, req, url, 302)
	}
}

func (s *groupService) AddMember(wr http.ResponseWriter, req *http.Request,
	ps httprouter.Params, user *core.User, group *core.Group) {
	mem := core.Membership{}
	for _, v := range user.Memberships {
		if v.Group.ID == group.ID {
			mem = v
			break
		}
	}
	canAdd := false
	for _, v := range mem.Roles {
		if v.Can(core.EditMembers) {
			canAdd = true
		}
	}
	ok := true
	msg := messages.Group{}
	if !canAdd {
		ok = false
		msg.Member = "You don't have permission to edit members"
	} else {
		uname := req.PostFormValue("username")
		userNew := core.User{Username: uname}
		s.us.GetUserByName(&userNew)
		memNew := core.Membership{User: &userNew, Group: group}
		if s.gs.AddMember(&memNew) == nil {
			group.Memberships = append(group.Memberships, memNew)
		} else {
			ok = false
			msg.Member = "Failed to add new member"
		}
	}
	if !ok {
		s.vs.EditGroupFail(wr, req, user, group, &msg)
	} else {
		url := fmt.Sprintf("/groups/%v", group.ID)
		http.Redirect(wr, req, url, 302)
	}
}

func (s *groupService) AddRole(wr http.ResponseWriter, req *http.Request,
	ps httprouter.Params, user *core.User, group *core.Group) {
	mem := core.Membership{}
	for _, m := range user.Memberships {
		if m.Group.ID == group.ID {
			mem = m
			break
		}
	}
	editRole := false
	for _, r := range mem.Roles {
		if r.Can(core.EditRoles) {
			editRole = true
			break
		}
	}
	if !editRole {
		http.Error(wr, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	name := req.PostFormValue("name")
	editGroup := req.PostFormValue("editgroup")
	editChores := req.PostFormValue("editchores")
	editMembers := req.PostFormValue("editmembers")
	editRoles := req.PostFormValue("editroles")
	getsChores := req.PostFormValue("getsChores")
	msg := messages.Group{}
	ok := validateGroupName(name, &msg)
	if ok {
		s.gs.GetRoles(group)
		for _, r := range group.Roles {
			if r.Name == name {
				ok = false
				msg.Name = "Role name already exists!"
			}
		}
	}
	role := core.Role{Name: name, GetsChores: getsChores == "true", Group: group}
	if ok {
		role.Set(core.EditGroup, editGroup == "true")
		role.Set(core.EditChores, editChores == "true")
		role.Set(core.EditMembers, editMembers == "true")
		role.Set(core.EditRoles, editRoles == "true")
		if s.gs.AddRole(&role) != nil {
			ok = false
			msg.General = "Unable to add role due to an unexpected error"
		}
	}
	if !ok {
		s.vs.NewRoleFail(wr, user, group, &msg)
	} else {
		url := fmt.Sprintf("/groups/%v", group.ID)
		http.Redirect(wr, req, url, 302)
	}
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
		edit := s.gs.CanEdit(&group, &user)
		if !edit {
			http.Error(wr, "You do not have permission to edit this group.", http.StatusUnauthorized)
			return
		}

		handler(wr, req, ps, &user, &group)
	}
}
