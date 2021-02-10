package types

// Role defines what a member has access to within a group
type Role struct {
	ID          string
	Name        string
	Permissions int
	GroupID     int
}

type PermBit int

const (
	EditMembers = 0
	EditChores  = 1
	EditGroup   = 2
	DeleteGroup = 3
)

/* Permission Bits
0 - EditMembers
1 - EditChores
2 - EditGroup
3 - DeleteGroup
*/

func (role *Role) CanEditMembers() bool {
	mask := 1 << EditMembers
	return role.Permissions&mask != 0
}

func (role *Role) CanEditChores() bool {
	mask := 1 << EditChores
	return role.Permissions&mask != 0
}

func (role *Role) CanEditGroup() bool {
	mask := 1 << EditGroup
	return role.Permissions&mask != 0
}

func (role *Role) CanDeleteGroup() bool {
	mask := 1 << DeleteGroup
	return role.Permissions&mask != 0
}
