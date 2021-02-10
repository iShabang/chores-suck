package types

// Group defines properties for a group
type Group struct {
	ID          int
	Name        string
	Memberships []Membership
	Roles       []Role
	Chores      []Chore
}
