package postgre

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
	db, err := sql.Open("postgres", "dbname=choressuck sslmode=disable")
	s := &Storage{
		Db: db,
	}
	if err != nil {
		log.Fatal(err)
	}
	return s
}

// GetUserByName fetches a user from the database by unique username
func (s *Storage) GetUserByName(user *types.User) error {
	query := `
	SELECT user.user_id, user.email, user.pword, user.created_at 
	FROM users 
	WHERE user.uname = $1`
	err := s.Db.QueryRow(query, user.Username).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt)
	return err
}

func (s *Storage) GetUserByEmail(user *types.User) error {
	query := `
	SELECT user.user_id, user.uname, user.pword, user.created_at 
	FROM users 
	WHERE user.email = $1`
	err := s.Db.QueryRow(query, user.Email).Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt)
	return err
}

// GetUserByID fetches a user from the database by unique ID
func (s *Storage) GetUserByID(ID string) (types.User, error) {
	user := types.User{}
	err := s.Db.QueryRow("SELECT id, uuid, name, email, created_at FROM users WHERE id = $1", ID).Scan(&user.ID, &user.UUID, &user.Username, &user.Email, &user.CreatedAt)
	return user, err
}

// GetUserChoreList fetches a list of chore data
func (s *Storage) GetUserChoreList(userID int) ([]types.ChoreListItem, error) {
	rows, err := s.Db.Query("SELECT chore_assignments.date_due, chores.name, group.name FROM chore_assignments WHERE chore_assignments.user_id = $1 INNER JOIN chores ON chores.chore_id = chore_assignments.chore_id INNER JOIN groups ON groups.group_id = chores.group_id", userID)

	if err != nil {
		return nil, err
	}

	var chores []types.ChoreListItem
	defer rows.Close()
	for rows.Next() {
		ch := types.ChoreListItem{}
		err = rows.Scan(&ch.DateDue, &ch.ChoreName, &ch.GroupName)
		if err != nil {
			return nil, err
		}
		chores = append(chores, ch)
	}

	return chores, err
}

func (s *Storage) GetUserMemberships(user *types.User) error {
	rows, err := s.Db.Query("SELECT memberships.joined_at, memberships.group_id, groups.name FROM memberships WHERE memberships.user_id = $1 INNER JOIN groups ON groups.group_id = memberships.group_id", user.ID)

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
			return err
		}
		user.Memberships = append(user.Memberships, mem)
	}

	return nil
}

// GetMemberChores fetches the list of chores assigned to a member
func (s *Storage) GetMemberChores(member *types.Membership) error {
	rows, err := s.Db.Query("SELECT chores.chore_id, chores.description, chores.name, chores.duration, chore_assignments.complete, chore_assignments.date_assigned, chore_assignments.date_complete FROM chores WHERE chores.group_id = $1 INNER JOIN chore_assignment ON chore_assignment.chore_id = chores.chore_id", member.Group.ID)

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
func (s *Storage) GetSession(ID string) (types.Session, error) {
	p := types.Session{}
	err := s.Db.QueryRow("SELECT id, ses_id, values, created, expires, modified FROM sessions WHERE ses_id = ?", ID).Scan(&p.ID, &p.SesID, &p.Values, &p.Created)

	return p, err
}

// DeleteSession removes a session from the database
func (s *Storage) DeleteSession(ID string) error {
	statement, err := s.Db.Prepare("DELETE FROM sessions where ses_id = $1")
	if err != nil {
		return err
	}
	defer statement.Close()
	_, err = statement.Exec(ID)
	return err
}

// UpsertSession inserts or updates a session in the database. If the session does not exist
// it is created otherwise the existing session is updated.
func (s *Storage) UpsertSession(ses *types.Session) error {
	statement, err := s.Db.Prepare("INSERT INTO sessions (ses_id, values, created) VALUES ($1,$2,$3) ON CONFLICT ses_id UPDATE SET values = $2")
	if err != nil {
		return err
	}
	defer statement.Close()
	_, err = statement.Exec(ses.SesID, ses.Values, ses.Created)
	return err
}
