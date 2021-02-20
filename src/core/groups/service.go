package groups

import (
	"chores-suck/core/types"
)

type Repository interface {
	CreateGroup(group *types.Group) error
	CreateRole(role *types.Role) error
}

type Service interface {
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
	// When creating the group, the user creating the group will be the owner role
	// create owner role, assign role to the current authenticated user
	// create the default role, assign to the current authenticated user
	// create the admin role, assign to the current authenticated user
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

	// Create a membership for the current authenticated user
}
