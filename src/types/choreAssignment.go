package types

import "time"

type ChoreAssignment struct {
	Complete     bool
	DateAssigned time.Time
	DateComplete time.Time
	ChoreID      int
	UserID       int
}
