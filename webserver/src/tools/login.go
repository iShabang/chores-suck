package tools

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
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

	fmt.Println("password good")

	if u.Attempts > 0 {
		h.c.UpdateUserAttempts(u.Username, 0)
	}

	tokenString, err, _ := GenToken(newUser.Name)
	if err != nil {
		fmt.Print("token failure")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println("got token")

	wt := WebToken{
		Token:  tokenString,
		Expire: "test",
	}

	js, err := json.Marshal(wt)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Print(js)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

	fmt.Println("login success")
}
