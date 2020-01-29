package tools

import (
	"fmt"
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
	fmt.Println("call to authorize")
	result := true
	cookie, err := r.Cookie("session")
	result = (err == nil) && (cookie != nil)

	if err != nil {
		fmt.Println(err)
	} else if cookie == nil {
		fmt.Println("cookie empty")
	}

	var sess *db.Session
	if result {
		fmt.Printf("got cookie. sid %v\n", cookie.Value)
		sess, err = h.conn.FindSession(cookie.Value)
		result = (err == nil) && (sess != nil)
	}

	if result {
		fmt.Printf("got session. user: %v\n", sess.UserId)
		result = (sess.SessionId != "" && sess.UserId != "")
	}

	var expTime int64
	if result {
		expTime, err := strconv.Atoi(sess.ExpireTime)
		result = (err == nil) && (expTime > 0)
	}

	if result {
		fmt.Printf("get expire time %v\n", expTime)
		currentTime := time.Now().Unix()
		result = (expTime > currentTime)
	}

	userId := ""
	if result {
		fmt.Println("auth success")
		userId = sess.UserId
	}

	return userId, err
}
