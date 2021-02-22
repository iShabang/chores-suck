package groups

import (
	"chores-suck/core/groups"
	"chores-suck/core/types"
	"chores-suck/core/users"
	ce "chores-suck/rest/errors"
	"chores-suck/rest/messages"
	"errors"
	"net/http"
)

var (
	// ErrInvalidFormData occurs when invalid form data is received when submitting a post request for groups
	ErrInvalidFormData = errors.New("Invalid form input")
)

type Service interface {
	CreateGroup(wr http.ResponseWriter, req *http.Request, uid uint64, msg *messages.CreateGroup) error
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
