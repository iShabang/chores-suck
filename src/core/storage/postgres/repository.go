package postgres

import (
	"chores-suck/core"
	"chores-suck/core/storage/errors"

	_ "github.com/lib/pq" //Required for compilation. Functionality is wrapped using database/sql

	"database/sql"
	"log"
	"os"
	"time"
)

type pgsession struct {
	ID       int
	SesID    string
	Values   string
	Created  time.Time
	Expires  time.Time
	Modified time.Time
}

// Storage defines properties of a storage object
type Storage struct {
	Db *sql.DB
}

// NewStorage creates and returns a new storage object
func NewStorage() *Storage {
	connString := os.Getenv("POSTGRES_CONN")
	db, err := sql.Open("postgres", connString)
	s := &Storage{
		Db: db,
	}
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return s
}

// GetUserByName fetches a user from the database by unique username
func (s *Storage) GetUserByName(user *core.User) error {
	query := `
	SELECT users.id, users.email, users.pword, users.created_at 
	FROM users 
	WHERE users.uname = $1`
	err := s.Db.QueryRow(query, user.Username).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return errors.ErrNotFound
	}
	return err
}

func (s *Storage) GetUserByEmail(user *core.User) error {
	query := `
	SELECT users.id, users.uname, users.pword, users.created_at 
	FROM users 
	WHERE users.email = $1`
	err := s.Db.QueryRow(query, user.Email).Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return errors.ErrNotFound
	}
	return err
}

// GetUserByID fetches a user from the database by unique ID
func (s *Storage) GetUserByID(user *core.User) error {
	err := s.Db.QueryRow("SELECT uname, email, pword, created_at FROM users WHERE id = $1", user.ID).Scan(&user.Username, &user.Email, &user.Password, &user.CreatedAt)
	return err
}

// CreateUser adds a new user to the database
func (s *Storage) CreateUser(user *core.User) error {
	ca := time.Now().UTC()
	query := `INSERT INTO users (uname, email, pword, created_at) VALUES ($1,$2,$3,$4) RETURNING id`
	err := s.Db.QueryRow(query, user.Username, user.Email, user.Password, ca).Scan(&user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetChores(t interface{}) error {
	switch v := t.(type) {
	case *core.User:
		return s.GetUserChores(v)
	// TODO: Implement GetGroupChores
	//case *core.Group:
	//	return s.GetGroupChores(v)
	default:
		return errors.ErrType
	}
}

func (s *Storage) GetUserChores(user *core.User) error {
	query := `
	SELECT ca.complete, ca.date_assigned, ca.date_complete, ca.date_due,
	c.id, c.name, c.description, c.duration, g.id, g.name
	FROM chore_assignments ca
	INNER JOIN chores c ON c.id = ca.chore_id
	INNER JOIN groups g ON g.id = c.group_id
	WHERE ca.user_id = $1`

	rows, err := s.Db.Query(query, user.ID)

	if err != nil {
		return err
	}

	for rows.Next() {
		ca := core.ChoreAssignment{User: user}
		c := core.Chore{User: user}
		g := core.Group{}

		err = rows.Scan(&ca.Complete, &ca.DateAssigned, &ca.DateComplete, &ca.DateDue,
			&c.ID, &c.Name, &c.Description, &c.Duration, &g.ID, &g.Name)

		if err != nil && err != sql.ErrNoRows {
			return err
		}

		ca.Chore = &c
		c.Assignment = &ca
		c.Group = &g
		user.Chores = append(user.Chores, c)
	}

	return nil
}

func (s *Storage) GetGroupByID(group *core.Group) error {
	query := `
	SELECT name FROM groups WHERE id = $1`
	e := s.Db.QueryRow(query, group.ID).Scan(&group.Name)
	return e
}

func (s *Storage) GetMembership(mem *core.Membership) error {
	query := `SELECT joined_at FROM memberships WHERE group_id = $1 AND user_id = $2`
	e := s.Db.QueryRow(query, mem.Group.ID, mem.User.ID).Scan(&mem.JoinedAt)
	return e
}

func (s *Storage) GetMemberships(t interface{}) error {
	switch v := t.(type) {
	case *core.User:
		return s.GetUserMemberships(v)
	case *core.Group:
		return s.GetGroupMemberships(v)
	case *core.Role:
		return s.GetRoleMemberships(v)
	default:
		return errors.ErrType
	}
}

func (s *Storage) GetUserMemberships(user *core.User) error {
	query := `
	SELECT m.joined_at, m.group_id, g.name
	FROM memberships m
	INNER JOIN groups g ON g.id = m.group_id
	WHERE m.user_id = $1`

	rows, err := s.Db.Query(query, user.ID)

	if err != nil {
		return err
	}

	user.Memberships = []core.Membership{}
	defer rows.Close()
	for rows.Next() {
		mem := core.Membership{
			User:  user,
			Group: &core.Group{},
		}
		err := rows.Scan(&mem.JoinedAt, &mem.Group.ID, &mem.Group.Name)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil
			}
			return err
		}
		user.Memberships = append(user.Memberships, mem)
	}

	return nil
}

