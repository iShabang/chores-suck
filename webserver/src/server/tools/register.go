package tools

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"server/tools/database"
)

/********************************************************
TYPES
********************************************************/
type RegisterHandler struct {
	c *db.Connection
}

type NewUser struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Username  string `json:"username"`
}

/********************************************************
INITIALIZER
********************************************************/
func NewRegister(c *db.Connection) *RegisterHandler {
	return &RegisterHandler{
		c: c,
	}
}

/********************************************************
HTTP HANDLERS
********************************************************/
func (h RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		h.handlePOST(w, r)
	default:
		http.Error(w, "Invlalid register command", http.StatusMethodNotAllowed)
	}
}

func (h RegisterHandler) handlePOST(w http.ResponseWriter, r *http.Request) {
	var nu NewUser
	json.NewDecoder(r.Body).Decode(&nu)

	// Ideally there should be another backend script that would take care of this before hand.
	if nu.FirstName == "" || nu.LastName == "" || nu.Email == "" || nu.Password == "" || nu.Username == "" {
		fmt.Println("fields missing")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Encrypt the password before saving to the database
	// There should be another backend script that validates the password before this point
	p := []byte(nu.Password)
	hp, err := bcrypt.GenerateFromPassword(p, 10)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// We have to convert between types since mongo uses bson and we expect json from the client
	ul := db.User{
		FirstName: nu.FirstName,
		LastName:  nu.LastName,
		Email:     nu.Email,
		Password:  string(hp),
		Username:  nu.Username,
		Attempts:  0,
	}

	// Add user to the database
	_, err = h.c.AddUser(&ul)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
