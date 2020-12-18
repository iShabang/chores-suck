package users

// Session type containing user session information
type Session struct {
	ID         uint64
	UserID     uint64
	ExpireTime int64
}

type SessionID string
