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
	GetMemberships(t interface{}) error
	GetMembership(mem *Membership) error
	GetRoles(t interface{}) error
	UpdateGroup(group *Group) error
	CanEdit(group *Group, user *User) bool
	DeleteMember(mem *Membership) error
	AddMember(mem *Membership) error
	AddRole(role *Role) error
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

	mem := Membership{JoinedAt: time.Now().UTC(), User: user, Group: &group}
	e = s.repo.CreateMembership(&mem)
	if e != nil {
		return e
	}
	owner := Role{Name: "Owner", Group: &group}
	owner.SetAll(true)
	e = s.repo.CreateRole(&owner)
	if e != nil {
		return e
	}

	e = s.repo.CreateRoleAssignment(owner.ID, user.ID)
	if e != nil {
		return e
	}

	admin := Role{Name: "Admin", Group: &group}
	admin.SetAll(true)
	e = s.repo.CreateRole(&admin)
	if e != nil {
		return e
	}

	e = s.repo.CreateRoleAssignment(admin.ID, user.ID)
	if e != nil {
		return e
	}

	def := Role{Name: "Default", Group: &group, GetsChores: true}
	e = s.repo.CreateRole(&def)
	if e != nil {
		return e
	}

	e = s.repo.CreateRoleAssignment(def.ID, user.ID)
	if e != nil {
		return e
	}

	return nil
}

func (s *groupService) GetGroup(group *Group) error {
	if e := s.repo.GetGroupByID(group); e != nil {
		return e
	}
	return nil
}

func (s *groupService) GetMemberships(t interface{}) error {
	return s.repo.GetMemberships(t)
}

func (s *groupService) GetRoles(t interface{}) error {
	e := s.repo.GetRoles(t)
	switch v := t.(type) {
	case *Membership:
		v.BuildSuperRole()
	}
	return e
}

func (s *groupService) UpdateGroup(group *Group) error {
	e := s.repo.UpdateGroup(group)
	return e
}

func (s *groupService) CanEdit(group *Group, user *User) bool {
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
		return false
	}
	log.Printf("CanEdit: Number of roles: %v", len(mem.Roles))
	if !mem.SuperRole.CanEdit() {
		log.Printf("User has insufficient privileges to edit this group: %v", mem.SuperRole.Permissions)
		return false
	}
	return true
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
	if e := s.repo.GetMembership(mem); e != nil {
		return e
	}
	if e := s.repo.GetRoles(mem); e != nil {
		return e
	}
	for _, v := range mem.Roles {
		mem.SuperRole.Permissions |= v.Permissions
	}
	//TODO: Get chore assignments
	return nil
}

func (s *groupService) AddMember(mem *Membership) error {
	mem.JoinedAt = time.Now().UTC()
	return s.repo.CreateMembership(mem)
}

func (s *groupService) AddRole(role *Role) error {
	return s.repo.CreateRole(role)
}
