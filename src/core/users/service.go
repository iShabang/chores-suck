package users

import (
	"errors"

	storagErr "chores-suck/core/storage/errors"
	"chores-suck/core/types"
)

var (
	ErrEmailExists = errors.New("Email already registered")
	ErrNameExists  = errors.New("Username already registered")
)

type Repository interface {
	GetUserByName(user *types.User) error
	GetUserByEmail(user *types.User) error
	GetUserByID(user *types.User) error
	CreateUser(user *types.User) error
	GetMemberships(t interface{}) error
}

type Service interface {
	GetUserByName(user *types.User) error
	GetUserByID(user *types.User) error
	CreateUser(user *types.User) error
	CheckEmailExists(email string) (bool, error)
	CheckUsernameExists(name string) (bool, error)
	GetMemberships(user *types.User) error
}

type service struct {
	repo Repository
}

func NewService(rep Repository) Service {
	return &service{
		repo: rep,
	}
}

func (s *service) GetUserByName(user *types.User) error {
	e := s.repo.GetUserByName(user)
	return e
}

func (s *service) GetUserByID(user *types.User) error {
	e := s.repo.GetUserByID(user)
	return e
}

func (s *service) CreateUser(user *types.User) error {
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

func (s *service) CheckEmailExists(email string) (bool, error) {
	user := types.User{Email: email}
	e := s.repo.GetUserByEmail(&user)
	if e == storagErr.ErrNotFound {
		return false, nil
	}
	return true, e
}

func (s *service) CheckUsernameExists(name string) (bool, error) {
	user := types.User{Username: name}
	e := s.repo.GetUserByName(&user)
	if e == storagErr.ErrNotFound {
		return false, nil
	}
	return true, e
}

func (s *service) GetMemberships(user *types.User) error {
	e := s.repo.GetMemberships(user)
	return e
}
