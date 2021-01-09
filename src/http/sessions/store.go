package sessions

import (
	"encoding/base32"
	"net/http"
	"strings"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

// Repository defines storage functionality for sessions
type Repository interface {
	GetSession(ses *sessions.Session) error
	DeleteSession(ses *sessions.Session) error
	UpsertSession(ses *sessions.Session) error
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
			if err = s.repo.GetSession(session); err != nil {
				session.IsNew = false
			}
		}
	}
	return session, err
}

// Save should persist session to the underlying store implementation.
func (s *Store) Save(req *http.Request, w http.ResponseWriter, ses *sessions.Session) error {
	if ses.Options.MaxAge < 0 {
		// TODO: Delete session from database
		if err := s.repo.DeleteSession(ses); err != nil {
			return err
		}
		http.SetCookie(w, sessions.NewCookie(ses.Name(), "", ses.Options))
		return nil
	}

	if ses.ID == "" {
		// Because the ID is used in the filename, encode it to
		// use alphanumeric characters only.
		ses.ID = strings.TrimRight(
			base32.StdEncoding.EncodeToString(
				securecookie.GenerateRandomKey(32)), "=")
	}

	if err := s.repo.UpsertSession(ses); err != nil {
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
