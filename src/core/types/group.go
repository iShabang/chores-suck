package types

// Group defines properties for a group
type Group struct {
	ID          uint64
	Name        string
	Memberships []Membership
	Roles       []Role
	Chores      []Chore
}
