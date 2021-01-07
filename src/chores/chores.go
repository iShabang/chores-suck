package chores

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

// Chore holds information for a chore
type Chore struct {
	Name     string
	Minutes  uint16
	Assignee uint64
}

// Person holds information for a person in the context of chores
type Person struct {
	Name  string
	Score uint16
	ID    uint64
}

// Randomize randomly distributes a set of chores to a set of people.
// Each person will have a minimum amount of chores to work on based
// on the time of each chore.
func Randomize(c []Chore, p []Person) {
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
		minScore += uint32(c[i].Minutes / 5)
	}
	minScore = minScore / uint32(len(p))

	var rIndex int
	allChecked := false

	for i := range c {
		if rIndex >= len(p) {
			allChecked = true
			rIndex = 0
		}
		if uint32(p[rIndex].Score) < minScore {
			c[i].Assignee = p[rIndex].ID
			p[rIndex].Score += c[i].Minutes / 5
		}
		if allChecked {
			c[i].Assignee = p[rIndex].ID
			p[rIndex].Score += c[i].Minutes / 5
		}
		rIndex++
	}
}

// Rotate rotates the assigned chores amongst the people.
func Rotate(c []Chore, p []Person) {
	// Rotate the chores amongst the roommates.
	// The roommates must already have chores assigned.
	// The first roommate in the list gets the last roommates chores
	// the second roommate gets the first roommates chores
	sort.Slice(p, func(i, j int) bool {
		return p[i].ID < p[j].ID
	})

	rmap := make(map[uint64]uint64)
	for i := range p {
		j := i + 1
		if j >= len(p) {
			j = 0
		}
		rmap[p[i].ID] = p[j].ID
	}

	for i := range c {
		c[i].Assignee = rmap[c[i].Assignee]
	}
}

func printChores(c []Chore) string {
	var s string
	for i := range c {
		s += fmt.Sprintf("%s\t%v\n", c[i].Name, c[i].Assignee)
	}
	return s
}

func printRoommates(r []Person) {
	for i := range r {
		fmt.Printf("%s\t%v\n", r[i].Name, r[i].Score)
	}
}
