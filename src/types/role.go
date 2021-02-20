package types

// Role defines what a member has access to within a group
type Role struct {
	ID          string
	Name        string
	Permissions int
	Group       *Group
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
