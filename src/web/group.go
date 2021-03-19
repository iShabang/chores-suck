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
	UpdateGroup(wr http.ResponseWriter, req *http.Request, ps httprouter.Params, user *core.User, group *core.Group)
	AddRole(wr http.ResponseWriter, req *http.Request, ps httprouter.Params, user *core.User, group *core.Group)
	GroupAccess(handler func(wr http.ResponseWriter, req *http.Request,
		ps httprouter.Params, user *core.User, group *core.Group)) authParamHandle
}

type groupService struct {
	gs core.GroupService
	us core.UserService
	cs core.ChoreService
}

func NewGroupService(g core.GroupService, u core.UserService, c core.ChoreService) GroupService {
	return &groupService{
		gs: g,
		us: u,
		cs: c,
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
		http.Redirect(wr, req, "/groups/create", 302)
		return
	}
	e = s.gs.CreateGroup(groupName, &user)
	if e != nil {
		handleError(internalError(e), wr)
		return
	}
	http.Redirect(wr, req, "/dashboard", 302)
}

func (s *groupService) UpdateGroup(wr http.ResponseWriter, req *http.Request,
	ps httprouter.Params, user *core.User, group *core.Group) {
	if submit := req.PostFormValue("submit_1"); submit != "" {
		s.updateName(wr, req, ps, user, group)
	} else if submit := req.PostFormValue("submit_2"); submit != "" {
		s.delMember(wr, req, ps, user, group)
	} else if submit := req.PostFormValue("submit_3"); submit != "" {
		s.addMember(wr, req, ps, user, group)
	} else if submit := req.PostFormValue("submit_4"); submit != "" {
		s.random(wr, req, ps, user, group)
	} else {
		log.Print("Failed to find submit value")
	}
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
		url := fmt.Sprintf("/roles/create/%v", group.ID)
		http.Redirect(wr, req, url, 302)
	} else {
		url := fmt.Sprintf("/groups/update/%v", group.ID)
		http.Redirect(wr, req, url, 302)
	}
}

func (s *groupService) updateName(wr http.ResponseWriter, req *http.Request, ps httprouter.Params, user *core.User, group *core.Group) {
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
	http.Redirect(wr, req, fmt.Sprintf("/groups/update/%v", group.ID), 302)
}

func (s *groupService) addMember(wr http.ResponseWriter, req *http.Request, ps httprouter.Params, user *core.User, group *core.Group) {
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
	url := fmt.Sprintf("/groups/update/%v", group.ID)
	http.Redirect(wr, req, url, 302)
}

func (s *groupService) delMember(wr http.ResponseWriter, req *http.Request, ps httprouter.Params, user *core.User, group *core.Group) {
	userID, e := strconv.ParseUint(req.PostFormValue("user_id"), 10, 64)
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
	http.Redirect(wr, req, fmt.Sprintf("/groups/update/%v", group.ID), 302)
}

func (s *groupService) random(wr http.ResponseWriter, req *http.Request, _ httprouter.Params, u *core.User, g *core.Group) {
	var msg string
	if e := s.gs.GetChores(g); e != nil {
		msg = e.Error()
	} else if e := s.cs.Randomize(g); e != nil {
		msg = e.Error()
	}
	if msg != "" {
		SetFlash(wr, "choreError", []byte(msg))
	}
	http.Redirect(wr, req, fmt.Sprintf("/groups/update/%v", g.ID), 302)
}

func (s *groupService) GroupAccess(handler func(wr http.ResponseWriter, req *http.Request,
	ps httprouter.Params, user *core.User, group *core.Group)) authParamHandle {
	return func(wr http.ResponseWriter, req *http.Request, ps httprouter.Params, uid uint64) {
		var groupID uint64
		groupID, e := strconv.ParseUint(ps.ByName("groupID"), 10, 64)
		if e != nil {
			if groupID, e = strconv.ParseUint(req.FormValue("group_id"), 10, 64); e != nil {
				http.Error(wr, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
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
