package core

import (
	"errors"
)

type ChoreRepository interface {
	CreateChore(*Chore) error
	GetChores(interface{}) error
	GetChore(*Chore) error
}

type ChoreService interface {
	Create(*Chore) error
	GetChore(*Chore) error
}

type choreService struct {
	repo ChoreRepository
}

func NewChoreService(r ChoreRepository) ChoreService {
	return &choreService{
		repo: r,
	}
}

func (s *choreService) Create(ch *Chore) error {
	if e := s.repo.GetChores(ch.Group); e != nil {
		return errors.New("An unexpected error occurred")
	}
	for _, v := range ch.Group.Chores {
		if v.Name == ch.Name {
			return errors.New("Chore already exists")
		}
	}
	if e := s.repo.CreateChore(ch); e != nil {
		return errors.New("An unexpected error occurred")
	}
	return nil
}

func (s *choreService) GetChore(ch *Chore) error {
	if e := s.repo.GetChore(ch); e != nil {
		return e
	}
	return nil
}
