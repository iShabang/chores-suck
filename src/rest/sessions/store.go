package sessions

import (
	"errors"
	"log"
	"net/http"
	"time"

	se "chores-suck/core/storage/errors"
	"chores-suck/core/types"

	"github.com/google/uuid"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

// Repository defines storage functionality for sessions
type Repository interface {
	GetSession(ses *types.Session) error
	DeleteSession(ID string) error
	UpsertSession(ses *types.Session) error
}

// Store defines properties of a session store
type Store struct {
	codecs []securecookie.Codec
	repo   Repository
	opts   *sessions.Options
}

// NewStore initializes a new store
func NewStore(rep Repository, keyPairs ...[]byte) *Store {
	store := &Store{
		codecs: securecookie.CodecsFromPairs(keyPairs...),
		repo:   rep,
		opts: &sessions.Options{
			Path:   "/",
			MaxAge: 86400 * 30,
		},
	}
	return store
}

// Get should return a cached session.
func (s *Store) Get(req *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(req).Get(s, name)
}

// New should create and return a new session.
//
// Note that New should never return a nil session, even in the case of
// an error if using the Registry infrastructure to cache the session.
func (s *Store) New(req *http.Request, name string) (*sessions.Session, error) {
	session := sessions.NewSession(s, name)
	opts := *s.opts
	session.Options = &opts
	session.IsNew = true
	var err error
	if cookie, errCookie := req.Cookie(name); errCookie == nil {
		if err = securecookie.DecodeMulti(name, cookie.Value, &session.ID, s.codecs...); err == nil {
			var ts = types.Session{}
			ts.UUID = session.ID
			if err := s.repo.GetSession(&ts); err == nil {
				session.IsNew = false
				err = securecookie.DecodeMulti(name, ts.Values, &session.Values, s.codecs...)
			} else if err != se.ErrNotFound {
				return session, err
			}
		}
	}
	return session, err
}

// Save should persist session to the underlying store implementation.
func (s *Store) Save(req *http.Request, w http.ResponseWriter, ses *sessions.Session) error {
	if ses.Options.MaxAge < 0 {
		if err := s.repo.DeleteSession(ses.ID); err != nil {
			return err
		}
		http.SetCookie(w, sessions.NewCookie(ses.Name(), "", ses.Options))
		return nil
	}

	if ses.ID == "" {
		uuid, err := uuid.NewRandom()
		if err != nil {
			return err
		}
		ses.ID = uuid.String()
	}

	var err error
	ts := types.Session{}
	ts.UUID = ses.ID
	ts.Created = time.Now().UTC()
	if ts.Values, err = securecookie.EncodeMulti(ses.Name(), ses.Values, s.codecs...); err != nil {
		return err
	}
	var userID uint64
	userID, ok := ses.Values["userid"].(uint64)
	if !ok || userID == 0 {
		log.Printf("Store.Save: user id failure. ok: %v, id: %v", ok, userID)
		return errors.New("Store.Save: Failed to get user id from session")
	}
	ts.UserID = userID

	if err = s.repo.UpsertSession(&ts); err != nil {
		return err
	}

	encoded, err := securecookie.EncodeMulti(ses.Name(), ses.ID, s.codecs...)
	if err != nil {
		return err
	}

	http.SetCookie(w, sessions.NewCookie(ses.Name(), encoded, ses.Options))
	return nil
}

// MaxAge sets the maximum age for the store and the underlying cookie
// implementation. Individual sessions can be deleted by setting Options.MaxAge
// = -1 for that session.
func (s *Store) MaxAge(age int) {
	s.opts.MaxAge = age

	// Set the maxAge for each securecookie instance.
	for _, codec := range s.codecs {
		if sc, ok := codec.(*securecookie.SecureCookie); ok {
			sc.MaxAge(age)
		}
	}
}
