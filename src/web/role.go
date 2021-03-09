package web

import (
	"chores-suck/core"
	"chores-suck/web/messages"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type RoleService interface {
	RemoveMember(wr http.ResponseWriter, req *http.Request,
		ps httprouter.Params, user *core.User, group *core.Group)
	AddMember(wr http.ResponseWriter, req *http.Request,
		ps httprouter.Params, user *core.User, group *core.Group)
}

type roleService struct {
	gs core.GroupService
	rs core.RoleService
	us core.UserService
	vs ViewService
}

func NewRoleService(g core.GroupService, r core.RoleService, u core.UserService, v ViewService) RoleService {
	return &roleService{
		gs: g,
		rs: r,
		us: u,
		vs: v,
	}
}

func (s *roleService) RemoveMember(wr http.ResponseWriter, req *http.Request,
	ps httprouter.Params, user *core.User, group *core.Group) {
	//check if they can edit roles
	var mem *core.Membership
	for i := range group.Memberships {
		if group.Memberships[i].User.ID == user.ID {
			mem = &group.Memberships[i]
			s.gs.GetRoles(mem)
		}
	}
	if !mem.SuperRole.Can(core.EditRoles) {
		log.Printf("UserID: %v, GroupID: %v, RoleService.RemoveUser: Insufficient privileges", user.ID, group.ID)
		http.Error(wr, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	//Check if the role exists in the group
	roleID, _ := strconv.ParseUint(ps.ByName("roleID"), 10, 64)
	s.gs.GetRoles(group)
	role := group.FindRole(roleID)
	if role == nil {
		log.Printf("RoleService.RemoveUser: Role not found")
		http.Error(wr, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	msg := messages.Group{}
	ok := true
	if role.Name == "Owner" {
		msg.General = "Cannot remove owner"
		ok = false
	} else {
		delID, _ := strconv.ParseUint(ps.ByName("userID"), 10, 64)
		//remove the user from the role
		if s.rs.RemoveMember(role.ID, delID) != nil {
			msg.General = "An unexpected error occurred"
			ok = false
		}
	}
	if !ok {
		s.gs.GetMemberships(role)
		s.vs.UpdateRoleFail(wr, user, group, role, &msg)
	} else {
		url := fmt.Sprintf("/groups/%v/update/role/%v", group.ID, role.ID)
		http.Redirect(wr, req, url, 302)
	}
}

func (s *roleService) AddMember(wr http.ResponseWriter, req *http.Request,
	ps httprouter.Params, user *core.User, group *core.Group) {
	mem := group.FindMember(user.ID)
	if !mem.SuperRole.Can(core.EditMembers) {
		log.Printf("RoleService.AddMember: Error: Insufficient permissions. Perm: %v", mem.SuperRole.Permissions)
		http.Error(wr, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	roleID, _ := strconv.ParseUint(ps.ByName("roleID"), 10, 64)
	s.gs.GetRoles(group)
	role := group.FindRole(roleID)
	if role == nil {
		log.Printf("RoleService.AddMember: Error: Role not found")
		http.Error(wr, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	ok := true
	username := req.PostFormValue("username")
	newUser := core.User{Username: username}
	msg := messages.Group{}
	if e := s.us.GetUserByName(&newUser); e != nil {
		log.Printf("RoleService.AddMember: Error getting user: %s", e.Error())
		msg.General = "Unable to add user"
		ok = false
	} else {
		newMem := group.FindMember(newUser.ID)
		if newMem == nil {
			log.Printf("RoleService.AddMember: User is not a member of this group")
			msg.General = "Unable to add user"
			ok = false
		} else {
			if e := s.rs.AddMember(role.ID, newUser.ID); e != nil {
				log.Printf("RoleService.AddMember: Error adding member: %s", e.Error())
				msg.General = "Failed to add member due to an unexpected error"
				ok = false
			}
		}
	}
	if !ok {
		role := core.Role{ID: roleID}
		s.gs.GetMemberships(&role)
		s.vs.UpdateRoleFail(wr, user, group, &role, &msg)
	} else {
		url := fmt.Sprintf("/groups/%v/update/role/%v", group.ID, roleID)
		http.Redirect(wr, req, url, 302)
	}
}
