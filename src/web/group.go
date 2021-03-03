package web

import (
	"chores-suck/core"
	"chores-suck/web/messages"
	"errors"
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
	EditGroup(wr http.ResponseWriter, req *http.Request, ps httprouter.Params, uid uint64, group *core.Group)
	GroupAccess(handler func(wr http.ResponseWriter, req *http.Request,
		ps httprouter.Params, uid uint64, group *core.Group)) authParamHandle
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

func (s *groupService) RemoveUser(p Params) {
}

func (s *groupService) EditGroup(wr http.ResponseWriter, req *http.Request,
	ps httprouter.Params, uid uint64, group *core.Group) {
	user := core.User{ID: uid}
	e := s.us.GetUserByID(&user)
	if e != nil {
		http.Error(wr, e.Error(), http.StatusInternalServerError)
		return
	}
	editGroup := false
	for _, r := range group.Roles {
		for _, u := range r.Users {
			if u.ID == user.ID && r.Can(core.EditGroup) {
				editGroup = true
				break
			}
		}
	}
	if !editGroup {
		http.Error(wr, "You do not have permission to edit this group", http.StatusUnauthorized)
		return
	}
	groupName := req.PostFormValue("groupname")
	group.Name = groupName
	msg := messages.Group{}
	if validateGroupName(groupName, &msg) {
		e = s.gs.UpdateGroup(group)
		if e != nil {
			handleError(internalError(e), wr)
		}
	} else {
		s.vs.EditGroupFail(wr, req, &user, group, &msg)
	}
}

func (s *groupService) GroupAccess(handler func(wr http.ResponseWriter, req *http.Request,
	ps httprouter.Params, uid uint64, group *core.Group)) authParamHandle {
	return func(wr http.ResponseWriter, req *http.Request, ps httprouter.Params, uid uint64) {
		groupID, e := strconv.ParseUint(ps.ByName("id"), 10, 64)
		if e != nil {
			http.Error(wr, e.Error(), http.StatusInternalServerError)
		}
		group := core.Group{ID: groupID}
		e = s.gs.GetGroup(&group)

		edit, e := s.gs.CanEdit(&group, uid)
		if e != nil {
			handleError(e, wr)
			return
		} else if !edit {
			http.Error(wr, "You do not have permission to edit this group.", http.StatusUnauthorized)
			return
		}

		handler(wr, req, ps, uid, &group)
	}
}
