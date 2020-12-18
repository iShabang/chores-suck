package users

import (
	"time"
)

// Service provides functionality for user types
type Service interface {
	Authenticate(string, string) (User, error)
	Authorize(string) (Session, error)
	StartSession(User) (Session, error)
}

type service struct {
	repo Repository
}

func (s service) Authenticate(username string, password string) (User, error) {

	u, e := s.repo.GetUser(username)

	if e != nil {
		//TODO: Handle error
	}

	r := checkpword(password, u.Password)

	if !r {
		// TODO: Handle wrong password
	}

	return u, nil
}

func (s service) Authorize(sid string) (Session, error) {
	ses, e := s.repo.GetSession(sid)

	if e != nil {
		// TODO: Handle error
		// Session not found or database issue
	}

	if ses.ExpireTime <= time.Now().Unix() {
		// TODO: Handle expired session
	}

	return ses, nil
}

func (s service) StartSession(user User) (Session, error) {
	// TODO: Generate a random unique ID
	id := uint64(0)

	et := time.Now().Add(24 * 7 * time.Hour)

	se := Session{
		ID:         id,
		UserID:     user.ID,
		ExpireTime: et.Unix(),
	}

	e := s.repo.AddSession(se)
	if e != nil {
		// TODO: Handle error
	}

	return se, nil
}
