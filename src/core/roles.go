package core

import (
	"errors"
)

type RoleRepository interface {
	RemoveMember(roleID uint64, userID uint64) error
	AddMember(roleID uint64, userID uint64) error
	GetRoles(t interface{}) error
	GetRole(role *Role) error
	UpdateRole(role *Role) error
}

type RoleService interface {
	RemoveMember(role *Role, userID uint64, user *User) error
	AddMember(role *Role, username string, user *User) error
	GetRole(role *Role) error
	Update(role *Role, newRole *Role, user *User) error
}

type roleService struct {
	repo RoleRepository
	us   UserService
}

func NewRoleService(re RoleRepository, u UserService) RoleService {
	return &roleService{
		repo: re,
		us:   u,
	}
}

func (s *roleService) RemoveMember(role *Role, userID uint64, user *User) error {
	mem := role.Group.FindMember(user.ID)
	if e := s.repo.GetRoles(mem); e != nil {
		return errors.New("An unexpected error occurred")
	}
	if !mem.SuperRole.Can(EditRoles) {
		return errors.New("You do not have permission to edit role members")
	}
	if e := s.repo.GetRoles(role.Group); e != nil {
		return errors.New("An unexpected error occurred")
	}
	oldRole := role.Group.FindRole(role.ID)
	if oldRole == nil {
		return errors.New("Invalid request")
	}
	if oldRole.Name == "Owner" {
		return errors.New("Cannot remove owner")
	}
	if e := s.repo.RemoveMember(role.ID, userID); e != nil {
		return errors.New("An unexpected error occurred")
	}
	return nil
}

func (s *roleService) AddMember(role *Role, username string, user *User) error {
	authMem := role.Group.FindMember(user.ID)
	if !authMem.SuperRole.Can(EditMembers) {
		return errors.New("You do not have permission to add members!")
	}
	if e := s.repo.GetRoles(role.Group); e != nil {
		return errors.New("An unexpected error occurred")
	}
	oldRole := role.Group.FindRole(role.ID)
	if oldRole == nil {
		return errors.New("Invalid request")
	}
	newUser := User{Username: username}
	if e := s.us.GetUserByName(&newUser); e != nil {
		return errors.New("Member not found")
	}
	mem := role.Group.FindMember(newUser.ID)
	if mem == nil {
		return errors.New("Member not found")
	}
	if e := s.repo.AddMember(role.ID, newUser.ID); e != nil {
		return errors.New("An unexpected error occurred")
	}
	return nil
}

func (s *roleService) GetRole(role *Role) error {
	return s.repo.GetRole(role)
}

func (s *roleService) Update(role *Role, newRole *Role, user *User) error {
	if role.Name == "Owner" || role.Name == "Admin" || role.Name == "Default" {
		return errors.New("Cannot make changes to Owner, Admin, or Default roles")
	}
	if e := s.repo.GetRoles(role.Group); e != nil {
		return errors.New("An unexpected error occurred")
	}
	if role.Name != newRole.Name {
		for i := range role.Group.Roles {
			if newRole.Name == role.Group.Roles[i].Name {
				return errors.New("Role name already exists")
			}
		}
	}
	if e := s.repo.UpdateRole(newRole); e != nil {
		return errors.New("An unexpected error occurred")
	}
	return nil

}
