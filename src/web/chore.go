package web

import (
	"chores-suck/core"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type ChoreService interface {
	Create(http.ResponseWriter, *http.Request, httprouter.Params, *core.User, *core.Group)
}

type choreService struct {
	cs core.ChoreService
}

func NewChoreService(c core.ChoreService) ChoreService {
	return &choreService{
		cs: c,
	}
}

func (s *choreService) Create(wr http.ResponseWriter, req *http.Request,
	ps httprouter.Params, user *core.User, group *core.Group) {
	var msg string
	mem := group.FindMember(user.ID)
	if !mem.SuperRole.Can(core.EditChores) {
		http.Error(wr, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	choreName := req.PostFormValue("chore_name")
	choreDesc := req.PostFormValue("chore_desc")
	choreTime, e := strconv.Atoi(req.PostFormValue("chore_dur"))
	chore := core.Chore{Group: group, Name: choreName, Description: choreDesc, Duration: choreTime}
	if e != nil {
		http.Error(wr, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if e := validateGroupName(choreName); e != nil {
		msg = e.Error()
	} else if e := s.cs.Create(&chore); e != nil {
		msg = e.Error()
	}
	if msg != "" {
		SetFlash(wr, "genError", []byte(msg))
	}
	url := fmt.Sprintf("/chores/create/%v", group.ID)
	http.Redirect(wr, req, url, 302)
}
