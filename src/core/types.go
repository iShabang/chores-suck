package core

import "time"

// Chore describes properties of a chore
type Chore struct {
	ID          int
	Description string
	Name        string
	Duration    int
	Group       *Group
	User        *User
	Assignment  *ChoreAssignment
}

type ChoreAssignment struct {
	Complete     bool
	DateAssigned time.Time
	DateComplete time.Time
	DateDue      time.Time
	Chore        *Chore
	User         *User
}

// Group defines properties for a group
type Group struct {
	ID          uint64
	Name        string
	Memberships []Membership
	Roles       []Role
	Chores      []Chore
}

// Membership relates a User to a specific group
type Membership struct {
	JoinedAt    time.Time
	User        *User
	Group       *Group
	Assignments []ChoreAssignment
	Roles       []Role
	SuperRole   Role
}

// User defines properties of a user
type User struct {
	ID          uint64
	Username    string
	Email       string
	Password    string
	CreatedAt   time.Time
	Memberships []Membership
	Chores      []Chore
}

// Session contains properties for a session pulled from the database
type Session struct {
	UUID    string
	Values  string
	Created time.Time
	UserID  uint64
}

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
	EditRoles   = 3
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
