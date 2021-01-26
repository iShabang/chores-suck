package types

// Chore describes properties of a chore
type Chore struct {
	ID          string
	Description string
	Name        string
	GroupID     string

	// Estimated time in minutes it takes to complete the task
	TimeToComplete int
}
