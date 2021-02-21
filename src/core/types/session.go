package types

import "time"

// Session contains properties for a session pulled from the database
type Session struct {
	UUID    string
	Values  string
	Created time.Time
	UserID  uint64
}
