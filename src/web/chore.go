package web

import (
	"chores-suck/core"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type ChoreService interface {
	Create(http.ResponseWriter, *http.Request, httprouter.Params, *core.User, *core.Group)
	Update(http.ResponseWriter, *http.Request, *core.User, *core.Chore)
	ChoreMW(handler func(http.ResponseWriter, *http.Request, *core.User, *core.Chore)) authParamHandle
}

type choreService struct {
	cs core.ChoreService
	gs core.GroupService
	us core.UserService
}

func NewChoreService(c core.ChoreService, g core.GroupService, u core.UserService) ChoreService {
	return &choreService{
		cs: c,
		gs: g,
		us: u,
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

func (s *choreService) Update(wr http.ResponseWriter, req *http.Request, us *core.User, ch *core.Chore) {
	if req.PostFormValue("submit_1") != "" {
		s.update(wr, req, us, ch)
	} else if req.PostFormValue("submit_2") != "" {
		s.delete(wr, req, us, ch)
	} else {
		http.Error(wr, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
}

func (s *choreService) delete(wr http.ResponseWriter, req *http.Request, us *core.User, ch *core.Chore) {
	if e := s.cs.Delete(ch); e != nil {
		SetFlash(wr, "genError", []byte(e.Error()))
		http.Redirect(wr, req, fmt.Sprintf("/chores/update/%v", ch.ID), 302)
		return
	}
	http.Redirect(wr, req, fmt.Sprintf("/groups/update/%v", ch.Group.ID), 302)
}

func (s *choreService) update(wr http.ResponseWriter, req *http.Request, us *core.User, ch *core.Chore) {
	var msg string
	choreName := req.PostFormValue("chore_name")
	choreDesc := req.PostFormValue("chore_desc")
	choreDur, e := strconv.Atoi(req.PostFormValue("chore_dur"))
	newChore := core.Chore{ID: ch.ID, Name: choreName, Description: choreDesc, Duration: choreDur}
	if e != nil {
		http.Error(wr, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if e := validateGroupName(choreName); e != nil {
		msg = e.Error()
	} else if e := s.cs.Update(ch, &newChore); e != nil {
		msg = e.Error()
	}
	if msg != "" {
		SetFlash(wr, "genError", []byte(msg))
	}
	http.Redirect(wr, req, fmt.Sprintf("/chores/update/%v", ch.ID), 302)
}

func (s *choreService) ChoreMW(handler func(http.ResponseWriter, *http.Request, *core.User, *core.Chore)) authParamHandle {
	return func(wr http.ResponseWriter, req *http.Request, ps httprouter.Params, userID uint64) {
		//Get Chore
		choreID, e := strconv.ParseUint(ps.ByName("choreID"), 10, 64)
		if e != nil {
			http.Error(wr, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		chore := core.Chore{ID: choreID}
		if e = s.cs.GetChore(&chore); e != nil {
			//Internal server error
			log.Printf("ChoreMW: Failed to grab chore: %s", e.Error())
			http.Error(wr, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		} else if chore.Name == "" {
			//Not found
			http.Error(wr, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		//Get Group
		if e = s.gs.GetGroup(chore.Group); e != nil {
			//internal server error
			log.Printf("ChoreMW: Failed to get group: %s", e.Error())
			http.Error(wr, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		//Get Group memberships
		if e = s.gs.GetMemberships(chore.Group); e != nil {
			//internal server error
			log.Printf("ChoreMW: Failed to get members: %s", e.Error())
			http.Error(wr, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		//Get user membership
		mem := chore.Group.FindMember(userID)
		if mem == nil {
			//Unauthorized
			http.Error(wr, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		if e = s.gs.GetRoles(mem); e != nil {
			http.Error(wr, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		//Check if member can edit chores
		if !mem.SuperRole.Can(core.EditChores) {
			//Unauthorized
			http.Error(wr, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		//Get User
		user := core.User{ID: userID}
		if e = s.us.GetUserByID(&user); e != nil {
			//Internal server error
			log.Printf("ChoreMW: Failed to get user: %s", e.Error())
			http.Error(wr, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		handler(wr, req, &user, &chore)
	}
}
