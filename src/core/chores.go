package core

import (
	"errors"
	"log"
)

type ChoreRepository interface {
	CreateChore(*Chore) error
	GetChores(interface{}) error
	GetChore(*Chore) error
	UpdateChore(*Chore) error
}

type ChoreService interface {
	Create(*Chore) error
	Update(ch *Chore, new *Chore) error
	GetChore(*Chore) error
}

type choreService struct {
	repo ChoreRepository
	gs   GroupService
}

func NewChoreService(r ChoreRepository, g GroupService) ChoreService {
	return &choreService{
		repo: r,
		gs:   g,
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

func (s *choreService) Update(ch *Chore, new *Chore) error {
	if ch.Name != new.Name {
		if e := s.gs.GetChores(ch.Group); e != nil {
			log.Printf("ChoreService: Update: Failed to get group chores: %s", e.Error())
			return errors.New("An unexpected error occurred")
		}
		if c := ch.Group.FindChore(new.Name); c != nil {
			return errors.New("Chore name already in use")
		}
	}
	if e := s.repo.UpdateChore(new); e != nil {
		log.Printf("ChoreService: Update: Failed to update: %s", e.Error())
		return errors.New("An unexpected error occurred")
	}
	return nil
}
