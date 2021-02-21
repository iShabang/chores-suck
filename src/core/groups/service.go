package groups

import (
	"chores-suck/core/types"
	"time"
)

type Repository interface {
	// CreateGroup adds a new group object to storage
	CreateGroup(group *types.Group) error

	// CreateRole adds a new role object to storage
	CreateRole(role *types.Role) error

	// CreateMembership adds a new membership object to storage
	CreateMembership(mem *types.Membership) error
}

type Service interface {
	// CreateGroup creates a new group with the default roles (owner, admin, default) and creates a new membership
	// for the owner (passed in user)
	CreateGroup(name string, user *types.User) error
}

type service struct {
	repo Repository
}

func (s *service) CreateGroup(name string, user *types.User) error {
	group := types.Group{Name: name}
	e := s.repo.CreateGroup(&group)
	if e != nil {
		return e
	}
	users := make([]*types.User, 1)
	users[0] = user

	owner := types.Role{Name: "Owner", Group: &group, Users: users}
	owner.SetAll(true)
	e = s.repo.CreateRole(&owner)
	if e != nil {
		return e
	}

	admin := types.Role{Name: "Admin", Group: &group, Users: users}
	admin.SetAll(true)
	e = s.repo.CreateRole(&admin)
	if e != nil {
		return e
	}

	def := types.Role{Name: "Default", Group: &group, Users: users}
	e = s.repo.CreateRole(&def)
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
