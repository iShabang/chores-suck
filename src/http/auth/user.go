package auth

// User defines properties of a user to be authorized/authenticated
type User struct {
	ID       uint64
	Name     string
	Password string
}