func (s *Storage) GetGroupMemberships(group *core.Group) error {
	query := `
	SELECT m.joined_at, m.user_id, u.uname
	FROM memberships m
	INNER JOIN users u ON m.user_id = u.id
	WHERE m.group_id = $1`

	rows, err := s.Db.Query(query, group.ID)
	if err != nil {
		return err
	}
	group.Memberships = []core.Membership{}
	defer rows.Close()
	for rows.Next() {
		mem := core.Membership{Group: group, User: &core.User{}}
		err = rows.Scan(&mem.JoinedAt, &mem.User.ID, &mem.User.Username)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil
			}
			return err
		}
		group.Memberships = append(group.Memberships, mem)
	}

	return nil
}

func (s *Storage) GetRoleMemberships(role *core.Role) error {
	query := `
	SELECT m.joined_at, m.user_id, u.uname
	FROM memberships m
	INNER JOIN role_assignments ra ON ra.role_id = $1
	INNER JOIN users u ON u.id = ra.user_id
	WHERE m.user_id = ra.user_id AND m.group_id = $2`
	rows, e := s.Db.Query(query, role.ID, role.Group.ID)
	if e != nil {
		return e
	}
	for rows.Next() {
		mem := core.Membership{User: &core.User{}, Group: role.Group}
		e = rows.Scan(&mem.JoinedAt, &mem.User.ID, &mem.User.Username)
		if e != nil {
			log.Print(e.Error())
			return e
		}
		role.Members = append(role.Members, mem)
	}
	return nil
}

func (s *Storage) CreateGroup(group *core.Group) error {
	query := `INSERT INTO groups (name) VALUES ($1) RETURNING id`
	e := s.Db.QueryRow(query, &group.Name).Scan(&group.ID)
	return e
}

func (s *Storage) UpdateGroup(group *core.Group) error {
	query := `UPDATE groups SET name = $1 WHERE id = $2`
	_, e := s.Db.Exec(query, group.Name, group.ID)
	return e
}

func (s *Storage) CreateRole(role *core.Role) error {
	query := `INSERT INTO roles (name, permissions, group_id, gets_chores) VALUES ($1,$2,$3,$4) RETURNING id`
	e := s.Db.QueryRow(query, role.Name, role.Permissions, role.Group.ID, role.GetsChores).Scan(&role.ID)
	return e
}

func (s *Storage) CreateRoleAssignment(roleID uint64, userID uint64) error {
	query := `INSERT INTO role_assignments (role_id, user_id) VALUES ($1,$2)`
	_, e := s.Db.Exec(query, roleID, userID)
	return e
}

func (s *Storage) CreateMembership(mem *core.Membership) error {
	query := `INSERT INTO memberships (joined_at, user_id, group_id) VALUES ($1,$2,$3)`
	_, e := s.Db.Exec(query, mem.JoinedAt, mem.User.ID, mem.Group.ID)
	return e
}

func (s *Storage) GetMemberChores(member *core.Membership) error {
	rows, err := s.Db.Query("SELECT chores.id, chores.description, chores.name, chores.duration, chore_assignments.complete, chore_assignments.date_assigned, chore_assignments.date_complete FROM chores WHERE chores.group_id = $1 INNER JOIN chore_assignment ON chore_assignment.chore_id = chores.id", member.Group.ID)

	if err != nil {
		return err
	}

	member.Assignments = []core.ChoreAssignment{}
	defer rows.Close()
	for rows.Next() {
		ca := core.ChoreAssignment{
			User:  member.User,
			Chore: &core.Chore{Group: member.Group},
		}
		err := rows.Scan(&ca.Chore.ID, &ca.Chore.Description, &ca.Chore.Name, &ca.Chore.Duration, &ca.Complete, &ca.DateAssigned, &ca.DateComplete)

		if err != nil {
			return err
		}

		member.Assignments = append(member.Assignments, ca)
	}

	return nil
}

func (s *Storage) GetRoles(t interface{}) error {
	switch v := t.(type) {
	case *core.Group:
		return s.GetGroupRoles(v)
	case *core.Membership:
		return s.GetMemberRoles(v)
	default:
		return errors.ErrType
	}
}

