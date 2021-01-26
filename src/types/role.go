package types

// Role defines what a member has access to within a group
type Role struct {
	ID           string
	Name         string
	Persmissions Permissions
}

// Permissions defines the actions a member with a specific role is allowed to perform in a group
type Permissions struct {
	EditMembers bool
	EditChores  bool
	EditGroup   bool
	DeleteGroup bool
}
