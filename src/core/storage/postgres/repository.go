package postgres

import (
	"chores-suck/core/storage/errors"
	"chores-suck/core/types"
	"database/sql"
	"time"

	"log"

	_ "github.com/lib/pq" //Required for compilation. Functionality is wrapped using database/sql
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
	db, err := sql.Open("postgres", "host=192.168.1.202 port=5432 user=pi password=CorkstCork12 dbname=choressuck")
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
func (s *Storage) GetUserByName(user *types.User) error {
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

func (s *Storage) GetUserByEmail(user *types.User) error {
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
func (s *Storage) GetUserByID(user *types.User) error {
	err := s.Db.QueryRow("SELECT uname, email, pword, created_at FROM users WHERE id = $1", user.ID).Scan(&user.Username, &user.Email, &user.Password, &user.CreatedAt)
	return err
}

// CreateUser adds a new user to the database
func (s *Storage) CreateUser(user *types.User) error {
	ca := time.Now().UTC()
	query := `INSERT INTO users (uname, email, pword, created_at) VALUES ($1,$2,$3,$4) RETURNING id`
	err := s.Db.QueryRow(query, user.Username, user.Email, user.Password, ca).Scan(&user.ID)
	if err != nil {
		return err
	}
	return nil
}

// GetUserChoreList fetches a list of chore data
func (s *Storage) GetUserChoreList(user *types.User) ([]types.ChoreListItem, error) {
	query := `
	SELECT ca.date_due, c.name, g.name
	FROM chore_assignments ca
	INNER JOIN chores c ON c.id = ca.chore_id
	INNER JOIN groups g ON g.id = c.group_id
	WHERE ca.user_id = $1`

	rows, err := s.Db.Query(query, user.ID)

	if err != nil {
		return nil, err
	}

	var chores []types.ChoreListItem
	defer rows.Close()
	for rows.Next() {
		ch := types.ChoreListItem{}
		err = rows.Scan(&ch.DateDue, &ch.ChoreName, &ch.GroupName)
		if err != nil {
			// An empty list is a valid value
			if err == sql.ErrNoRows {
				return chores, nil
			}

			return nil, err
		}
		chores = append(chores, ch)
	}

	return chores, err
}

func (s *Storage) GetGroupByID(group *types.Group) error {
	query := `
	SELECT name FROM groups WHERE id = $1`
	e := s.Db.QueryRow(query, group.ID).Scan(&group.Name)
	return e
}

func (s *Storage) GetMemberships(t interface{}) error {
	switch v := t.(type) {
	case *types.User:
		return s.GetUserMemberships(v)
	case *types.Group:
		return s.GetGroupMemberships(v)
	default:
		return errors.ErrType
	}
}

func (s *Storage) GetUserMemberships(user *types.User) error {
	query := `
	SELECT m.joined_at, m.group_id, g.name
	FROM memberships m
	INNER JOIN groups g ON g.id = m.group_id
	WHERE m.user_id = $1`

	rows, err := s.Db.Query(query, user.ID)

	if err != nil {
		return err
	}

	user.Memberships = []types.Membership{}
	defer rows.Close()
	for rows.Next() {
		mem := types.Membership{
			User:  user,
			Group: &types.Group{},
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

func (s *Storage) GetGroupMemberships(group *types.Group) error {
	query := `
	SELECT m.joined_at, m.user_id, u.uname
	FROM memberships m
	INNER JOIN users u ON m.user_id = u.id
	WHERE m.group_id = $1`

	rows, err := s.Db.Query(query, group.ID)
	if err != nil {
		return err
	}
	group.Memberships = []types.Membership{}
	defer rows.Close()
	for rows.Next() {
		mem := types.Membership{Group: group, User: &types.User{}}
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

func (s *Storage) CreateGroup(group *types.Group) error {
	query := `INSERT INTO groups (name) VALUES ($1) RETURNING id`
	e := s.Db.QueryRow(query, &group.Name).Scan(&group.ID)
	return e
}

func (s *Storage) CreateRole(role *types.Role) error {
	query := `INSERT INTO roles (name, permissions, group_id, gets_chores) VALUES ($1,$2,$3,$4) RETURNING id`
	e := s.Db.QueryRow(query, role.Name, role.Permissions, role.Group.ID, role.GetsChores).Scan(&role.ID)
	return e
}

func (s *Storage) CreateRoleAssignment(roleID uint64, userID uint64) error {
	query := `INSERT INTO role_assignments (role_id, user_id) VALUES ($1,$2)`
	_, e := s.Db.Exec(query, roleID, userID)
	return e
}

func (s *Storage) CreateMembership(mem *types.Membership) error {
	query := `INSERT INTO memberships (joined_at, user_id, group_id) VALUES ($1,$2,$3)`
	_, e := s.Db.Exec(query, mem.JoinedAt, mem.User.ID, mem.Group.ID)
	return e
}

func (s *Storage) GetMemberChores(member *types.Membership) error {
	rows, err := s.Db.Query("SELECT chores.id, chores.description, chores.name, chores.duration, chore_assignments.complete, chore_assignments.date_assigned, chore_assignments.date_complete FROM chores WHERE chores.group_id = $1 INNER JOIN chore_assignment ON chore_assignment.chore_id = chores.id", member.Group.ID)

	if err != nil {
		return err
	}

	member.Assignments = []types.ChoreAssignment{}
	defer rows.Close()
	for rows.Next() {
		ca := types.ChoreAssignment{
			User:  member.User,
			Chore: &types.Chore{Group: member.Group},
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
	case *types.Group:
		return s.GetGroupRoles(v)
	case *types.Membership:
		return s.GetMemberRoles(v)
	default:
		return errors.ErrType
	}
}

func (s *Storage) GetGroupRoles(group *types.Group) error {
	query := `
	SELECT id, name, permissions, gets_chores
	FROM roles
	WHERE group_id = $1`
	rows, e := s.Db.Query(query, group.ID)
	if e != nil {
		return e
	}
	group.Roles = []types.Role{}
	defer rows.Close()
	for rows.Next() {
		role := types.Role{Group: group}
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

func (s *Storage) GetMemberRoles(member *types.Membership) error {
	query := `
	SELECT r.id, r.name, r.permissions, r.gets_chores
	FROM role_assignments ra
	INNER JOIN roles r on r.id = ra.role_id
	WHERE ra.user_id = $1`

	rows, e := s.Db.Query(query, member.User.ID)
	if e != nil {
		return e
	}
	member.Roles = []types.Role{}
	defer rows.Close()
	for rows.Next() {
		role := types.Role{Group: member.Group}
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

// GetSession fetches a session frm the database by session id
func (s *Storage) GetSession(ses *types.Session) error {
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
func (s *Storage) UpsertSession(ses *types.Session) error {
	statement, err := s.Db.Prepare("INSERT INTO sessions (uuid, values, created, user_id) VALUES ($1,$2,$3,$4) ON CONFLICT (uuid) DO UPDATE SET values = $2")
	if err != nil {
		return err
	}
	defer statement.Close()
	_, err = statement.Exec(ses.UUID, ses.Values, ses.Created, ses.UserID)
	return err
}
