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
	var nu UserLarge
	json.NewDecoder(r.Body).Decode(&nu)
	if nu.FirstName == "" || nu.LastName == "" || nu.Email == "" || nu.Password == "" || nu.Username == "" {
		fmt.Println("fields missing")
		return
	}

	p := []byte(nu.Password)
	hp, err := bcrypt.GenerateFromPassword(p, 10)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	nu.Password = string(hp)
	_, err = h.c.AddUser(&nu)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
