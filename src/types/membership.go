package types

import (
	"time"
)

// Membership relates a User to a specific group
type Membership struct {
	JoinedAt time.Time
	UserID   int
	GroupID  int
}
