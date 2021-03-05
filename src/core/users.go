package core

import (
	"errors"

	storagErr "chores-suck/core/storage/errors"
)

var (
	ErrEmailExists = errors.New("Email already registered")
	ErrNameExists  = errors.New("Username already registered")
)

type UserRepository interface {
	GetUserByName(user *User) error
	GetUserByEmail(user *User) error
	GetUserByID(user *User) error
	CreateUser(user *User) error
	GetMemberships(t interface{}) error
	GetChores(t interface{}) error
}

type UserService interface {
	GetUserByName(user *User) error
	GetUserByID(user *User) error
	CreateUser(user *User) error
	CheckEmailExists(email string) (bool, error)
	CheckUsernameExists(name string) (bool, error)
	GetMemberships(user *User) error
	GetChores(user *User) error
}

type userService struct {
	repo UserRepository
}

func NewUserService(rep UserRepository) UserService {
	return &userService{
		repo: rep,
	}
}

func (s *userService) GetUserByName(user *User) error {
	if e := s.repo.GetUserByName(user); e != nil {
		return e
	}
	return s.getUserInternal(user)
}

func (s *userService) GetUserByID(user *User) error {
	if e := s.repo.GetUserByID(user); e != nil {
		return e
	}
	return s.getUserInternal(user)
}

func (s *userService) getUserInternal(user *User) error {
	return s.repo.GetMemberships(user)
	// TODO: get chores
}

func (s *userService) CreateUser(user *User) error {
	exists, e := s.CheckEmailExists(user.Email)
	if e != nil {
		return e
	} else if exists {
		return ErrEmailExists
	}

	exists, e = s.CheckUsernameExists(user.Username)
	if e != nil {
		return e
	} else if exists {
		return ErrNameExists
	}

	e = s.repo.CreateUser(user)
	return e

}

func (s *userService) CheckEmailExists(email string) (bool, error) {
	user := User{Email: email}
	e := s.repo.GetUserByEmail(&user)
	if e == storagErr.ErrNotFound {
		return false, nil
	}
	return true, e
}

func (s *userService) CheckUsernameExists(name string) (bool, error) {
	user := User{Username: name}
	e := s.repo.GetUserByName(&user)
	if e == storagErr.ErrNotFound {
		return false, nil
	}
	return true, e
}

func (s *userService) GetMemberships(user *User) error {
	e := s.repo.GetMemberships(user)
	return e
}

func (s *userService) GetChores(user *User) error {
	e := s.repo.GetChores(user)
	return e
}
