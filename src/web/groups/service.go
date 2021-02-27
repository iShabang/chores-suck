package groups

import (
	"chores-suck/core/groups"
	"chores-suck/core/types"
	"chores-suck/core/users"
	ce "chores-suck/web/errors"
	"chores-suck/web/messages"
	"errors"
	"log"
	"net/http"
)

var (
	// ErrInvalidFormData occurs when invalid form data is received when submitting a post request for groups
	ErrInvalidFormData = errors.New("Invalid form input")
)

type Service interface {
	CreateGroup(wr http.ResponseWriter, req *http.Request, uid uint64, msg *messages.CreateGroup) error
	CanEdit(group *types.Group, uid uint64) (bool, error)
	GetGroup(group *types.Group) error
	UpdateGroup(group *types.Group) error
}

type service struct {
	gs groups.Service
	us users.Service
}

func NewService(g groups.Service, u users.Service) Service {
	return &service{
		gs: g,
		us: u,
	}
}

func (s *service) CreateGroup(wr http.ResponseWriter, req *http.Request, uid uint64, msg *messages.CreateGroup) error {
	groupName := req.PostFormValue("groupname")

	if !validateName(groupName, msg) {
		return ErrInvalidFormData
	}

	user := types.User{ID: uid}
	e := s.us.GetUserByID(&user)
	if e != nil {
		return ce.StatusError{Code: http.StatusInternalServerError, Err: e}
	}

	e = s.gs.CreateGroup(groupName, &user)
	if e != nil {
		return ce.StatusError{Code: http.StatusInternalServerError, Err: e}
	}

	return nil

}

func (s *service) CanEdit(group *types.Group, uid uint64) (bool, error) {

	isMember := false
	var mem types.Membership
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

func (s *service) GetGroup(group *types.Group) error {
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

func (s *service) UpdateGroup(group *types.Group) error {
	e := s.gs.UpdateGroup(group)
	return ce.StatusError{Code: http.StatusInternalServerError, Err: e}
}

func internalError(e error) error {
	return ce.StatusError{Code: http.StatusInternalServerError, Err: e}
}
