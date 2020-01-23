package tools

import (
	"encoding/json"
	"fmt"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type WebToken struct {
	Token  string `json:"token"`
	Expire string `json:"expire"`
}

func (h LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		h.handlePOST(w, r)
	default:
		http.Error(w, "Invlalid login command", http.StatusMethodNotAllowed)
	}
}

type LoginHandler struct {
	c *Connection
}

// FYI: This is how you do dependency injection in Go
func NewLogin(conn *Connection) *LoginHandler {
	return &LoginHandler{
		c: conn,
	}
}

func (h LoginHandler) handlePOST(w http.ResponseWriter, r *http.Request) {
	fmt.Println("starting login")
	var newUser User
	json.NewDecoder(r.Body).Decode(&newUser)
	fmt.Printf("username: %v password: %v\n", newUser.Name, newUser.Password)

	// run query for the user.
	u, err := h.c.GetUser(newUser.Name)
	fmt.Printf("id: %v\nfirst_name: %v\nlast_name: %v\nemail: %v\npassword: %v\nusername: %v\nattempts: %v\n", u.Id, u.FirstName, u.LastName, u.Email, u.Password, u.Username, u.Attempts)

	if err == ErrNotFound {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Print(ErrNotFound)
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
	np := []byte(newUser.Password)
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

	// store session id and expire time in database
	h.c.AddSession(&u, id.String(), expireTime)

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
