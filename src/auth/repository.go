package auth

// Repository defines functionality for a session repository
type Repository interface {
	Add(Session) error
	Get(ID string) (Session, error)
	Delete(ID string) error
	GC() error
}
