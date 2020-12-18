package users

// Repository functionality for user database implementations
type Repository interface {
	AddUser(User) error
	AddSession(Session) error
	Update(User, string, string) error
	GetUser(string) (User, error)
	GetSession(string) (Session, error)
}
