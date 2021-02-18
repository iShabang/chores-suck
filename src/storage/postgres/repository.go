package postgres

import (
	"chores-suck/types"
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
	return err
}

func (s *Storage) GetUserByEmail(user *types.User) error {
	query := `
	SELECT users.id, users.uname, users.pword, users.created_at 
	FROM users 
	WHERE users.email = $1`
	err := s.Db.QueryRow(query, user.Email).Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt)
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

// GetMemberChores fetches the list of chores assigned to a member
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

// GetSession fetches a session frm the database by session id
func (s *Storage) GetSession(ses *types.Session) error {
	err := s.Db.QueryRow("SELECT values, created FROM sessions WHERE uuid = $1", ses.UUID).Scan(&ses.Values, &ses.Created)

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
	statement, err := s.Db.Prepare("INSERT INTO sessions (uuid, values, created) VALUES ($1,$2,$3) ON CONFLICT (uuid) DO UPDATE SET values = $2")
	if err != nil {
		return err
	}
	defer statement.Close()
	_, err = statement.Exec(ses.UUID, ses.Values, ses.Created)
	return err
}
