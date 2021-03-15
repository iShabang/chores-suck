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
	UpdateRole(role *Role) error
	GetChores(t interface{}) error
}

type GroupService interface {
	// CreateGroup creates a new group with the default roles (owner, admin, default) and creates a new membership
	// for the owner (passed in user)
	CreateGroup(name string, user *User) error
	GetGroup(group *Group) error
	GetMemberships(t interface{}) error
	GetMembership(mem *Membership) error
	GetRoles(t interface{}) error
	UpdateGroup(group *Group, user *User) error
	CanEdit(group *Group, user *User) bool
	DeleteMember(mem *Membership, user *User) error
	AddMember(mem *Membership, user *User) error
	AddRole(role *Role, user *User) error
	UpdateRole(role *Role, user *User) error
	GetChores(group *Group) error
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

func (s *groupService) UpdateGroup(group *Group, user *User) error {
	mem := group.FindMember(user.ID)
	e := s.GetRoles(mem)
	if e != nil {
		return e
	}
	if !mem.SuperRole.Can(EditGroup) {
		return errors.New("Insufficient permissions")
	}
	return s.repo.UpdateGroup(group)
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

func (s *groupService) DeleteMember(mem *Membership, user *User) error {
	authMem := mem.Group.FindMember(user.ID)
	if e := s.GetRoles(authMem); e != nil {
		return errors.New("An unexpected error occurred")
	}
	if !authMem.SuperRole.Can(EditMembers) {
		return errors.New("You do not have permission to remove members!")
	}
	if e := s.GetRoles(mem); e != nil {
		return errors.New("An unexpected error occurred")
	}
	for _, v := range mem.Roles {
		if v.Name == "Owner" {
			return errors.New("Cannot delete owner")
		}
	}
	if e := s.repo.DeleteMember(mem); e != nil {
		return errors.New("An unexpected error occurred")
	}
	return nil
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

func (s *groupService) AddMember(mem *Membership, user *User) error {
	authMem := mem.Group.FindMember(user.ID)
	if e := s.GetRoles(authMem); e != nil {
		return errors.New("An unexpected error occurred")
	}
	if !authMem.SuperRole.Can(EditMembers) {
		return errors.New("You do not have permission to add members!")
	}
	mem.JoinedAt = time.Now().UTC()
	if e := s.repo.CreateMembership(mem); e != nil {
		return errors.New("An unexpected error occurred")
	}
	return nil
}

func (s *groupService) AddRole(role *Role, user *User) error {
	mem := role.Group.FindMember(user.ID)
	if e := s.GetRoles(mem); e != nil {
		return errors.New("An unexpected error occurred")
	}
	if !mem.SuperRole.Can(EditRoles) {
		return errors.New("You do not have permission to add roles!")
	}
	if e := s.GetRoles(role.Group); e != nil {
		return errors.New("An unexpected error occurred")
	}
	for _, r := range role.Group.Roles {
		if r.Name == role.Name {
			return errors.New("Role already exists")
		}
	}
	if e := s.repo.CreateRole(role); e != nil {
		return errors.New("An unexpected error occurred")
	}
	return nil
}

func (s *groupService) UpdateRole(role *Role, user *User) error {
	mem := role.Group.FindMember(user.ID)
	if e := s.GetRoles(mem); e != nil {
		return errors.New("An unexpected error occurred")
	}
	if !mem.SuperRole.Can(EditRoles) {
		return errors.New("You do not have permission to update roles")
	}
	s.GetRoles(role.Group)
	oldRole := role.Group.FindRole(role.ID)
	if oldRole == nil {
		return errors.New("Invalid request")
	} else if oldRole.Name == "Owner" || oldRole.Name == "Admin" || oldRole.Name == "Default" {
		return errors.New("Cannot make changes to Owner, Admin, or Default roles")
	}
	if e := s.repo.UpdateRole(role); e != nil {
		return errors.New("An unexpected error occurred")
	}
	return nil
}

func (s *groupService) GetChores(group *Group) error {
	if e := s.repo.GetChores(group); e != nil {
		return errors.New("An unexpected error occurred")
	}
	return nil
}
