package chores

import (
	"testing"
)

var (
	c []Chore = []Chore{
		{Name: "Bathroom Floor", Minutes: 10},
		{Name: "Kitchen Floor", Minutes: 15},
		{Name: "Kitchen Counter", Minutes: 5},
		{Name: "Front Floor", Minutes: 5},
		{Name: "Brawl Station", Minutes: 5},
		{Name: "Stove", Minutes: 5},
		{Name: "Toilet/Trash", Minutes: 5},
		{Name: "Fridge", Minutes: 5},
		{Name: "Bathtub", Minutes: 5},
	}

	r []Roommate = []Roommate{
		{Name: "Shannon", Score: 0, ID: 12345},
		{Name: "Shaun", Score: 0, ID: 12346},
		{Name: "Tyler", Score: 0, ID: 12347},
		{Name: "Preston", Score: 0, ID: 12348},
	}
)

func TestRandomize(t *testing.T) {
	Randomize(c, r)
	out1 := printChores(c)

	Randomize(c, r)
	out2 := printChores(c)

	if out1 == out2 {
		t.Error("Randomize produced equal output")
	}
}

func TestRotate(t *testing.T) {
	Randomize(c, r)
	out1 := c[0].Assignee

	Rotate(c, r)
	out2 := c[0].Assignee

	if out1 == out2 {
		t.Error("Assignee's the same after rotation")
	}
}
