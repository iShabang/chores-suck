package chores

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

type Chore struct {
	Name     string
	Minutes  uint16
	Assignee uint64
}

type Roommate struct {
	Name  string
	Score uint16
	ID    uint64
}

func Randomize(c []Chore, r []Roommate) {
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
	minScore = minScore / uint32(len(r))

	var rIndex uint
	allChecked := false

	for i := range c {
		if int(rIndex) >= len(r) {
			allChecked = true
			rIndex = 0
		}
		if uint32(r[rIndex].Score) < minScore {
			c[i].Assignee = r[rIndex].ID
			r[rIndex].Score += c[i].Minutes / 5
		}
		if allChecked {
			c[i].Assignee = r[rIndex].ID
			r[rIndex].Score += c[i].Minutes / 5
		}
		rIndex++
	}
}

func Rotate(c []Chore, r []Roommate) {
	// Rotate the chores amongst the roommates.
	// The roommates must already have chores assigned.
	// The first roommate in the list gets the last roommates chores
	// the second roommate gets the first roommates chores
	sort.Slice(r, func(i, j int) bool {
		return r[i].ID < r[j].ID
	})

	rmap := make(map[uint64]uint64)
	for i := range r {
		j := i + 1
		if j >= len(r) {
			j = 0
		}
		rmap[r[i].ID] = r[j].ID
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

func printRoommates(r []Roommate) {
	for i := range r {
		fmt.Printf("%s\t%v\n", r[i].Name, r[i].Score)
	}
}
