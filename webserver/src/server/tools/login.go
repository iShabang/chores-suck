package tools

import (
	"encoding/json"
	"fmt"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"server/tools/database"
	"time"
)

/********************************************************
TYPES
********************************************************/
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginHandler struct {
	c *db.Connection
}

/********************************************************
INITIALIZER
********************************************************/
// FYI: This is how you do dependency injection in Go
func NewLogin(conn *db.Connection) *LoginHandler {
	return &LoginHandler{
		c: conn,
	}
}

/********************************************************
HTTP
********************************************************/
func (h LoginHandler) Login(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Invlalid login command", http.StatusMethodNotAllowed)
		return
	}

	var cred Credentials
	json.NewDecoder(r.Body).Decode(&cred)

	u, err := h.c.GetUser(cred.Username)

	if err == db.ErrNotFound {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Print(db.ErrNotFound)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("database error")
		return
	}

	if u.Attempts > 2 {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Printf("attempts = %v", u.Attempts)
		return
	}

	hp := []byte(u.Password)
	np := []byte(cred.Password)
	err = bcrypt.CompareHashAndPassword(hp, np)

	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			h.c.UpdateUserAttempts(u.Username, u.Attempts+1)
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Print("wrong password")
		} else {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	if u.Attempts > 0 {
		h.c.UpdateUserAttempts(u.Username, 0)
	}

	// Generate a session id
	id := xid.New()

	// Calculate expire time
	expireTime := time.Now().Add(24 * 7 * time.Hour)

	session := db.Session{
		SessionId:  id.String(),
		UserId:     u.Id,
		ExpireTime: fmt.Sprintf("%v", expireTime.Unix()),
	}

	// store session id and expire time in database
	h.c.AddSession(&session)

	// store id in a cookie
	cookie := http.Cookie{
		Name:     "session",
		Value:    id.String(),
		Expires:  expireTime,
		Secure:   true,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	// send the response
	w.WriteHeader(http.StatusOK)
}

func (h LoginHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid login command", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session")
	result := (err != nil)

	if result {
		err = h.c.DeleteSession(cookie.Value)
		result = (err != nil)
	}

	if result {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}
