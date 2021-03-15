package chores

import (
	"chores-suck/core"
	"fmt"
	"math/rand"
	"sort"
	"time"
)

// Person holds information for a person in the context of chores
type Person struct {
	Name  string
	Score uint16
	ID    uint64
}

// Randomize randomly distributes a set of chores to a set of people.
// Each person will have a minimum amount of chores to work on based
// on the time of each chore.
func Randomize(c []core.Chore, p []core.User) {
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
		if uint32(scores[p[rIndex].ID]) < minScore || allChecked {
			c[i].Assignment.User = &p[rIndex]
			scores[p[rIndex].ID] += c[i].Duration / 5
		}
		rIndex++
	}
}

// Rotate rotates the assigned chores amongst the people.
func Rotate(c []core.Chore, u []core.User) {
	// Rotate the chores amongst the roommates.
	// The roommates must already have chores assigned.
	// The first roommate in the list gets the last roommates chores
	// the second roommate gets the first roommates chores
	sort.Slice(u, func(i, j int) bool {
		return u[i].ID < u[j].ID
	})

	rmap := make(map[uint64]*core.User)
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

func printChores(c []core.Chore) string {
	var s string
	for i := range c {
		s += fmt.Sprintf("%s\t%v\n", c[i].Name, c[i].Assignment.User.Username)
	}
	return s
}

func printRoommates(r []Person) {
	for i := range r {
		fmt.Printf("%s\t%v\n", r[i].Name, r[i].Score)
	}
}
