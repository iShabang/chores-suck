package core

import (
	"errors"
	"fmt"
	"log"
)

type RoleRepository interface {
	RemoveMember(roleID uint64, userID uint64) error
	AddMember(roleID uint64, userID uint64) error
	GetRoles(t interface{}) error
	GetRole(role *Role) error
	UpdateRole(role *Role) error
	DeleteRole(role *Role) error
}

type RoleService interface {
	RemoveMember(role *Role, userID uint64, user *User) error
	AddMember(role *Role, username string, user *User) error
	GetRole(role *Role) error
	Update(role *Role, newRole *Role, user *User) error
	Delete(role *Role) error
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
	if role.Name == "Owner" {
		return errors.New("Cannot remove owner")
	}
	if e := s.repo.RemoveMember(role.ID, userID); e != nil {
		return errors.New("An unexpected error occurred")
	}
	return nil
}

func (s *roleService) AddMember(role *Role, username string, user *User) error {
	if role.Name == "Owner" {
		return errors.New("There can only be one owner")
	}
	mem := role.Group.FindMember(username)
	if mem == nil {
		return errors.New("Member not found")
	}
	if e := s.repo.AddMember(role.ID, mem.User.ID); e != nil {
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

func (s *roleService) Delete(role *Role) error {
	if role.Name == "Owner" || role.Name == "Admin" || role.Name == "Default" {
		msg := fmt.Sprintf("Cannot delete %s role", role.Name)
		return errors.New(msg)
	}
	if e := s.repo.DeleteRole(role); e != nil {
		log.Printf("Core: RoleService: Delete: %s", e.Error())
		return errors.New("An unexpected error occurred")
	}
	return nil
}
