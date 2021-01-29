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

// GetUserByID fetches a user from the database by unique ID
func (s *Storage) GetUserByID(ID string) (types.User, error) {
	user := types.User{}
	err := s.Db.QueryRow("SELECT id, uuid, name, email, created_at FROM users WHERE id = $1", ID).Scan(&user.ID, &user.UUID, &user.Username, &user.Email, &user.CreatedAt)
	return user, err
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
