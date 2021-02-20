package types

import "time"

// User defines properties of a user
type User struct {
	ID          uint64
	Username    string
	Email       string
	Password    string
	CreatedAt   time.Time
	Memberships []Membership
}
