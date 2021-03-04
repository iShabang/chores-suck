package core

import (
	"errors"
	"log"
	"time"
)

type GroupRepository interface {
	CreateGroup(group *Group) error
	CreateRole(role *Role) error
	CreateRoleAssignment(roleID uint64, userID uint64) error
	CreateMembership(mem *Membership) error
	GetGroupByID(group *Group) error
	GetMemberships(t interface{}) error
	GetMembership(mem *Membership) error
	GetRoles(t interface{}) error
	UpdateGroup(group *Group) error
	DeleteMember(mem *Membership) error
}

type GroupService interface {
	// CreateGroup creates a new group with the default roles (owner, admin, default) and creates a new membership
	// for the owner (passed in user)
	CreateGroup(name string, user *User) error
	GetGroup(group *Group) error
	GetMemberships(group *Group) error
	GetMembership(mem *Membership) error
	GetRoles(t interface{}) error
	UpdateGroup(group *Group) error
	CanEdit(group *Group, user *User) (bool, error)
	DeleteMember(mem *Membership) error
	AddMember(mem *Membership) error
}

type groupService struct {
	repo GroupRepository
}

func NewGroupService(r GroupRepository) GroupService {
	return &groupService{
		repo: r,
	}
}

func (s *groupService) CreateGroup(name string, user *User) error {
	group := Group{Name: name}
	e := s.repo.CreateGroup(&group)
	if e != nil {
		return e
	}
	users := make([]User, 1)
	users[0] = *user

	owner := Role{Name: "Owner", Group: &group, Users: users}
	owner.SetAll(true)
	e = s.repo.CreateRole(&owner)
	if e != nil {
		return e
	}

	e = s.repo.CreateRoleAssignment(owner.ID, user.ID)
	if e != nil {
		return e
	}

	admin := Role{Name: "Admin", Group: &group, Users: users}
	admin.SetAll(true)
	e = s.repo.CreateRole(&admin)
	if e != nil {
		return e
	}

	e = s.repo.CreateRoleAssignment(admin.ID, user.ID)
	if e != nil {
		return e
	}

	def := Role{Name: "Default", Group: &group, Users: users, GetsChores: true}
	e = s.repo.CreateRole(&def)
	if e != nil {
		return e
	}

	e = s.repo.CreateRoleAssignment(def.ID, user.ID)
	if e != nil {
		return e
	}

	mem := Membership{JoinedAt: time.Now().UTC(), User: user, Group: &group}
	e = s.repo.CreateMembership(&mem)
	if e != nil {
		return e
	}

	return nil
}

func (s *groupService) GetGroup(group *Group) error {
	e := s.repo.GetGroupByID(group)
	return e
}

func (s *groupService) GetMemberships(group *Group) error {
	e := s.repo.GetMemberships(group)
	return e
}

func (s *groupService) GetRoles(t interface{}) error {
	e := s.repo.GetRoles(t)
	return e
}

func (s *groupService) UpdateGroup(group *Group) error {
	e := s.repo.UpdateGroup(group)
	return e
}

func (s *groupService) CanEdit(group *Group, user *User) (bool, error) {
	isMember := false
	var mem Membership
	for _, v := range group.Memberships {
		if v.User.ID == user.ID {
			isMember = true
			mem = v
			break
		}
	}
	if !isMember {
		log.Print("User is not a member of the group to edit")
		return false, nil
	}
	user.Memberships = append(user.Memberships, mem)
	e := s.GetRoles(&user.Memberships[0])
	if e != nil {
		return false, e
	}

	canEdit := false
	for _, v := range user.Memberships[0].Roles {
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

func (s *groupService) DeleteMember(mem *Membership) error {
	for _, v := range mem.Roles {
		if v.Name == "Owner" {
			log.Print("Cannot delete owner")
			return errors.New("Cannot delete owner")
		}
	}
	return s.repo.DeleteMember(mem)
}

func (s *groupService) GetMembership(mem *Membership) error {
	return s.repo.GetMembership(mem)
}

func (s *groupService) AddMember(mem *Membership) error {
	mem.JoinedAt = time.Now().UTC()
	return s.repo.CreateMembership(mem)
}
