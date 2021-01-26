package types

import (
	"time"
)

// Membership relates a User to a specific group
type Membership struct {
	ID         string
	DateJoined time.Time
	UserID     string
	GroupID    string
}
