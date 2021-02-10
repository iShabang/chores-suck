package types

// Role defines what a member has access to within a group
type Role struct {
	ID           string
	Name         string
	Persmissions int
	GroupID      int
}

/* Permission Bits
0 - EditMembers
1 - EditChores
2 - EditGroup
3 - DeleteGroup
*/
