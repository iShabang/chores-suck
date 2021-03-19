package core

import (
	"errors"
	"log"
	"math/rand"
	"sort"
	"time"
)

type ChoreRepository interface {
	CreateChore(*Chore) error
	GetChores(interface{}) error
	GetChore(*Chore) error
	UpdateChore(*Chore) error
	DeleteChore(*Chore) error
	InsertAssignments([]ChoreAssignment) error
	DeleteAssignments([]ChoreAssignment) error
}

type ChoreService interface {
	Create(*Chore) error
	Update(ch *Chore, new *Chore) error
	Delete(ch *Chore) error
	GetChore(*Chore) error
	Randomize(g *Group) error
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

func (s *choreService) Delete(ch *Chore) error {
	if e := s.repo.DeleteChore(ch); e != nil {
		log.Printf("ChoreService: Delete: Operation Failed: %s", e.Error())
		return errors.New("An unexpected error occurred")
	}
	return nil
}

func (s *choreService) Randomize(g *Group) error {
	//Store old assignments
	oldCa := make([]ChoreAssignment, 0, len(g.Chores))
	newCa := make([]ChoreAssignment, 0, len(g.Chores))
	for i := range g.Chores {
		if g.Chores[i].Assignment != nil {
			oldCa = append(oldCa, *g.Chores[i].Assignment)
		}
		// TODO: determine the due date
		g.Chores[i].Assignment = &ChoreAssignment{
			Chore:        &g.Chores[i],
			DateAssigned: time.Now().UTC(),
		}
	}
	// TODO: only pass members that get chores
	randomize(g.Chores, g.Memberships)
	failed := false
	for i := range g.Chores {
		newCa = append(newCa, *g.Chores[i].Assignment)
		if newCa[i].User == nil {
			log.Printf("Missing user for chore: %v", g.Chores[i].ID)
			failed = true
		}
	}
	if failed {
		return errors.New("Failed to distribute chores")
	}
	if e := s.repo.DeleteAssignments(oldCa); e != nil {
		// return error message
		log.Printf("Core: ChoreService: Randomize: failed to delete assignments: %s", e.Error())
		return errors.New("An unexpected error occurred")
	}
	//Update the database
	if e := s.repo.InsertAssignments(newCa); e != nil {
		//return error message
		log.Printf("Core: ChoreService: Randomize: failed to insert assignments: %s", e.Error())
		return errors.New("An unexpected error occurred")
	}
	return nil
}

// Randomize randomly distributes a set of chores to a set of people.
// Each person will have a minimum amount of chores to work on based
// on the time of each chore.
func randomize(c []Chore, p []Membership) {
	/* Fisher-Yates Shuffle Algorithm
	for i = n-1; i>0; i--
		j = random number from 0 <= j <= i
		swap c[j] with c[i]
	*/
	// Randomize the list of chores
	rand.Seed(time.Now().UnixNano())
	for i := len(c) - 1; i > 0; i-- {
		j := rand.Intn(i)
		c[i], c[j] = c[j], c[i]
	}

	// Find the average score of all the chores
	var minScore uint32
	for i := range c {
		minScore += uint32(c[i].Duration / 5)
	}
	minScore = minScore / uint32(len(p))

	// Generate a map of users to scores
	scores := make(map[uint64]int)

	var rIndex int
	allChecked := false

	for i := range c {
		if rIndex >= len(p) {
			allChecked = true
			rIndex = 0
		}
		if uint32(scores[p[rIndex].User.ID]) < minScore || allChecked {
			c[i].Assignment.User = p[rIndex].User
			scores[p[rIndex].User.ID] += c[i].Duration / 5
		}
		rIndex++
	}
}

// Rotate rotates the assigned chores amongst the people.
func rotate(c []Chore, u []User) {
	// Rotate the chores amongst the roommates.
	// The roommates must already have chores assigned.
	// The first roommate in the list gets the last roommates chores
	// the second roommate gets the first roommates chores
	sort.Slice(u, func(i, j int) bool {
		return u[i].ID < u[j].ID
	})

	rmap := make(map[uint64]*User)
	for i := range u {
		j := i + 1
		if j >= len(u) {
			j = 0
		}
		rmap[u[i].ID] = &u[j]
	}

	for i := range c {
		c[i].Assignment.User = rmap[c[i].Assignment.User.ID]
	}
}
