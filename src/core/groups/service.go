package groups

import (
	"chores-suck/core/types"
	"time"
)

type Repository interface {
	CreateGroup(group *types.Group) error
	CreateRole(role *types.Role) error
	CreateRoleAssignment(roleID uint64, userID uint64) error
	CreateMembership(mem *types.Membership) error
	GetGroupByID(group *types.Group) error
	GetMemberships(t interface{}) error
	GetRoles(t interface{}) error
}

type Service interface {
	// CreateGroup creates a new group with the default roles (owner, admin, default) and creates a new membership
	// for the owner (passed in user)
	CreateGroup(name string, user *types.User) error
	GetGroup(group *types.Group) error
	GetMemberships(group *types.Group) error
	GetRoles(t interface{}) error
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{
		repo: r,
	}
}

func (s *service) CreateGroup(name string, user *types.User) error {
	group := types.Group{Name: name}
	e := s.repo.CreateGroup(&group)
	if e != nil {
		return e
	}
	users := make([]types.User, 1)
	users[0] = *user

	owner := types.Role{Name: "Owner", Group: &group, Users: users}
	owner.SetAll(true)
	e = s.repo.CreateRole(&owner)
	if e != nil {
		return e
	}

	e = s.repo.CreateRoleAssignment(owner.ID, user.ID)
	if e != nil {
		return e
	}

	admin := types.Role{Name: "Admin", Group: &group, Users: users}
	admin.SetAll(true)
	e = s.repo.CreateRole(&admin)
	if e != nil {
		return e
	}

	e = s.repo.CreateRoleAssignment(admin.ID, user.ID)
	if e != nil {
		return e
	}

	def := types.Role{Name: "Default", Group: &group, Users: users, GetsChores: true}
	e = s.repo.CreateRole(&def)
	if e != nil {
		return e
	}

	e = s.repo.CreateRoleAssignment(def.ID, user.ID)
	if e != nil {
		return e
	}

	mem := types.Membership{JoinedAt: time.Now().UTC(), User: user, Group: &group}
	e = s.repo.CreateMembership(&mem)
	if e != nil {
		return e
	}

	return nil
}

func (s *service) GetGroup(group *types.Group) error {
	e := s.repo.GetGroupByID(group)
	return e
}

func (s *service) GetMemberships(group *types.Group) error {
	e := s.repo.GetMemberships(group)
	return e
}

func (s *service) GetRoles(t interface{}) error {
	e := s.repo.GetRoles(t)
	return e
}
