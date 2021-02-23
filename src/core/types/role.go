package types

// Role defines what a member has access to within a group and if they get chores assigned to them
type Role struct {
	ID          uint64
	Name        string
	Permissions int
	GetsChores  bool
	Group       *Group
	Users       []User
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

func (role *Role) Can(bit PermBit) bool {
	mask := 1 << bit
	return role.Permissions&mask != 0
}

func (role *Role) Set(bit PermBit, value bool) {
	mask := 1 << bit
	if value {
		role.Permissions |= mask
	} else {
		mask = ^mask
		role.Permissions &= mask
	}
}

func (role *Role) SetAll(value bool) {
	role.Permissions = 0
	if value {
		role.Permissions = ^role.Permissions
	}
}

func (role *Role) CanEdit() bool {
	mask := 1 << EditMembers
	mask |= 1 << EditChores
	mask |= 1 << EditGroup
	return role.Permissions&mask != 0
}
