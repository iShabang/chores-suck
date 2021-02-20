package groups

import (
	"chores-suck/core/groups"
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
}

func (s *service) CreateGroup(wr http.ResponseWriter, req *http.Request, uid uint64, msg *messages.CreateGroup) error {
	groupName := req.PostFormValue("groupName")

	if !validateName(groupName, msg) {
		return ErrInvalidFormData
	}

	e := s.gs.CreateGroup(groupName, uid)
	if e != nil {
		return ce.StatusError{Code: http.StatusInternalServerError, Err: e}
	}

	return nil

}
