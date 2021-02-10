package types

import "time"

// Session contains properties for a session pulled from the database
type Session struct {
	ID      int
	SesID   string
	Values  string
	Created time.Time
}
