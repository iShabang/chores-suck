package tools

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type RegisterHandler struct {
	c *Connection
}

type NewUser struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Username  string `json:"username"`
}

func NewRegister(c *Connection) *RegisterHandler {
	return &RegisterHandler{
		c: c,
	}
}

func (h RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		h.handlePOST(w, r)
	default:
		http.Error(w, "Invlalid login command", http.StatusMethodNotAllowed)
	}
}

func (h RegisterHandler) handlePOST(w http.ResponseWriter, r *http.Request) {
	var nu NewUser
	json.NewDecoder(r.Body).Decode(&nu)
	fmt.Println(nu)
	if nu.FirstName == "" || nu.LastName == "" || nu.Email == "" || nu.Password == "" || nu.Username == "" {
		fmt.Println("fields missing")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	p := []byte(nu.Password)
	hp, err := bcrypt.GenerateFromPassword(p, 10)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	ul := User{
		FirstName: nu.FirstName,
		LastName:  nu.LastName,
		Email:     nu.Email,
		Password:  string(hp),
		Username:  nu.Username,
		Attempts:  0,
	}

	id, err := h.c.AddUser(&ul)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Printf("Register: got an id %v\n", id)

	w.WriteHeader(http.StatusOK)
}
