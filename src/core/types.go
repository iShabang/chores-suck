package core

import "time"

// Chore describes properties of a chore
type Chore struct {
	ID          uint64
	Description string
	Name        string
	Duration    int
	Group       *Group
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

func (g *Group) FindRole(id uint64) *Role {
	for i := range g.Roles {
		if g.Roles[i].ID == id {
			return &g.Roles[i]
		}
	}
	return nil
}

func (g *Group) FindMember(value interface{}) *Membership {
	switch v := value.(type) {
	case uint64:
		for i := range g.Memberships {
			if g.Memberships[i].User.ID == v {
				return &g.Memberships[i]
			}
		}
	case string:
		for i := range g.Memberships {
			if g.Memberships[i].User.Username == v {
				return &g.Memberships[i]
			}
		}
	}
	return nil
}

func (g *Group) FindChore(value interface{}) *Chore {
	switch v := value.(type) {
	case uint64:
		for i, c := range g.Chores {
			if v == c.ID {
				return &g.Chores[i]
			}
		}
	case string:
		for i, c := range g.Chores {
			if v == c.Name {
				return &g.Chores[i]
			}
		}
	}
	return nil
}

// Membership relates a User to a specific group
type Membership struct {
	JoinedAt    time.Time
	User        *User
	Group       *Group
	Assignments []ChoreAssignment
	Roles       []Role
	//SuperRole is a combination of all the roles assigned to the member
	//This is a convenience for checking permissions for a member without
	//checking all of their roles individually. See Membership.BuildSuperRole()
	//for an in-depth look on what information is available in the super role
	SuperRole Role
}

func (m *Membership) BuildSuperRole() {
	m.SuperRole.Name = "SuperRole"
	m.SuperRole.Group = m.Group
	for i := range m.Roles {
		m.SuperRole.Permissions |= m.Roles[i].Permissions
		m.SuperRole.GetsChores = m.SuperRole.GetsChores || m.Roles[i].GetsChores
	}
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
	Members     []Membership
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
