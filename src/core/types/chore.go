package types

import "time"

// Chore describes properties of a chore
type Chore struct {
	ID          int
	Description string
	Name        string
	Duration    int
	Group       *Group
	User        *User
	Assignment  *ChoreAssignment
}

type ChoreListItem struct {
	GroupName string
	ChoreName string
	DateDue   time.Time
}
