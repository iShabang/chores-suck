package web

import (
	"chores-suck/core"
	"chores-suck/web/messages"
	"errors"
	"log"
	"net/http"
)

var (
	// ErrInvalidFormData occurs when invalid form data is received when submitting a post request for groups
	ErrInvalidFormData = errors.New("Invalid form input")
)

type GroupService interface {
	CreateGroup(wr http.ResponseWriter, req *http.Request, uid uint64, msg *messages.CreateGroup) error
	CanEdit(group *core.Group, uid uint64) (bool, error)
	GetGroup(group *core.Group) error
	UpdateGroup(group *core.Group) error
}

type groupService struct {
	gs core.GroupService
	us core.UserService
}

func NewGroupService(g core.GroupService, u core.UserService) GroupService {
	return &groupService{
		gs: g,
		us: u,
	}
}

func (s *groupService) CreateGroup(wr http.ResponseWriter, req *http.Request, uid uint64, msg *messages.CreateGroup) error {
	groupName := req.PostFormValue("groupname")

	if !validateGroupName(groupName, msg) {
		return ErrInvalidFormData
	}

	user := core.User{ID: uid}
	e := s.us.GetUserByID(&user)
	if e != nil {
		return StatusError{Code: http.StatusInternalServerError, Err: e}
	}

	e = s.gs.CreateGroup(groupName, &user)
	if e != nil {
		return StatusError{Code: http.StatusInternalServerError, Err: e}
	}

	return nil

}

func (s *groupService) CanEdit(group *core.Group, uid uint64) (bool, error) {

	isMember := false
	var mem core.Membership
	for _, v := range group.Memberships {
		if v.User.ID == uid {
			isMember = true
			mem = v
			break
		}
	}
	if !isMember {
		log.Print("User is not a member of the group to edit")
		return false, nil
	}

	e := s.gs.GetRoles(&mem)
	if e != nil {
		return false, internalError(e)
	}

	canEdit := false
	for _, v := range mem.Roles {
		if v.CanEdit() {
			canEdit = true
			break
		}
	}
	if !canEdit {
		log.Print("User has insufficient privileges to edit this group")
		return false, nil
	}

	return true, nil
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