func (s *Storage) GetGroupRoles(group *core.Group) error {
	query := `
	SELECT id, name, permissions, gets_chores
	FROM roles
	WHERE group_id = $1`
	rows, e := s.Db.Query(query, group.ID)
	if e != nil {
		return e
	}
	group.Roles = []core.Role{}
	defer rows.Close()
	for rows.Next() {
		role := core.Role{Group: group}
		e := rows.Scan(&role.ID, &role.Name, &role.Permissions, &role.GetsChores)
		if e != nil {
			if e == sql.ErrNoRows {
				return nil
			}
			return e
		}
		group.Roles = append(group.Roles, role)
	}
	return nil
}

func (s *Storage) GetMemberRoles(member *core.Membership) error {
	query := `
	SELECT r.id, r.name, r.permissions, r.gets_chores
	FROM role_assignments ra
	INNER JOIN roles r on r.id = ra.role_id
	WHERE ra.user_id = $1 AND r.group_id = $2`

	rows, e := s.Db.Query(query, member.User.ID, member.Group.ID)
	if e != nil {
		return e
	}
	member.Roles = []core.Role{}
	defer rows.Close()
	for rows.Next() {
		role := core.Role{Group: member.Group}
		e := rows.Scan(&role.ID, &role.Name, &role.Permissions, &role.GetsChores)
		if e != nil {
			if e == sql.ErrNoRows {
				return nil
			}
			return e
		}
		member.Roles = append(member.Roles, role)
	}
	return nil

}

func (s *Storage) GetRole(role *core.Role) error {
	query := `
	SELECT name, permissions, gets_chores, group_id
	FROM roles WHERE id = $1`
	role.Group = &core.Group{}
	e := s.Db.QueryRow(query, role.ID).Scan(&role.Name, &role.Permissions, &role.GetsChores, &role.Group.ID)
	if e == sql.ErrNoRows {
		return nil
	} else {
		return e
	}
}

func (s *Storage) UpdateRole(role *core.Role) error {
	query := `UPDATE roles SET (name, permissions, gets_chores) = ($1, $2, $3) WHERE id = $4`
	_, e := s.Db.Exec(query, role.Name, role.Permissions, role.GetsChores, role.ID)
	return e
}

func (s *Storage) DeleteMember(mem *core.Membership) error {
	query := `DELETE FROM memberships WHERE group_id = $1 AND user_id = $2`
	_, e := s.Db.Exec(query, mem.Group.ID, mem.User.ID)
	if e != nil {
		return e
	}
	query = `DELETE FROM role_assignments WHERE user_id = $1 AND role_id = $2`
	for _, v := range mem.Roles {
		_, e = s.Db.Exec(query, mem.User.ID, v.ID)
		if e != nil {
			return e
		}
	}
	return nil
}

func (s *Storage) RemoveMember(roleID uint64, userID uint64) error {
	query := `DELETE FROM role_assignments WHERE role_id = $1 AND user_id = $2`
	_, e := s.Db.Exec(query, roleID, userID)
	return e
}

func (s *Storage) AddMember(roleID uint64, userID uint64) error {
	query := `INSERT INTO role_assignments (user_id, role_id) VALUES($1,$2)`
	_, e := s.Db.Exec(query, userID, roleID)
	return e
}

// GetSession fetches a session frm the database by session id
func (s *Storage) GetSession(ses *core.Session) error {
	err := s.Db.QueryRow("SELECT values, created, user_id FROM sessions WHERE uuid = $1", ses.UUID).Scan(&ses.Values, &ses.Created, &ses.UserID)

	if err == sql.ErrNoRows {
		return errors.ErrNotFound
	}

	return err
}

// DeleteSession removes a session from the database
func (s *Storage) DeleteSession(UUID string) error {
	statement, err := s.Db.Prepare("DELETE FROM sessions WHERE uuid = $1")
	if err != nil {
		return err
	}
	defer statement.Close()
	_, err = statement.Exec(UUID)
	return err
}

// UpsertSession inserts or updates a session in the database. If the session does not exist
// it is created otherwise the existing session is updated.
func (s *Storage) UpsertSession(ses *core.Session) error {
	statement, err := s.Db.Prepare("INSERT INTO sessions (uuid, values, created, user_id) VALUES ($1,$2,$3,$4) ON CONFLICT (uuid) DO UPDATE SET values = $2")
	if err != nil {
		return err
	}
	defer statement.Close()
	_, err = statement.Exec(ses.UUID, ses.Values, ses.Created, ses.UserID)
	return err
}
