package tools

import (
	"net/http"
	"server/tools/database"
	"strconv"
	"time"
)

/********************************************************
TYPES
********************************************************/

type AuthHandler struct {
	conn *db.Connection
}

/********************************************************
INITIALIZER
********************************************************/

func NewAuthHandler(c *db.Connection) *AuthHandler {
	return &AuthHandler{
		conn: c,
	}
}

/********************************************************
EXPORTED METHODS
********************************************************/
func (h *AuthHandler) AuthorizeRequest(r *http.Request) (string, error) {
	result := true
	cookie, err := r.Cookie("session")
	result = (err != nil)

	var sess *db.Session
	if result {
		sess, err = h.conn.FindSession(cookie.Value)
		result = (err != nil) && (sess.SessionId != "" && sess.UserId != "")
	}

	var expTime int64
	if result {
		expTime, err := strconv.Atoi(sess.ExpireTime)
		result = (err != nil) && (expTime > 0)
	}

	if result {
		currentTime := time.Now().Unix()
		result = (expTime > currentTime)
	}

	userId := ""
	if result {
		userId = sess.UserId
	}

	return userId, err
}
