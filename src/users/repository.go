package users

// Repository functionality for user database implementations
type Repository interface {
	Add(User) error
	Update(User, string, string) error
	Get(string) (User, error)
}
