package types

import "time"

type ChoreAssignment struct {
	Complete     bool
	DateAssigned time.Time
	DateComplete time.Time
	DateDue      time.Time
	Chore        *Chore
	User         *User
}
