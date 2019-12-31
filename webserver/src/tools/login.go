package tools

import (
	//"golang.org/x/crypto/bcrypt"
	"encoding/json"
	"fmt"
	"net/http"
)

func (h LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.handleGET(w, r)
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

func (h LoginHandler) handleGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Status", "200")
}

func (h LoginHandler) handlePOST(w http.ResponseWriter, r *http.Request) {
	fmt.Println("starting login")
	var newUser User
	json.NewDecoder(r.Body).Decode(&newUser)
	fmt.Printf("username: %v password: %v\n", newUser.Name, newUser.Password)

	// run query for the user.
	u, err := h.c.GetUser(newUser.Name)

	if err == ErrNotFound {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Print(ErrNotFound)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("database error")
		return
	}

	fmt.Println("got user")

	if u.Attempts > 2 {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Printf("attempts = %v", u.Attempts)
		return
	}

	fmt.Println("attempts good")

	if u.Password != newUser.Password {
		h.c.UpdateUserAttempts(u.Username, u.Attempts+1)
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Print("wrong password")
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println("password good")

	if u.Attempts > 0 {
		h.c.UpdateUserAttempts(u.Username, 0)
	}

	tokenString, err, expireTime := GenToken(newUser.Name)
	if err != nil {
		fmt.Print("token failure")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println("got token")

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expireTime,
	})
	fmt.Println("login success")
}
